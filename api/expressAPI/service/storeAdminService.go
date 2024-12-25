package service

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
)

type StoreAdminServiceInterface interface {
	Save(te *domain.StoreAdmin) (*domain.StoreAdmin, error)
	// BatchSave 批量新增，返回成功的条数
	BatchSave(list []*domain.StoreAdmin) (int64, error)
	// PageList 分页查询
	//
	//	参数
	//		query: 查询参数对象
	//		page: 当前页码
	//		size: 每页条数
	//
	//	返回
	//		[]: 店铺管理数组
	//      int64: 总条数
	//      int64: 总页数
	PageList(query *param.StoreAdminSearchParam) ([]*domain.StoreAdmin, int64, int64, error)
	// BatchDelete 批量删除
	BatchDelete(ids []any) (rows int64, err error)
	SelectById(id int64) (*domain.StoreAdmin, error)
	SelectCountById(id int64) (int64, error)
	Update(storeAdmin *domain.StoreAdmin) (int64, error)
}
