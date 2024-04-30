package main

import (
	"apiProject/email/common"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {

	config := common.MailConfig{
		Host:     "smtp.qq.com",
		Port:     "25",
		Email:    "745876299@qq.com",
		Password: "mxexfejfdcmhbfbb",
	}
	content := common.MailContent{
		To:      []string{"745876299@qq.com"},
		Cc:      []string{"745876299@qq.com"},
		Bcc:     []string{},
		Subject: "Test Subject Smtp",
		Body:    "这是来自go smtp发送的邮件",
		AttachmentPath: []string{
			"/Users/yangge/Downloads/测试中文小图.jpg",
			"/Users/yangge/Downloads/c548a7d37d4f27e4d14ca6941d11392c.mp4",
		},
	}
	//
	//"/Users/yangge/Downloads/雪落黄山 _ 当霜染一半山头, 风也不再轻柔｜8K超清.mp4",
	// /Users/yangge/Downloads/20240215005635413-Screenrecorder-2024-02-15-00-54-00-961.mp4
	sendMail(config, content)
}

func buildHeaders(w *bytes.Buffer, headers map[string]string) {
	for key, value := range headers {
		w.WriteString(key + ": " + value + "\r\n")
	}
	w.WriteString("\r\n")
}

func sendMail(config common.MailConfig, content common.MailContent) {
	// 连接到SMTP服务器
	auth := smtp.PlainAuth("", config.Email, config.Password, config.Host)

	var message bytes.Buffer
	writer := multipart.NewWriter(&message)

	// 设置邮件头部
	headers := map[string]string{
		"From":         config.Email,
		"To":           strings.Join(content.To, ","),
		"Cc":           strings.Join(content.Cc, ","),
		"Bcc":          strings.Join(content.Bcc, ","),
		"Subject":      content.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "multipart/mixed; boundary=" + writer.Boundary(),
	}

	// 设置邮件头部
	buildHeaders(&message, headers)

	// 正文部分
	//message.WriteString("--" + writer.Boundary() + "\r\n")
	//message.WriteString("Content-Type:text/html;charset=utf-8\r\n")
	//message.WriteString("\r\n")
	//message.WriteString("<html><body>")
	//message.WriteString(content.Body + "<h3>" + time.Now().Format("2006-01-02 15:04:05") + "</h3>")
	//message.WriteString("</body></html>\r\n")
	//message.WriteString("\r\n")

	// 读取HTML模板文件

	tmpl, err := template.ParseFiles("./email/email_template.html")
	if err != nil {
		log.Printf("Failed to parse template file: %v", err)
		return
	}

	imagePath := "/Users/yangge/Pictures/vlcsnap-2024-04-07-20h31m31s596.png"
	fileContent, err := os.ReadFile(imagePath)
	if err != nil {
		log.Printf("Failed to read image file: %v", err)
		return
	}

	// 检测图片的MIME类型
	imageMine, err := mimetype.DetectFile(imagePath)
	if err != nil {
		log.Printf("Failed to detect MIME type for imagePath %s: %v", imagePath, err)
		return
	}
	// 获取图片base64编码
	encodedImage := base64.StdEncoding.EncodeToString(fileContent)

	// 创建一个buffer来保存渲染后的HTML内容
	var bodyBuffer bytes.Buffer
	data := struct {
		Body      string
		Timestamp string
		ImageMine string
		ImageData string
	}{
		Body:      content.Body,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		ImageMine: imageMine.String(),
		ImageData: encodedImage,
	}

	// 渲染HTML模板
	err = tmpl.Execute(&bodyBuffer, data)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		return
	}

	// 添加HTML版本的邮件正文
	message.WriteString("--" + writer.Boundary() + "\r\n")
	message.WriteString("Content-Type:text/html;charset=utf-8\r\n")
	message.WriteString("\r\n")
	message.WriteString(bodyBuffer.String() + "\r\n")

	// 附件部分
	//for _, attachmentPath := range content.AttachmentPath {
	//	mimeType, err := mimetype.DetectFile(attachmentPath)
	//	if err != nil {
	//		log.Printf("Failed to detect MIME type for %s: %v", attachmentPath, err)
	//		continue
	//	}
	//
	//	fileContent, err := os.ReadFile(attachmentPath)
	//	if err != nil {
	//		log.Printf("Failed to read attachment %s: %v", attachmentPath, err)
	//		continue
	//	}
	//
	//	_, filename := filepath.Split(attachmentPath)
	//	// 进行URL编码
	//	encodedFilename := url.QueryEscape(filename)
	//
	//	message.WriteString("--" + writer.Boundary() + "\r\n")
	//	message.WriteString("Content-Disposition: attachment; filename*=utf-8''" + encodedFilename + "\r\n") // 使用utf-8编码
	//	message.WriteString("Content-Type: " + mimeType.String() + "\r\n")
	//	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	//	message.WriteString("\r\n")
	//	message.WriteString(base64.StdEncoding.EncodeToString(fileContent))
	//	message.WriteString("\r\n")
	//}

	// 附件部分
	for _, attachmentPath := range content.AttachmentPath {
		err := appendAttachToMessage(&message, writer, attachmentPath)
		if err != nil {
			log.Printf("Failed to attach file %s: %v", attachmentPath, err)
			continue
		}
	}

	message.WriteString("--" + writer.Boundary() + "--\r\n")

	err5 := writer.Close()
	if err5 != nil {
		log.Printf("Failed to close multipart writer: %v", err5)
		return
	}

	startTime := time.Now()
	err3 := smtp.SendMail(config.Host+":"+config.Port, auth, config.Email, append(content.To, content.Cc...), message.Bytes())
	if err3 != nil {
		log.Printf("发送邮件失败: %v", err3)
		//log.Fatal(err) // 发送邮件失败时，整个程序会直接退出
		//return err
	}

	elapsedTime := time.Since(startTime).Seconds()
	log.Printf("邮件发送耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
	//return nil
}

// appendAttachToMessage 添加附件到邮箱信息体中
//
// 参数
//
//	message (*bytes.Buffer): 消息buff对象
//	writer (*multipart.Writer): 附件写入处理器
//	filePath (string): 附件路径
func appendAttachToMessage(message *bytes.Buffer, writer *multipart.Writer, filePath string) (err error) {
	mimeType, err := mimetype.DetectFile(filePath)
	if err != nil {
		message.WriteString("Content-Type:application/octet-stream" + "\r\n")
	}

	//fileContent, err := os.ReadFile(filePath)
	//if err != nil {
	//	log.Printf("Failed to read attachment %s: %v", filePath, err)
	//	return
	//}

	openFile, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open attachment %s: %v", filePath, err)
		return
	}
	defer openFile.Close()
	// 读取附件方式一
	//var buffer bytes.Buffer
	//_, err = io.Copy(&buffer, openFile)
	//if err != nil {
	//	log.Printf("Failed to read attachment %s: %v", filePath, err)
	//	return
	//}
	//
	//fileContent := buffer.Bytes()

	// 读取附件方式二
	var fileContent []byte
	buffer := make([]byte, 2048)

	for {
		n, err := io.ReadFull(openFile, buffer)
		if n > 0 {
			fileContent = append(fileContent, buffer[:n]...)
		}
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				log.Printf("Failed to read attachment %s: %v", filePath, err)
			}
			break
		}
	}

	// 读取附件方式三
	//bufferedReader := bufio.NewReader(openFile)
	//var fileContent bytes.Buffer
	//
	//buffer := make([]byte, 2048)
	//for {
	//	n, err := bufferedReader.Read(buffer)
	//	if n > 0 {
	//		fileContent.Write(buffer[:n])
	//	}
	//	if err != nil {
	//		if err != io.EOF {
	//			log.Printf("Failed to read file %s: %v", filePath, err)
	//		}
	//		break
	//	}
	//}

	_, file := filepath.Split(filePath)
	encodedFilename := url.QueryEscape(file)

	message.WriteString("--" + writer.Boundary() + "\r\n")
	message.WriteString("Content-Disposition: attachment; filename*=utf-8''" + encodedFilename + "\r\n") // 使用utf-8编码
	message.WriteString("Content-Type: " + mimeType.String() + "\r\n")
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("\r\n")
	message.WriteString(base64.StdEncoding.EncodeToString(fileContent))
	message.WriteString("\r\n")

	return nil
}
