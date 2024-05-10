package impl

import (
	"apiProject/api/expressAPI/types/domain"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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
		log.Printf("快递单号：%s与取件码：%s不能重复", te.ExpressNumber, te.PickupCode)
		return nil, errors.New(fmt.Sprintf("快递单号：%s与取件码：%s不能重复", te.ExpressNumber, te.PickupCode))
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
		log.Printf("快递新增错误===%v", err)
		return nil, err
	}

	// postgresql不能返回新增的ID，只能通过其他方式查询新增后的那条数据
	rowCount, err := result.RowsAffected()
	if err != nil {
		log.Printf("获取插入行数错误===%v", err)
		return nil, err
	}

	var data *domain.TestExpress
	if rowCount == 1 {
		data, err = pg.SelectByNumAndPickupCode(te.ExpressNumber, te.PickupCode)
		if err != nil {
			log.Printf("查询快递信息失败===%v", err)
			return nil, err
		}
	}
	return data, nil
}

//goland:noinspection SqlResolve,SqlCaseVsIf
func (pg *TestDictTypeDb) SelectByNumAndPickupCode(expressNumber, pickupCode string) (*domain.TestExpress, error) {
	// 返回值这里不需要加上*号
	var result domain.TestExpress
	err := pg.Db.QueryRow(
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
				AND pickup_code = $2`, expressNumber, pickupCode).
		Scan(&result.ID,
			&result.ExpressName,
			&result.ExpressNumber,
			&result.PickupCode,
			&result.FromUsername,
			&result.FromUserPhone,
			&result.FromUserAddress,
			&result.FromUserIdNumber,
			&result.CreateBy,
			&result.CreateTime,
			&result.UpdateBy,
			&result.UpdateTime,
			&result.Remarks,
			&result.DelFlag)
	if err != nil {
		log.Printf("通过快递单号与取件码查询异常===%v", err)
		return nil, err
	}

	// 这里需要加上&
	return &result, nil
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
	log.Println("生成的批量sql===", batchSql)
	result, err := pg.Db.Exec(batchSql, valueArgList...)
	if err != nil {
		log.Printf("批量新增执行错误===%v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("批量新增执行失败===%v", err)
		return 0, err
	}

	return rowsAffected, nil
}

func (pg *TestDictTypeDb) PageList(te *domain.TestExpress, page, size int64) ([]*domain.TestExpress, int64, int64, error) {
	return nil, 0, 0, nil
}

//goland:noinspection SqlResolve
func (pg *TestDictTypeDb) CheckByNumAndPickupCode(expressNumber, pickupCode string) bool {
	var selectCount int64
	err := pg.Db.QueryRow("SELECT COUNT(*) FROM tb_test_express WHERE express_number = $1 AND pickup_code = $2", expressNumber, pickupCode).Scan(&selectCount)
	if err != nil {
		log.Printf("通过快递单号与取件码查询异常===%v", err)
		return false
	}
	if selectCount >= 1 {
		return true
	}
	return false
}

//goland:noinspection SqlResolve
func (pg *TestDictTypeDb) BatchDelete(ids []string) (rows int64, err error) {
	deleteSql := fmt.Sprintf("DELETE FROM tb_test_express WHERE id IN (%s)", strings.Join(ids, ", "))
	log.Println("测试快递批量删除sql===", deleteSql)
	result, err := pg.Db.Exec(deleteSql)
	if err != nil {
		log.Printf("测试快递批量删除异常===%v", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("测试快递批量删除执行失败===%v", err)
		return 0, err
	}

	return rowsAffected, nil
}
