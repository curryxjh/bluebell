package mysql

import (
	"bluebell/models"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "SELECT community_id, community_name FROM community"
	if err := db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no community found")
			err = nil
		}
	}
	return
}

// GetCommunityListWithMemberCount 获取社区列表（包含成员数量）
func GetCommunityListWithMemberCount() ([]*models.CommunityWithMembers, error) {
	sqlStr := `
		SELECT 
			c.community_id, 
			c.community_name,
			COUNT(DISTINCT uc.user_id) as member_count
		FROM community c
		LEFT JOIN user_community uc ON c.community_id = uc.community_id
		GROUP BY c.community_id, c.community_name
		ORDER BY member_count DESC
	`
	var communityList []*models.CommunityWithMembers
	if err := db.Select(&communityList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("no community found")
			return []*models.CommunityWithMembers{}, nil
		}
		return nil, err
	}
	return communityList, nil
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := "SELECT community_id, community_name, introduction, create_time FROM community WHERE community_id = ?"
	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			return nil, err
		}
	}
	fmt.Println(community.ID, community.Name, community.Introduction, community.CreateTime)

	return community, nil
}
