package mysql

import "bluebell/models"

func CreatePost(p *models.Post) (err error) {
	sqlStr := "INSERT INTO post (post_id, title, content, author_id, community_id) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}
