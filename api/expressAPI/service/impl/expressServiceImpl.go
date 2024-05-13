package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type ExpressDB struct {
	Db *sql.DB
}

func NewExpressDB(db *sql.DB) *ExpressDB {
	return &ExpressDB{
		Db: db,
	}
}

func (e *ExpressDB) CreateExpress(express *domain.Express) (*domain.Express, error) {
	rows, err := e.Db.Exec("INSERT INTO tool_express_manage(user_id, express_name, express_number, from_name, from_phone, from_address, pickup_code, create_by, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())",
		express.UserId, express.ExpressName, express.ExpressNumber, express.FromName, express.FromPhone, express.FromAddress, express.PickupCode, express.CreateBy, time.Now())
	if err != nil {
		return nil, err
	}
	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}
	//express.ID = id
	//express.CreateTime = FormatTime(timeNow)
	express, err = e.GetExpress(id)
	if err != nil {
		log.Printf("获取快递异常===%v", err)
		return nil, err
	}
	return express, nil
}

func (e *ExpressDB) GetExpress(id int64) (*domain.Express, error) {
	// 执行查询操作
	row := e.Db.QueryRow("SELECT id, IFNULL(user_id, ''), IFNULL(express_name, ''), IFNULL(express_number, ''), IFNULL(from_name, ''), IFNULL(from_phone, ''), IFNULL(from_address, ''), IFNULL(pickup_code, ''), IFNULL(create_by, ''), IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s' ), ''), IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s' ), '') FROM tool_express_manage WHERE id = ?", id)
	// 创建 Express 对象
	express := &domain.Express{}

	// 从查询结果中扫描数据到 Express 对象
	err := row.Scan(
		&express.ID,
		&express.UserId,
		&express.ExpressName,
		&express.ExpressNumber,
		&express.FromName,
		&express.FromPhone,
		&express.FromAddress,
		&express.PickupCode,
		&express.CreateBy,
		&express.CreateTime,
		&express.UpdateTime,
	)
	if err != nil {
		return nil, err
	}

	// 返回查询结果
	return express, nil
}

func (e *ExpressDB) SelectExpressPage(expressName string, pageStr, sizeStr string) ([]*domain.Express, int64, int64, error) {
	// 查询总记录数
	var totalRecords int64
	countSql := "SELECT COUNT(*) FROM tool_express_manage" + buildWhereClause(expressName)
	log.Printf("查询条数sql===%s", countSql)
	err := e.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}

	// 计算LIMIT的偏移量
	//offset := (page - 1) * size
	// 计算总页数
	//totalPages := (totalRecords + size - 1) / size
	//转换string strconv.FormatInt(page, 10)
	size, offset, totalPages := BuildPageOffset(pageStr, sizeStr, totalRecords)

	// 执行查询操作
	// 构建 SQL 查询语句
	limitQuerySql := "SELECT id, IFNULL(user_id, ''), IFNULL(express_name, ''), IFNULL(express_number,''), IFNULL(from_name, ''), IFNULL(from_phone, ''), IFNULL(from_address, ''), IFNULL(pickup_code, ''), IFNULL(create_by, ''), IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s' ), ''), IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s' ), '') FROM tool_express_manage" + buildWhereClause(expressName) + " LIMIT ?, ?"
	log.Printf("查询分页sql===%s, offset===%d, size====%d", limitQuerySql, offset, size)
	rows, err := e.Db.Query(limitQuerySql, offset, size)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}
	// defer rows.Close()

	// 创建 Express 对象数组
	//var expressesArray []*Express

	// 遍历查询结果
	/*for rows.Next() {
		// 创建 Express 对象
		express := &Express{}

		// 从查询结果中扫描数据到 Express 对象
		err := rows.Scan(
			&express.ID,
			&express.UserId,
			&express.ExpressName,
			&express.ExpressNumber,
			&express.FromName,
			&express.FromPhone,
			&express.FromAddress,
			&express.PickupCode,
			&express.CreateBy,
			&express.CreateTime,
		)
		if err != nil {
			log.Print(err)
			return nil, 0, 0, err
		}

		// 将 Express 对象添加到数组中
		expressesArray = append(expressesArray, express)
	}*/

	// 检查遍历过程中是否有错误
	//if err := rows.Err(); err != nil {
	//	log.Print(err)
	//	return nil, 0, 0, err
	//}
	expressesList, err := e.buildPageData(rows)
	if err != nil {
		return expressesList, 0, 0, err
	}
	RowsClose(rows, "快递分页查询")
	// 返回查询结果数组
	return expressesList, totalRecords, totalPages, nil
}

