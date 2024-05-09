package service

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
)

// ExpressServiceInterface 快递服务接口
type ExpressServiceInterface interface {
	// CreateExpress 新增快递信息
	//
	// 参数
	//	express (Express): 快递对象
	CreateExpress(express *domain.Express) (*domain.Express, error)
	// GetExpress 查询快递信息
	//
	// 参数
	//	id (int64): 主键ID
	GetExpress(id int64) (*domain.Express, error)
	// SelectExpressPage 分页查询
	//
	// 参数
	//	expressName (string): 快递名称
	//	page (string): 当前页码
	//	size (string): 每页条数
	SelectExpressPage(expressName string, page, size string) ([]*domain.Express, int64, int64, error)
	// DeleteById 删除
	//
	// 参数
	//	id (int64): 主键ID
	DeleteById(id int64) (rows int64, err error)
	// SelectExpressPageByParam 分页查询
	//
	// 参数
	//	searchParam (ExpressSearchParam): 快递查询参数对象
	SelectExpressPageByParam(searchParam *param.ExpressSearchParam) ([]*domain.Express, int64, int64, error)
	// UpdateExpress 修改快递信息
	//
	// 参数
	//	express (Express): 快递对象
	UpdateExpress(express *domain.Express) (int64, error)
	// BatchDeleteByIds 批量删除
	//
	// 参数
	//	ids ([]string): 主键ID集合
	BatchDeleteByIds(ids []string) (rows int64, err error)
	// BatchCreateExpress 批量新增
	//
	// 参数
	//	list ([]Express): 快递切片
	BatchCreateExpress(list []*domain.Express) (int64, error)
}
