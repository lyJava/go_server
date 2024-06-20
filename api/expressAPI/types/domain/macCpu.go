package domain

// MacCpu 苹果处理器
type MacCpu struct {
	Id                    int64  `json:"id,omitempty"`                    // 主键ID
	CpuName               string `json:"cpuName,omitempty"`               // 处理器名称
	CpuType               string `json:"cpuType,omitempty"`               // 处理器类型枚举(Intel或者Apple Silicon)
	CpuBasicBoost         string `json:"cpuBasicBoost,omitempty"`         // 处理器基本频率
	CpuTruboBoost         string `json:"cpuTruboBoost,omitempty"`         // 处理器Trubo频率
	CpuCoreNumber         int32  `json:"cpuCoreNumber,omitempty"`         // 处理器核心数
	CpuThreadNumber       int32  `json:"cpuThreadNumber,omitempty"`       // 处理器线程数
	CpuCache              string `json:"cpuCache,omitempty"`              // 处理器缓存
	CpuTdp                string `json:"cpuTdp,omitempty"`                // 处理器TDP功耗
	MemoryWidth           string `json:"memoryWidth,omitempty"`           // 内存带宽
	MediaProcessingEngine string `json:"mediaProcessingEngine,omitempty"` // 媒体处理引擎
}
