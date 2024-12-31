package service

import (
	"apiProject/api/expressAPI/types/domain"
)

// TestUserForGormService 用户服务接口（gorm方式）
type TestUserForGormService interface {
	GetUserById(id int64) (*domain.TestUser, error)
	// CreateUser 用户创建
	CreateUser(user *domain.TestUser) (*domain.TestUser, error)
	// UserLogin 用户登录
	UserLogin(user *domain.TestUser) (*domain.TestUser, error)
	// UpdateUser 用户修改
	UpdateUser(user *domain.TestUser) (int64, error)
}
