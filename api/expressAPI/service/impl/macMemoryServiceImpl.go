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

type MacMemoryDb struct {
	Db *sql.DB
}

func NewMacMemoryDb(pg *sql.DB) *MacMemoryDb {
	return &MacMemoryDb{
		Db: pg,
	}
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) Save(memory *domain.MacMemory) (*domain.MacMemory, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存新增开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果内存新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果内存新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果内存新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果内存新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果内存新增事务提交完成")
			}
		}
	}()

	var lastInsertId int64

	// 加上RETURNING id，然后使用Scan可以返回新增数据的主键ID
	if err = tx.QueryRow(fmt.Sprintf(`INSERT INTO tb_mac_memory(%s) VALUES ($1, $2, $3, $4, $5) RETURNING id`, MacCpuCommonColumn),
		memory.MemorySize,
		memory.MemorySpeed,
		memory.IntegrationFlag,
		memory.MemoryType,
		memory.EccCheck,
	).Scan(&lastInsertId); err != nil {
		zap.L().Sugar().Errorf("苹果内存新增执行错误===%+v", err)
		return nil, err
	}

	zap.L().Sugar().Infof("苹果内存新增返回自增主键ID: %d", lastInsertId)

	return selectMemoryDetail(tx, lastInsertId)
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) BatchSave(list []*domain.MacMemory) (int64, error) {
	// 占位符切片
	var placeholderList []string
	// 对应的值
	var valueArgList []interface{}
	// 占位符的数量
	numPlaceholders := 5

	for i := 0; i < len(list); i += numPlaceholders {
		end := i + numPlaceholders
		if end > len(list) {
			end = len(list)
		}
		for j := i; j < end; j++ {
			memory := list[j]
			placeholderList = append(placeholderList, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
				len(valueArgList)+1, len(valueArgList)+2, len(valueArgList)+3, len(valueArgList)+4, len(valueArgList)+5))
			valueArgList = append(valueArgList,
				memory.MemorySize,
				memory.MemorySpeed,
				memory.IntegrationFlag,
				memory.MemoryType,
				memory.EccCheck,
			)
		}
	}

	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量新增开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Infof("苹果内存批量新增事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Infof("苹果内存批量新增事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Infof("苹果内存批量新增事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果内存批量新增提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Infof("苹果内存批量新增事务提交完成")
			}
		}
	}()

	batchSql := fmt.Sprintf(`INSERT INTO tb_mac_memory (%s) VALUES %s`, MacMemoryCommonColumn, strings.Join(placeholderList, ","))
	zap.L().Sugar().Infof("苹果内存批量新增sql===%s", batchSql)
	result, err := tx.Exec(batchSql, valueArgList...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量新增执行错误===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量新增执行失败===%+v", err)
		return 0, err
	}

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) SelectById(id int64) (*domain.MacMemory, error) {
	memory := &domain.MacMemory{}
	querySql := fmt.Sprintf(`SELECT id, %s FROM tb_mac_memory WHERE id = $1`, MacMemoryCommonColumn)
	if err := pg.Db.QueryRow(querySql, id).
		Scan(
			&memory.Id,
			&memory.MemorySize,
			&memory.MemorySpeed,
			&memory.IntegrationFlag,
			&memory.MemoryType,
			&memory.EccCheck,
		); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Sugar().Errorf("通过ID查询苹果内存错误===%+v", err)
			return nil, fmt.Errorf("no MacCpu found with id %d", id)
		}
		zap.L().Sugar().Errorf("通过ID查询苹果内存错误===%+v", err)
		return nil, err
	}
	return memory, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) Update(memory *domain.MacMemory) (*domain.MacMemory, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存修改开启事务失败===%+v", err)
		return nil, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果内存修改事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果内存修改事务回滚,发生错误===%+v", err)
		} else {
			zap.L().Sugar().Info("苹果内存修改事务正在提交")
			// 正常结束则提交事务
			if err = tx.Commit(); err != nil {
				zap.L().Sugar().Errorf("苹果内存修改提交事务失败===%+v", err)
			} else {
				zap.L().Sugar().Info("苹果内存修改事务提交完成")
			}
		}
	}()

	setClause, args, placeholderIndex, err := memoryDynamicUpdate(memory)
	zap.L().Sugar().Infof("苹果内存修改动态参数:\n%v,\nargsCount:%d", utils.ToJsonFormat(args), placeholderIndex)

	if err != nil {
		return nil, err
	}

	updateSql := fmt.Sprintf("UPDATE tb_mac_memory SET %s WHERE id = $%d RETURNING id", setClause, placeholderIndex)
	zap.L().Sugar().Infof("苹果内存修改动态sql:%s", updateSql)

	var cpuId int64

	args = append(args, memory.Id)
	if err = tx.QueryRow(updateSql, args...).Scan(&cpuId); err != nil {
		zap.L().Sugar().Errorf("苹果内存修改执行错误===%+v", err)
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("未查询到数据，请确认参数有效性")
		}
		return nil, err
	}

	zap.L().Sugar().Infof("苹果内存修改返回ID:%d", cpuId)

	return selectMemoryDetail(tx, cpuId)
}

