package models

import "time"

// UserCommunity 用户社区关系
type UserCommunity struct {
	ID          int64     `json:"id" db:"id"`
	UserID      int64     `json:"user_id" db:"user_id"`
	CommunityID int64     `json:"community_id" db:"community_id"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// ParamJoinCommunity 加入社区请求参数
type ParamJoinCommunity struct {
	CommunityID int64 `json:"community_id" binding:"required"`
}
