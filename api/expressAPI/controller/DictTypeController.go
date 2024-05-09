package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/response"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type DictController struct {
	service service.DictTypeService
}

func DictControllerInit(s service.DictTypeService) *DictController {
	return &DictController{
		service: s,
	}
}

func (d *DictController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/dictType/list", d.handlePageList).Methods("POST")
}

func (d *DictController) handlePageList(w http.ResponseWriter, r *http.Request) {
	var searchParam map[string]string

	if err := json.NewDecoder(r.Body).Decode(&searchParam); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("查询失败"))
		return
	}
	list, totalPages, totalRecords, err := d.service.GetDictList(nil, searchParam["page"], searchParam["size"])
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("分页查询失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(response.NewPageData(totalRecords, totalPages, list)))
	return
}
