package controller

import (
	"apiProject/api/response"
	"apiProject/api/utils"
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// DownloadController 下载控制器
type DownloadController struct{}

// DownloadControllerInit 下载控制器初始化
func DownloadControllerInit() *DownloadController {
	return &DownloadController{}
}

func (*DownloadController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/download/file", DownloadHandler3).Methods("GET")
	r.HandleFunc("/create/zip", handler4).Methods("GET")
}

// maxFileSize 文件大小10MB
const maxFileSize = 10 << 20

// DownloadHandler 文件下载处理器
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	isZip := r.URL.Query().Get("zip")
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}
	var useZip = false
	if isZip == "true" {
		useZip = true
	}
	// 文件路径拼接
	filePath := filepath.Join("./upload", fileName)

	if useZip {
		zipName := utils.GetPrefix(fileName) + ".zip"
		//var fileList []string
		//fileList = append(fileList, filePath)
		//createAndWriteZipFileToResponse(zipName, fileList, w)
		// 创建临时目录
		/*tempDir, err := os.MkdirTemp("", "temp-download")
		if err != nil {
			log.Println("创建临时目录失败", err.Error())
			response.WriteJson(w, response.FailMessageResp("创建临时目录失败"))
			return
		}
		defer os.RemoveAll(tempDir)*/

		//zipName := utils.GetPrefix(fileName) + ".zip"
		/*log.Println("zip文件名", zipName)
		//zipName := filepath.Join(tempDir, "测试.zip")
		//zipName := filepath.Join(tempDir, url.QueryEscape(strings.TrimSuffix(fileName, filepath.Ext(fileName)))) + ".zip"
		// 创建ZIP文件路径
		zipOutPath, err := os.Create(zipName)
		log.Println("zip文件路径", zipOutPath.Name())
		if err != nil {
			log.Println("创建ZIP文件异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("创建zip文件失败"))
			return
		}
		defer zipOutPath.Close()
		//zw := zip.NewWriter(zipOutPath)*/

		/*// 创建内存缓冲区
		buf := new(bytes.Buffer)
		// 创建 ZIP 编写器
		zw := zip.NewWriter(buf)

		// 将文件内容添加到ZIP文件中
		fileInZip, err := zw.Create(fileName)
		if err != nil {
			log.Println("创建ZIP文件内部异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("创建ZIP文件内部失败"))
			return
		}

		log.Println("fileInZip====", fileInZip)

		// todo 该方式文件比较大的时候比较耗内存
		_, err = io.Copy(fileInZip, file)
		if err != nil {
			log.Println("复制文件内容到ZIP文件内部异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("复制文件内容到ZIP文件内部失败"))
			return
		}

		// 关闭zip.Writer
		if err := zw.Close(); err != nil {
			log.Println("关闭ZIP编写器异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("关闭ZIP编写器失败"))
			return
		}*/

		// 将文件内容逐块复制到 ZIP 文件中
		/*buf := make([]byte, 1024)
		for {
			n, err := file.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println("读取文件异常", err.Error())
				response.WriteJson(w, response.FailMessageResp("读取文件失败"))
				return
			}

			_, err = fileInZip.Write(buf[:n])
			if err != nil {
				log.Println("写入ZIP文件内部异常", err.Error())
				http.Error(w, "写入ZIP文件内部失败", http.StatusInternalServerError)
				return
			}
		}*/
		var files []string
		files = append(files, filePath)
		buf, err := utils.CreatZipBuffer(files)
		if err != nil {
			response.WriteJson(w, response.FailMessageResp(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+filepath.Base(zipName))
		w.Header().Set("Content-Transfer-Encoding", "binary")

		// 提供 ZIP 文件给客户端下载
		//http.ServeContent(w, r, zipFileInfo.Name(), zipFileInfo.ModTime(), zipOutPath)
		// 提供 ZIP 文件给客户端下载
		//http.ServeFile(w, r, zipName)
		//http.ServeContent(w, r, zipFileInfo.Name(), zipFileInfo.ModTime(), zipFile)
		buf.WriteTo(w)
	} else {

		file, err := utils.OpenFile(filePath)
		if err != nil {
			response.WriteJson(w, response.FailMessageResp(err.Error()))
			return
		}

		defer file.Close()

		fileInfo, err := utils.GetFileInfo(file)
		if err != nil {
			log.Println("获取文件信息异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("获取文件信息失败"))
			return
		}

		fileSize := fileInfo.Size()

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

		if fileSize > maxFileSize {
			http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
		} else {
			http.ServeFile(w, r, filePath)
		}
	}
}

// 创建 ZIP 文件并写入到 HTTP 响应中
func createAndWriteZipFileToResponse(zipName string, files []string, w http.ResponseWriter) error {
	// 创建 ZIP 文件
	zipBuffer := new(bytes.Buffer)
	zw := zip.NewWriter(zipBuffer)

	// 逐个文件添加到 ZIP 文件中
	for _, filePath := range files {
		// 打开要添加到 ZIP 文件的文件
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将文件添加到 ZIP 文件中
		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}

		// 设置zip文件名
		header.Name = filepath.Base(filePath)

		// 创建zip文件头
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}

	// 关闭 ZIP Writer
	err := zw.Close()
	if err != nil {
		return err
	}

	// 设置 HTTP 响应头
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+zipName)
	w.Header().Set("Content-Length", strconv.Itoa(len(zipBuffer.Bytes())))

	// 将 ZIP 文件写入 HTTP 响应
	_, err = w.Write(zipBuffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func DownloadHandler2(w http.ResponseWriter, r *http.Request) {
	// 获取文件名和是否需要压缩
	isZip := r.URL.Query().Get("zip")
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}
	var useZip = false
	if isZip == "true" {
		useZip = true
	}

	// 文件路径拼接
	filePath := filepath.Join("./upload", fileName)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("打开文件异常", err.Error())
		response.WriteJson(w, response.FailMessageResp("文件不存在"))
		return
	}
	defer file.Close()

	if useZip {
		// zip文件名
		zipName := utils.GetPrefix(fileName) + ".zip"
		// 创建内存缓冲区
		buf := new(bytes.Buffer)
		// 创建 ZIP 编写器
		zw := zip.NewWriter(buf)

		// 将文件内容添加到 ZIP 文件中
		fileInZip, err := zw.Create(zipName)
		if err != nil {
			log.Println("创建ZIP文件内部异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("创建ZIP文件内部失败"))
			return
		}

		log.Println("生成zip文件路径", zipName)

		// 将文件内容复制到 ZIP 文件中
		_, err = io.Copy(fileInZip, file)
		if err != nil {
			log.Println("复制文件内容到ZIP文件内部异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("复制文件内容到ZIP文件内部失败"))
			return
		}

		// 关闭 ZIP 编写器
		if err := zw.Close(); err != nil {
			log.Println("关闭ZIP编写器异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("关闭ZIP编写器失败"))
			return
		}

		// 将缓冲区的内容写入 HTTP 响应体
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(zipName))
		w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
		if _, err := buf.WriteTo(w); err != nil {
			log.Println("写入HTTP响应体异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("写入HTTP响应体失败"))
			return
		}
	} else {
		// 如果不需要压缩，则直接提供文件下载
		http.ServeFile(w, r, filePath)
	}
}

func DownloadHandler3(w http.ResponseWriter, r *http.Request) {
	// 获取文件名和是否需要压缩
	isZip := r.URL.Query().Get("zip")
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}
	var useZip = false
	if isZip == "true" {
		useZip = true
	}

	// 文件路径拼接
	filePath := filepath.Join("./upload", fileName)

	if useZip {
		var files []string
		files = append(files, filePath)
		tempFile, err := utils.CreateZipTemp(files, "temp-zip-*.zip")
		if err != nil {
			response.WriteJson(w, response.FailMessageResp(err.Error()))
			return
		}

		tempInfo, err := utils.GetFileInfo(tempFile)
		if err != nil {
			response.WriteJson(w, response.FailMessageResp(err.Error()))
			return
		}
		contentLength := strconv.FormatInt(tempInfo.Size(), 10)
		zipDownloadName := filepath.Base(tempInfo.Name())

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Disposition", "attachment; filename*=utf-8''"+zipDownloadName)
		//http.ServeFile(w, r, filePath)
		// 将临时文件的内容写入 HTTP 响应体
		http.ServeContent(w, r, zipDownloadName, tempInfo.ModTime(), tempFile)
	} else {
		// 打开文件，使用http.ServeContent方式下载
		/*file, err := os.Open(filePath)
		if err != nil {
			log.Println("打开文件异常", err.Error())
			response.WriteJson(w, response.FailMessageResp("文件不存在"))
			return
		}
		defer file.Close()

		fileInfo, err := utils.GetFileInfo(file)
		if err != nil {
			response.WriteJson(w, response.FailMessageResp(err.Error()))
			return
		}
		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)*/
		// 如果不需要压缩，则直接提供文件下载
		http.ServeFile(w, r, filePath)
		return
	}
}

const (
	volumeSize = 1024 * 1024 * 10 // 每个分卷的大小，这里以10MB为例
)

func handler4(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}
	// 文件路径拼接
	filePath := filepath.Join("./upload", fileName)

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()

	// 创建压缩文件
	zipFileName := "multipart.zip"
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// 创建压缩写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 将文件添加到压缩文件中
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		panic(err)
	}
	header.Name = fileName
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(writer, file)
	if err != nil {
		panic(err)
	}

	// 计算分卷数量
	numVolumes := (fileSize + volumeSize - 1) / volumeSize

	// 拆分压缩文件为分卷
	for i := int64(1); i <= numVolumes; i++ {
		volumeFileName := fmt.Sprintf("multipart_%d.zip", i)
		volumeFile, err := os.Create(volumeFileName)
		if err != nil {
			panic(err)
		}
		defer volumeFile.Close()

		// 创建分卷压缩文件
		volumeWriter := zip.NewWriter(volumeFile)

		// 添加分卷文件
		header.Name = fmt.Sprintf("%d_%s", i, fileName)
		writer, err := volumeWriter.CreateHeader(header)
		if err != nil {
			panic(err)
		}
		_, err = file.Seek((i-1)*volumeSize, io.SeekStart)
		if err != nil {
			panic(err)
		}
		_, err = io.CopyN(writer, file, volumeSize)
		if err != nil && err != io.EOF {
			panic(err)
		}

		// 关闭分卷文件
		err = volumeWriter.Close()
		if err != nil {
			panic(err)
		}

		fmt.Printf("第 %d 个分卷创建成功\n", i)
	}

	fmt.Println("所有分卷创建完成")
}
