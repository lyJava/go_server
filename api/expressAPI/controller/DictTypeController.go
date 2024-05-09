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
	router.HandleFunc("/dictType/save", d.handleSave).Methods("POST")
	router.HandleFunc("/dictType/{dictTypeId}", d.handleDetail).Methods("GET")
}

// handlePageList 处理分页查询
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

// handleSave 处理保存
func (d *DictController) handleSave(w http.ResponseWriter, r *http.Request) {
	var dictType *domain.DictType
	if err := json.NewDecoder(r.Body).Decode(&dictType); err != nil {
		log.Printf("字典类型新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("新增参数解析失败"))
		return
	}

	saveDictType, err := d.service.SaveDictType(dictType)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkDataResp(saveDictType))
	return
}

// handleDetail 处理查询
func (d *DictController) handleDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["dictTypeId"]
	if idStr == "" {
		response.WriteJson(w, response.FailMessageResp("ID不能为空"))
		return
	}

	detail, err := d.service.SelectDictTypeById(utils.ConvertToInt64(idStr))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("查询字典类型失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(detail))
	return
}
