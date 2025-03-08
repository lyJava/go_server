package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
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

	/*dateFormat, err := utils.DateFormat(user.Birthday)
	if err != nil {
		return nil, errors.New("用户出生日期错误")
	}
	user.Birthday = dateFormat*/
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

func (u *TestUserGormDB) BatchCreateUser(list []*domain.TestUser) (int64, error) {
	tx := u.Db.CreateInBatches(list, 100)
	if tx.Error != nil {
		zap.L().Sugar().Errorf("用户批量创建异常: %+v", tx.Error)
		return 0, errors.New("用户批量创建异常")
	}
	var batchIds []int64
	for i := 0; i < len(list); i++ {
		batchIds = append(batchIds, list[i].Id)
	}
	zap.L().Sugar().Infof("用户批量创建返回的ID切片: %+v", batchIds)
	return tx.RowsAffected, nil
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

// BatchUpdateUser 批量更新
//
//goland:noinspection SqlResolve,SqlCaseVsIf,SqlError
func (u *TestUserGormDB) BatchUpdateUser(list []*domain.TestUser) (int64, error) {
	var rowsAffected int64
	// 用来存储每个字段的更新条件
	var usernameUpdates, passwordUpdates, emailUpdates, birthdayUpdates, phoneUpdates, addressUpdates []string
	var args []any
	var ids []any

	commonSqlFragment := "WHEN id = %d THEN '%s'"
	// 遍历 list 生成动态 SQL 更新
	for _, user := range list {
		userId := user.Id
		ids = append(ids, userId)
		// 构造动态更新的字段
		if user.Username != "" {
			// 这里直接将 id 和 username 填入 SQL 语句
			usernameUpdates = append(usernameUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Username))
		}
		if user.Password != "" {
			passwordUpdates = append(passwordUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Password))
		}
		if user.Email != "" {
			emailUpdates = append(emailUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Email))
		}
		//if user.Birthday !=  {
		//	birthdayUpdates = append(birthdayUpdates, fmt.Sprintf(commonSqlFragment, userId, formatDateForSQL(user.Birthday)))
		//}
		if user.Birthday != nil {
			birthdayUpdates = append(birthdayUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Birthday))
		}
		if user.Phone != "" {
			phoneUpdates = append(phoneUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Phone))
		}
		if user.Address != "" {
			addressUpdates = append(addressUpdates, fmt.Sprintf(commonSqlFragment, userId, user.Address))
		}
	}

	// 动态生成 SQL 语句
	sql := fmt.Sprintf(`
        UPDATE tb_test_user SET
            username = CASE %s ELSE username END, 
            password = CASE %s ELSE password END, 
            email = CASE %s ELSE email END, 
            birthday = CASE %s ELSE birthday END, 
            phone = CASE %s ELSE phone END, 
            address = CASE %s ELSE address END
        WHERE id IN (%s)`,
		// 为每个字段生成相应的条件
		strings.Join(usernameUpdates, " "),
		strings.Join(passwordUpdates, " "),
		strings.Join(emailUpdates, " "),
		strings.Join(birthdayUpdates, " "),
		strings.Join(phoneUpdates, " "),
		strings.Join(addressUpdates, " "),
		// 将所有的 IDs 作为条件传入
		strings.Join(generatePlaceholderArray(len(ids)), ","),
	)

	// 将 ids 作为条件传入
	args = append(args, ids...)

	// 执行 SQL
	tx := u.Db.Exec(sql, args...)
	if tx.Error != nil {
		zap.L().Sugar().Errorf("批量更新异常: %+v", tx.Error)
		return rowsAffected, errors.New("批量更新异常")
	}
	rowsAffected = tx.RowsAffected
	/*tx := u.Db.Begin()

	for _, user := range list {
		result := tx.Model(&domain.TestUser{}).Where("id = ?", user.Id).Updates(&domain.TestUser{
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			Birthday: user.Birthday,
			Phone:    user.Phone,
			Address:  user.Address,
		})
		if result.Error != nil {
			tx.Rollback()
			zap.L().Sugar().Errorf("用户更批量新异常: %+v", tx.Error)
			return 0, errors.New("用户批量更新异常")
		}

		// 检查更新影响的行数
		if result.RowsAffected > 0 {
			rowsAffected += result.RowsAffected
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		zap.L().Sugar().Errorf("用户批量更新事务提交失败: %+v", err)
		return 0, errors.New("用户批量更新失败")
	}
	*/
	zap.L().Sugar().Infof("用户批量更新成功条数: %d", rowsAffected)

	return rowsAffected, nil
}

