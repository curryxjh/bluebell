package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `INSERT INTO comment(post_id, author_id, content, parent_id, status)
				VALUES (?, ?, ?, ?, 1)`
	result, err := db.Exec(sqlStr, comment.PostID, comment.AuthorID, comment.Content, comment.ParentID)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		zap.L().Error("get comment id failed", zap.Error(err))
		return err
	}
	comment.ID = id
	return nil
}

// GetCommentsByPostID 根据帖子ID获取评论列表
func GetCommentsByPostID(postID int64) (comments []*models.CommentDetail, err error) {
	sqlStr := `SELECT 
		c.id, c.post_id, c.author_id, c.content, c.parent_id, c.status, c.create_time,
		u.username as author_name
	FROM comment c
	LEFT JOIN user u ON c.author_id = u.user_id
	WHERE c.post_id = ? AND c.status = 1
	ORDER BY c.create_time DESC`

	if err = db.Select(&comments, sqlStr, postID); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no comments found for post", zap.Int64("post_id", postID))
			return []*models.CommentDetail{}, nil
		}
		zap.L().Error("query comments failed", zap.Error(err))
		return nil, err
	}
	return comments, nil
}

// GetCommentByID 根据ID获取评论
func GetCommentByID(commentID int64) (comment *models.Comment, err error) {
	comment = new(models.Comment)
	sqlStr := `SELECT id, post_id, author_id, content, parent_id, status, create_time
				FROM comment
				WHERE id = ?`
	if err = db.Get(comment, sqlStr, commentID); err != nil {
		zap.L().Error("query comment failed", zap.Error(err))
		return nil, err
	}
	return comment, nil
}

// DeleteComment 删除评论（软删除）
func DeleteComment(commentID int64) (err error) {
	sqlStr := `UPDATE comment SET status = 0 WHERE id = ?`
	_, err = db.Exec(sqlStr, commentID)
	if err != nil {
		zap.L().Error("delete comment failed", zap.Error(err))
		return err
	}
	return nil
}

// GetCommentCount 获取帖子的评论数量
func GetCommentCount(postID int64) (count int64, err error) {
	sqlStr := `SELECT COUNT(*) FROM comment WHERE post_id = ? AND status = 1`
	err = db.Get(&count, sqlStr, postID)
	if err != nil {
		zap.L().Error("get comment count failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}
