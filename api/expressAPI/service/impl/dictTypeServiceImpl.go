package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type DictTypeDb struct {
	Db *sql.DB
}

func NewDictTypeDb(pg *sql.DB) *DictTypeDb {
	return &DictTypeDb{
		Db: pg,
	}
}

//goland:noinspection SqlResolve
func (pg *DictTypeDb) GetDictList(d *domain.DictType, page, sizeStr string) ([]*domain.DictType, int64, int64, error) {
	var dictTypeList []*domain.DictType
	var totalRecords int64

	var err error

	countSql := "SELECT COUNT(*) FROM tb_sys_dict_type" + buildWhereParam(d)
	log.Println("字典类型分页countSql===", countSql)
	err = pg.Db.QueryRow(countSql).Scan(&totalRecords)
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
	rows, err := pg.Db.Query(querySql)
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

//goland:noinspection SqlResolve
func (pg *DictTypeDb) SelectDictTypeById(id int64) (*domain.DictType, error) {
	var dictType = &domain.DictType{}
	queryRow := pg.Db.QueryRow(`
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
					tb_sys_dict_type WHERE id = $1`, id)

	err := queryRow.Scan(
		&dictType.Id,
		&dictType.DictName,
		&dictType.DictType,
		&dictType.TypeStatus,
		&dictType.CreateBy,
		&dictType.CreateTime,
		&dictType.UpdateBy,
		&dictType.UpdateTime,
		&dictType.Remark,
		&dictType.DelFlag)
	if err != nil {
		log.Printf("字典详情查询错误===%v", err)
		return nil, err
	}
	return dictType, nil
}

// SaveDictType 保存字典类型
//
//goland:noinspection SqlResolve
func (pg *DictTypeDb) SaveDictType(dt *domain.DictType) (*domain.DictType, error) {
	result, err := pg.Db.Exec(
		`INSERT INTO
			  		tb_sys_dict_type (dict_name, dict_type, type_status, create_by, create_time, update_by, update_time, remark)
		 	   VALUES
			   		($1, $2, $3, $4, CURRENT_TIMESTAMP, $5, CURRENT_TIMESTAMP, $6)`,
		dt.DictName, dt.DictType, dt.TypeStatus, dt.CreateBy, dt.UpdateBy, dt.Remark)

	if err != nil {
		log.Printf("字典新增错误===%v", err)
		return nil, err
	}

	// postgresql不能返回新增的ID，只能通过其他方式查询新增后的那条数据
	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("获取插入行数错误===%v", err)
		return nil, err
	}

	var dictType *domain.DictType
	if rowCount == 1 {
		dictType, err = pg.SelectDetailByObj(domain.NewDictTypeDetail(dt.DictType))
		if err != nil {
			log.Printf("字典新增失败===%v", err)
			return nil, errors.New("字典新增失败")
		}
	} else {
		log.Println("字典新增失败")
		return nil, errors.New("字典新增失败")
	}

	return dictType, nil
}

//goland:noinspection SqlResolve
func (pg *DictTypeDb) UpdateDictType(d *domain.DictType) (int64, error) {
	result, err := pg.Db.Exec(
		`UPDATE tb_sys_dict_type
		  SET
			dict_name = $1,
			dict_type = $2,
			type_status = $3,
			update_by = $4,
			update_time = CURRENT_TIMESTAMP,
			remark = $5
		WHERE
			id = $6`,
		d.DictName, d.DictType, d.TypeStatus, d.UpdateBy, d.Remark, d.Id)

	if err != nil {
		log.Printf("字典修改错误===%v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("字典修改失败===%v", err)
		return 0, errors.New("字典修改失败")
	}
	return rowsAffected, nil
}

//goland:noinspection SqlResolve
func (pg *DictTypeDb) DeleteDictType(id int64) bool {
	result, err := pg.Db.Exec(`DELETE FROM tb_sys_dict_type WHERE id = $1`, id)
	if err != nil {
		log.Printf("字典删除错误===%v", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("字典删除失败===%v", err)
		return false
	}

	if rowsAffected != 1 {
		return false
	}
	return true
}

//goland:noinspection SqlResolve
func (pg *DictTypeDb) SelectDetailByObj(dt domain.DictType) (*domain.DictType, error) {
	var dictType = &domain.DictType{}
	queryRow := pg.Db.QueryRow(`
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
					tb_sys_dict_type WHERE dict_type = $1`, dt.DictType)

	err := queryRow.Scan(
		&dictType.Id,
		&dictType.DictName,
		&dictType.DictType,
		&dictType.TypeStatus,
		&dictType.CreateBy,
		&dictType.CreateTime,
		&dictType.UpdateBy,
		&dictType.UpdateTime,
		&dictType.Remark,
		&dictType.DelFlag)
	if err != nil {
		log.Printf("字典详情查询错误===%v", err)
		return nil, err
	}
	return dictType, nil
}

//goland:noinspection SqlResolve
func (pg *DictTypeDb) CheckTypeIsExist(id *int64, typeStr string) bool {
	var count int64
	err := pg.Db.QueryRow("SELECT COUNT(*) FROM tb_sys_dict_type WHERE dict_type = $1 AND id != $2", typeStr, id).Scan(&count)
	if err != nil {
		log.Printf("字典类型验证错误===%v", err)
		return false
	}
	if count == 0 {
		return false
	}
	return true
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
