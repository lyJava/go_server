package domain

// MacMemory 苹果内存结构体
type MacMemory struct {
	Id              int64  `json:"id,omitempty"`              // 主键ID
	MemorySize      string `json:"memorySize,omitempty"`      // 内存大小(GB)
	MemorySpeed     string `json:"memorySpeed,omitempty"`     // 内存频率(MHz)
	IntegrationFlag string `json:"integrationFlag,omitempty"` // 是否集成
	MemoryType      string `json:"memoryType,omitempty"`      // 内存类型
	EccCheck        int32  `json:"eccCheck,omitempty"`        // 是否支持ECC校验
}
