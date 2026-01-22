package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

var (
	ErrorInvalidParam = mysql.ErrorInvalidParam
	ErrorNoPermission = mysql.ErrorNoPermission
)

// CreateComment 创建评论
func CreateComment(authorID int64, p *models.ParamCreateComment) (err error) {
	// 转换post_id
	postID, err := strconv.ParseInt(p.PostID, 10, 64)
	if err != nil {
		zap.L().Error("invalid post_id", zap.Error(err))
		return ErrorInvalidParam
	}

	// 生成评论ID
	commentID := snowflake.GenID()

	comment := &models.Comment{
		ID:       commentID,
		PostID:   postID,
		AuthorID: authorID,
		Content:  p.Content,
		ParentID: p.ParentID,
		Status:   1,
	}

	// 保存到数据库
	return mysql.CreateComment(comment)
}

// GetCommentsByPostID 获取帖子的评论列表
func GetCommentsByPostID(postIDStr string) (comments []*models.CommentDetail, err error) {
	// 转换post_id
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("invalid post_id", zap.Error(err))
		return nil, ErrorInvalidParam
	}

	// 从数据库查询
	comments, err = mysql.GetCommentsByPostID(postID)
	if err != nil {
		zap.L().Error("mysql.GetCommentsByPostID failed", zap.Error(err))
		return nil, err
	}

	return comments, nil
}

// DeleteComment 删除评论
func DeleteComment(commentID int64, userID int64) (err error) {
	// 获取评论信息
	comment, err := mysql.GetCommentByID(commentID)
	if err != nil {
		zap.L().Error("mysql.GetCommentByID failed", zap.Error(err))
		return err
	}

	// 检查权限（只有作者本人可以删除）
	if comment.AuthorID != userID {
		return ErrorNoPermission
	}

	// 软删除
	return mysql.DeleteComment(commentID)
}
