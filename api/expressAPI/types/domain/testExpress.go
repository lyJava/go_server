package domain

// TestExpress 测试快递结构(postgresql的表)
//
// 设置了omitempty,字段的值为空值时，JSON序列化时将忽略该字段，因此不会显示空值。
type TestExpress struct {
	ID               *int64 `json:"id,omitempty"`               // 主键ID
	ExpressName      string `json:"expressName,omitempty"`      // 快递名称
	ExpressNumber    string `json:"expressNumber,omitempty"`    // 快递单号
	PickupCode       string `json:"pickupCode,omitempty"`       // 取件码
	FromUsername     string `json:"fromUsername,omitempty"`     // 发件人姓名
	FromUserPhone    string `json:"fromUserPhone,omitempty"`    // 发件人手机
	FromUserAddress  string `json:"fromUserAddress,omitempty"`  // 发件人地址
	FromUserIdNumber string `json:"fromUserIdNumber,omitempty"` // 发件人身份证号
	CreateBy         string `json:"createBy,omitempty"`         // 创建人
	CreateTime       string `json:"createTime,omitempty"`       // 创建时间
	UpdateBy         string `json:"updateBy,omitempty"`         // 修改人
	UpdateTime       string `json:"updateTime,omitempty"`       // 更新时间
	Remarks          string `json:"remarks,omitempty"`          // 备注
	DelFlag          int    `json:"delFlag,omitempty"`          // 是否删除(0:正常；1:删除)
}
