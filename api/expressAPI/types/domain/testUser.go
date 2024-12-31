package domain

// TestUser 测试用户信息
type TestUser struct {
	Id       int64  `json:"id,omitempty"`       // 主键ID
	Username string `json:"username,omitempty"` // 用户名
	Password string `json:"password,omitempty"` // 密码
	Email    string `json:"email,omitempty"`    // 邮箱
	Birthday string `json:"birthday,omitempty"` // 出生日期
	Phone    string `json:"phone,omitempty"`    // 手机
	Address  string `json:"address,omitempty"`  //地址

}
