package param

// ExpressSearchParam 分页查询参数结构
type ExpressSearchParam struct {
	ExpressName   string `json:"expressName"`   // 快递名称
	ExpressNumber string `json:"expressNumber"` // 快递单号
	UserId        string `json:"userId"`        // 用户ID
	FromName      string `json:"fromName"`      // 发件人姓名
	FromPhone     string `json:"fromPhone"`     // 发件人手机
	FromAddress   string `json:"fromAddress"`   // 发件人地址
	PickupCode    string `json:"pickupCode"`    // 取件码
	CreateBy      string `json:"createBy"`      // 创建人
	Page          int64  `json:"page"`          // 当前页码
	Size          int64  `json:"size"`          // 每页条数
	Column        string `json:"column"`        // 排序的列
	Order         string `json:"order"`         // 降序或者升序
}

// StoreAdminSearchParam 分页查询参数结构
type StoreAdminSearchParam struct {
	UserName     string `json:"userName"`
	Mobile       string `json:"mobile"`
	StoreName    string `json:"storeName"`
	StatusValue  string `json:"statusValue"`
	RealName     string `json:"realName"`
	MerchantId   int64  `json:"merchantId"`
	MerchantName string `json:"merchantName"`
	Page         int64  `json:"page"`   // 当前页码
	Size         int64  `json:"size"`   // 每页条数
	Column       string `json:"column"` // 排序的列
	Order        string `json:"order"`  // 降序或者升序
}

// Payment 支付结构体
type Payment struct {
	Id       string `json:"id"`       // 支付ID
	Quantity int64  `json:"quantity"` // 支付金额(分)
}

// Order 订单结构体
type Order struct {
	Id    string    `json:"id"`    // 订单ID
	Items []Payment `json:"items"` // 支付数组
}

// MacBookPageParam 苹果本分页查询参数
type MacBookPageParam struct {
	Page      int64  `json:"page"`      // 当前页码
	Size      int64  `json:"size"`      // 每页条数
	Column    string `json:"column"`    // 排序的列
	Order     string `json:"order"`     // 降序或者升序
	ProName   string `json:"proName"`   // 名称
	ProColor  string `json:"proColor"`  // 颜色
	ProYear   string `json:"proYear"`   // 年份
	ProType   string `json:"proType"`   // 类型
	MoldModel string `json:"moldModel"` // 型号
}
