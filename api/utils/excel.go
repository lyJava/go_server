package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/xuri/excelize/v2"
)

//goland:noinspection GoUnhandledErrorResult
func CreateExcel() {
	f := excelize.NewFile()

	defer func() {
		// 关闭文件
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 创建工作表，刚创建的excel默认是没有工作表的
	_, err := f.NewSheet("sheet1")

	if err != nil {
		log.Printf("设置工作表名称错误===%v", err)
		return
	}
	// 继续创建工作表，这里macOS预览的时候才会默认激活了sheet1，不需要的可以删除该代码
	_, _ = f.NewSheet("sheet2")

	f.SetCellValue("Sheet1", "A1", "快递名称")
	f.SetCellValue("Sheet1", "B1", "快递单号")
	f.SetCellValue("Sheet1", "C1", "取件码")
	f.SetCellValue("Sheet1", "D1", "发件人姓名")
	f.SetCellValue("Sheet1", "E1", "发件人手机")
	f.SetCellValue("Sheet1", "F1", "发件人地址")
	f.SetCellValue("Sheet1", "G1", "发件人身份证号")
	f.SetCellValue("Sheet1", "H1", "创建人")
	f.SetCellValue("Sheet1", "I1", "创建时间")
	f.SetCellValue("Sheet1", "J1", "修改人")
	f.SetCellValue("Sheet1", "K1", "修改时间")
	// 设置活动的工作表
	f.SetActiveSheet(0)

	currentPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录错误===%v", err)
	}
	log.Println("获取当前工作目录路径:", currentPath)

	// 按指定路径保存文件
	if err := f.SaveAs(currentPath + "/excel/Book3.xlsx"); err != nil {
		fmt.Println(err)
	}
}

// WriteTestExpressToExcel 将测试快递数据写入Excel
//
//goland:noinspection GoUnhandledErrorResult
func WriteTestExpressToExcel[T any](fileName string, headers []string, list []T) string {
	f := excelize.NewFile()

	// 关闭文件
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 写入表头
	for i, header := range headers {
		cell := indexToColumn(i) + "1"
		log.Println("表头", cell)
		f.SetCellValue("Sheet1", cell, header)
		width := len(header) * 2 // 根据表头内容设置宽度，这里简单的以字符数乘以2来设置宽度
		f.SetColWidth("Sheet1", cell, cell, float64(width))
	}

	// 写入数据
	for rowIndex, record := range list {
		rowIndex += 2 // 跳过表头，从第二行开始写入数据
		v := reflect.ValueOf(record)
		if v.Kind() == reflect.Ptr {
			v = v.Elem() // 如果是指针类型，则获取其指向的值
		}
		for colIndex := 0; colIndex < v.NumField(); colIndex++ {
			cell := indexToColumn(colIndex) + strconv.Itoa(rowIndex)
			// 获取字段值并转换为字符串
			var value string
			if v.Type().Field(colIndex).Name == "Id" {
				idPtr := v.Field(colIndex).Interface().(*int64)
				if idPtr != nil {
					value = ConvertInt64ToStr(*idPtr)
				} else {
					value = "" // 如果指针为 nil，则将值设为空字符串
				}
			} else {
				value = fmt.Sprintf("%v", v.Field(colIndex).Interface())
			}

			log.Println("获取的value", value)

			f.SetCellValue("Sheet1", cell, value)

			// 写入数据到单元格
			f.SetCellValue("Sheet1", cell, value)

			// 根据写入内容的长度动态设置单元格的宽度加上一些额外的空间
			colWidth := len(value) + 5
			colName, _ := excelize.ColumnNumberToName(colIndex + 1)
			f.SetColWidth("Sheet1", colName, colName, float64(colWidth))
		}
	}

	currentPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录错误===%v", err)
	}
	log.Println("获取当前工作目录路径:", currentPath)

	filePath := filepath.Join(currentPath, fileName)
	// 保存文件
	if err := f.SaveAs(currentPath + fileName); err != nil {
		fmt.Println(err)
	}
	return filePath
}

// indexToColumn 将索引转换为 Excel 列名
func indexToColumn(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = (index-index%26)/26 - 1
	}
	return result
}
