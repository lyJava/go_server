package impl

import (
	"database/sql"
	"log"

	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
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
		log.Printf(message+"关闭row结果错误===%+v", err)
		zap.L().Sugar().Errorf(message+"关闭row结果错误===%+v", err)
		return err
	}
	log.Println(message + "结果成功关闭")
	zap.L().Sugar().Infoln(message+"结果成功关闭")
	return nil
}