func (u *TestUserGormDB) DeleteUser(id int64) (int64, error) {
	tx := u.Db.Delete(&domain.TestUser{}, id)
	if tx.Error != nil {
		zap.L().Sugar().Errorf("用户删除异常: %+v", tx.Error)
		return 0, errors.New("用户更新异常")
	}
	return tx.RowsAffected, nil
}

func (u *TestUserGormDB) BatchDeleteUser(ids []any) (int64, error) {
	tx := u.Db.Delete(&domain.TestUser{}, ids)
	if tx.Error != nil {
		zap.L().Sugar().Errorf("用户批量删除异常: %+v", tx.Error)
		return 0, errors.New("用户批量删除失败")
	}
	return tx.RowsAffected, nil
}

// SelectPage 查询分页
func (u *TestUserGormDB) SelectPage(param *param.TestUserPageParam) ([]*domain.TestUser, int64, int64, error) {
	setPageDefault(param)

	// 计算分页参数
	offset := (param.Page - 1) * param.Size

	// 构建查询条件
	db := u.Db.Model(&domain.TestUser{})
	// 动态构建查询条件
	if param.Username != "" {
		db = db.Where(fmt.Sprintf("username LIKE CONCAT('%%', '%s', '%%')", param.Username))
	}
	if param.Email != "" {
		db = db.Where(fmt.Sprintf("email LIKE CONCAT('%%', '%s', '%%')", param.Email))
	}
	if param.Birthday != "" {
		db = db.Where("birthday = ?", param.Birthday)
	}
	if param.Phone != "" {
		db = db.Where("phone LIKE ?", fmt.Sprintf("%%%s%%", param.Phone))
	}
	if param.Address != "" {
		db = db.Where("address LIKE ?", fmt.Sprintf("%%%s%%", param.Address))
	}

	// 处理排序（如果 param.Column 和 param.Order 为空，则不进行排序）
	if param.Column != "" && param.Order != "" {
		db = db.Order(fmt.Sprintf("%s %s", param.Column, param.Order))
	}

	// 查询总记录数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		zap.L().Sugar().Errorf("获取总记录数异常: %+v", err)
		return nil, 0, 0, errors.New("获取总记录数失败")
	}

	// 查询用户列表
	var testUserList []*domain.TestUser
	if err := db.Offset(int(offset)).Limit(int(param.Size)).Find(&testUserList).Error; err != nil {
		zap.L().Sugar().Errorf("获取用户列表异常: %+v", err)
		return nil, 0, 0, errors.New("获取用户列表失败")
	}

	/*if len(testUserList) > 0 {
		for _, user := range testUserList {
			birthday := user.Birthday
			if birthday != "" {
				dateFormat, err := utils.DateFormat(birthday)
				if err != nil {
					return nil, 0, 0, errors.New("获取用户列表失败")
				}
				user.Birthday = dateFormat
			}
		}
	}*/

	return testUserList, total, getTotalPage(total, param.Size), nil
}

func (u *TestUserGormDB) SelectAll() ([]*domain.TestUser, error) {
	var list []*domain.TestUser
	if err := u.Db.Find(&list).Error; err != nil {
		zap.L().Sugar().Errorf("查询所有用户异常: %+v", err)
		return nil, errors.New("查询所有用户失败")
	}
	/*if len(list) > 0 {
		for _, user := range list {
			user.Birthday = user.Birthday[:10]
		}
	}*/
	return list, nil
}

// formatDateForSQL 格式化日期
//
// 参数:
//
//	birthday: 出生日期字符串
//
// 返回:
//
//	string: 格式化后的日期
func formatDateForSQL(birthday string) string {
	parsed, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		zap.L().Sugar().Errorf("日期格式化失败: %v", err)
		return birthday // 如果格式化失败，返回原始日期
	}
	return parsed.Format("2006-01-02") // 返回标准的日期格式
}

// generatePlaceholders 生成SQL占位符(例如：`?, ?, ?`)
//
// 参数:
//   - n: 占位符个数
//
// 返回:
//   - string 占位符字符串
func generatePlaceholders(n int) string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ",")
}

// generatePlaceholderArray 生成SQL占位符切片 (例如：{?, ?, ?})
//
// 参数:
//   - n: 占位符个数
//
// 返回:
//   - string 占位符切片
func generatePlaceholderArray(count int) []string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return placeholders
}

// setPageDefault 分页参数默认设置
//
// 参数:
//
//	param：分页查询参数结构体
//
// 返回：
//
//	*：分页查询参数结构体本身
func setPageDefault(param *param.TestUserPageParam) *param.TestUserPageParam {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Size <= 0 {
		param.Size = 10
	}
	return param
}

// getTotalPage 计算总页数
func getTotalPage(total, size int64) int64 {
	return int64(int(math.Ceil(float64(total) / float64(size))))
}
