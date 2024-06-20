package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
	"log"
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
		log.Printf("苹果处理器新增开启事务失败===%+v", err)
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
				log.Printf("苹果处理器新增提交事务失败===%+v", err)
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
	macCpu := &domain.MacCpu{}
	if err = tx.QueryRow(fmt.Sprintf(`SELECT id, %s FROM tb_mac_cpu WHERE id = $1`, CommonColumn), lastInsertId).
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
		log.Printf("通过ID查询苹果处理器异常===%+v", err)
		zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
		return nil, err
	}

	// 显式提交事务
	/*if err = tx.Commit(); err != nil {
		log.Printf("苹果处理器新增提交事务失败===%+v", err)
		zap.L().Sugar().Errorf("苹果处理器新增提交事务失败===%+v", err)
		return nil, err
	}*/

	return macCpu, nil
}

func (pg *MacCpuDb) BatchSave(list []*domain.MacCpu) (int64, error) {
	return 0, nil
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
			log.Printf("通过ID查询苹果处理器异常===%+v", err)
			zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
			return nil, fmt.Errorf("no MacCpu found with id %d", id)
		}
		log.Printf("通过ID查询苹果处理器异常===%+v", err)
		zap.L().Sugar().Errorf("通过ID查询苹果处理器错误===%+v", err)
		return nil, err
	}
	return cpu, nil
}

func (pg *MacCpuDb) Update(cpu *domain.MacCpu) (int64, error) {
	return 0, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlNoDataSourceInspection
func (pg *MacCpuDb) DeleteById(id int64) (int64, error) {
	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		log.Printf("苹果处理器删除开启事务失败===%+v", err)
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
				log.Printf("苹果处理器删除提交事务失败===%+v", err)
				zap.L().Sugar().Errorf("苹果处理器删除提交事务失败===%+v", err)
			}
		}
	}()

	result, err := tx.Exec(`DELETE FROM tb_mac_cpu WHERE id = $1`, id)
	if err != nil {
		log.Printf("苹果处理器删除异常===%+v", err)
		zap.L().Sugar().Errorf("苹果处理器删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("苹果处理器删除执行获取条数===%+v", err)
		zap.L().Sugar().Errorf("苹果处理器删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("苹果处理器删除执行获取条数===%d", rowsAffected)

	return rowsAffected, nil
}
