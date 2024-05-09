package domain

// Express 快递对象结构
//
// 设置了omitempty,字段的值为空值时，JSON序列化时将忽略该字段，因此不会显示空值。
type Express struct {
	ID            int64  `json:"id"`                      // 主键ID
	UserId        string `json:"userId,omitempty"`        // 用户ID
	ExpressName   string `json:"expressName,omitempty"`   // 快递名称
	ExpressNumber string `json:"expressNumber,omitempty"` // 快递单号
	FromName      string `json:"fromName,omitempty"`      // 发件人姓名
	FromPhone     string `json:"fromPhone,omitempty"`     // 发件人手机
	FromAddress   string `json:"fromAddress,omitempty"`   // 发件人地址
	PickupCode    string `json:"pickupCode,omitempty"`    // 取件码
	CreateBy      string `json:"createBy,omitempty"`      // 创建人
	CreateTime    string `json:"createTime,omitempty"`    // 创建时间
	UpdateTime    string `json:"updateTime,omitempty"`    // 更新时间
}
