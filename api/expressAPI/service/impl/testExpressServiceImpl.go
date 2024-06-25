package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap"
)

type TestDictTypeDb struct {
	Db *sql.DB
}

func NewTestDictTypeDb(pg *sql.DB) *TestDictTypeDb {
	return &TestDictTypeDb{
		Db: pg,
	}
}

//goland:noinspection SqlResolve
func (pg *TestDictTypeDb) Save(te *domain.TestExpress) (*domain.TestExpress, error) {
	// 先检查是否存在
	isExist := pg.CheckByNumAndPickupCode(te.ExpressNumber, te.PickupCode)
	if isExist {
		log.Printf("快递单号:%s与取件码:%s不能重复", te.ExpressNumber, te.PickupCode)
		zap.L().Sugar().Errorf("快递单号:%s与取件码:%s不能重复", te.ExpressNumber, te.PickupCode)
		return nil, fmt.Errorf("快递单号:%s与取件码:%s不能重复", te.ExpressNumber, te.PickupCode)
	}

	result, err := pg.Db.Exec(`INSERT INTO tb_test_express(express_name, express_number, pickup_code, from_username, 
                            from_user_phone, from_user_address, from_user_id_number, create_by, remarks, del_flag)
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		te.ExpressName,
		te.ExpressNumber,
		te.PickupCode,
		te.FromUsername,
		te.FromUserPhone,
		te.FromUserAddress,
		te.FromUserIdNumber,
		te.CreateBy,
		te.Remarks,
		te.DelFlag,
	)

	if err != nil {
		log.Printf("测试快递新增错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递新增错误===%+v", err)
		return nil, err
	}

	// postgresql不能返回新增的ID，只能通过其他方式查询新增后的那条数据
	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Printf("测试快递新增获取行数错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递新增获取行数错误===%+v", err)
		return nil, err
	}

	zap.L().Sugar().Infof("测试快递新增获取行数===%d", rowCount)

	var data *domain.TestExpress
	if rowCount == 1 {
		data, err = pg.SelectByNumAndPickupCode(te.ExpressNumber, te.PickupCode)
		if err != nil {
			log.Printf("测试查询快递信息失败===%+v", err)
			zap.L().Sugar().Errorf("测试查询快递信息失败===%+v", err)
			return nil, err
		}
	}
	return data, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *TestDictTypeDb) SelectByNumAndPickupCode(expressNumber, pickupCode string) (*domain.TestExpress, error) {
	// 返回值这里不需要加上*号
	testExpress := &domain.TestExpress{}
	row := pg.Db.QueryRow(
		`SELECT
				id,
				express_name,
				express_number,
				pickup_code,
				from_username,
				from_user_phone,
				from_user_address,
				from_user_id_number,
				create_by,
				CASE
					WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS create_time,
				CASE
					WHEN update_by IS NULL THEN ''
					ELSE update_by
				END AS update_by,
				CASE
					WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
					ELSE ''
				END AS update_time,
				remarks,
				CASE
					WHEN del_flag = 0 THEN '正常'
					ELSE '删除'
				END AS del_flag
			FROM
				tb_test_express
			WHERE
				express_number = $1
				AND pickup_code = $2`, expressNumber, pickupCode)
	err := rowScan(row, testExpress)
	if err != nil {
		log.Printf("通过快递单号与取件码查询错误===%+v", err)
		zap.L().Sugar().Errorf("通过快递单号与取件码查询错误===%+v", err)
		return nil, err
	}

	// 这里需要加上&
	return testExpress, nil
}

//goland:noinspection SqlResolve, SqlError
func (pg *TestDictTypeDb) BatchSave(list []*domain.TestExpress) (int64, error) {
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
			testExpress := list[j]
			placeholderList = append(placeholderList, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				len(valueArgList)+1, len(valueArgList)+2, len(valueArgList)+3, len(valueArgList)+4, len(valueArgList)+5,
				len(valueArgList)+6, len(valueArgList)+7, len(valueArgList)+8, len(valueArgList)+9, len(valueArgList)+10))
			valueArgList = append(valueArgList,
				testExpress.ExpressName,
				testExpress.ExpressNumber,
				testExpress.PickupCode,
				testExpress.FromUsername,
				testExpress.FromUserPhone,
				testExpress.FromUserAddress,
				testExpress.FromUserIdNumber,
				testExpress.CreateBy,
				testExpress.Remarks,
				testExpress.DelFlag,
			)
		}
	}

	batchSql := fmt.Sprintf(`INSERT INTO tb_test_express (express_name, express_number, pickup_code, from_username, 
                             		from_user_phone,from_user_address,from_user_id_number,create_by,remarks, 
                             		del_flag) VALUES %s %s`, strings.Join(placeholderList, ","), "\r\n")
	log.Println("测试快递批量新增sql===", batchSql)
	zap.L().Sugar().Infof("测试快递批量新增sql===%s", batchSql)
	result, err := pg.Db.Exec(batchSql, valueArgList...)
	if err != nil {
		log.Printf("测试快递批量新增执行错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递批量新增执行错误===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("测试快递批量新增执行失败===%+v", err)
		zap.L().Sugar().Errorf("测试快递批量新增执行失败===%+v", err)
		return 0, err
	}

	return rowsAffected, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf,SqlCaseVsLimit
func (pg *TestDictTypeDb) PageList(te *domain.TestExpress, page, size int64) ([]*domain.TestExpress, int64, int64, error) {
	// 查询总记录数
	var totalRecords int64
	countSql := "SELECT COUNT(*) FROM tb_test_express" + buildCountByEntity(te)
	zap.L().Sugar().Infof("测试快递分页查询count的sql===%s", countSql)

	err := pg.Db.QueryRow(countSql).Scan(&totalRecords)
	if err != nil {
		zap.L().Sugar().Errorf("测试快递分页查询总条数错误===%+v", err)
		return nil, 0, 0, err
	}

	zap.L().Sugar().Infof("测试快递分页查询总条数===%d", totalRecords)

	size, offset, totalPages := BuildPageOffset(page, size, totalRecords)

	rows, err := pg.Db.Query(fmt.Sprintf(
		`SELECT
					id,
					express_name,
					express_number,
					pickup_code,
					from_username,
					from_user_phone,
					from_user_address,
					from_user_id_number,
					create_by,
					CASE
						WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS create_time,
					CASE
						WHEN update_by IS NULL THEN ''
						ELSE update_by
					END AS update_by,
					CASE
						WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
						ELSE ''
					END AS update_time,
					remarks,
					CASE
						WHEN del_flag = 0 THEN '正常'
						ELSE '删除'
					END AS del_flag
				FROM
					tb_test_express %s
				ORDER BY id DESC	
				OFFSET $1 LIMIT $2`, buildCountByEntity(te)), offset, size)

	if err != nil {
		log.Printf("测试快递分页查询行数错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递分页查询行数错误===%+v", err)
		return nil, 0, 0, err
	}

	defer RowsClose(rows, "测试快递分页查询")

	var list []*domain.TestExpress

	for rows.Next() {
		testExpress := &domain.TestExpress{}
		err := rowsScan(rows, testExpress)
		if err != nil {
			log.Printf("测试快递分页返回异常===%+v", err.Error())
			zap.L().Sugar().Errorf("测试快递分页返回错误===%+v", err)
			return nil, 0, 0, err
		}

		list = append(list, testExpress)
	}
	return list, totalRecords, totalPages, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *TestDictTypeDb) SelectById(id int64) (*domain.TestExpress, error) {
	row := pg.Db.QueryRow(
		`SELECT
			id,
			express_name,
			express_number,
			pickup_code,
			from_username,
			from_user_phone,
			from_user_address,
			from_user_id_number,
			CASE
				WHEN create_by IS NULL THEN ''
				ELSE create_by
			END AS create_by,
			CASE
				WHEN create_time IS NOT NULL THEN to_char(create_time, 'YYYY-MM-DD HH24:MI:SS')
				ELSE ''
			END AS create_time,
			CASE
				WHEN update_by IS NULL THEN ''
				ELSE update_by
			END AS update_by,
			CASE
				WHEN update_time IS NOT NULL THEN to_char(update_time, 'YYYY-MM-DD HH24:MI:SS')
				ELSE ''
			END AS update_time,
			remarks,
			CASE
				WHEN del_flag = 0 THEN '正常'
				ELSE '删除'
			END AS del_flag
		FROM
			tb_test_express
		WHERE id = $1`, id)

	textExpress := &domain.TestExpress{}
	err := rowScan(row, textExpress)
	if err != nil {
		log.Printf("通过ID查询测试快递异常===%+v", err)
		zap.L().Sugar().Errorf("通过ID查询测试快递错误===%+v", err)
		return nil, err
	}
	return textExpress, nil
}

//goland:noinspection SqlResolve
func (pg *TestDictTypeDb) CheckByNumAndPickupCode(expressNumber, pickupCode string) bool {
	var selectCount int64
	err := pg.Db.QueryRow("SELECT COUNT(*) FROM tb_test_express WHERE express_number = $1 AND pickup_code = $2", expressNumber, pickupCode).Scan(&selectCount)
	if err != nil {
		log.Printf("通过快递单号与取件码查询异常===%+v", err)
		return false
	}
	if selectCount >= 1 {
		return true
	}
	return false
}

//goland:noinspection SqlResolve
func (pg *TestDictTypeDb) BatchDelete(ids []any) (rows int64, err error) {
	//deleteSql := fmt.Sprintf("DELETE FROM tb_test_express WHERE id IN (%s)", strings.Join(ids, ", "))

	// 开启事务
	tx, err := pg.Db.Begin()
	if err != nil {
		log.Printf("开启事务失败===%+v", err)
		zap.L().Sugar().Errorf("开启事务失败===%+v", err)
		return 0, err
	}

	// 确保事务的提交或回滚
	defer func() {
		if p := recover(); p != nil {
			zap.L().Sugar().Info("字典类型批量删除事务即将回滚（panic恢复）")
			tx.Rollback()
			panic(p) // 重新panic以便外层捕获
		} else if err != nil {
			tx.Rollback() // 发生错误则回滚事务
			zap.L().Sugar().Errorf("字典类型批量删除事务回滚,发生错误===%+v", err)
		} else {
			err = tx.Commit() // 正常结束则提交事务
			if err != nil {
				log.Printf("字典类型批量删除提交事务失败===%+v", err)
				zap.L().Sugar().Errorf("字典类型批量删除提交事务失败===%+v", err)
			}
		}
	}()

	delSql := fmt.Sprintf("DELETE FROM tb_test_express WHERE id IN %s", utils.BuildArgsWithBrackets(ids))
	log.Printf("测试快递批量删除sql===%s", delSql)
	zap.L().Sugar().Debugf("测试快递批量删除sql===%s", delSql)

	//result, err := pg.Db.Exec(fmt.Sprintf("DELETE FROM tb_test_express WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids))), ids...)
	result, err := tx.Exec(fmt.Sprintf("DELETE FROM tb_test_express WHERE id IN (%s)", utils.GeneratePlaceholders(len(ids))), ids...)
	if err != nil {
		log.Printf("测试快递批量删除异常===%+v", err)
		zap.L().Sugar().Errorf("测试快递批量删除异常===%+v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("测试快递批量删除执行获取条数===%+v", err)
		zap.L().Sugar().Errorf("测试快递批量删除执行获取条数===%+v", err)
		return 0, err
	}

	zap.L().Sugar().Infof("测试快递批量删除执行获取条数===%d", rowsAffected)

	return rowsAffected, nil
}

// 构建 WHERE 子句
func buildCountByEntity(te *domain.TestExpress) string {
	if te != nil {
		var clauses []string
		expressName := te.ExpressName
		expressNumber := te.ExpressNumber
		pickupCode := te.PickupCode
		fromUsername := te.FromUsername
		fromUserPhone := te.FromUserPhone
		fromUserAddress := te.FromUserAddress
		fromUserIdNumber := te.FromUserIdNumber
		createBy := te.CreateBy
		delFlag := te.DelFlag

		if expressName != "" {
			clauses = append(clauses, "express_name LIKE CONCAT('%', '"+expressName+"', '%')")
		}
		if expressNumber != "" {
			clauses = append(clauses, "express_number LIKE CONCAT('%', '"+expressNumber+"', '%')")
		}
		if pickupCode != "" {
			clauses = append(clauses, "pickup_code = '"+pickupCode+"'")
		}

		if fromUsername != "" {
			clauses = append(clauses, "from_username LIKE CONCAT('%', '"+fromUsername+"', '%')")
		}
		if fromUserPhone != "" {
			clauses = append(clauses, "from_user_phone LIKE CONCAT('%', '"+fromUserPhone+"', '%')")
		}
		if fromUserAddress != "" {
			clauses = append(clauses, "from_user_address LIKE CONCAT('%', '"+fromUserAddress+"', '%')")
		}
		if fromUserIdNumber != "" {
			clauses = append(clauses, "from_user_id_number ='"+fromUserIdNumber+"'")
		}

		if createBy != "" {
			clauses = append(clauses, "create_by LIKE CONCAT('%', '"+createBy+"', '%')")
		}

		if delFlag != "" {
			clauses = append(clauses, fmt.Sprintf("del_flag = %d", utils.ConvertToInt64(delFlag)))
		}

		if len(clauses) > 0 {
			return " WHERE " + strings.Join(clauses, " AND ")
		}
	}

	return ""
}

// rowScan 将数据行row转换到结构体
func rowScan(row *sql.Row, testExpress *domain.TestExpress) error {
	err := row.Scan(&testExpress.Id,
		&testExpress.ExpressName,
		&testExpress.ExpressNumber,
		&testExpress.PickupCode,
		&testExpress.FromUsername,
		&testExpress.FromUserPhone,
		&testExpress.FromUserAddress,
		&testExpress.FromUserIdNumber,
		&testExpress.CreateBy,
		&testExpress.CreateTime,
		&testExpress.UpdateBy,
		&testExpress.UpdateTime,
		&testExpress.Remarks,
		&testExpress.DelFlag)
	if err != nil {
		log.Printf("测试快递数据转换错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递数据转换错误===%+v", err)
		return err
	}
	return nil
}

// rowsScan 将数据行rows转换到结构体
func rowsScan(rows *sql.Rows, testExpress *domain.TestExpress) error {
	err := rows.Scan(&testExpress.Id,
		&testExpress.ExpressName,
		&testExpress.ExpressNumber,
		&testExpress.PickupCode,
		&testExpress.FromUsername,
		&testExpress.FromUserPhone,
		&testExpress.FromUserAddress,
		&testExpress.FromUserIdNumber,
		&testExpress.CreateBy,
		&testExpress.CreateTime,
		&testExpress.UpdateBy,
		&testExpress.UpdateTime,
		&testExpress.Remarks,
		&testExpress.DelFlag)
	if err != nil {
		log.Printf("测试快递数据多条转换错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递数据多条转换错误===%+v", err)
		return err
	}
	return nil
}
