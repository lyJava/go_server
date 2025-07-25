package mailSend

import (
	"apiProject/email/common"
	"bytes"
	"encoding/base64"
	"fmt"
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

	"github.com/gabriel-vasile/mimetype"
	"github.com/jordan-wright/email"
	"gopkg.in/gomail.v2"
)

type User struct {
	Name     string
	Position string
	Email    string
	Status   string // 或者使用一个自定义类型，但在模板中我们使用字符串
}

type ServiceStatus struct {
	Name       string
	Status     string // 可以是“运行中”、“警告”、“停止”等
	Usage      int    // 使用率百分比
	LastUpdate string // 最后更新时间
}

type TemplateData struct {
	Body        string
	Timestamp   string
	ImageMine   string
	ImageData   string
	ImageUrl    string
	ServiceList []ServiceStatus
	UserList    []User
}

// SendMailBySmtp 发送邮件
func SendMailBySmtp(config common.MailConfig, content common.MailContent) error {
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

	tpl, err := ParseTemplate("./email_template.html")
	if err != nil {
		log.Printf("Failed to parse template file: %+v", err)
		return err
	}

	imagePath := "/Users/yangge/Pictures/vlcsnap-2024-04-07-20h31m31s596.png"
	imageMine, base64Data, err := GetImgMineAndContent(imagePath)
	if err != nil {
		log.Printf("Failed to detect MIME type for imagePath %s: %+v", imagePath, err)
		return err
	}

	// 创建一个buffer来保存渲染后的HTML内容
	var bodyBuffer bytes.Buffer
	data := &TemplateData{
		Body:        content.Body,
		Timestamp:   time.Now().Format("2006-01-02 15:04:05"),
		ImageMine:   imageMine,
		ImageData:   base64Data, // 获取图片base64编码
		ImageUrl:    "",
		UserList:    nil,
		ServiceList: nil,
	}

	// 渲染HTML模板
	if err = tpl.Execute(&bodyBuffer, data); err != nil {
		log.Printf("Failed to execute bodyBuffer: %+v", err)
		return err
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

	if err := writer.Close(); err != nil {
		log.Printf("Failed to close multipart writer: %v", err)
		return err
	}

	startTime := time.Now()
	port := strconv.Itoa(config.Port)

	if err := smtp.SendMail(config.Host+":"+port, auth, config.Email, append(content.To, content.Cc...), message.Bytes()); err != nil {
		if strings.Contains(err.Error(), "short response") {
			log.Printf("⚠️ 服务器短响应但邮件可能已发送成功: %v", err)
			// 选择性返回 nil，表示你要“容忍”这类错误
			elapsedTime := time.Since(startTime).Seconds()
			log.Printf("smtp邮件发送成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
			return nil
		}
		log.Printf("发送邮件失败: %+v", err)
		return err
	}
	elapsedTime := time.Since(startTime).Seconds()
	log.Printf("smtp邮件发送成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
	return nil
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

func SendMailByGmail(config common.MailConfig, content common.MailContent) error {

	d := gomail.NewDialer(config.Host, int(config.Port), config.Email, config.Password)

	// 创建邮件
	mail := gomail.NewMessage()

	mail.SetHeader("From", config.Email)
	mail.SetHeader("To", content.To...)
	mail.SetHeader("Cc", content.Cc...)
	mail.SetHeader("Bcc", content.Bcc...)
	mail.SetHeader("Subject", content.Subject)

	// 添加 HTML 内容
	mail.AddAlternative("text/html;charset=utf-8", content.Body+"<h3>"+time.Now().Format("2006-01-02 15:04:05")+"</h3>")

	imageMime, base64Data, err := GetImgMineAndContent("/Users/yangge/Pictures/vlcsnap-2024-04-07-20h31m31s596.png")
	if err != nil {
		return err
	}

	// 拼接img标签
	// mail.AddAlternative("text/html;charset=utf-8", "<img src=\"data:"+imageMime.String()+";base64,"+encodedImage+"\" alt=\"img\"/>")
	imgTag := fmt.Sprintf("<div style=\"margin: 0 auto\"><img src=\"data:%s;base64,%s\" alt=\"img\"/></div>", imageMime, base64Data)
	mail.AddAlternative("text/html;charset=utf-8", imgTag)
	// 附件部分
	for _, attachmentPath := range content.AttachmentPath {
		mail.Attach(attachmentPath)
	}

	startTime := time.Now()

	// 发送邮件
	if err := d.DialAndSend(mail); err != nil {
		log.Fatalf("Failed to mailSend email: %v", err)
		return err
	}
	elapsedTime := time.Since(startTime).Seconds()
	log.Printf("gmail方式发送邮件成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
	return nil
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
	//添加附件
	m.Attach("/Users/yangge/Downloads/测试中文小图.jpg")
	mail := gomail.NewDialer("smtp.qq.com", 25, "745876299@qq.com", "mxexfejfdcmhbfbb")
	// Send the email to Bob, Cora and Dan.
	if err := mail.DialAndSend(m); err != nil {
		panic(err)
	}
}

func SendEmailByJordan(config common.MailConfig, content common.MailContent) error {
	e := email.NewEmail()
	e.From = "745876299@qq.com"
	e.To = content.To
	e.Cc = content.Cc
	e.Subject = "来自Jordan的邮件"
	e.Text = []byte(`这是来自的邮件=====你好`)
	//e.HTML = []byte(`<h1>这是来自 <b>Go-Jordan</b> 的邮件</h1><p>你好！</p>`)

	// 1. 解析HTML模板
	tpl, err := ParseTemplate("./email_template.html")
	if err != nil {
		return err
	}

	imgMine, base64Data, err := GetImgMineAndContent("/Users/yangge/Downloads/wallpaper-5045169.jpg")
	if err != nil {
		log.Printf("解析模板失败: %v", err)
		return err
	}
	// 2. 准备模板数据
	data := &TemplateData{
		Body:      "尊贵的用户",
		Timestamp: time.Now().Format("2006年1月2日 15:04"),
		ImageMine: imgMine,
		ImageData: base64Data,
		ImageUrl:  "https://cdn.pixabay.com/photo/2025/06/09/16/27/animal-9650392_1280.jpg",
		ServiceList: []ServiceStatus{
			{Name: "邮件发送服务", Status: "运行中", Usage: 75, LastUpdate: time.Now().Format("2006-01-02 15:04:05")},
			{Name: "数据库服务", Status: "警告", Usage: 92, LastUpdate: "2025-07-24 14:23:10"},
			{Name: "文件存储服务", Status: "运行中", Usage: 42, LastUpdate: "2025-07-25 08:45:32"},
			{Name: "任务调度服务", Status: "停止", Usage: 0, LastUpdate: "2025-07-20 19:12:57"},
		},
		UserList: []User{
			{Name: "张云", Position: "系统架构师", Email: "zhangyun@example.com", Status: "在线"},
			{Name: "李思雨", Position: "前端工程师", Email: "lisyu@example.com", Status: "在线"},
			{Name: "王建国", Position: "后端工程师", Email: "wangjg@example.com", Status: "休假"},
			{Name: "陈婷婷", Position: "产品经理", Email: "chentt@example.com", Status: "在线"},
		},
	}

	// 3. 渲染HTML内容
	var htmlBody bytes.Buffer
	if err := tpl.Execute(&htmlBody, data); err != nil {
		log.Printf("渲染模板失败: %v", err)
		return err
	}

	// 4. 设置HTML内容
	e.HTML = htmlBody.Bytes()

	// 5. 纯文本备用内容（可选）
	e.Text = fmt.Appendf(nil, "你好！%s\n%s\n发送时间：%s",
		content.Subject,
		content.Body,
		data.Timestamp)

	// 添加附件
	_, err1 := e.AttachFile("/Users/yangge/Downloads/123-small.jpg")
	_, err2 := e.AttachFile("/Users/yangge/Downloads/c548a7d37d4f27e4d14ca6941d11392c.mp4")
	if err1 != nil || err2 != nil {
		log.Printf("附件添加失败: %v %v", err1, err2)
		return fmt.Errorf("jordan方式添加邮件失败,err1==%+v, err2==%+v", err1, err2)
	}

	startTime := time.Now()
	// 发送邮件（使用 StartTLS）
	if err = e.Send("smtp.qq.com:587", smtp.PlainAuth("", "745876299@qq.com", "mxexfejfdcmhbfbb", "smtp.qq.com")); err != nil {
		if strings.Contains(err.Error(), "short response") {
			log.Printf("⚠️ 服务器短响应但邮件可能已发送成功: %v", err)
			// 选择性返回 nil，表示你要“容忍”这类错误
			elapsedTime := time.Since(startTime).Seconds()
			log.Printf("jordan方式发送邮件成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
			return nil
		}
		log.Printf("jordan方式发送邮件失败: %v", err)
		return err
	}
	elapsedTime := time.Since(startTime).Seconds()
	log.Printf("jordan方式发送邮件成功！耗时： %s 秒", fmt.Sprintf("%.2f", elapsedTime))
	return nil
}

func ParseTemplate(tplPath string) (*template.Template, error) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("解析模板失败: %v", err)
		return nil, err
	}
	return tpl, nil
}

// GetImgMine 检测图片的MIME类型与内容
func GetImgMineAndContent(imagePath string) (string, string, error) {
	if imagePath == "" {
		return "", "", nil
	}

	imageMime, err := mimetype.DetectFile(imagePath)
	if err != nil {
		log.Printf("Failed to detect MIME type for imagePath %s: %v", imagePath, err)
		return "", "", err
	}

	content, err := os.ReadFile(imagePath)
	if err != nil {
		log.Printf("Failed to read image file: %+v", err)
		return "", "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(content)

	return imageMime.String(), base64Data, nil
}
