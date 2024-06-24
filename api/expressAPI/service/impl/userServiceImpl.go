package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"database/sql"
	"errors"

	"go.uber.org/zap"
)

type UserDB struct {
	Db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{
		Db: db,
	}
}

func (u *UserDB) GetUserById(id int64) (*domain.User, error) {
	// 创建 Express 对象
	user := &domain.User{}
	// 执行查询操作
	err := u.Db.QueryRow(
		`SELECT
				user_id,
				IFNULL(username, ''),
				IFNULL(nick_name, ''),
				IFNULL(phone, ''),
				IFNULL(email, ''),
				IFNULL(password, ''),
				IFNULL(create_by, ''),
				IFNULL(update_by, ''),
				IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'), ''),
				IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), '')
			FROM
				sys_user
			WHERE
				user_id = ?`, id).Scan(
		&user.UserId,
		&user.Username,
		&user.Nickname,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreateBy,
		&user.UpdateBy,
		&user.CreateTime,
		&user.UpdateTime,
	)
	if err != nil {
		zap.L().Sugar().Errorf("get user by id error ==%+v", err)
		return nil, errors.New("获取用户失败")
	}
	return user, nil
}

func (u *UserDB) CreateUser(user *domain.User) (*domain.User, error) {

	// 开启事务
	tx, err := u.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("用户创建开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("用户创建事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("用户创建事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("用户创建事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("用户创建提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("用户创建事务提交完成")
			}
		}
	}()
	// 执行查询操作
	row, err := tx.Exec(
		`INSERT INTO
			   	sys_user(username, nick_name, phone, email, password, create_by, create_time, update_time, token)
			   VALUES(?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)`,
		user.Username, user.Nickname, user.Phone, user.Email, user.Password, user.CreateBy, user.Token)

	if err != nil {
		zap.L().Sugar().Errorf("用户创建执行错误==%+v", err)
		return nil, errors.New("创建用户失败")
	}
	id, err := row.LastInsertId()
	if err != nil {
		zap.L().Sugar().Errorf("用户创建获取最后ID错误==%+v", err)
		return nil, errors.New("创建用户失败")
	}
	/* userById, err := u.GetUserById(id)
	if err != nil {
		log.Printf("用户通过ID查询错误==%+v", err)
		return nil, errors.New("创建用户失败")
	} */

	// 创建 Express 对象
	userNew := &domain.User{}
	// 执行查询操作
	err = tx.QueryRow(
		`SELECT
				user_id,
				IFNULL(username, ''),
				IFNULL(nick_name, ''),
				IFNULL(phone, ''),
				IFNULL(email, ''),
				IFNULL(password, ''),
				IFNULL(create_by, ''),
				IFNULL(update_by, ''),
				IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s'), ''),
				IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s'), '')
			FROM
				sys_user
			WHERE
				user_id = ?`, id).Scan(
		&userNew.UserId,
		&userNew.Username,
		&userNew.Nickname,
		&userNew.Phone,
		&userNew.Email,
		&userNew.Password,
		&userNew.CreateBy,
		&userNew.UpdateBy,
		&userNew.CreateTime,
		&userNew.UpdateTime,
	)
	if err != nil {
		zap.L().Sugar().Errorf("get user by id error ==%+v", err)
		return nil, errors.New("获取用户失败")
	}

	return userNew, nil
}

func (u *UserDB) UserLogin(us *domain.User) (*domain.User, error) {
	// 创建 Express 对象
	user := &domain.User{}
	err := u.Db.QueryRow(
		`SELECT
					user_id,
					IFNULL(username, ''),
					IFNULL(nick_name, ''),
					IFNULL(password, '')
				FROM
					sys_user
				WHERE
					username = ?`, us.Username).Scan(
		&user.UserId,
		&user.Username,
		&user.Nickname,
		&user.Password,
	)
	if err != nil {
		zap.L().Sugar().Errorf("用户登录查询错误==%+v", err)
		return nil, err
	}
	return user, nil
}
