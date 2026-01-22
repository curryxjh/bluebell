package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateCommentHandler(c *gin.Context) {

	p := new(models.ParamCreateComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	if err := logic.CreateComment(userID, p); err != nil {
		zap.L().Error("logic.CreateComment failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, nil)
}

func GetCommentsByPostIDHandler(c *gin.Context) {
	postIDStr := c.Param("id")
	if postIDStr == "" {
		ResponseError(c, CodeInvalidParam)
		return
	}

	comments, err := logic.GetCommentsByPostID(postIDStr)
	if err != nil {
		zap.L().Error("logic.GetCommentsByPostID failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, comments)
}

func DeleteCommentHandler(c *gin.Context) {
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 删除评论
	if err := logic.DeleteComment(commentID, userID); err != nil {
		zap.L().Error("logic.DeleteComment failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, nil)
}
