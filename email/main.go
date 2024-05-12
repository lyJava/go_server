package main

import (
	"apiProject/email/common"
	"apiProject/email/send"
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
	send.SendMail(config, content)

	config2 := common.MailConfig{
		Host:     "smtp.qq.com",
		Port:     "25",
		Email:    "745876299@qq.com",
		Password: "mxexfejfdcmhbfbb",
	}

	content2 := common.MailContent{
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

	send.SendMail2(config2, content2)
}
