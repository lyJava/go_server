package config

// RabbitmqConfig rabbitmq配置结构
type RabbitmqConfig struct {
	Rabbitmq RabbitmqConfigItem `mapstructure:"rabbitmq"`
}

// RabbitmqConfigItem rabbitmq配置项结构体
type RabbitmqConfigItem struct {
	Enable            bool   `mapstructure:"enable"`             // 是否启用
	Host              string `mapstructure:"host"`               // 主机
	Port              int    `mapstructure:"port"`               // 端口
	Username          string `mapstructure:"username"`           // 用户
	Password          string `mapstructure:"password"`           // 密码
	VirtualHost       string `mapstructure:"virtual-host"`       // 虚拟主机
	ConnectionTimeout int64  `mapstructure:"connection-timeout"` // 超时时间
	PublisherReturns  bool   `mapstructure:"publisher-returns"`  // 是否返回
}