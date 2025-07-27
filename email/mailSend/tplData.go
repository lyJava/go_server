package mailSend

// MailConfig 存储邮件相关的配置信息
type MailConfig struct {
	From     string
	Host     string
	Port     int
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
