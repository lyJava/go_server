package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type DictTypeDb struct {
	Db *sql.DB
}

func NewDictTypeDb(db *sql.DB) *DictTypeDb {
	return &DictTypeDb{
		Db: db,
	}
}

func (db *DictTypeDb) GetDictList(d *domain.DictType, page, sizeStr string) ([]*domain.DictType, int64, int64, error) {
	var dictTypeList []*domain.DictType
	var totalRecords int64

	var err error

	countSql := "SELECT COUNT(*) FROM tb_sys_dict_type" + buildWhereParam(d)
	log.Println("字典类型分页countSql===", countSql)
	err = db.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}

	pageSize, offset, totalPages := BuildPageOffset(page, sizeStr, totalRecords)
	querySql := fmt.Sprintf(`
				SELECT
					id,
					dict_name,
					dict_type,
					CASE
						WHEN type_status = '0' THEN '正常'
						ELSE '停用'
					END AS type_status,
					create_by,
					CASE
						WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS create_time,
					update_by,
					CASE
						WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS update_time,
					remark,
					CASE
						WHEN del_flag = '0' THEN '正常'
						ELSE '删除'
					END AS del_flag
				FROM
					tb_sys_dict_type`+buildWhereParam(d)+" LIMIT %d OFFSET %d", pageSize, offset)
	log.Println("字典类型分页查询sql===", querySql)
	rows, err := db.Db.Query(querySql)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		dictType := domain.DictType{}
		err = rows.Scan(
			&dictType.Id,
			&dictType.DictName,
			&dictType.DictType,
			&dictType.TypeStatus,
			&dictType.CreateBy,
			&dictType.CreateTime,
			&dictType.UpdateBy,
			&dictType.UpdateTime,
			&dictType.Remark,
			&dictType.DelFlag,
		)

		if err != nil {
			log.Print(err)
			return nil, 0, 0, err
		}

		dictTypeList = append(dictTypeList, &dictType)
	}

	return dictTypeList, totalPages, totalRecords, nil
}

// buildWhereParam 构建多条件动态查询
func buildWhereParam(dict *domain.DictType) string {
	if dict != nil {
		var clauses []string
		dictName := dict.DictName
		dictType := dict.DictType
		remark := dict.Remark

		if dictName != "" {
			clauses = append(clauses, "dict_name LIKE CONCAT('%', '"+dictName+"', '%')")
		}
		if dictType != "" {
			clauses = append(clauses, "dict_type = '"+dictType+"'")
		}

		if remark != "" {
			clauses = append(clauses, "remark LIKE CONCAT('%', '"+remark+"', '%')")
		}

		if len(clauses) > 0 {
			return " WHERE " + strings.Join(clauses, " AND ")
		}
	}
	return ""
}
