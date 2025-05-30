package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	userID := snowflake.GenID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
