package service

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
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
	// SelectPage 分页查询
	//
	// 参数:
	//     param: 包含分页查询的条件参数，类型为 *param.TestUserPageParam，包含以下字段：
	//     - Page: 当前页码
	//     - Size: 每页条数
	//     - Column: 排序的列（可选）
	//     - Order: 排序方式（可选）
	//     - Username: 用户名（可选）
	//     - Email: 邮箱（可选）
	//     - Birthday: 出生日期（可选）
	//     - Phone: 电话号码（可选）
	//     - Address: 地址（可选）
	//
	// 返回:
	//     -[]*domain.TestUser: 符合条件的用户列表，包含每个用户的详细信息
	//     -int64: 查询结果的总记录数
	//     -int64: 总页数，基于查询的总记录数和每页条数计算
	//     -error: 如果执行过程中出现任何错误，则返回相应的错误信息。如果没有错误，返回nil
	SelectPage(param *param.TestUserPageParam) ([]*domain.TestUser, int64, int64, error)
}
