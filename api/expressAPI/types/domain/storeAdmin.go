package domain

// StoreAdmin 商店管理员
type StoreAdmin struct {
	Id           int64  `json:"id,omitempty"`           // 主键ID
	UserName     string `json:"userName,omitempty"`     // 名称
	Mobile       string `json:"mobile,omitempty"`       // 手机
	RealName     string `json:"realName,omitempty"`     // 真实姓名
	StatusValue  string `json:"statusValue,omitempty"`  // 状态值
	StoreName    string `json:"storeName,omitempty"`    // 店铺名称
	MerchantId   int64  `json:"merchantId,omitempty"`   // 租户ID
	MerchantName string `json:"merchantName,omitempty"` // 租户名称
	CreateTime   string `json:"createTime,omitempty"`   // 创建时间
	UpdateTime   string `json:"updateTime,omitempty"`   // 修改时间
}
