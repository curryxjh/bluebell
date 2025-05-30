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
