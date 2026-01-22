package mysql

import (
	"bluebell/models"
	"database/sql"
	"go.uber.org/zap"
)

// JoinCommunity 用户加入社区
func JoinCommunity(userID, communityID int64) error {
	sqlStr := "INSERT INTO user_community (user_id, community_id) VALUES (?, ?)"
	_, err := db.Exec(sqlStr, userID, communityID)
	if err != nil {
		zap.L().Error("insert user_community failed", zap.Error(err))
		return err
	}
	return nil
}

// CheckUserInCommunity 检查用户是否已加入社区
func CheckUserInCommunity(userID, communityID int64) (bool, error) {
	sqlStr := "SELECT id FROM user_community WHERE user_id = ? AND community_id = ?"
	var id int64
	err := db.Get(&id, sqlStr, userID, communityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		zap.L().Error("check user in community failed", zap.Error(err))
		return false, err
	}
	return true, nil
}

// LeaveCommunity 用户退出社区
func LeaveCommunity(userID, communityID int64) error {
	sqlStr := "DELETE FROM user_community WHERE user_id = ? AND community_id = ?"
	_, err := db.Exec(sqlStr, userID, communityID)
	if err != nil {
		zap.L().Error("delete user_community failed", zap.Error(err))
		return err
	}
	return nil
}

// GetUserCommunities 获取用户加入的所有社区
func GetUserCommunities(userID int64) ([]*models.Community, error) {
	sqlStr := `
		SELECT c.community_id, c.community_name 
		FROM community c
		INNER JOIN user_community uc ON c.community_id = uc.community_id
		WHERE uc.user_id = ?
		ORDER BY uc.create_time DESC
	`
	var communities []*models.Community
	err := db.Select(&communities, sqlStr, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.Community{}, nil
		}
		zap.L().Error("get user communities failed", zap.Error(err))
		return nil, err
	}
	return communities, nil
}

// GetCommunityMemberCount 获取社区成员数量
func GetCommunityMemberCount(communityID int64) (int64, error) {
	sqlStr := "SELECT COUNT(*) FROM user_community WHERE community_id = ?"
	var count int64
	err := db.Get(&count, sqlStr, communityID)
	if err != nil {
		zap.L().Error("get community member count failed", zap.Error(err))
		return 0, err
	}
	return count, nil
}
