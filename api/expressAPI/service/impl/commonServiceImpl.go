package impl

import (
	"database/sql"

	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	// 苹果笔记本列
	MacBookCommonColumn = `b.pro_name, b.pro_color, b.pro_year, b.pro_type, b.mold_model, b.monitor_id, b.cpu_id, b.memory_id, b.storage_info, b.size_info, b.weight_info, b.gpu_id, b.chip_type, b.memory_max, 
							CASE WHEN b.create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS') ELSE '' END AS create_time, 
							CASE WHEN b.update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS') ELSE '' END AS update_time, 
							c.id, c.cpu_name, c.cpu_type, c.cpu_basic_boost, c.cpu_trubo_boost, c.cpu_core_number, c.cpu_thread_number, c.cpu_cache, c.cpu_tdp, c.memory_width, c.media_processing_engine,
							m.id, m.memory_size, m.memory_speed, m.integration_flag, m.memory_type, m.ecc_check`
	// 苹果处理器列
	MacCpuCommonColumn = "cpu_name, cpu_type, cpu_basic_boost, cpu_trubo_boost, cpu_core_number, cpu_thread_number, cpu_cache, cpu_tdp, memory_width, media_processing_engine"
	// 苹果内存列
	MacMemoryCommonColumn = "memory_size, memory_speed, integration_flag, memory_type, ecc_check"
)

// BuildPageOffset 计算总页数，偏移量
func BuildPageOffset(pageStr, sizeStr interface{}, totalRecords int64) (int64, int64, int64) {
	var page, size int64

	switch v := pageStr.(type) {
	case int64:
		page = v
	case string:
		page = cast.ToInt64(v)
	default:
		return 0, 0, 0
	}

	switch v := sizeStr.(type) {
	case int64:
		size = v
	case string:
		size = cast.ToInt64(v)

	default:
		return 0, 0, 0
	}

	// 计算LIMIT的偏移量
	offset := (page - 1) * size
	// 计算总页数
	totalPages := (totalRecords + size - 1) / size
	return size, offset, totalPages
}

// RowsClose 关闭rows
func RowsClose(rows *sql.Rows, message string) error {
	err := rows.Close()
	if err != nil {
		zap.L().Sugar().Errorf(message+"关闭row结果错误===%+v", err)
		return err
	}
	zap.L().Sugar().Infoln(message + "结果成功关闭")
	return nil
}

func buildOrderBy(column, order, tableAlias string) string {
	var orderBy string

	if tableAlias == "" {
		if column == "" && order == "" {
			orderBy = " ORDER BY id DESC"
		} else {
			orderBy = " ORDER BY " + column + " " + order
		}
	} else {
		if column == "" && order == "" {
			orderBy = " ORDER BY " + tableAlias + ".id DESC"
		} else {
			orderBy = " ORDER BY " + tableAlias + "." + column + " " + order
		}
	}

	return orderBy
}
