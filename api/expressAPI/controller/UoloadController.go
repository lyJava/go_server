package controller

import (
	"apiProject/api/response"
	"bufio"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// UploadController 上传控制器
type UploadController struct{}

// UploadControllerInit 上传控制器初始化
func UploadControllerInit() *UploadController {
	return &UploadController{}
}

func (*UploadController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
}

// 10 << 20 是位运算符 << 的用法，表示左移操作。在这里，10 << 20 表示将整数 10 左移 20 位，等于 10 * 2^20，
// 即 10 乘以 1024 * 1024，结果是 10485760。因此，10 << 20 等价于 10485760，表示 10 MB。

// maxUploadSize 上传文件大小50MB
const maxUploadSize = 50 << 20

// bufferSize 缓冲区大小
const bufferSize = 4096 // 设置缓冲区大小为 4096 字节

// UploadHandler 文件上传处理器
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析表单，设置内存限制为10MB
	// 表示设置内存限制为 10 MB。也就是说，如果请求体的大小超过了 10 MB，剩余的数据会被写入临时文件中。这个限制是为了防止恶意或不当使用大量内存导致服务器资源耗尽，从而影响其他请求的处理
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println("表单内存超出限制异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("表单数据超过最大内存设置"))
		return
	}

	// multipartFile 获取上传的文件
	multipartFile, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("检索文件出错", err.Error())
		response.WriteJson(w, response.FailMessageResp("检索文件失败"))
		return
	}

	defer multipartFile.Close()

	// 检查文件大小
	if handler.Size > maxUploadSize {
		log.Printf("文件大小超出限制:%s", "50MB")
		// http.Error(w, fmt.Sprintf("File size exceeds the limit of %d bytes", maxUploadSize), http.StatusBadRequest)
		response.WriteJson(w, response.FailMessageResp("文件大小不能超过50M"))
		return
	}

	// 设置文件上传保存的目录
	uploadDir := "./upload/"

	// 创建上传目录，如果不存在的话
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Println("文件上传创建目录异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("上传文件失败"))
		return
	}

	// 拼接文件路径(直接拼接可能会因为系统差异创建目录时产生异常)
	filePath := filepath.Join(uploadDir, handler.Filename)
	// 创建文件
	newFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("文件上传打开文件异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("上传文件失败"))
		return
	}
	defer newFile.Close()

	// 将文件内容复制到目标文件
	/*_, err = io.Copy(f, file)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("上传文件失败"))
		return
	}*/

	// 创建带有指定缓冲区大小的写入器
	writer := bufio.NewWriterSize(newFile, bufferSize)
	defer writer.Flush()

	// 将上传的文件内容复制到目标文件
	_, err = io.Copy(writer, multipartFile)
	if err != nil {
		log.Println("文件上传内容复制到目标文件异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("上传文件失败"))
		return
	}
	response.WriteJson(w, response.OkMessageResp("上传文件成功"))
}
