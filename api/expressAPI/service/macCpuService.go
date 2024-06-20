package service

import "apiProject/api/expressAPI/types/domain"

type MacCpuService interface {
	// Save 保存新增
	Save(cpu *domain.MacCpu) (*domain.MacCpu, error)
	// BatchSave 批量保存新增
	BatchSave(cpuList []*domain.MacCpu) (int64, error)
	// SelectById 通过主键ID查询
	SelectById(id int64) (*domain.MacCpu, error)
	// Update 修改/更新
	Update(cpu *domain.MacCpu) (int64, error)
	// DeleteById 通过主键ID删除
	DeleteById(id int64) (int64, error)
}
