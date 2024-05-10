package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
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
}

func (td *TestExpressController) handleTestExpressSave(w http.ResponseWriter, r *http.Request) {
	var testExpress *domain.TestExpress
	err := json.NewDecoder(r.Body).Decode(&testExpress)
	if err != nil {
		log.Printf("快递新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("快递新增参数解析失败"))
		return
	}
	r.Body.Close()

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
