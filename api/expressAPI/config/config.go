package config

import (
	"apiProject/api/expressAPI/types"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var EnvConfig = InitConfig()

func InitConfig() *types.MysqlConfig {
	serverConfig, sqlConfigItem, RsaKeyInfo, jwtSecret, rabbitmqConfig := buildConfig()

	// 设置时区为亚洲/上海
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}

	return &types.MysqlConfig{
		ServerPort: serverConfig.Port,
		DbUser:     sqlConfigItem.Username,
		DbPass:     sqlConfigItem.Password,
		DbAddress:  sqlConfigItem.Url,
		DbName:     sqlConfigItem.Database,
		Loc:        loc,
		PublicKey:  RsaKeyInfo.Public,
		PrivateKey: RsaKeyInfo.Private,
		JWTSecret:  jwtSecret.Secret,
		Rabbitmq:   rabbitmqConfig,
	}
	/*return &types.MysqlConfig{
		ServerPort: GetEnv("SERVER_PORT", "8086"),
		DbUser:     GetEnv("DB_USER", "root"),
		DbPass:     GetEnv("DB_PASS", "root244112311"),
		DbAddress:  fmt.Sprintf("%s:%s", GetEnv("DB_HOST", "localhost"), GetEnv("DB_PORT", "3306")),
		DbName:     GetEnv("DB_NAME", "eladmin"),
		Loc:        loc,
		JWTSecret:  GetEnv("JWT_SECRET", "testjwtsecretadmin"),
		PublicKey:  "-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA3JnmWnHSYY+D703RzGS6\nPK/YvslCz/I8K+vr7PZq5ND0OhoLP+pexhBnDtHhDw4vFhFsosNl70YZkOImQiX+\nglF0PJqZld1TAtzSl4UtTGo/TI5fHuxy4q31Q+b/kun8lvfvU7ZJe/uWUbKx6dC7\nw6Pum7k5C9oqaOKXHllp78e1Wu3+qnvU2KnNI8gO2PyGUnLxJMJCIfvwqZPwRh07\nKi0yqsUFWuoIpvyLZePR/R35MGf6FQyavlZx7rNFUpw9prHbzC1nu9iSVi4A7B6l\nikjoEThi+uXIQHptQseIxosVhs2z04SCEobRv1H86vs/ciJKE7Y8nUIcwC0H8a/h\nTQIDAQAB\n-----END RSA PUBLIC KEY-----",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA3JnmWnHSYY+D703RzGS6PK/YvslCz/I8K+vr7PZq5ND0OhoL\nP+pexhBnDtHhDw4vFhFsosNl70YZkOImQiX+glF0PJqZld1TAtzSl4UtTGo/TI5f\nHuxy4q31Q+b/kun8lvfvU7ZJe/uWUbKx6dC7w6Pum7k5C9oqaOKXHllp78e1Wu3+\nqnvU2KnNI8gO2PyGUnLxJMJCIfvwqZPwRh07Ki0yqsUFWuoIpvyLZePR/R35MGf6\nFQyavlZx7rNFUpw9prHbzC1nu9iSVi4A7B6likjoEThi+uXIQHptQseIxosVhs2z\n04SCEobRv1H86vs/ciJKE7Y8nUIcwC0H8a/hTQIDAQABAoIBAF29xEZAwd6VRsJE\n9lb9oqoxK1B/Y7XLwMgFO775Q5kyNeYOtSMW6+kMhU6l3xYvt9CP3PMZR1KzHiAU\nCZ/oV0t3Y4ZxR7yITUMVJSQgAozLRVS51y/j2Dn9JBETsxzx81UPzJJtDrLxyQG0\nhqfN/Ev5eGaSAezIa2cgiojqA/tQvpP42qNvp+cDJGgRyLVmYfXiRco8pGrnW9k/\n0AFe5N9GgvjFvKBtYhxMLdnZLbgF/a0piYcas6ja3Wnb6K1HQXzfE4fy66WCh2Xn\nWoU3FfP9PTqWVT7JSCrSR0ot6ctY/2mKNn/AfoK556PZHJ7BuW4ylwXRPiVAhRrZ\nYPgmYyECgYEA/gmpjM8Ge4YQK33nOZ3qOOiuGQ6Pvzm5i6+VZmq4+6M14AFEeQH/\noQqxHcABlHym9+eWKHOv0nAy61p8qIkg2OQVH0nOdQSvd51N09c0ZNGsOX5XDJlU\n/6ls5VEdI0FPgKpEk53RD+eUCQ6L+xOGM6NFe8eQzqtPYCN4jl8dnWUCgYEA3k4e\nnL8Y16XpJGF5wp+ixsmdwsh9XlWygfILw5onYgGnaSv5FbxjCDSJedGJ0CM6jJXL\niLp7tTwpqe0oBOYW53KkYXn0posEU4lfGe27OQkbvzRXEQE3AGMcjOmz/mNepx8z\n3ZxmG43+NO8Eh6ElKYTKK7dkULMOgPo5bTs/yckCgYBtyHsvUOB6TUt7oCNm8Omh\nwlxKk9JnT2jyBuVHp2NdzACiV6nhqY1xaQ91zd5g7yWxCLIJtUUMalR3BVnN88Tw\nNlEyflDsnSO/S4mwvNX1o+8LwZ+Y4EKtYeifiVhQPg8/iVWtfYw1lVySNWklDiD2\n+94xSeM4jSv2Xh3hWRWRSQKBgD8wy4jY1THvak86GgdVo0qIYvzMSr6282/2optu\nRUWZnMHLixk/nJLnhDCJfHgam3j814c9Iw8IU/uGezqxQM93ifxfU0jH+WnZgZv4\nNKDo0udN9HXT95N3mNUBVXW5P12YBAE5hNjOSvU2//2hs9OSeHlmvvAlhbjp58sB\n7YbpAoGBALJrHmV9VXJ+nZH3GqSvLtscs6fnnOevToeEfukr56lBGylhfq8iUqea\nrzcJSDpUq+Ead38uBba5iapvFLeTSkFBYFNY4yHnPW4Uv/VuuIGIH9NVI0PTOK7Y\n5yeZ8oGYJBbjhTbrDF0YpMdcZOsvhjwp7D44TpcV9DDmsdkcLccY\n-----END RSA PRIVATE KEY-----",
	}*/
}

