package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"apiProject/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type MacBookDb struct {
	Db *sql.DB
}

func NewMacBookDb(pg *sql.DB) *MacBookDb {
	return &MacBookDb{
		Db: pg,
	}
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) Save(macBook *domain.MacBook) (*domain.MacBook, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本新增开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果笔记本新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果笔记本新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果笔记本新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果笔记本新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果笔记本新增事务提交完成")
			}
		}
	}()

	zap.L().Sugar().Info("苹果笔记本新增参数===%+v", macBook)

	var lastInsertId int64

	// 加上RETURNING id，然后使用Scan可以返回新增数据的主键ID
	if err = tx.QueryRow(`INSERT INTO tb_mac_book(
							pro_name, pro_color, pro_year, pro_type, mold_model, monitor_id, cpu_id, memory_id, storage_info, size_info, weight_info, gpu_id, chip_type, memory_max, create_time, update_time) 
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`,
		macBook.ProName,
		macBook.ProColor,
		macBook.ProYear,
		macBook.ProType,
		macBook.MoldModel,
		macBook.MonitorId,
		macBook.CpuId,
		macBook.MemoryId,
		macBook.StorageInfo,
		macBook.SizeInfo,
		macBook.WeightInfo,
		macBook.GpuId,
		macBook.ChipType,
		macBook.MemoryMax,
	).Scan(&lastInsertId); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本新增执行错误===%+v", err)
		return nil, err
	}

	zap.L().Sugar().Infof("苹果笔记本新增返回自增主键ID: %d", lastInsertId)
	return selectBookDetail(tx, lastInsertId)
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) BatchSave(list []*domain.MacBook) (int64, error) {
	// 占位符切片
	var placeholderList []string
	// 对应的值
	var valueArgList []interface{}
	// 占位符的数量
	numPlaceholders := 14

	for i := 0; i < len(list); i += numPlaceholders {
		end := i + numPlaceholders
		if end > len(list) {
			end = len(list)
		}
		for j := i; j < end; j++ {
			macBook := list[j]
			placeholderList = append(placeholderList, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
				len(valueArgList)+1, len(valueArgList)+2, len(valueArgList)+3, len(valueArgList)+4, len(valueArgList)+5,
				len(valueArgList)+6, len(valueArgList)+7, len(valueArgList)+8, len(valueArgList)+9, len(valueArgList)+10,
				len(valueArgList)+11, len(valueArgList)+12, len(valueArgList)+13, len(valueArgList)+14))
			valueArgList = append(valueArgList,
				macBook.ProName,
				macBook.ProColor,
				macBook.ProYear,
				macBook.ProType,
				macBook.MoldModel,
				macBook.MonitorId,
				macBook.CpuId,
				macBook.MemoryId,
				macBook.StorageInfo,
				macBook.SizeInfo,
				macBook.WeightInfo,
				macBook.GpuId,
				macBook.ChipType,
				macBook.MemoryMax,
			)
		}
	}

	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量新增开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Infof("苹果笔记本批量新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Infof("苹果笔记本批量新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Infof("苹果笔记本批量新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果笔记本批量新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Infof("苹果笔记本批量新增事务提交完成")
			}
		}
	}()

	batchSql := fmt.Sprintf(`INSERT INTO tb_mac_book (pro_name, pro_color, pro_year, pro_type, mold_model, monitor_id, cpu_id, memory_id, storage_info, size_info, weight_info, gpu_id, chip_type, memory_max, create_time, update_time) VALUES %s`, strings.Join(placeholderList, ","))
	zap.L().Sugar().Infof("苹果笔记本批量新增sql===%s", batchSql)
	result, err := tx.Exec(batchSql, valueArgList...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量新增执行错误===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量新增执行失败===%+v", err)
		return 0, err
	}

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) SelectById(id int64) (*domain.MacBook, error) {
	macBook := &domain.MacBook{}
	querySql := fmt.Sprintf(`SELECT 
								b.id, %s 
							FROM tb_mac_book b 
							LEFT JOIN tb_mac_cpu c ON c.id = b.cpu_id 
							LEFT JOIN tb_mac_memory m ON m.id = b.memory_id 
							WHERE b.id = $1`, MacBookCommonColumn)
	if err := pg.Db.QueryRow(querySql, id).
		Scan(
			&macBook.Id, &macBook.ProName, &macBook.ProColor, &macBook.ProYear, &macBook.ProType, &macBook.MoldModel, &macBook.MonitorId, &macBook.CpuId, &macBook.MemoryId, &macBook.StorageInfo, &macBook.SizeInfo, &macBook.WeightInfo, &macBook.GpuId, &macBook.ChipType, &macBook.MemoryMax, &macBook.CreateTime, &macBook.UpdateTime,
			&macBook.CpuInfo.Id, &macBook.CpuInfo.CpuName, &macBook.CpuInfo.CpuType, &macBook.CpuInfo.CpuBasicBoost, &macBook.CpuInfo.CpuTruboBoost, &macBook.CpuInfo.CpuCoreNumber, &macBook.CpuInfo.CpuThreadNumber, &macBook.CpuInfo.CpuCache, &macBook.CpuInfo.CpuTdp, &macBook.CpuInfo.MemoryWidth, &macBook.CpuInfo.MediaProcessingEngine,
			&macBook.MemoryInfo.Id, &macBook.MemoryInfo.MemorySize, &macBook.MemoryInfo.MemorySpeed, &macBook.MemoryInfo.IntegrationFlag, &macBook.MemoryInfo.MemoryType, &macBook.MemoryInfo.EccCheck,
		); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Sugar().Errorf("通过ID查询苹果笔记本错误===%+v", err)
			return nil, fmt.Errorf("no MacCpu found with id %d", id)
		}
		zap.L().Sugar().Errorf("通过ID查询苹果笔记本错误===%+v", err)
		return nil, err
	}
	return macBook, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) Update(macBook *domain.MacBook) (*domain.MacBook, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本修改开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果笔记本修改事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果笔记本修改事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果笔记本修改事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果笔记本修改提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果笔记本修改事务提交完成")
			}
		}
	}()

	setClause, args, placeholderIndex, err := bookDynamicUpdate(macBook)
	zap.L().Sugar().Infof("苹果笔记本修改动态参数:\n%v,\nargsCount:%d", utils.ToJsonFormat(args), placeholderIndex)

	if err != nil {
		return nil, err
	}

	updateSql := fmt.Sprintf("UPDATE tb_mac_book SET %s, update_time = CURRENT_TIMESTAMP WHERE id = $%d RETURNING id", setClause, placeholderIndex)
	zap.L().Sugar().Infof("苹果笔记本修改动态sql:%s", updateSql)

	var cpuId int64

	args = append(args, macBook.Id)
	if err = tx.QueryRow(updateSql, args...).Scan(&cpuId); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本修改执行错误===%+v", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("未查询到数据，请确认参数有效性")
		}
		return nil, err
	}

	zap.L().Sugar().Infof("苹果笔记本修改返回ID:%d", cpuId)

	return selectBookDetail(tx, cpuId)
}

