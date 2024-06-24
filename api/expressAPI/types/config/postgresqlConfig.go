package config


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