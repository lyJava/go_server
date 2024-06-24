package config

import "github.com/golang-jwt/jwt"

// ServerConfig 服务器配置
type ServerConfig struct {
	Server ServerConfigItem `mapstructure:"server"`
}

// ServerConfigItem 服务项配置
type ServerConfigItem struct {
	Port       int    `mapstructure:"port"`
	Version    string `mapstructure:"version"`
	RequestLog bool   `mapstructure:"request-log"`
}

// MyClaims 自定义jwt的token返回字段
type MyClaims struct {
	Username           string `json:"username"` // 用户名
	UserId             string `json:"userId"`   //用户ID
	jwt.StandardClaims        //jwt的Claims
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