package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/utils"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TestUserGormDB struct {
	Db *gorm.DB
}

func NewTestUserGormDB(db *gorm.DB) *TestUserGormDB {
	return &TestUserGormDB{
		Db: db,
	}
}

func (u *TestUserGormDB) GetUserById(id int64) (*domain.TestUser, error) {
	var user domain.TestUser
	if err := u.Db.First(&user, id).Error; err != nil {
		zap.L().Sugar().Errorf("GetUserById error: %+v", err)
		return nil, errors.New("获取用户失败")
	}

	dateFormat, err := utils.DateFormat(user.Birthday)
	if err != nil {
		return nil, errors.New("用户出生日期错误")
	}
	user.Birthday = dateFormat
	return &user, nil
}

func (u *TestUserGormDB) CreateUser(user *domain.TestUser) (*domain.TestUser, error) {
	// GORM 会自动使用事务，确保插入操作的原子性
	if err := u.Db.Create(user).Error; err != nil {
		zap.L().Sugar().Errorf("CreateUser error: %+v", err)
		return nil, errors.New("创建用户失败")
	}
	return u.GetUserById(user.Id)
}

func (u *TestUserGormDB) UserLogin(user *domain.TestUser) (*domain.TestUser, error) {
	var loginUser domain.TestUser
	if err := u.Db.Where("username = ?", user.Username).First(&loginUser).Error; err != nil {
		zap.L().Sugar().Errorf("UserLogin Search By username error: %+v", err)
		return nil, errors.New("当前用户不存在")
	}

	if err := u.Db.Where("username = ? AND password = ?", user.Username, user.Password).First(&loginUser).Error; err != nil {
		zap.L().Sugar().Errorf("UserLogin Search By username And Password error: %+v", err)
		return nil, errors.New("当前用户名与密码不匹配")
	}

	return &loginUser, nil
}

func (u *TestUserGormDB) UpdateUser(user *domain.TestUser) (int64, error) {
	tx := u.Db.Debug().Model(&domain.TestUser{}).Where("id = ?", user.Id).Updates(&domain.TestUser{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Birthday: user.Birthday,
		Phone:    user.Phone,
		Address:  user.Address,
	},
	)

	// 检查是否发生了错误
	if tx.Error != nil {
		zap.L().Sugar().Errorf("用户更新异常: %+v", tx.Error)
		return 0, errors.New("用户更新异常")
	}
	return tx.RowsAffected, nil
}
