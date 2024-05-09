package impl

import "apiProject/api/utils"

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
