package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type TestExpressController struct {
	service service.TestExpressService
}

func TestExpressControllerInit(s service.TestExpressService) *TestExpressController {
	return &TestExpressController{
		service: s,
	}
}

func (td *TestExpressController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/testExpress/save", td.handleTestExpressSave).Methods("POST")
	router.HandleFunc("/testExpress/batchSave", td.handleTestExpressBatchSave).Methods("POST")
	router.HandleFunc("/testExpress/batchDelete", td.handleTestExpressBatchDelete).Methods("POST")
	router.HandleFunc("/testExpress/page", td.handleTestExpressPage).Methods("POST")
	router.HandleFunc("/testExpress/{id}", td.handleTestExpressGetById).Methods("GET")
	router.HandleFunc("/testExpress/export/excel", td.handleTestExpressExport).Methods("GET")
}

// handleTestExpressSave 新增
func (td *TestExpressController) handleTestExpressSave(w http.ResponseWriter, r *http.Request) {
	var testExpress *domain.TestExpress
	err := json.NewDecoder(r.Body).Decode(&testExpress)
	if err != nil {
		log.Printf("快递新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递新增参数解析失败"))
		return
	}

	utils.CloseBodyError("测试快递新增失败", w, r)

	marshal, _ := json.MarshalIndent(testExpress, "", "    ")
	log.Printf("测试快递新增数据===\r\n%s", string(marshal))
	result, err := td.service.Save(testExpress)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("新增成功", result))
}

// handleTestExpressBatchSave 批量新增
func (td *TestExpressController) handleTestExpressBatchSave(w http.ResponseWriter, r *http.Request) {
	var testExpressList []*domain.TestExpress
	err := json.NewDecoder(r.Body).Decode(&testExpressList)
	if err != nil {
		log.Printf("快递新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递新增参数解析失败"))
		return
	}

	utils.CloseBodyError("测试快递批量新增失败", w, r)

	if len(testExpressList) <= 0 {
		response.WriteJson(w, response.FailMessageResp("测试快递批量新增不能低少一条数据"))
		return
	}

	result, err := td.service.BatchSave(testExpressList)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("测试快递批量新增成功", result))
}

// handleTestExpressBatchDelete 批量删除
func (td *TestExpressController) handleTestExpressBatchDelete(w http.ResponseWriter, r *http.Request) {
	var ids []any
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Printf("快递新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递批量删除参数解析失败"))
		return
	}

	// 加上defer会在请求结束后关闭
	defer utils.CloseBodyError("测试快递批量删除", w, r)

	if len(ids) == 0 {
		response.WriteJson(w, response.FailMessageResp("测试快递批量删除参数验证失败"))
		return
	}

	result, err := td.service.BatchDelete(ids)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if result == 0 {
		response.WriteJson(w, response.FailMessageResp("没有可执行的数据"))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("测试快递批量删除成功", result))
}

// handleTestExpressPage 分页查询
func (td *TestExpressController) handleTestExpressPage(w http.ResponseWriter, r *http.Request) {
	var queryMap = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&queryMap)
	if err != nil {
		log.Printf("测试快递分页查询参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递快递分页查询参数解析失败"))
		return
	}
	utils.CloseBodyError("测试快递分页查询失败", w, r)

	page, _ := queryMap["page"].(string)
	size, _ := queryMap["size"].(string)
	log.Printf("测试快递分页查询分页参数===page=%v,size=%v", page, size)

	if utils.ConvertToInt64(size) > 500 {
		response.WriteJson(w, response.FailMessageResp("测试快递分页查询单次不能超过500条"))
		return
	}

	var testExpress = &domain.TestExpress{}

	queryObj := queryMap["obj"]
	if queryObj != nil {
		testExpressObj, ok := queryObj.(map[string]interface{})
		if !ok {
			log.Println("测试快递分页查询参数obj类型错误")
			response.WriteJson(w, response.FailMessageResp("测试快递分页查询参数obj类型错误"))
			return
		}

		testExpress = &domain.TestExpress{
			ExpressName:      utils.GetStrFromMap(testExpressObj, "expressName"),
			ExpressNumber:    utils.GetStrFromMap(testExpressObj, "expressNumber"),
			PickupCode:       utils.GetStrFromMap(testExpressObj, "pickupCode"),
			FromUsername:     utils.GetStrFromMap(testExpressObj, "fromUsername"),
			FromUserPhone:    utils.GetStrFromMap(testExpressObj, "fromUserPhone"),
			FromUserAddress:  utils.GetStrFromMap(testExpressObj, "fromUserAddress"),
			FromUserIdNumber: utils.GetStrFromMap(testExpressObj, "fromUserIdNumber"),
			CreateBy:         utils.GetStrFromMap(testExpressObj, "createBy"),
			Remarks:          utils.GetStrFromMap(testExpressObj, "remarks"),
			DelFlag:          utils.GetStrFromMap(testExpressObj, "delFlag"),
		}
	}

	marshal, _ := json.MarshalIndent(testExpress, "", "    ")
	log.Printf("测试快递分页查询对象参数===\r\n%s", string(marshal))

	list, totalRecord, totalPage, err := td.service.PageList(testExpress, utils.ConvertToInt64(page), utils.ConvertToInt64(size))
	if err != nil {
		log.Printf("快递分页查询失败===%v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}
	// 设置表头
	headers := []string{
		"主键ID",
		"快递名称",
		"快递单号",
		"取件码",
		"发件人姓名",
		"发件人手机",
		"发件人地址",
		"发件人身份证号",
		"创建人",
		"创建时间",
		"修改人",
		"修改时间",
		"备注",
		"是否删除",
	}
	excelName := "data_" + time.Now().Format("20060102150405") + ".xlsx"
	filePath := utils.WriteTestExpressToExcel("/excel/"+excelName, headers, list)
	log.Println("生成excel路径===", filePath)

	downloadUrl := "http://localhost:3000/testExpress/export/excel?fileName=" + excelName
	log.Println("Excel下载链接===", downloadUrl)

	dataMap := map[string]interface{}{
		"pageData":    response.NewPageData(totalRecord, totalPage, list),
		"downloadUrl": downloadUrl,
	}
	response.WriteJson(w, response.OkDataResp(dataMap))
}

// handleTestExpressGetById 查询
func (td *TestExpressController) handleTestExpressGetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	detail, err := td.service.SelectById(utils.ConvertToInt64(id))
	log.Println("测试快递详情", detail)
	if err != nil {
		log.Printf("测试快递通过ID查询失败===%v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	// 设置表头
	/* headers := []string{
		"主键ID",
		"快递名称",
		"快递单号",
		"取件码",
		"发件人姓名",
		"发件人手机",
		"发件人地址",
		"发件人身份证号",
		"创建人",
		"创建时间",
		"修改人",
		"修改时间",
		"备注",
		"是否删除",
	}

	var list []*domain.TestExpress
	list = append(list, detail)
	utils.WriteTestExpressToExcel("/excel/data_"+time.Now().Format("20060102150405")+".xlsx", headers, list) */
	response.WriteJson(w, response.OkDataResp(detail))
}

// handleTestExpressExport 导出Excel
func (td *TestExpressController) handleTestExpressExport(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")
	if fileName == "" {
		response.WriteJson(w, response.FailMessageResp("文件名不能为空"))
		return
	}
	currentPath, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录错误===%v", err)
	}
	log.Println("获取当前目录:", currentPath)

	filePath := filepath.Join(currentPath+"/excel", fileName)
	utils.DownloadFile(filePath, w, r)
}
