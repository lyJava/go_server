package domain

type MacBook struct {
	Id          int64     `json:"id,omitempty"`          // 主键ID
	ProName     string    `json:"proName,omitempty"`     // 名称
	ProColor    string    `json:"proColor,omitempty"`    // 外观颜色
	ProYear     string    `json:"proYear,omitempty"`     // 年份
	ProType     string    `json:"proType,omitempty"`     // 类型
	MoldModel   string    `json:"moldModel,omitempty"`   // 模具型号
	MonitorId   int64     `json:"monitorId,omitempty"`   // 显示器ID
	CpuId       int64     `json:"cpuId,omitempty"`       // 处理器ID
	MemoryId    int64     `json:"memoryId,omitempty"`    // 内存ID
	StorageInfo string    `json:"storageInfo,omitempty"` // 存储信息
	SizeInfo    string    `json:"sizeInfo,omitempty"`    // 尺寸
	WeightInfo  string    `json:"weightInfo,omitempty"`  // 重量
	GpuId       int64     `json:"gpuId,omitempty"`       // 图形处理器ID
	ChipType    string    `json:"chipType,omitempty"`    // 芯片类型
	MemoryMax   int64     `json:"memoryMax,omitempty"`   // 最大支持内存
	CreateTime  string    `json:"createTime,omitempty"`  // 创建时间
	UpdateTime  string    `json:"updateTime,omitempty"`  // 更新时间
	CpuInfo     MacCpu    `json:"cpuInfo,omitempty"`     // 处理器信息
	MemoryInfo  MacMemory `json:"memoryInfo,omitempty"`  // 内存信息
}
