package common

// MailConfig 存储邮件相关的配置信息
type MailConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
}

// MailContent 存储邮件的内容信息
type MailContent struct {
	To             []string
	Cc             []string
	Bcc            []string
	Subject        string
	Body           string
	AttachmentPath []string // 你的附件路径
}
