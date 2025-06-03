package redis

// redis keys 注意使用命名空间的方式，方便查询和拆分

const (
	Prefix             = "bluebell:"  // 项目key前缀
	KeyPostTimeZSet    = "post:time"  // zset；帖子及发布时间
	KeyPostScoreZSet   = "post:score" // zset；帖子及投票分数
	KeyPostVotedZSetPF = "post:voted" // zset；记录用户及投票类型；参数是post id

	KeyCommunitySetPF = "community:" // set保存每个分区下的帖子id
)

func getRedisKey(key string) string {
	return Prefix + key
}
