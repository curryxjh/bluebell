package redis

import (
	"bluebell/models"
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	return client.ZRevRange(context.Background(), key, start, end).Result()
}

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从redis获取id
	// 根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	// 确定查询的索引起点
	return getIDsFormKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPF + id)
		pipeline.ZCount(context.Background(), key, "1", "1")
	}
	cmders, err := pipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}
	data = make([]int64, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}

	// 使用zinterstore把分区的帖子set和帖子分数的zset生成一个新的zset
	// 针对新的zset按之前的逻辑取数据

	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + strconv.Itoa(int(p.CommunityID))

	if client.Exists(context.Background(), key).Val() < 1 {
		// 不存在
		pipeline := client.Pipeline()
		pipeline.ZInterStore(context.Background(), key, &redis.ZStore{
			Keys:      []string{cKey, orderKey},
			Aggregate: "Max",
		})
		pipeline.Expire(context.Background(), key, time.Second*60)
		_, err := pipeline.Exec(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
