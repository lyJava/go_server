package service

import (
	"apiProject/api/expressAPI/types"
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetUserById(id int64) (*types.User, error)
	// CreatUser 创建用户
	CreatUser(user *types.User) (*types.User, error)
	// UserLogin 用户登录
	UserLogin(user *types.User) (*types.User, error)
}
