package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "xjh021022.gg",
		DB:           "bluebell",
		Port:         3307,
		MaxIdleConns: 10,
		MaxOpenConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
	}
	panic(err)
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    12,
		CommunityID: 1,
		Title:       "title",
		Content:     "content",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("create post success")
}
