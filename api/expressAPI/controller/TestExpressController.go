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
	return
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
	return
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
	return
}
