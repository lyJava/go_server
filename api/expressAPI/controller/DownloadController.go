package controller

import (
	"apiProject/api/response"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// DownloadController 下载控制器
type DownloadController struct{}

// DownloadControllerInit 下载控制器初始化
func DownloadControllerInit() *DownloadController {
	return &DownloadController{}
}

func (*DownloadController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/download/file", DownloadHandler).Methods("GET")
}

// maxFileSize 文件大小10MB
const maxFileSize = 10 << 20

// DownloadHandler 文件下载处理器
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}

	// 文件路径拼接
	filePath := filepath.Join("./upload", fileName)

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("打开文件异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("文件不存在"))
		return
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("获取文件信息异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("获取文件信息失败"))
		return
	}

	// 读取文件的前512字节
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		log.Println("读取文件异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("读取文件失败"))
		return
	}
	// 获取文件的MIME类型
	fileMineType := http.DetectContentType(buffer)
	log.Printf("当前文件MIME:%s", fileMineType)

	w.Header().Set("Content-Type", fileMineType)
	// 对文件名进行编码处理，避免在firefox或postman中请求下载变成response.bin
	w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+url.QueryEscape(fileName))
	w.Header().Set("Content-Transfer-Encoding", "binary")

	if fileInfo.Size() > maxFileSize {
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	} else {
		http.ServeFile(w, r, filePath)
	}
}
