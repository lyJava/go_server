package controller

import (
	"apiProject/api/response"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	if fileInfo.Size() > maxFileSize {
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
	} else {
		http.ServeFile(w, r, filePath)
	}
}