func (e *ExpressDB) SelectExpressPageByParam(param *param.ExpressSearchParam) ([]*domain.Express, int64, int64, error) {
	// 查询总记录数
	var totalRecords int64
	countSql := "SELECT COUNT(*) FROM tool_express_manage" + buildWhereClauseByParam(param)
	//log.Printf("查询条数sql===%s", countSql)
	err := e.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}

	size, offset, totalPages := BuildPageOffset(param.Page, param.Size, totalRecords)

	// 构建 SQL 查询语句
	limitQuerySql := "SELECT id, IFNULL(user_id, ''), IFNULL(express_name, ''), IFNULL(express_number,''), IFNULL(from_name, ''), IFNULL(from_phone, ''), IFNULL(from_address, ''), IFNULL(pickup_code, ''), IFNULL(create_by, ''), IFNULL(DATE_FORMAT(create_time, '%Y-%m-%d %H:%i:%s' ), ''), IFNULL(DATE_FORMAT(update_time, '%Y-%m-%d %H:%i:%s' ), '') FROM tool_express_manage" + buildOrderBy(param) + buildWhereClauseByParam(param) + " LIMIT ?, ?"
	log.Printf("查询分页sql===%s, offset===%d, size====%d", limitQuerySql, offset, size)
	rows, err := e.Db.Query(limitQuerySql, offset, size)
	if err != nil {
		log.Print(err)
		return nil, 0, 0, err
	}
	// todo 这里的rows必须关闭，使用defer方式时候，rows,Close()位置不重要，但是使用自定义方法关闭rows，就得注意下关闭调用的顺序
	// defer rows.Close()

	expressesArray, err := e.buildPageData(rows)
	if err != nil {
		return expressesArray, 0, 0, err
	}

	RowsClose(rows, "快递分页查询")

	// 返回查询结果数组
	return expressesArray, totalRecords, totalPages, nil
}

func (e *ExpressDB) buildPageData(rows *sql.Rows) ([]*domain.Express, error) {

	// 创建 Express 对象数组
	var expressesList []*domain.Express

	// 遍历查询结果
	for rows.Next() {
		// 创建 Express 对象
		express := &domain.Express{}

		// 从查询结果中扫描数据到 Express 对象
		err := rows.Scan(
			&express.ID,
			&express.UserId,
			&express.ExpressName,
			&express.ExpressNumber,
			&express.FromName,
			&express.FromPhone,
			&express.FromAddress,
			&express.PickupCode,
			&express.CreateBy,
			&express.CreateTime,
			&express.UpdateTime,
		)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		// 将 Express 对象添加到数组中
		expressesList = append(expressesList, express)
	}

	// 检查遍历过程中是否有错误
	if err := rows.Err(); err != nil {
		log.Print(err)
		return nil, err
	}
	return expressesList, nil
}

func (e *ExpressDB) DeleteById(id int64) (int64, error) {
	rows, err := e.Db.Exec("DELETE FROM tool_express_manage WHERE id = ?", id)
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	if rowsAffected == 0 {
		return 0, errors.New("未找到要删除的记录")
	}
	return rowsAffected, nil
}

