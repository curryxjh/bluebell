package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"time"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票值多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostTimeZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(context.Background(), getRedisKey(KeyPostScoreZSet), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(context.Background(), cKey, postID)
	_, err := pipeline.Exec(context.Background())
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	// 1.判断投票限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(context.Background(), getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 2, 3需要放到一个pipeline事务中操作
	// 2.更新帖子分数
	// 先查当前用户给当前帖子的投票记录
	ov := client.ZScore(context.Background(), getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 更新: 如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(value - ov) // 计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(context.Background(), getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)
	// 3.记录用户为该帖子投票的数据
	if value == 0 {
		pipeline.ZRem(context.Background(), getRedisKey(KeyPostVotedZSetPF+postID), userID)
	} else {
		pipeline.ZAdd(context.Background(), getRedisKey(KeyPostScoreZSet+postID), &redis.Z{
			Score:  value,
			Member: userID,
		})
	}
	_, err := pipeline.Exec(context.Background())
	return err
}
