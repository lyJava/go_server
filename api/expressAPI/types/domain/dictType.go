package domain

type DictType struct {
	Id         int64  `json:"userId,omitempty"`
	DictName   string `json:"dictName,omitempty"`
	DictType   string `json:"dictType,omitempty"`
	TypeStatus string `json:"typeStatus,omitempty"`
	CreateBy   string `json:"createBy,omitempty"`
	CreateTime string `json:"createTime,omitempty"`
	UpdateBy   string `json:"updateBy,omitempty"`
	UpdateTime string `json:"updateTime,omitempty"`
	Remark     string `json:"remark,omitempty"`
	DelFlag    string `json:"delFlag,omitempty"`
}
