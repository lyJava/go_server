package impl

import (
	"apiProject/api/expressAPI/types"
	"database/sql"
	"errors"
	"log"
)

type UserDB struct {
	Db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{
		Db: db,
	}
}

func (u *UserDB) GetUserById(id int64) (*types.User, error) {
	// 执行查询操作
	row := u.Db.QueryRow("SELECT user_id,  IFNULL(username, ''), IFNULL(nick_name, ''), IFNULL(phone, ''), IFNULL(email, ''), IFNULL(password, ''), IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s' ), ''), IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s' ), '') FROM sys_user WHERE user_id = ?", id)
	// 创建 Express 对象
	user := &types.User{}

	// 从查询结果中扫描数据到 Express 对象
	err := row.Scan(
		&user.UserId,
		&user.Username,
		&user.Nickname,
		&user.Phone,
		&user.Email,
		&user.Password,
		&user.CreateTime,
		&user.UpdateTime,
	)
	if err != nil {
		log.Printf("get user by id error ==%v", err)
		return nil, errors.New("获取用户失败")
	}
	return user, nil
}

func (u *UserDB) CreatUser(user *types.User) (*types.User, error) {
	// 执行查询操作
	row, err := u.Db.Exec("INSERT INTO sys_user(username, nick_name, phone, email, password, create_by, create_time, update_time, token) VALUES(?, ?, ?, ?, ?, ?, NOW(), NOW(), ?)",
		user.Username, user.Nickname, user.Phone, user.Email, user.Password, user.CreateBy, user.Token)

	if err != nil {
		log.Printf("exec error ==%v", err)
		return nil, errors.New("创建用户失败")
	}
	id, err := row.LastInsertId()
	if err != nil {
		log.Printf("row get lastInsertId error ==%v", err)
		return nil, errors.New("创建用户失败")
	}
	userById, err := u.GetUserById(id)
	if err != nil {
		return nil, errors.New("创建用户失败")
	}
	return userById, nil
}

func (u *UserDB) UserLogin(us *types.User) (*types.User, error) {
	row := u.Db.QueryRow("SELECT user_id, IFNULL(username, ''), IFNULL(nick_name, ''),  IFNULL(password, '') FROM sys_user WHERE username = ?", us.Username)
	// 创建 Express 对象
	user := &types.User{}

	err := row.Scan(
		&user.UserId,
		&user.Username,
		&user.Nickname,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
