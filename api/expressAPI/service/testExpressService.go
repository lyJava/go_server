package service

import "apiProject/api/expressAPI/types/domain"

type TestExpressService interface {
	Save(te *domain.TestExpress) (*domain.TestExpress, error)
	// BatchSave 批量新增，返回成功的条数
	BatchSave(list []*domain.TestExpress) (int64, error)
	PageList(te *domain.TestExpress, page, size int64) ([]*domain.TestExpress, int64, int64, error)
	// BatchDelete 批量删除
	BatchDelete(ids []string) (rows int64, err error)
}