func selectBookDetail(tx *sql.Tx, bookId int64) (*domain.MacBook, error) {
	zap.L().Sugar().Infof("通过ID查询苹果笔记本参数===%d", bookId)
	macBook := &domain.MacBook{}

	// 使用tx的查询，保证与插入操作在同一个事务中
	if err := tx.QueryRow(fmt.Sprintf(`SELECT 
										b.id, %s 
									FROM tb_mac_book b 
									LEFT JOIN tb_mac_cpu c ON c.id = b.cpu_id 
									LEFT JOIN tb_mac_memory m ON m.id = b.memory_id 
									WHERE b.id = $1`, MacBookCommonColumn), bookId).
		Scan(
			&macBook.Id, &macBook.ProName, &macBook.ProColor, &macBook.ProYear, &macBook.ProType, &macBook.MoldModel, &macBook.MonitorId, &macBook.CpuId, &macBook.MemoryId, &macBook.StorageInfo, &macBook.SizeInfo, &macBook.WeightInfo, &macBook.GpuId, &macBook.ChipType, &macBook.MemoryMax, &macBook.CreateTime, &macBook.UpdateTime,
			&macBook.CpuInfo.Id, &macBook.CpuInfo.CpuName, &macBook.CpuInfo.CpuType, &macBook.CpuInfo.CpuBasicBoost, &macBook.CpuInfo.CpuTruboBoost, &macBook.CpuInfo.CpuCoreNumber, &macBook.CpuInfo.CpuThreadNumber, &macBook.CpuInfo.CpuCache, &macBook.CpuInfo.CpuTdp, &macBook.CpuInfo.MemoryWidth, &macBook.CpuInfo.MediaProcessingEngine,
			&macBook.MemoryInfo.Id, &macBook.MemoryInfo.MemorySize, &macBook.MemoryInfo.MemorySpeed, &macBook.MemoryInfo.IntegrationFlag, &macBook.MemoryInfo.MemoryType, &macBook.MemoryInfo.EccCheck,
		); err != nil {
		zap.L().Sugar().Errorf("通过ID查询苹果笔记本错误===%+v", err)
		return nil, err
	}

	return macBook, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) DeleteById(id int64) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本删除开启事务失败===%+v", err)
		return 0, err
	}
	zap.L().Sugar().Info("苹果笔记本删除事务开启成功")
	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果笔记本删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果笔记本删除事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果笔记本删除事务正在提交")
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果笔记本删除提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果笔记本删除事务提交成功")
			}
		}
	}()

	result, err := tx.Exec(`DELETE FROM tb_mac_book WHERE id = $1`, id)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果笔记本删除执行获取条数===%d", rowsAffected)

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacBookDb) BatchDeleteByIds(ids []any) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量删除开启事务失败===%+v", err)
		return 0, err
	}
	zap.L().Sugar().Info("苹果笔记本批量删除事务开启成功")

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果笔记本批量删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果笔记本批量删除事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果笔记本批量删除事务正在提交")
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果笔记本批量删除提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果笔记本批量删除事务提交成功")
			}
		}
	}()

	deleteSql := fmt.Sprintf("DELETE FROM tb_mac_book WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids)))
	zap.L().Sugar().Infof("苹果笔记本批量删除执行sql===%s", deleteSql)

	result, err := tx.Exec(deleteSql, ids...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量删除执行获取条数错误===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果笔记本批量删除执行成功条数===%d", rowsAffected)

	return rowsAffected, nil
}

func (pg *MacBookDb) SelectPage(pageParam *param.MacBookPageParam) ([]*domain.MacBook, int64, int64, error) {
	whereParam := buildWhereClauseForMacBookPage(pageParam)
	// 查询总记录数
	var totalRecords int64
	countSql := `SELECT COUNT(*) FROM tb_mac_book b` + whereParam
	zap.L().Sugar().Infof("测试快递分页查询count的sql===%s", countSql)

	err := pg.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本分页查询总条数错误===%+v", err)
		return nil, 0, 0, err
	}

	zap.L().Sugar().Infof("苹果笔记本分页查询总条数===%d", totalRecords)

	size, offset, totalPages := BuildPageOffset(pageParam.Page, pageParam.Size, totalRecords)

	zap.L().Sugar().Infof("苹果笔记本分页查询offset=%d, size=%d", offset, size)

	pageSql := fmt.Sprintf(`SELECT 
								b.id, %s 
							FROM tb_mac_book b 
							LEFT JOIN tb_mac_cpu c ON c.id = b.cpu_id 
							LEFT JOIN tb_mac_memory m ON m.id = b.memory_id 
							%s 
							%s 
							OFFSET $1 LIMIT $2`, MacBookCommonColumn, whereParam, buildOrderBy(pageParam.Column, pageParam.Order, "b"))
	zap.L().Sugar().Infof("苹果笔记本分页查询sql===%s", pageSql)
	rows, err := pg.Db.Query(pageSql, offset, size)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本分页查询行数错误===%+v", err)
		return nil, 0, 0, err
	}

	defer RowsClose(rows, "苹果笔记本分页查询")

	var list []*domain.MacBook

	for rows.Next() {
		macBook := &domain.MacBook{}
		if err = rows.Scan(&macBook.Id, &macBook.ProName, &macBook.ProColor, &macBook.ProYear, &macBook.ProType, &macBook.MoldModel, &macBook.MonitorId, &macBook.CpuId, &macBook.MemoryId, &macBook.StorageInfo, &macBook.SizeInfo, &macBook.WeightInfo, &macBook.GpuId, &macBook.ChipType, &macBook.MemoryMax, &macBook.CreateTime, &macBook.UpdateTime,
			&macBook.CpuInfo.Id, &macBook.CpuInfo.CpuName, &macBook.CpuInfo.CpuType, &macBook.CpuInfo.CpuBasicBoost, &macBook.CpuInfo.CpuTruboBoost, &macBook.CpuInfo.CpuCoreNumber, &macBook.CpuInfo.CpuThreadNumber, &macBook.CpuInfo.CpuCache, &macBook.CpuInfo.CpuTdp, &macBook.CpuInfo.MemoryWidth, &macBook.CpuInfo.MediaProcessingEngine,
			&macBook.MemoryInfo.Id, &macBook.MemoryInfo.MemorySize, &macBook.MemoryInfo.MemorySpeed, &macBook.MemoryInfo.IntegrationFlag, &macBook.MemoryInfo.MemoryType, &macBook.MemoryInfo.EccCheck); err != nil {
			zap.L().Sugar().Errorf("苹果笔记本分页查询错误===%+v", err)
			return nil, 0, 0, err
		}
		list = append(list, macBook)
	}

	return list, totalRecords, totalPages, nil
}

