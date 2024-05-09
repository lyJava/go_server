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

// PostgresqlConfig postgresql连接配置
type PostgresqlConfig struct {
	Postgresql PostgresqlConfigItem `mapstructure:"postgresql"` // postgresql连接配置项
}

// PostgresqlConfigItem postgresql连接配置项目
type PostgresqlConfigItem struct {
	Host     string `mapstructure:"host"`     // 连接URL
	Port     int    `mapstructure:"port"`     // 连接端口
	User     string `mapstructure:"user"`     // 用户名
	Password string `mapstructure:"password"` //密码
	DbName   string `mapstructure:"db-name"`  // 数据库名称
}
