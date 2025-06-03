package mysql

import (
	"bluebell/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := "INSERT INTO post (post_id, title, content, author_id, community_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostById(pid int64) (p *models.Post, err error) {
	p = new(models.Post)
	sqlStr := "SELECT post_id, title, content, author_id, community_id, create_time FROM post WHERE post_id = ?"
	err = db.Get(p, sqlStr, pid)
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select
	post_id, title, content, author_id, community_id, create_time
	from post
	order by create_time desc
	limit ?, ?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

func GetPostByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)

	err = db.Select(&postList, query, args...)
	return
}
