package service

import (
	"apiProject/api/expressAPI/types/domain"
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetUserById(id int64) (*domain.User, error)
	// CreatUser 创建用户
	CreatUser(user *domain.User) (*domain.User, error)
	// UserLogin 用户登录
	UserLogin(user *domain.User) (*domain.User, error)
}
