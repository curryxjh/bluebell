package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"
	"go.uber.org/zap"
	"strconv"
)

func CreatePost(p *models.Post) (err error) {
	p.ID = snowflake.GenID()
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
}

func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}

	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}

	// 获取帖子投票数
	var voteCount int64 = 0
	voteNum, err := redis.GetPostVoteData([]string{strconv.FormatInt(pid, 10)})
	if err != nil {
		zap.L().Warn("redis.GetPostVoteData failed, using default 0", zap.Error(err))
	} else if len(voteNum) > 0 {
		voteCount = voteNum[0]
	}

	zap.L().Info("GetPostById success", 
		zap.Int64("post_id", pid), 
		zap.String("title", post.Title),
		zap.Int64("vote_count", voteCount))

	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		VoteNum:         voteCount,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return nil, err
	}
	
	if len(posts) == 0 {
		zap.L().Info("No posts found in database")
		return []*models.ApiPostDetail{}, nil
	}
	
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}

		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         0, // 从MySQL获取时，投票数默认为0
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	
	zap.L().Info("GetPostList success", zap.Int("count", len(data)))
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 去redis中查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil || len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data, fallback to mysql", zap.Error(err))
		// Redis中没有数据，直接从MySQL获取
		return GetPostList(p.Page, p.Size)
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 根据id去mysql查询帖子的详细信息
	// 返回数据按照给定id的顺序返回
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetPostByIDs failed", zap.Error(err))
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	
	if len(posts) == 0 {
		zap.L().Warn("GetPostByIDs return 0 posts")
		return []*models.ApiPostDetail{}, nil
	}
	
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData failed", zap.Error(err))
		// 如果获取投票数失败，使用0填充
		voteData = make([]int64, len(posts))
	}
	
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		
		voteNum := int64(0)
		if idx < len(voteData) {
			voteNum = voteData[idx]
		}
		
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteNum,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		if len(ids) == 0 {
			zap.L().Warn("redis.GetCommunityPostIDsInOrder(p) return 0 data", zap.Error(err))
			return
		}
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))
	// 根据id去mysql查询帖子的详细信息
	// 返回数据按照给定id的顺序返回
	posts, err := mysql.GetPostByIDs(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts", posts))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	for idx, post := range posts {
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {
		// 查所有的
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("logic.GetCommunityPostList failed", zap.Error(err))
		return nil, err
	}
	return
}
