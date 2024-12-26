package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"apiProject/api/utils"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"strings"
	"time"
)

type StoreAdminDb struct {
	Db *sql.DB
}

func NewStoreAdminDb(pg *sql.DB) *StoreAdminDb {
	return &StoreAdminDb{
		Db: pg,
	}
}

//goland:noinspection SqlResolve
func (pg *StoreAdminDb) Save(te *domain.StoreAdmin) (*domain.StoreAdmin, error) {

	// var err error
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("商户管理员新增新增开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("商户管理员新增新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
			zap.L().Sugar().Errorf("商户管理员新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("商户管理员新增事务正在提交")
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("商户管理员新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹商户管理员新增事务提交完成")
			}
		}
	}()

	// 新增数据返回的主键ID
	var lastInsertId int64

	if err = tx.QueryRow(`INSERT INTO tb_store_admin(user_name, mobile, real_name, status_value, store_name, merchant_id, merchant_name)
							VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		te.UserName,
		te.Mobile,
		te.RealName,
		te.StatusValue,
		te.StoreName,
		te.MerchantId,
		te.MerchantName,
	).Scan(&lastInsertId); err != nil {
		zap.L().Sugar().Errorf("商户管理员新增执行错误===%+v", err)
		return nil, err
	}

	if err != nil {
		zap.L().Sugar().Errorf("商户管理员新增错误===%+v", err)
		return nil, err
	}

	storeAdmin := &domain.StoreAdmin{}
	// 使用tx的查询，保证与插入操作在同一个事务中
	if err := tx.QueryRow(
		`SELECT 
    			id, 
    			user_name, 
    			mobile, 
    			real_name, 
    			status_value, 
    			store_name, 
    			merchant_id, 
    			merchant_name, 
       			CASE
					WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS create_time,
				CASE
					WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS update_time 
			FROM tb_store_admin WHERE id = $1`, lastInsertId).Scan(
		&storeAdmin.Id,
		&storeAdmin.UserName,
		&storeAdmin.Mobile,
		&storeAdmin.RealName,
		&storeAdmin.StatusValue,
		&storeAdmin.StoreName,
		&storeAdmin.MerchantId,
		&storeAdmin.MerchantName,
		&storeAdmin.CreateTime,
		&storeAdmin.UpdateTime,
	); err != nil {
		zap.L().Sugar().Errorf("通过ID查询商户管理员错误===%+v", err)
		return nil, err
	}

	return storeAdmin, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *StoreAdminDb) SelectByStoreAdminId(id int64) (*domain.StoreAdmin, error) {
	// 返回值这里不需要加上*号
	storeAdmin := &domain.StoreAdmin{}
	row := pg.Db.QueryRow(
		`SELECT
				id,
				user_name,
				mobile,
				real_name,
				status_value,
				store_name,
				merchant_id,
				merchant_name,
				CASE
					WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS create_time,
				CASE
					WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS update_time
			FROM
				tb_store_admin
			WHERE
				id = $1`,
		id)
	err := rowScanForStoreAdmin(row, storeAdmin)
	if err != nil {
		zap.L().Sugar().Errorf("商户管理查询错误===%+v", err)
		return nil, err
	}

	// 这里需要加上&
	return storeAdmin, nil
}

//goland:noinspection SqlResolve, SqlError
func (pg *StoreAdminDb) BatchSave(list []*domain.StoreAdmin) (int64, error) {
	// 占位符切片
	var placeholderList []string
	// 对应的值
	var valueArgList []interface{}
	// 占位符的数量
	numPlaceholders := 7

	for i := 0; i < len(list); i += numPlaceholders {
		end := i + numPlaceholders
		if end > len(list) {
			end = len(list)
		}
		for j := i; j < end; j++ {
			storeAdmin := list[j]
			placeholderList = append(placeholderList, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				len(valueArgList)+1, len(valueArgList)+2, len(valueArgList)+3, len(valueArgList)+4, len(valueArgList)+5,
				len(valueArgList)+6, len(valueArgList)+7))
			valueArgList = append(valueArgList,
				storeAdmin.UserName,
				storeAdmin.Mobile,
				storeAdmin.RealName,
				storeAdmin.StatusValue,
				storeAdmin.StoreName,
				storeAdmin.MerchantId,
				storeAdmin.MerchantName,
			)
		}
	}

	batchSql := fmt.Sprintf(`INSERT INTO tb_store_admin (user_name, mobile, real_name, status_value, store_name, merchant_id, merchant_name) VALUES %s %s`, strings.Join(placeholderList, ","), "\r\n")
	zap.L().Sugar().Infof("商户管理批量新增sql===%s", batchSql)
	result, err := pg.Db.Exec(batchSql, valueArgList...)
	if err != nil {
		zap.L().Sugar().Errorf("商户管理批量新增执行错误===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("商户管理批量新增执行失败===%+v", err)
		return 0, err
	}

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlCaseVsLimit
func (pg *StoreAdminDb) PageList(storeAdmin *param.StoreAdminSearchParam) ([]*domain.StoreAdmin, int64, int64, error) {
	// 查询总记录数
	var totalRecords int64
	countSql := "SELECT COUNT(*) FROM tb_store_admin" + buildCountForStoreAdmin(storeAdmin)
	zap.L().Sugar().Infof("商户分页查询count的sql===%s", countSql)

	err := pg.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		zap.L().Sugar().Errorf("商户分页查询总条数错误===%+v", err)
		return nil, 0, 0, err
	}

	zap.L().Sugar().Infof("商户分页查询总条数===%d", totalRecords)

	size, offset, totalPages := BuildPageOffset(storeAdmin.Page, storeAdmin.Size, totalRecords)

	rows, err := pg.Db.Query(fmt.Sprintf(
		`SELECT
					id,
					user_name,
					mobile,
					real_name,
					status_value,
					store_name,
					merchant_id,
					merchant_name,
					CASE
						WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS create_time,
					CASE
						WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS update_time
				FROM
					tb_store_admin %s
				ORDER BY id DESC	
				OFFSET $1 LIMIT $2`, buildCountForStoreAdmin(storeAdmin)), offset, size)

	if err != nil {
		zap.L().Sugar().Errorf("商户分页查询行数错误===%+v", err)
		return nil, 0, 0, err
	}

	defer RowsClose(rows, "商户分页查询")

	var list []*domain.StoreAdmin

	for rows.Next() {
		storeAdmin := &domain.StoreAdmin{}
		err := rowsScanForStoreAdminList(rows, storeAdmin)
		if err != nil {
			zap.L().Sugar().Errorf("商户分页返回错误===%+v", err)
			return nil, 0, 0, err
		}

		list = append(list, storeAdmin)
	}
	return list, totalRecords, totalPages, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *StoreAdminDb) SelectById(id int64) (*domain.StoreAdmin, error) {
	row := pg.Db.QueryRow(
		`SELECT
					id,
					user_name,
					mobile,
					real_name,
					status_value,
					store_name,
					merchant_id,
					merchant_name,
					CASE
						WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS create_time,
					CASE
						WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS update_time
				FROM
					tb_store_admin
			WHERE id = $1`, id)

	data := &domain.StoreAdmin{}
	err := rowScanForStoreAdmin(row, data)
	if err != nil {
		zap.L().Sugar().Errorf("通过ID查询商户错误===%+v", err)
		return nil, err
	}
	return data, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *StoreAdminDb) SelectCountById(id int64) (int64, error) {
	row := pg.Db.QueryRow(`SELECT COUNT(*) FROM tb_store_admin WHERE id = $1`, id)
	count := int64(0)
	// 将查询结果扫描到 count 变量
	err := row.Scan(&count)
	if err != nil {
		zap.L().Sugar().Errorf("通过ID查询商户错误===%+v", err)
		return 0, err
	}
	return count, nil
}

//goland:noinspection SqlResolve
func (pg *StoreAdminDb) BatchDelete(ids []any) (rows int64, err error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("商户管理批量删除开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("商户管理批量删除事务即将回滚（panic恢复）")
			tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("商户管理批量删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("商户管理批量删除提交事务失败===%+v", err)
			}
		}
	}()

	result, err := tx.Exec(fmt.Sprintf("DELETE FROM tb_store_admin WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids))), ids...)
	if err != nil {
		zap.L().Sugar().Errorf("商户管理批量删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("商户管理批量删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("商户管理批量删除执行获取条数===%d", rowsAffected)
	return rowsAffected, nil
}

func (pg *StoreAdminDb) Update(storeAdmin *domain.StoreAdmin) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("商户修改开启事务失败===%+v", err)
		return 0, err
	}
	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("商户修改事务即将回滚（panic恢复）")
			tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("商户修改事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("商户修改提交事务失败===%+v", err)
			}
		}
	}()

	rows, err := tx.Exec("UPDATE tb_store_admin SET user_name = $1, mobile = $2, real_name = $3, status_value = $4, store_name = $5, merchant_id = $6, merchant_name = $7, update_time= $8  WHERE id = $9",
		storeAdmin.UserName, storeAdmin.Mobile, storeAdmin.RealName, storeAdmin.StatusValue, storeAdmin.StoreName, storeAdmin.MerchantId, storeAdmin.MerchantName, time.Now(), storeAdmin.Id)
	if err != nil {
		zap.L().Sugar().Errorf("商户修改执行错误===%+v", err)
		return 0, err
	}
	count, err := rows.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("商户修改获取执行条数错误===%+v", err)
		return 0, err
	}
	return count, nil
}

// 构建 WHERE 子句
func buildCountForStoreAdmin(sa *param.StoreAdminSearchParam) string {
	if sa != nil {
		var clauses []string
		userName := sa.UserName
		mobile := sa.Mobile
		realName := sa.RealName
		statusValue := sa.StatusValue
		storeName := sa.StoreName
		merchantId := sa.MerchantId
		merchantName := sa.MerchantName
		if userName != "" {
			clauses = append(clauses, "user_name LIKE CONCAT('%', '"+userName+"', '%')")
		}
		if mobile != "" {
			clauses = append(clauses, "mobile LIKE CONCAT('%', '"+mobile+"', '%')")
		}
		if realName != "" {
			clauses = append(clauses, "real_name LIKE CONCAT('%', '"+realName+"', '%')")
		}
		if statusValue != "" {
			clauses = append(clauses, "status_value = '"+statusValue+"'")
		}

		if storeName != "" {
			clauses = append(clauses, "store_name = '"+storeName+"'")
		}

		if merchantId != 0 {
			clauses = append(clauses, "merchant_id = '"+utils.ConvertInt64ToStr(merchantId)+"'")
		}
		if merchantName != "" {
			clauses = append(clauses, "merchant_name LIKE CONCAT('%', '"+merchantName+"', '%')")
		}

		if len(clauses) > 0 {
			return " WHERE " + strings.Join(clauses, " AND ")
		}
	}

	return ""
}

// rowScan 将数据行row转换到结构体
func rowScanForStoreAdmin(row *sql.Row, storeAdmin *domain.StoreAdmin) error {
	err := row.Scan(&storeAdmin.Id,
		&storeAdmin.UserName,
		&storeAdmin.Mobile,
		&storeAdmin.RealName,
		&storeAdmin.StatusValue,
		&storeAdmin.StoreName,
		&storeAdmin.MerchantId,
		&storeAdmin.MerchantName,
		&storeAdmin.CreateTime,
		&storeAdmin.UpdateTime)
	if err != nil {
		zap.L().Sugar().Errorf("商户数据转换错误===%+v", err)
		return err
	}
	return nil
}

// rowsScan 将数据行rows转换到结构体
func rowsScanForStoreAdminList(rows *sql.Rows, storeAdmin *domain.StoreAdmin) error {
	err := rows.Scan(&storeAdmin.Id,
		&storeAdmin.UserName,
		&storeAdmin.Mobile,
		&storeAdmin.RealName,
		&storeAdmin.StatusValue,
		&storeAdmin.StoreName,
		&storeAdmin.MerchantId,
		&storeAdmin.MerchantName,
		&storeAdmin.CreateTime,
		&storeAdmin.UpdateTime)
	if err != nil {
		zap.L().Sugar().Errorf("商户数据多条转换错误===%+v", err)
		return err
	}
	return nil
}
