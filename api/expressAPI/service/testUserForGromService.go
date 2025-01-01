package service

import (
	"apiProject/api/expressAPI/types/domain"
)

// TestUserForGormService 用户服务接口（gorm方式）
type TestUserForGormService interface {
	GetUserById(id int64) (*domain.TestUser, error)
	// CreateUser 用户创建
	CreateUser(user *domain.TestUser) (*domain.TestUser, error)
	// BatchCreateUser 用户批量创建
	BatchCreateUser(list []*domain.TestUser) (int64, error)
	// UserLogin 用户登录
	UserLogin(user *domain.TestUser) (*domain.TestUser, error)
	// UpdateUser 用户修改
	UpdateUser(user *domain.TestUser) (int64, error)
	// BatchUpdateUser 用户批量修改
	BatchUpdateUser(list []*domain.TestUser) (int64, error)
	// DeleteUser 用户删除
	DeleteUser(id int64) (int64, error)
	// BatchDeleteUser 用户批量删除
	BatchDeleteUser(ids []any) (int64, error)
}
