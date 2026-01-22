package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"errors"
)

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

// GetCommunityListWithMemberCount 获取社区列表（包含成员数量）
func GetCommunityListWithMemberCount() ([]*models.CommunityWithMembers, error) {
	return mysql.GetCommunityListWithMemberCount()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}

// JoinCommunity 加入社区
func JoinCommunity(userID, communityID int64) error {
	// 1. 检查社区是否存在
	_, err := mysql.GetCommunityDetailByID(communityID)
	if err != nil {
		return errors.New("社区不存在")
	}

	// 2. 检查用户是否已加入该社区
	joined, err := mysql.CheckUserInCommunity(userID, communityID)
	if err != nil {
		return err
	}
	if joined {
		return errors.New("已经加入该社区")
	}

	// 3. 加入社区
	return mysql.JoinCommunity(userID, communityID)
}

// LeaveCommunity 退出社区
func LeaveCommunity(userID, communityID int64) error {
	// 检查用户是否已加入该社区
	joined, err := mysql.CheckUserInCommunity(userID, communityID)
	if err != nil {
		return err
	}
	if !joined {
		return errors.New("未加入该社区")
	}

	return mysql.LeaveCommunity(userID, communityID)
}

// GetUserCommunities 获取用户加入的所有社区
func GetUserCommunities(userID int64) ([]*models.Community, error) {
	return mysql.GetUserCommunities(userID)
}
