package types

import (
	"github.com/golang-jwt/jwt"
	"time"
)

// MysqlConfig mysql的结构
type MysqlConfig struct {
	ServerPort string         // 服务端口号
	DbUser     string         // 数据库用户名
	DbPass     string         // 数据库密码
	DbAddress  string         // 数据库地址
	DbName     string         // 数据库名称
	JWTSecret  string         // jwt的token加密
	Loc        *time.Location // 设置时区
	PublicKey  string         // rsa公钥
	PrivateKey string         // rsa私钥
}

// ExpressSearchParam 分页查询参数结构
type ExpressSearchParam struct {
	ExpressName   string `json:"expressName"`   // 快递名称
	ExpressNumber string `json:"expressNumber"` // 快递单号
	UserId        string `json:"userId"`        // 用户ID
	FromName      string `json:"fromName"`      // 发件人姓名
	FromPhone     string `json:"fromPhone"`     // 发件人手机
	FromAddress   string `json:"fromAddress"`   // 发件人地址
	PickupCode    string `json:"pickupCode"`    // 取件码
	CreateBy      string `json:"createBy"`      // 创建人
	Page          int64  `json:"page"`          // 当前页码
	Size          int64  `json:"size"`          // 每页条数
}

// Express 快递对象结构
//
// 设置了omitempty,字段的值为空值时，JSON序列化时将忽略该字段，因此不会显示空值。
type Express struct {
	ID            int64  `json:"id"`                      // 主键ID
	UserId        string `json:"userId,omitempty"`        // 用户ID
	ExpressName   string `json:"expressName,omitempty"`   // 快递名称
	ExpressNumber string `json:"expressNumber,omitempty"` // 快递单号
	FromName      string `json:"fromName,omitempty"`      // 发件人姓名
	FromPhone     string `json:"fromPhone,omitempty"`     // 发件人手机
	FromAddress   string `json:"fromAddress,omitempty"`   // 发件人地址
	PickupCode    string `json:"pickupCode,omitempty"`    // 取件码
	CreateBy      string `json:"createBy,omitempty"`      // 创建人
	CreateTime    string `json:"createTime,omitempty"`    // 创建时间
	UpdateTime    string `json:"updateTime,omitempty"`    // 更新时间
}

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

// MyClaims 自定义jwt的token返回字段
type MyClaims struct {
	Username           string `json:"username"` // 用户名
	UserId             string `json:"userId"`   //用户ID
	jwt.StandardClaims        //jwt的Claims
}
