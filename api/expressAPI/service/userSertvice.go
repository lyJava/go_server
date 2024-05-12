package service

import (
	"apiProject/api/expressAPI/types/domain"
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetUserById(id int64) (*domain.User, error)
	// CreateUser 创建用户
	CreateUser(user *domain.User) (*domain.User, error)
	// UserLogin 用户登录
	UserLogin(user *domain.User) (*domain.User, error)
}
