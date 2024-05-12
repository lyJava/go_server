package send

import (
	"apiProject/email/common"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/gomail.v2"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SendMail(config common.MailConfig, content common.MailContent) {
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
	BuildHeaders(&message, headers)

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
		err := AppendAttachToMessage(&message, writer, attachmentPath)
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

// AppendAttachToMessage 添加附件到邮箱信息体中
//
// 参数
//
//	message (*bytes.Buffer): 消息buff对象
//	writer (*multipart.Writer): 附件写入处理器
//	filePath (string): 附件路径
func AppendAttachToMessage(message *bytes.Buffer, writer *multipart.Writer, filePath string) (err error) {
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

func BuildHeaders(w *bytes.Buffer, headers map[string]string) {
	for key, value := range headers {
		w.WriteString(key + ": " + value + "\r\n")
	}
	w.WriteString("\r\n")
}

func SendMail2(config common.MailConfig, content common.MailContent) {
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	d := gomail.NewDialer(config.Host, port, config.Email, config.Password)

	// 创建邮件
	mail := gomail.NewMessage()

	mail.SetHeader("From", config.Email)
	mail.SetHeader("To", content.To...)
	mail.SetHeader("Cc", content.Cc...)
	mail.SetHeader("Bcc", content.Bcc...)
	mail.SetHeader("Subject", content.Subject)

	// 添加 HTML 内容
	mail.AddAlternative("text/html;charset=utf-8", content.Body+"<h3>"+time.Now().Format("2006-01-02 15:04:05")+"</h3>")

	imagePath := "/Users/yangge/Pictures/vlcsnap-2024-04-07-20h31m31s596.png"
	fileContent, err := os.ReadFile(imagePath)
	if err != nil {
		log.Printf("Failed to read image file: %v", err)
		return
	}

	// 检测图片的MIME类型
	imageMime, err := mimetype.DetectFile(imagePath)
	if err != nil {
		log.Printf("Failed to detect MIME type for imagePath %s: %v", imagePath, err)
		return
	}
	// 获取图片base64编码
	encodedImage := base64.StdEncoding.EncodeToString(fileContent)
	// 拼接img标签
	// mail.AddAlternative("text/html;charset=utf-8", "<img src=\"data:"+imageMime.String()+";base64,"+encodedImage+"\" alt=\"img\"/>")
	imgTag := fmt.Sprintf("<div style=\"margin: 0 auto\"><img src=\"data:%s;base64,%s\" alt=\"img\"/></div>", imageMime, encodedImage)
	mail.AddAlternative("text/html;charset=utf-8", imgTag)
	// 附件部分
	for _, attachmentPath := range content.AttachmentPath {
		mail.Attach(attachmentPath)
	}

	startTime := time.Now()

	// 发送邮件
	if err := d.DialAndSend(mail); err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}
	elapsedTime := time.Since(startTime).Seconds()
	log.Printf("gmail方式发送邮件成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
}

func TestSend() {
	m := gomail.NewMessage()
	m.SetHeader("From", "745876299@qq.com")
	//m.SetHeader("To", "qiujiahongde@163.com", "mail12@163.com")  //发送多个人
	m.SetHeader("To", "745876299@qq.com")
	//m.SetHeader("Cc", "qiujiahongde@163.com") //抄送
	//m.SetHeader("Bcc", "309284701@qq.com")    // 密送
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	//发送html格式邮件。
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>! <p style='color:red'>red </p>")
	m.Attach("/Users/yangge/Downloads/测试中文小图.jpg") //添加附件
	mail := gomail.NewDialer("smtp.qq.com", 25, "745876299@qq.com", "mxexfejfdcmhbfbb")
	// Send the email to Bob, Cora and Dan.
	if err := mail.DialAndSend(m); err != nil {
		panic(err)
	}
}
