package config

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