func (e *ExpressDB) UpdateExpress(express *domain.Express) (int64, error) {
	rows, err := e.Db.Exec("UPDATE tool_express_manage SET user_id = ?, express_name = ?, express_number = ?, from_name = ?, from_phone = ?, from_address= ?, pickup_code = ?, create_by = ?, update_time = ? WHERE id = ?",
		express.UserId, express.ExpressName, express.ExpressNumber, express.FromName, express.FromPhone, express.FromAddress, express.PickupCode, express.CreateBy, time.Now(), express.ID)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	count, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (e *ExpressDB) BatchDeleteByIds(ids []string) (int64, error) {
	// 使用逗号分隔的字符串构建 IN 子句
	deleteSql := fmt.Sprintf("DELETE FROM tool_express_manage WHERE id IN (%s)", strings.Join(ids, ","))
	log.Printf("批量删除sql===%s", deleteSql)
	rows, err := e.Db.Exec(deleteSql)
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		log.Println(err)
		return 0, nil
	}
	if rowsAffected == 0 {
		return 0, errors.New("未找到要删除的记录")
	}
	return rowsAffected, nil
}

func (e *ExpressDB) BatchCreateExpress(list []*domain.Express) (int64, error) {
	var valueList []string
	var valueArgs []interface{}

	for _, express := range list {
		valueList = append(valueList, "(?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())")
		valueArgs = append(valueArgs, express.UserId, express.ExpressName, express.ExpressNumber, express.FromName, express.FromPhone, express.FromAddress, express.PickupCode, express.CreateBy)
	}

	stmt := fmt.Sprintf("INSERT INTO tool_express_manage(user_id, express_name, express_number, from_name, from_phone, from_address, pickup_code, create_by, create_time, update_time) VALUES %s", strings.Join(valueList, ","))
	// 执行批量新增sql
	rows, err := e.Db.Exec(stmt, valueArgs...)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	count, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// 构建 WHERE 子句
func buildWhereClause(expressName string) string {
	if expressName != "" {
		//return fmt.Sprintf(" WHERE express_name LIKE CONCAT('%%', '%s', '%%')", expressName)
		return " WHERE express_name LIKE CONCAT('%', '" + expressName + "', '%')"
	}
	return ""
}

// buildWhereClauseByParam 构建多条件动态查询
func buildWhereClauseByParam(param *param.ExpressSearchParam) string {
	if param != nil {
		var clauses []string
		if param.ExpressName != "" {
			clauses = append(clauses, "express_name LIKE CONCAT('%', '"+param.ExpressName+"', '%')")
		}
		if param.ExpressNumber != "" {
			clauses = append(clauses, "express_number = '"+param.ExpressNumber+"'")
		}
		if param.UserId != "" {
			clauses = append(clauses, "user_id = '"+param.UserId+"'")
		}
		if param.FromName != "" {
			clauses = append(clauses, "from_name LIKE CONCAT('%', '"+param.FromName+"', '%')")
		}
		if param.FromPhone != "" {
			clauses = append(clauses, "from_phone LIKE CONCAT('%', '"+param.FromPhone+"', '%')")
		}
		if param.FromAddress != "" {
			clauses = append(clauses, "from_address LIKE CONCAT('%', '"+param.FromAddress+"', '%')")
		}
		if param.PickupCode != "" {
			clauses = append(clauses, "pickup_code LIKE CONCAT('%', '"+param.PickupCode+"', '%')")
		}
		if param.CreateBy != "" {
			clauses = append(clauses, "create_by LIKE CONCAT('%', '"+param.CreateBy+"', '%')")
		}

		if len(clauses) > 0 {
			return " WHERE " + strings.Join(clauses, " AND ")
		}
	}
	return ""
}

func buildOrderBy(param *param.ExpressSearchParam) string {
	if param != nil {
		var orderBy string
		field := param.Field
		order := param.Order

		if field == "" && order == "" {
			orderBy = " ORDER BY id DESC"
		} else {
			orderBy = " ORDER BY " + field + " " + order
		}
		return orderBy
	}
	return ""
}
