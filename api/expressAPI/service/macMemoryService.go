package service

import "apiProject/api/expressAPI/types/domain"

type MacMemoryService interface {
	// Save 保存新增
	Save(memory *domain.MacMemory) (*domain.MacMemory, error)
	// BatchSave 批量保存新增
	BatchSave(memoryList []*domain.MacMemory) (int64, error)
	// SelectById 通过主键ID查询
	SelectById(id int64) (*domain.MacMemory, error)
	// Update 修改/更新
	Update(memory *domain.MacMemory) (*domain.MacMemory, error)
	// DeleteById 通过主键ID删除
	DeleteById(id int64) (int64, error)
	// BatchDeleteByIds 通过主键ID切片批量删除
	BatchDeleteByIds(ids []any) (int64, error)
}