package main

import (
	"apiProject/email/common"
	"apiProject/email/mailSend"
	"log"
)

func main() {

	config := common.MailConfig{
		Host:     "smtp.qq.com",
		Port:     587,
		Email:    "745876299@qq.com",
		Password: "mxexfejfdcmhbfbb",
	}
	content := common.MailContent{
		To:      []string{"745876299@qq.com"},
		Cc:      []string{"745876299@qq.com"},
		Bcc:     []string{},
		Subject: "使用smtp发送的邮件",
		Body:    "这是来自go smtp发送的邮件",
		AttachmentPath: []string{
			"/Users/yangge/Downloads/测试中文小图.jpg",
			"/Users/yangge/Downloads/c548a7d37d4f27e4d14ca6941d11392c.mp4",
		},
	} 
	//"/Users/yangge/Downloads/雪落黄山 _ 当霜染一半山头, 风也不再轻柔｜8K超清.mp4",
	// /Users/yangge/Downloads/20240215005635413-Screenrecorder-2024-02-15-00-54-00-961.mp4
	
	if err := mailSend.SendMailBySmtp(config, content); err!= nil {
		log.Fatalf("smtp发送邮件失败===%+v", err)
	}

	if err := mailSend.SendEmailByJordan(config, content); err!= nil {
		log.Fatalf("jordan发送邮件失败===%+v", err)
	}

	config2 := common.MailConfig{
		Host:     "smtp.qq.com",
		Port:     25,
		Email:    "745876299@qq.com",
		Password: "mxexfejfdcmhbfbb",
	}

	content2 := common.MailContent{
		To:      []string{"745876299@qq.com"},
		Cc:      []string{"745876299@qq.com"},
		Bcc:     []string{},
		Subject: "使用gomail发送邮件",
		Body:    "这是来自gomail发送的邮件",
		AttachmentPath: []string{
			//"/Users/yangge/Downloads/1713032174133.jpg",
			"/Users/yangge/Downloads/测试中文小图.jpg",
			"/Users/yangge/Downloads/c548a7d37d4f27e4d14ca6941d11392c.mp4",
		},
	}

	if err := mailSend.SendMailByGmail(config2, content2); err!= nil {
		log.Fatalf("gomail发送邮件失败===%+v", err)
	} 
}
