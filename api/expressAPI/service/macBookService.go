package service

import "apiProject/api/expressAPI/types/domain"

type MacBookService interface {
	// Save 保存新增
	Save(macBook *domain.MacBook) (*domain.MacBook, error)
	// BatchSave 批量保存新增
	BatchSave(macBookList []*domain.MacBook) (int64, error)
	// SelectById 通过主键ID查询
	SelectById(id int64) (*domain.MacBook, error)
	// Update 修改/更新
	Update(macBook *domain.MacBook) (*domain.MacBook, error)
	// DeleteById 通过主键ID删除
	DeleteById(id int64) (int64, error)
	// BatchDeleteByIds 通过主键ID切片批量删除
	BatchDeleteByIds(ids []any) (int64, error)
}