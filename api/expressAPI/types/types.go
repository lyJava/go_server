package types

import (
	"github.com/golang-jwt/jwt"
	"time"
)

// MysqlConfig mysql的结构
type MysqlConfig struct {
	ServerPort int                // 服务端口号
	DbUser     string             // 数据库用户名
	DbPass     string             // 数据库密码
	DbAddress  string             // 数据库地址
	DbName     string             // 数据库名称
	JWTSecret  string             // jwt的token加密
	Loc        *time.Location     // 设置时区
	PublicKey  string             // rsa公钥
	PrivateKey string             // rsa私钥
	Rabbitmq   RabbitmqConfigItem //rabbitmq配置
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

// Payment 支付结构体
type Payment struct {
	Id       string `json:"id"`       // 支付ID
	Quantity int64  `json:"quantity"` // 支付金额(分)
}

// Order 订单结构体
type Order struct {
	Id    string    `json:"id"`    // 订单ID
	Items []Payment `json:"items"` // 支付数组
}

// MyClaims 自定义jwt的token返回字段
type MyClaims struct {
	Username           string `json:"username"` // 用户名
	UserId             string `json:"userId"`   //用户ID
	jwt.StandardClaims        //jwt的Claims
}

// SqlConfig mysql配置
type SqlConfig struct {
	Mysql SqlConfigItem `mapstructure:"mysql"` // mysql
}

// SqlConfigItem mysql配置项
type SqlConfigItem struct {
	Url      string `mapstructure:"url"`      // 连接URL
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
	Database string `mapstructure:"database"` // 数据库名称
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Server ServerConfigItem `mapstructure:"server"`
}

// ServerConfigItem 服务项配置
type ServerConfigItem struct {
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}

// JwtConfig JWT的配置
type JwtConfig struct {
	Jwt JwtSecret `mapstructure:"jwt"`
}

// JwtSecret JWT的配置
type JwtSecret struct {
	Secret string `mapstructure:"secret"`
}

// RsaConfig rsa的配置
type RsaConfig struct {
	Rsa RsaKeyConfig `mapstructure:"rsa"` //密钥
}

// RsaKeyConfig rsa密钥的配置
type RsaKeyConfig struct {
	Key RsaKey `mapstructure:"key"` //密钥
}

// RsaKey rsa密钥结构
type RsaKey struct {
	Public  string `mapstructure:"public"`  // 公钥
	Private string `mapstructure:"private"` //私钥
}

// RedisConfig redis连接结构
type RedisConfig struct {
	Redis RedisConfigItem `mapstructure:"redis"`
}

// RedisConfigItem redis连接结构
type RedisConfigItem struct {
	Address  string `mapstructure:"address"`  // 连接地址
	Password string `mapstructure:"password"` // 密码
	Db       int    `mapstructure:"db"`       // 连接数据库索引，默认0
}

// RabbitmqConfig rabbitmq配置结构
type RabbitmqConfig struct {
	Rabbitmq RabbitmqConfigItem `mapstructure:"rabbitmq"`
}

// RabbitmqConfigItem rabbitmq配置项结构体
type RabbitmqConfigItem struct {
	Host              string `mapstructure:"host"`               // 主机
	Port              int    `mapstructure:"port"`               // 端口
	Username          string `mapstructure:"username"`           // 用户
	Password          string `mapstructure:"password"`           // 密码
	VirtualHost       string `mapstructure:"virtual-host"`       // 虚拟主机
	ConnectionTimeout int64  `mapstructure:"connection-timeout"` // 超时时间
	PublisherReturns  bool   `mapstructure:"publisher-returns"`  // 是否返回
}