// cpuDynamicUpdate 构建动态更新
func bookDynamicUpdate(macBook *domain.MacBook) (string, []any, int64, error) {
	var setClauses []string
	var args []interface{}
	placeholderIndex := int64(1)

	addClause := func(field string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, placeholderIndex))
		args = append(args, value)
		placeholderIndex++
	}

	if macBook.ProName != "" {
		addClause(utils.CamelToSnakeCase("ProName"), macBook.ProName)
	}
	if macBook.ProColor != "" {
		addClause(utils.CamelToSnakeCase("ProColor"), macBook.ProColor)
	}
	if macBook.ProYear != "" {
		addClause(utils.CamelToSnakeCase("ProYear"), macBook.ProYear)
	}
	if macBook.MoldModel != "" {
		addClause(utils.CamelToSnakeCase("MoldModel"), macBook.MoldModel)
	}

	if macBook.MonitorId != 0 {
		addClause(utils.CamelToSnakeCase("MonitorId"), macBook.MonitorId)
	}
	if macBook.CpuId != 0 {
		addClause(utils.CamelToSnakeCase("CpuId"), macBook.CpuId)
	}
	if macBook.MemoryId != 0 {
		addClause(utils.CamelToSnakeCase("MemoryId"), macBook.MemoryId)
	}

	if macBook.StorageInfo != "" {
		addClause(utils.CamelToSnakeCase("StorageInfo"), macBook.StorageInfo)
	}
	if macBook.SizeInfo != "" {
		addClause(utils.CamelToSnakeCase("SizeInfo"), macBook.SizeInfo)
	}
	if macBook.WeightInfo != "" {
		addClause(utils.CamelToSnakeCase("WeightInfo"), macBook.WeightInfo)
	}

	if macBook.GpuId != 0 {
		addClause(utils.CamelToSnakeCase("GpuId"), macBook.GpuId)
	}

	if macBook.ChipType != "" {
		addClause(utils.CamelToSnakeCase("ChipType"), macBook.ChipType)
	}

	if macBook.MemoryMax != 0 {
		addClause(utils.CamelToSnakeCase("MemoryMax"), macBook.MemoryMax)
	}

	if len(setClauses) == 0 {
		zap.L().Sugar().Info("no fields to update")
		return "", nil, 0, errors.New("没有需要更新的列信息")
	}

	setClause := strings.Join(setClauses, ", ")
	return setClause, args, placeholderIndex, nil
}

// buildWhereClauseForMacBookPage 构建多条件动态查询
func buildWhereClauseForMacBookPage(param *param.MacBookPageParam) string {
	if param != nil {
		var clauses []string
		if param.ProName != "" {
			clauses = append(clauses, "b.pro_name LIKE CONCAT('%', '"+param.ProName+"', '%')")
		}
		if param.ProColor != "" {
			clauses = append(clauses, "b.pro_color = '"+param.ProColor+"'")
		}
		if param.ProYear != "" {
			clauses = append(clauses, "b.pro_year LIKE CONCAT('%', '"+param.ProYear+"', '%')")
		}
		if param.ProType != "" {
			clauses = append(clauses, "b.pro_type  = '"+param.ProType+"'")
		}
		if param.MoldModel != "" {
			clauses = append(clauses, "b.mold_model LIKE CONCAT('%', '"+param.MoldModel+"', '%')")
		}

		if len(clauses) > 0 {
			return " WHERE " + strings.Join(clauses, " AND ")
		}
	}
	return ""
}
