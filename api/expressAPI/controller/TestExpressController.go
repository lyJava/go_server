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
}

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

func (td *TestExpressController) handleTestExpressBatchDelete(w http.ResponseWriter, r *http.Request) {
	var ids []string
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Printf("快递新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("测试快递批量删除参数解析失败"))
		return
	}

	utils.CloseBodyError("测试快递批量删除失败", w, r)

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
			ExpressName:      getStrFromMap(testExpressObj, "expressName"),
			ExpressNumber:    getStrFromMap(testExpressObj, "expressNumber"),
			PickupCode:       getStrFromMap(testExpressObj, "pickupCode"),
			FromUsername:     getStrFromMap(testExpressObj, "fromUsername"),
			FromUserPhone:    getStrFromMap(testExpressObj, "fromUserPhone"),
			FromUserAddress:  getStrFromMap(testExpressObj, "fromUserAddress"),
			FromUserIdNumber: getStrFromMap(testExpressObj, "fromUserIdNumber"),
			CreateBy:         getStrFromMap(testExpressObj, "createBy"),
			Remarks:          getStrFromMap(testExpressObj, "remarks"),
			DelFlag:          getStrFromMap(testExpressObj, "delFlag"),
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

	response.WriteJson(w, response.OkDataResp(response.NewPageData(totalRecord, totalPage, list)))
}

func (td* TestExpressController) handleTestExpressGetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	detail, err := td.service.SelectById(utils.ConvertToInt64(id))
	if err != nil {
		log.Printf("测试快递通过ID查询失败===%v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkDataResp(detail))
}


func getStrFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}
