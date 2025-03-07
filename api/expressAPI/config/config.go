package config

import (
	"apiProject/api/expressAPI/types/config"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"
)

var EnvConfig = InitConfig()

func InitConfig() *config.MysqlConfig {
	serverConfig, sqlConfigItem, RsaKeyInfo, jwtSecret, rabbitmqConfig := buildConfig()

	// 设置时区为亚洲/上海
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}

	return &config.MysqlConfig{
		DbUser:       sqlConfigItem.Username,
		DbPass:       sqlConfigItem.Password,
		DbAddress:    sqlConfigItem.Url,
		DbName:       sqlConfigItem.Database,
		Loc:          loc,
		PublicKey:    RsaKeyInfo.Public,
		PrivateKey:   RsaKeyInfo.Private,
		JWTSecret:    jwtSecret.Secret,
		Rabbitmq:     rabbitmqConfig,
		ServerConfig: serverConfig,
	}
}

// GetCurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// buildConfig 构建并返回配置对象
func buildConfig() (config.ServerConfigItem, config.SqlConfigItem, config.RsaKey, config.JwtSecret, config.RabbitmqConfigItem) {
	log.Println("当前执行文件的路径:", GetCurrentAbPathByCaller())
	currentPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录错误===%v", err)
	}
	log.Println("获取当前工作目录路径:", currentPath)
	// /Users/yangge/GolandProjects/apiProject/api/expressAPI/config/application.yml
	configAbsolutePath := filepath.Join(currentPath, "api/expressAPI/config")
	log.Println("获取配置文件绝对路径:", configAbsolutePath)
	// todo 需要配置Working directory与当前的main.go的目录保持一直，不然会出现go run main.go启动与golang的debug工具启动获取目录不一致情况
	viperConfig := ReadConfig(configAbsolutePath, "application", "yml")

	var serverConfig config.ServerConfig
	//serverMap := viperConfig.Get("server").(map[string]interface{})
	//mapstructure.Decode(serverMap, &serverConfig)
	viperConfig.Unmarshal(&serverConfig)
	log.Println("服务的端口号:", serverConfig.Server.Port)
	log.Println("服务版本信息:", serverConfig.Server.Version)

	var sqlConfig config.SqlConfig
	//mysqlMap := viperConfig.Get("mysql").(map[string]interface{})
	//fmt.Println("数据库用户：", mysqlMap["username"])
	//mapstructure.Decode(mysqlMap, &sqlConfig)
	viperConfig.Unmarshal(&sqlConfig)
	log.Println("数据库url：", viperConfig.Get("mysql.url"))

	var rsaConfig config.RsaConfig
	//rsaMap := viperConfig.Get("rsa").(map[string]interface{})
	//mapstructure.Decode(rsaMap, &rsaConfig)
	viperConfig.Unmarshal(&rsaConfig)

	var jwtConfig config.JwtConfig
	//jwtMap := viperConfig.Get("jwt").(map[string]interface{})
	//mapstructure.Decode(jwtMap, &jwtConfig)
	viperConfig.Unmarshal(&jwtConfig)

	var rabbitmqConfig config.RabbitmqConfig
	viperConfig.Unmarshal(&rabbitmqConfig)
	// 将json格式化输出
	rabbitmqItem, _ := json.MarshalIndent(rabbitmqConfig.Rabbitmq, "", "    ")
	log.Printf("rabbitmq配置信息==\r\n%s", rabbitmqItem)

	return serverConfig.Server, sqlConfig.Mysql, rsaConfig.Rsa.Key, jwtConfig.Jwt, rabbitmqConfig.Rabbitmq
}

// ReadConfig 读取配置
//
// 参数
//
//		path：配置路径
//		name：配置文件名(不包括后缀)
//	 configType：配置文件路径
//
// 返回
//
//	viper：*viper.Viper
func ReadConfig(path, name, configType string) *viper.Viper {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType(configType)

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
