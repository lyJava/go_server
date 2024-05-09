package service

import "apiProject/api/expressAPI/types/domain"

// DictTypeService 字典类型接口
type DictTypeService interface {
	// GetDictList 分页查询
	GetDictList(d *domain.DictType, page, size string) ([]*domain.DictType, int64, int64, error)

	// SelectDictTypeById 查询字典类型
	SelectDictTypeById(id int64) (*domain.DictType, error)

	// SaveDictType 新增字典类型
	SaveDictType(d *domain.DictType) (*domain.DictType, error)
}
