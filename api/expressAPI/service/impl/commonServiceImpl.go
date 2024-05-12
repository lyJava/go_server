package impl

import (
	"apiProject/api/utils"
	"database/sql"
	"log"
)

// BuildPageOffset 计算总页数，偏移量
func BuildPageOffset(pageStr, sizeStr interface{}, totalRecords int64) (int64, int64, int64) {
	var page, size int64

	switch v := pageStr.(type) {
	case int64:
		page = v
	case string:
		page = utils.ConvertToInt64(v)
	default:
		return 0, 0, 0
	}

	switch v := sizeStr.(type) {
	case int64:
		size = v
	case string:
		size = utils.ConvertToInt64(v)

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
		log.Printf(message+"关闭row结果错误===%v", err)
		return err
	} 
	log.Println(message + "结果成功关闭")
	return nil
}
