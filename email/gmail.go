package main

import (
	"apiProject/email/common"
	"encoding/base64"
	"fmt"
	_ "fmt"
	"github.com/gabriel-vasile/mimetype"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	_ "os"
	"strconv"
	"time"
)

// https://github.com/go-gomail/gomail
func testSend() {
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

func main() {
	config := common.MailConfig{
		Host:     "smtp.qq.com",
		Port:     "25",
		Email:    "745876299@qq.com",
		Password: "mxexfejfdcmhbfbb",
	}

	content := common.MailContent{
		To:      []string{"1179028989@qq.com"},
		Cc:      []string{"745876299@qq.com"},
		Bcc:     []string{},
		Subject: "Test Subject Gmail",
		Body:    "这是来自gmail发送的邮件",
		AttachmentPath: []string{
			//"/Users/yangge/Downloads/1713032174133.jpg",
			"/Users/yangge/Downloads/测试中文小图.jpg",
			"/Users/yangge/Downloads/c548a7d37d4f27e4d14ca6941d11392c.mp4",
		},
	}

	sendMail2(config, content)
}

func sendMail2(config common.MailConfig, content common.MailContent) {
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
