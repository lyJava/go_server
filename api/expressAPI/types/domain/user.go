package domain

// User 用户信息
type User struct {
	UserId     int64  `json:"userId,omitempty"`     // 主键ID
	Username   string `json:"username,omitempty"`   // 用户名
	Nickname   string `json:"nickname,omitempty"`   // 昵称
	Phone      string `json:"phone,omitempty"`      // 手机
	Email      string `json:"email,omitempty"`      // 邮箱
	Password   string `json:"password,omitempty"`   // 密码
	CreateBy   string `json:"createBy,omitempty"`   // 创建人
	CreateTime string `json:"createTime,omitempty"` // 创建时间
	UpdateTime string `json:"updateTime,omitempty"` // 更新时间
	Token      string `json:"token,omitempty"`      // 令牌
}
