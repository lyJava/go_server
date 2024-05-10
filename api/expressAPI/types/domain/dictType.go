package domain

type DictType struct {
	Id         *int64 `json:"id,omitempty"`         // 主键ID
	DictName   string `json:"dictName,omitempty"`   // 字典类型名称
	DictType   string `json:"dictType,omitempty"`   // 字典类型
	TypeStatus string `json:"typeStatus,omitempty"` // 是否启用(0:正常；1:停用)
	CreateBy   string `json:"createBy,omitempty"`   // 创建人
	CreateTime string `json:"createTime,omitempty"` // 创建时间
	UpdateBy   string `json:"updateBy,omitempty"`   // 修改人
	UpdateTime string `json:"updateTime,omitempty"` // 修改时间
	Remark     string `json:"remark,omitempty"`     // 备注
	DelFlag    string `json:"delFlag,omitempty"`    // 是否删除(0:正常；1:删除)
}

// NewDictTypeDetail 新增字典类型详情
func NewDictTypeDetail(typeStr string) DictType {
	return DictType{
		DictType: typeStr,
	}
}
