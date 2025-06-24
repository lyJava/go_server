package datasource

import (
	"apiProject/api/expressAPI/config"
	cfg "apiProject/api/expressAPI/types/config"
	"apiProject/api/utils"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

// InitMinioClient 初始化minio客户端
//
// 返回
//   - *minio.Client: minio客户端对象
//   - error: 可能存在的错误
func InitMinioClient() (*minio.Client, error) {

	viperConfig := config.ReadConfig("api/expressAPI/config", "application", "yml")
	if viperConfig == nil {
		log.Println("未读取到viper配置信息")
		return nil, errors.New("未读取到viper配置信息")
	}
	var minioConfig cfg.MinioConfig
	if err := viperConfig.Unmarshal(&minioConfig); err != nil {
		log.Printf("获取minio配置错误:%+v", err)
		return nil, errors.New("获取minio配置失败")
	}
	minConfigItem := minioConfig.MinioConfigItem
	log.Printf("minio客户端配置项===%v", utils.ToJsonFormat(minConfigItem))

	// 初始化 MinIO 客户端
	client, err := minio.New(minConfigItem.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minConfigItem.AccessKey, minConfigItem.SecretKey, ""),
		Secure: minConfigItem.UseSSL,
	})
	if err != nil {
		log.Printf("无法连接到 MinIO: %v", err)
		return nil, errors.New("无法连接到minio")
	}

	offline := client.IsOffline()
	log.Printf("minio客户端是否在线====:%v", offline)

	if !offline {
		return nil, errors.New("minio客户端不在线")
	}

	return client, nil
}
