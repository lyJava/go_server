package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"go.uber.org/zap"
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
	router.HandleFunc("/testExpress/save", td.handlerTestExpressSave).Methods("POST")
	router.HandleFunc("/testExpress/batchSave", td.handlerTestExpressBatchSave).Methods("POST")
	router.HandleFunc("/testExpress/batchDelete", td.handlerTestExpressBatchDelete).Methods("POST")
	router.HandleFunc("/testExpress/page", td.handlerTestExpressPage).Methods("POST")
	router.HandleFunc("/testExpress/{id}", td.handlerTestExpressGetById).Methods("GET")
	router.HandleFunc("/testExpress/export/excel", td.handlerTestExpressExport).Methods("GET")
}

// handlerTestExpressSave 新增
func (td *TestExpressController) handlerTestExpressSave(w http.ResponseWriter, r *http.Request) {
	var testExpress *domain.TestExpress
	err := json.NewDecoder(r.Body).Decode(&testExpress)
	if err != nil {
		log.Printf("测试快递新增参数解析错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试快递新增请求", w, r)

	marshal := utils.ToJsonFormat(testExpress)
	log.Printf("测试快递新增数据===\r\n%s", marshal)
	result, err := td.service.Save(testExpress)
	if err != nil {
		zap.L().Sugar().Errorf("测试快递新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("新增成功", result))
}

// handlerTestExpressBatchSave 批量新增
func (td *TestExpressController) handlerTestExpressBatchSave(w http.ResponseWriter, r *http.Request) {
	var testExpressList []*domain.TestExpress
	err := json.NewDecoder(r.Body).Decode(&testExpressList)
	if err != nil {
		log.Printf("测试快递新增参数解析错误===%v", err)
		zap.L().Sugar().Errorf("测试快递批量新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("测试快递批量新增请求", w, r)

	if len(testExpressList) <= 0 {
		zap.L().Sugar().Infof("测试快递批量新增数量:%d不能低少一条数据", len(testExpressList))
		response.WriteJson(w, response.FailMessageResp("测试快递批量新增不能低少一条数据"))
		return
	}

	result, err := td.service.BatchSave(testExpressList)
	if err != nil {
		zap.L().Sugar().Errorf("快递新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("测试快递批量新增成功", result))
}

// handlerTestExpressBatchDelete 批量删除
func (td *TestExpressController) handlerTestExpressBatchDelete(w http.ResponseWriter, r *http.Request) {
	var ids []any
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Printf("测试快递批量删除参数解析失败===%v", err)
		zap.L().Sugar().Errorf("快递新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递批量删除参数解析失败"))
		return
	}

	// 加上defer会在请求结束后关闭
	defer utils.CloseBodyError("测试快递批量删除请求", w, r)

	if len(ids) == 0 {
		zap.L().Sugar().Errorf("测试快递批量删除参数验证失败===%+v", ids)
		response.WriteJson(w, response.FailMessageResp("测试快递批量删除参数验证失败"))
		return
	}

	result, err := td.service.BatchDelete(ids)
	if err != nil {
		zap.L().Sugar().Errorf("测试快递批量删除错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	if result == 0 {
		response.WriteJson(w, response.FailMessageResp("没有可执行的数据"))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("测试快递批量删除成功", result))
}

// handlerTestExpressPage 分页查询
func (td *TestExpressController) handlerTestExpressPage(w http.ResponseWriter, r *http.Request) {
	var queryMap = make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&queryMap)
	if err != nil {
		log.Printf("测试快递分页查询参数解析错误===%+v", err)
		zap.L().Sugar().Errorf("测试快递分页查询参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递快递分页查询参数解析失败"))
		return
	}
	utils.CloseBodyError("测试快递分页查询失败", w, r)

	page := queryMap["page"].(string)
	if page == "" {
		zap.L().Sugar().Errorf("测试快递快递分页参数page验证失败===%s", page)
		response.WriteJson(w, response.FailMessageResp("测试快递快递分页参数page不能为空"))
		return
	}
	size := queryMap["size"].(string)
	if size == "" {
		zap.L().Sugar().Errorf("测试快递快递分页参数size验证失败===%s", page)
		response.WriteJson(w, response.FailMessageResp("测试快递快递分页参数size不能为空"))
		return
	}
	log.Printf("测试快递分页查询分页参数===page=%s,size=%s", page, size)
	zap.L().Sugar().Debugf("测试快递分页查询分页参数===page=%s,size=%s", page, size)

	if utils.ConvertToInt64(size) > 500 {
		zap.L().Sugar().Errorf("测试快递快递分页单次查询超过500条了,当前条数:%s", size)
		response.WriteJson(w, response.FailMessageResp("测试快递分页查询单次不能超过500条"))
		return
	}

	excel := queryMap["excel"].(bool)
	zap.L().Sugar().Infof("测试快递分页查询是否导出excel:%t", excel)

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

	marshal := utils.ToJsonFormat(testExpress)
	log.Printf("测试快递分页查询对象参数===\r\n%s", marshal)

	list, totalRecord, totalPage, err := td.service.PageList(testExpress, utils.ConvertToInt64(page), utils.ConvertToInt64(size))
	if err != nil {
		log.Printf("快递分页查询失败===%+v", err)
		zap.L().Sugar().Errorf("测试快递分页查询是否导出excel:%+v", err)
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

	downloadUrl := ""

	if excel {
		excelName := "data_" + time.Now().Format("20060102150405") + ".xlsx"
		filePath := utils.WriteTestExpressToExcel("/excel/"+excelName, headers, list)
		log.Println("生成excel路径===", filePath)

		downloadUrl = "http://localhost:3000/testExpress/export/excel?fileName=" + excelName
		log.Println("Excel下载链接===", downloadUrl)
		zap.L().Sugar().Infof("生成Excel路径:%s,下载链接:%s", filePath, downloadUrl)
	}

	dataMap := map[string]interface{}{
		"pageData":    response.NewPageData(totalRecord, totalPage, list),
		"downloadUrl": downloadUrl,
	}
	response.WriteJson(w, response.OkDataResp(dataMap))
}

// handlerTestExpressGetById 查询
func (td *TestExpressController) handlerTestExpressGetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	detail, err := td.service.SelectById(cast.ToInt64(id))
	log.Println("测试快递详情", detail)
	if err != nil {
		log.Printf("测试快递通过ID查询失败===%+v", err)
		zap.L().Sugar().Errorf("测试快递通过ID查询错误:%+v", err)
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

// handlerTestExpressExport 导出Excel
func (td *TestExpressController) handlerTestExpressExport(w http.ResponseWriter, r *http.Request) {
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
