package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/utils"
	"database/sql"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
)

type MacCpuDb struct {
	Db *sql.DB
}

func NewMacCpuDb(pg *sql.DB) *MacCpuDb {
	return &MacCpuDb{
		Db: pg,
	}
}

const (
	CommonColumn = "cpu_name, cpu_type, cpu_basic_boost, cpu_trubo_boost, cpu_core_number, cpu_thread_number, cpu_cache, cpu_tdp, memory_width, media_processing_engine"
)

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) Save(cpu *domain.MacCpu) (*domain.MacCpu, error) {
	// var err error
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器新增开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果处理器新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果处理器新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果处理器新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果处理器新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果处理器新增事务提交完成")
			}
		}
	}()

	var lastInsertId int64

	// 加上RETURNING id，然后使用Scan可以返回新增数据的主键ID
	if err = tx.QueryRow(fmt.Sprintf(`INSERT INTO tb_mac_cpu(%s) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, CommonColumn),
		cpu.CpuName,
		cpu.CpuType,
		cpu.CpuBasicBoost,
		cpu.CpuTruboBoost,
		cpu.CpuCoreNumber,
		cpu.CpuThreadNumber,
		cpu.CpuCache,
		cpu.CpuTdp,
		cpu.MemoryWidth,
		cpu.MediaProcessingEngine,
	).Scan(&lastInsertId); err != nil {
		zap.L().Sugar().Errorf("苹果处理器新增执行错误===%+v", err)
		// 发生错误则回滚
		//_ = tx.Rollback()
		return nil, err
	}

	zap.L().Sugar().Infof("苹果处理器新增返回自增主键ID: %d", lastInsertId)

	// 显式提交事务
	/*if err = tx.Commit(); err != nil {
		log.Printf("苹果处理器新增提交事务失败===%+v", err)
		zap.L().Sugar().Errorf("苹果处理器新增提交事务失败===%+v", err)
		return nil, err
	}*/

	//return macCpu, nil
	return selectDetail(err, tx, lastInsertId)
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) BatchSave(list []*domain.MacCpu) (int64, error) {
	// 占位符切片
	var placeholderList []string
	// 对应的值
	var valueArgList []interface{}
	// 占位符的数量
	numPlaceholders := 10

	for i := 0; i < len(list); i += numPlaceholders {
		end := i + numPlaceholders
		if end > len(list) {
			end = len(list)
		}
		for j := i; j < end; j++ {
			macCpu := list[j]
			placeholderList = append(placeholderList, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				len(valueArgList)+1, len(valueArgList)+2, len(valueArgList)+3, len(valueArgList)+4, len(valueArgList)+5,
				len(valueArgList)+6, len(valueArgList)+7, len(valueArgList)+8, len(valueArgList)+9, len(valueArgList)+10))
			valueArgList = append(valueArgList,
				macCpu.CpuName,
				macCpu.CpuType,
				macCpu.CpuBasicBoost,
				macCpu.CpuTruboBoost,
				macCpu.CpuCoreNumber,
				macCpu.CpuThreadNumber,
				macCpu.CpuCache,
				macCpu.CpuTdp,
				macCpu.MemoryWidth,
				macCpu.MediaProcessingEngine,
			)
		}
	}

	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量新增开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Debugf("苹果处理器批量新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果处理器批量新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Debugf("苹果处理器批量新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果处理器批量新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Debugf("苹果处理器批量新增事务提交完成")
			}
		}
	}()

	batchSql := fmt.Sprintf(`INSERT INTO tb_mac_cpu (%s) VALUES %s`, CommonColumn, strings.Join(placeholderList, ","))
	zap.L().Sugar().Infof("苹果处理器批量新增sql===%s", batchSql)
	result, err := tx.Exec(batchSql, valueArgList...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量新增执行错误===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量新增执行失败===%+v", err)
		return 0, err
	}

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) SelectById(id int64) (*domain.MacCpu, error) {
	cpu := &domain.MacCpu{}
	querySql := fmt.Sprintf(`SELECT id, %s FROM tb_mac_cpu WHERE id = $1`, CommonColumn)
	if err := pg.Db.QueryRow(querySql, id).
		Scan(
			&cpu.Id,
			&cpu.CpuName,
			&cpu.CpuType,
			&cpu.CpuBasicBoost,
			&cpu.CpuTruboBoost,
			&cpu.CpuCoreNumber,
			&cpu.CpuThreadNumber,
			&cpu.CpuCache,
			&cpu.CpuTdp,
			&cpu.MemoryWidth,
			&cpu.MediaProcessingEngine,
		); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
			return nil, fmt.Errorf("no MacCpu found with id %d", id)
		}
		zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
		return nil, err
	}
	return cpu, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) Update(cpu *domain.MacCpu) (*domain.MacCpu, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器修改开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果处理器修改事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果处理器修改事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果处理器修改事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果处理器修改提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果处理器修改事务提交完成")
			}
		}
	}()

	setClause, args, placeholderIndex, err := cpuDynamicUpdate(cpu)
	zap.L().Sugar().Infof("苹果处理器修改动态参数:\n%v,\nargsCount:%d", utils.ToJsonFormat(args), placeholderIndex)

	if err != nil {
		return nil, err
	}

	updateSql := fmt.Sprintf("UPDATE tb_mac_cpu SET %s WHERE id = $%d RETURNING id", setClause, placeholderIndex)
	zap.L().Sugar().Infof("苹果处理器修改动态sql:%s", updateSql)

	var cpuId int64

	args = append(args, cpu.Id)
	if err = tx.QueryRow(updateSql, args...).Scan(&cpuId); err != nil {
		zap.L().Sugar().Errorf("部门负责人修改执行错误===%+v", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("未查询到数据，请确认参数有效性")
		}
		return nil, err
	}

	zap.L().Sugar().Infof("部门负责人修改返回ID:%d", cpuId)

	return selectDetail(err, tx, cpuId)
}

func selectDetail(err error, tx *sql.Tx, cpuId int64) (*domain.MacCpu, error) {
	macCpu := &domain.MacCpu{}
	// 使用tx的查询，保证与插入操作在同一个事务中
	if err = tx.QueryRow(fmt.Sprintf(`SELECT id, %s FROM tb_mac_cpu WHERE id = $1`, CommonColumn), cpuId).
		Scan(
			&macCpu.Id,
			&macCpu.CpuName,
			&macCpu.CpuType,
			&macCpu.CpuBasicBoost,
			&macCpu.CpuTruboBoost,
			&macCpu.CpuCoreNumber,
			&macCpu.CpuThreadNumber,
			&macCpu.CpuCache,
			&macCpu.CpuTdp,
			&macCpu.MemoryWidth,
			&macCpu.MediaProcessingEngine,
		); err != nil {
		zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
		return nil, err
	}

	return macCpu, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) DeleteById(id int64) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器删除开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果处理器删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果处理器删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果处理器删除提交事务失败===%+v", err)
			}
		}
	}()

	result, err := tx.Exec(`DELETE FROM tb_mac_cpu WHERE id = $1`, id)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果处理器删除执行获取条数===%d", rowsAffected)

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) BatchDeleteByIds(ids []any) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量删除开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果处理器批量删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果处理器批量删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果处理器批量删除提交事务失败===%+v", err)
			}
		}
	}()

	deleteSql := fmt.Sprintf("DELETE FROM tb_mac_cpu WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids)))
	zap.L().Sugar().Infof("苹果处理器批量删除执行sql===%s", deleteSql)

	result, err := tx.Exec(deleteSql, ids...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量删除执行获取条数错误===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果处理器批量删除执行成功条数===%d", rowsAffected)

	return rowsAffected, nil
}

// cpuDynamicUpdate 构建动态更新
func cpuDynamicUpdate(cpu *domain.MacCpu) (string, []any, int64, error) {
	var setClauses []string
	var args []interface{}
	placeholderIndex := int64(1)

	addClause := func(field string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, placeholderIndex))
		args = append(args, value)
		placeholderIndex++
	}

	if cpu.CpuName != "" {
		addClause(utils.CamelToSnakeCase("CpuName"), cpu.CpuName)
	}
	if cpu.CpuType != "" {
		addClause(utils.CamelToSnakeCase("CpuType"), cpu.CpuType)
	}
	if cpu.CpuBasicBoost != "" {
		addClause(utils.CamelToSnakeCase("CpuBasicBoost"), cpu.CpuBasicBoost)
	}
	if cpu.CpuTruboBoost != "" {
		addClause(utils.CamelToSnakeCase("CpuTruboBoost"), cpu.CpuTruboBoost)
	}
	if cpu.CpuCoreNumber != 0 {
		addClause(utils.CamelToSnakeCase("CpuCoreNumber"), cpu.CpuCoreNumber)
	}
	if cpu.CpuThreadNumber != 0 {
		addClause(utils.CamelToSnakeCase("CpuThreadNumber"), cpu.CpuThreadNumber)
	}
	if cpu.CpuCache != "" {
		addClause(utils.CamelToSnakeCase("CpuCache"), cpu.CpuCache)
	}
	if cpu.CpuTdp != "" {
		addClause(utils.CamelToSnakeCase("CpuTdp"), cpu.CpuTdp)
	}
	if cpu.MemoryWidth != "" {
		addClause(utils.CamelToSnakeCase("MemoryWidth"), cpu.MemoryWidth)
	}
	if cpu.MediaProcessingEngine != "" {
		addClause(utils.CamelToSnakeCase("MediaProcessingEngine"), cpu.MediaProcessingEngine)
	}

	if len(setClauses) == 0 {
		zap.L().Sugar().Info("no fields to update")
		return "", nil, 0, errors.New("没有需要更新的列信息")
	}

	setClause := strings.Join(setClauses, ", ")
	return setClause, args, placeholderIndex, nil
}
