package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func CommunityHandler(c *gin.Context) {
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// CommunityListWithMembersHandler 获取社区列表（包含成员数量）
func CommunityListWithMembersHandler(c *gin.Context) {
	data, err := logic.GetCommunityListWithMemberCount()
	if err != nil {
		zap.L().Error("logic.GetCommunityListWithMemberCount failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

func CommunityDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(id)

	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}
	ResponseSuccess(c, data)
}

// JoinCommunityHandler 加入社区
func JoinCommunityHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamJoinCommunity)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("JoinCommunityHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 3. 加入社区
	if err := logic.JoinCommunity(userID, p.CommunityID); err != nil {
		zap.L().Error("logic.JoinCommunity failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServeBusy, err.Error())
		return
	}

	ResponseSuccess(c, nil)
}

// LeaveCommunityHandler 退出社区
func LeaveCommunityHandler(c *gin.Context) {
	// 1. 获取参数
	p := new(models.ParamJoinCommunity)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LeaveCommunityHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2. 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 3. 退出社区
	if err := logic.LeaveCommunity(userID, p.CommunityID); err != nil {
		zap.L().Error("logic.LeaveCommunity failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeServeBusy, err.Error())
		return
	}

	ResponseSuccess(c, nil)
}

// GetUserCommunitiesHandler 获取用户加入的所有社区
func GetUserCommunitiesHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 获取用户加入的社区列表
	data, err := logic.GetUserCommunities(userID)
	if err != nil {
		zap.L().Error("logic.GetUserCommunities failed", zap.Error(err))
		ResponseError(c, CodeServeBusy)
		return
	}

	ResponseSuccess(c, data)
}
