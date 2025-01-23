package config

// MinioConfig minio连接配置
type MinioConfig struct {
	MinioConfigItem MinioConfigConfigItem `mapstructure:"minio"` // minio连接配置项
}

// MinioConfigConfigItem minio连接配置连接配置项目
type MinioConfigConfigItem struct {
	Endpoint  string `mapstructure:"endpoint"`  // 连接端点
	AccessKey string `mapstructure:"accessKey"` // 访问密钥
	SecretKey string `mapstructure:"secretKey"` // 私密密钥
	UseSSL    bool   `mapstructure:"useSSL"`    // 是否启用SSL/TLS，一般情况不启用
}