// buildConfig 构建并返回配置对象
func buildConfig() (types.ServerConfigItem, types.SqlConfigItem, types.RsaKey, types.JwtSecret, types.RabbitmqConfigItem) {
	viperConfig := ReadConfig("config", "application", "yml")

	var serverConfig types.ServerConfig
	//serverMap := viperConfig.Get("server").(map[string]interface{})
	//mapstructure.Decode(serverMap, &serverConfig)
	viperConfig.Unmarshal(&serverConfig)

	log.Println("服务的端口号:", serverConfig.Server.Port)
	log.Println("服务版本信息:", serverConfig.Server.Version)

	var sqlConfig types.SqlConfig
	//mysqlMap := viperConfig.Get("mysql").(map[string]interface{})
	//fmt.Println("数据库用户：", mysqlMap["username"])
	//mapstructure.Decode(mysqlMap, &sqlConfig)
	viperConfig.Unmarshal(&sqlConfig)
	log.Println("数据库url：", viperConfig.Get("mysql.url"))

	var rsaConfig types.RsaConfig
	//rsaMap := viperConfig.Get("rsa").(map[string]interface{})
	//mapstructure.Decode(rsaMap, &rsaConfig)
	viperConfig.Unmarshal(&rsaConfig)

	var jwtConfig types.JwtConfig
	//jwtMap := viperConfig.Get("jwt").(map[string]interface{})
	//mapstructure.Decode(jwtMap, &jwtConfig)
	viperConfig.Unmarshal(&jwtConfig)

	var rabbitmqConfig types.RabbitmqConfig
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
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	return viper.GetViper()
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
