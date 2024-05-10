package service

import "apiProject/api/expressAPI/types/domain"

type TestExpressService interface {
	Save(te *domain.TestExpress) (*domain.TestExpress, error)
	BatchSave([]*domain.TestExpress) (int64, error)
	PageList(te *domain.TestExpress, page, size int64) ([]*domain.TestExpress, int64, int64, error)
}