func selectMemoryDetail(tx *sql.Tx, cpuId int64) (*domain.MacMemory, error) {
	memory := &domain.MacMemory{}
	// 使用tx的查询，保证与插入操作在同一个事务中
	if err := tx.QueryRow(fmt.Sprintf(`SELECT id, %s FROM tb_mac_memory WHERE id = $1`, MacMemoryCommonColumn), cpuId).
		Scan(
			&memory.Id,
			&memory.MemorySize,
			&memory.MemorySpeed,
			&memory.IntegrationFlag,
			&memory.MemoryType,
			&memory.EccCheck,
		); err != nil {
		zap.L().Sugar().Errorf("通过ID查询苹果内存错误===%+v", err)
		return nil, err
	}

	return memory, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) DeleteById(id int64) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存删除开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果内存删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果内存删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果内存删除提交事务失败===%+v", err)
			}
		}
	}()

	result, err := tx.Exec(`DELETE FROM tb_mac_memory WHERE id = $1`, id)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果内存删除执行获取条数===%d", rowsAffected)

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacMemoryDb) BatchDeleteByIds(ids []any) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量删除开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("苹果内存批量删除事务即将回滚（panic恢复）")
			_ = tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			_ = tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("苹果内存批量删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				zap.L().Sugar().Errorf("苹果内存批量删除提交事务失败===%+v", err)
			}
		}
	}()

	deleteSql := fmt.Sprintf("DELETE FROM tb_mac_memory WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids)))
	zap.L().Sugar().Infof("苹果内存批量删除执行sql===%s", deleteSql)

	result, err := tx.Exec(deleteSql, ids...)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量删除执行获取条数错误===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果内存批量删除执行成功条数===%d", rowsAffected)

	return rowsAffected, nil
}

// memoryDynamicUpdate 构建动态更新
func memoryDynamicUpdate(memory *domain.MacMemory) (string, []any, int64, error) {
	var setClauses []string
	var args []interface{}
	placeholderIndex := int64(1)

	addClause := func(field string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, placeholderIndex))
		args = append(args, value)
		placeholderIndex++
	}

	if memory.MemorySize != "" {
		addClause(utils.CamelToSnakeCase("MemorySize"), memory.MemorySize)
	}
	if memory.MemorySpeed != "" {
		addClause(utils.CamelToSnakeCase("MemorySpeed"), memory.MemorySpeed)
	}
	if memory.IntegrationFlag != "" {
		addClause(utils.CamelToSnakeCase("IntegrationFlag"), memory.IntegrationFlag)
	}
	if memory.MemoryType != "" {
		addClause(utils.CamelToSnakeCase("MemoryType"), memory.MemoryType)
	}
	if memory.EccCheck != "" {
		addClause(utils.CamelToSnakeCase("EccCheck"), memory.EccCheck)
	}

	if len(setClauses) == 0 {
		zap.L().Sugar().Info("no fields to update")
		return "", nil, 0, errors.New("没有需要更新的列信息")
	}

	setClause := strings.Join(setClauses, ", ")
	return setClause, args, placeholderIndex, nil
}
