package models

import "time"

// Comment 评论结构
type Comment struct {
	ID         int64     `json:"id" db:"id"`
	PostID     int64     `json:"post_id" db:"post_id"`
	AuthorID   int64     `json:"author_id" db:"author_id"`
	Content    string    `json:"content" db:"content"`
	ParentID   int64     `json:"parent_id" db:"parent_id"`
	Status     int8      `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

// CommentDetail 评论详情（包含作者信息）
type CommentDetail struct {
	ID         int64     `json:"id" db:"id"`
	PostID     int64     `json:"post_id" db:"post_id"`
	AuthorID   int64     `json:"author_id" db:"author_id"`
	AuthorName string    `json:"author_name" db:"author_name"`
	Content    string    `json:"content" db:"content"`
	ParentID   int64     `json:"parent_id" db:"parent_id"`
	Status     int8      `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

// ParamCreateComment 创建评论的参数
type ParamCreateComment struct {
	PostID   string `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required"`
	ParentID int64  `json:"parent_id"`
}
