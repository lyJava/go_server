package utils

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"log"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func CheckBuckets(client *minio.Client) ([]minio.BucketInfo, error) {
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		log.Printf("无法列出存储桶: %+v", err)
		return nil, err
	}

	log.Println("MinIO 存储桶:")
	for _, bucket := range buckets {
		log.Printf(" - %s\n", bucket.Name)
	}

	return buckets, nil
}

// UploadObj 上传
//
// 参数
//
//   - client (*minio.Client): minio客户端对象
//   - bucketName (string): 桶名称
//   - objectName (string): 上传后显示的文件名
//   - filePath (string): 上传文件路径
//
// 返回
//
//   - string: 文件上传后的临时分享URL
//   - error: 如果上传失败，返回错误信息
func UploadObj(client *minio.Client, bucketName, objectName, filePath string) (string, error) {
	// 验证存储桶是否存在
	exists, err := client.BucketExists(context.Background(), bucketName)

	if err != nil {
		return "", errors.New("获取桶失败")
	}
	// 桶不存在则创建
	if !exists {
		err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("创建桶异常: %+v", err)
			return "", errors.New("创建桶失败")
		}
	}

	// 获取文件的 MIME 类型
	fileExt := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(fileExt)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	log.Printf("上传文件MIME类型为: %s", mimeType)

	uploadInfo, err := client.FPutObject(context.Background(), bucketName, objectName,
		filePath, minio.PutObjectOptions{
			ContentType: mimeType, // 设置文件类型为图片
		})
	if err != nil {
		log.Printf("文件上传失败: %+v", err)
		return "", errors.New("文件上传失败")
	}
	log.Printf("文件上传成功! ===%v", ToJsonFormat(&uploadInfo))
	// inline:预览， attachment:下载
	previewUrl, err := client.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*300,
		url.Values{
			"response-content-type":        []string{mimeType},
			"response-content-disposition": []string{"inline; filename=" + objectName},
		},
	)
	if err != nil {
		return "", errors.New("返回临时分享URL失败")
	}

	temporarySharingUrl := previewUrl.String()
	log.Printf("返回临时分享URL===%s", temporarySharingUrl)
	return temporarySharingUrl, nil
}

// DownloadObj 下载
//
// 参数
//
//   - client (*minio.Client): minio客户端对象
//   - bucketName (string): 桶名称
//   - objectName (string): 需要下载的文件名
//   - filePath (string): 下载文件存放路径
//
// 返回
//
//   - string: 文件下载成功返回的绝对路径
//   - error: 下载失败返回错误
func DownloadObj(client *minio.Client, bucketName, objectName, filePath string) (string, error) {
	// 获取当前根目录
	workingDir, err := os.Getwd()
	if err != nil {
		log.Printf("获取工作目录失败: %v", err)
		return "", err
	}

	// 拼接 download 文件夹路径和文件名
	downloadDir := filepath.Join(workingDir, "download")
	// 创建 download 文件夹（如果不存在的话）
	if _, err = os.Stat(downloadDir); os.IsNotExist(err) {
		err = os.Mkdir(downloadDir, os.ModePerm)
		if err != nil {
			log.Printf("创建 download 文件夹失败: %v", err)
			return "", err
		}
	}

	fileSuffix := filepath.Ext(objectName)
	log.Printf("下载文件fileSufix: %v", fileSuffix)
	// 拼接最终的文件路径
	filePath = filepath.Join(downloadDir, uuid.NewString()+fileSuffix)

	err = client.FGetObject(context.Background(), bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("文件下载失败: %v", err)
		return "", err
	}
	log.Printf("文件下载成功!,路径为====%s", filePath)
	return filePath, nil
}

// RemoveObj 删除文件
//
// 参数
//
//   - client (*minio.Client): minio客户端对象
//   - bucketName (string): 桶名称
//   - objectName (string): 需要删除的文件名
//
// 返回
//
//   - error: 删除失败返回错误
func RemoveObj(client *minio.Client, bucketName, objectName string) error {
	if err := client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{}); err != nil {
		log.Printf("文件删除错误: %v", err)
		return errors.New("文件删除失败")
	}
	return nil
}
