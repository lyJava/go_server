package config

import "time"

// MysqlConfig mysql的结构
type MysqlConfig struct {
	DbUser       string             // 数据库用户名
	DbPass       string             // 数据库密码
	DbAddress    string             // 数据库地址
	DbName       string             // 数据库名称
	JWTSecret    string             // jwt的token加密
	Loc          *time.Location     // 设置时区
	PublicKey    string             // rsa公钥
	PrivateKey   string             // rsa私钥
	Rabbitmq     RabbitmqConfigItem //rabbitmq配置
	ServerConfig ServerConfigItem   //服务器配置
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
