package mysql

import "errors"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
	ErrorInvalidParam    = errors.New("请求参数错误")
	ErrorNoPermission    = errors.New("没有操作权限")
)
