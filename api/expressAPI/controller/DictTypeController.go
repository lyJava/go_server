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
	router.HandleFunc("/dictType/update", d.handleUpdate).Methods("PUT")
	router.HandleFunc("/dictType/{id}", d.handleDelete).Methods("DELETE")
}

// handlePageList 处理分页查询
func (d *DictController) handlePageList(w http.ResponseWriter, r *http.Request) {
	var searchParam map[string]string
	if err := json.NewDecoder(r.Body).Decode(&searchParam); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("查询失败"))
		return
	}

	utils.CloseBodyError("字典类型查询失败", w, r)

	list, totalPages, totalRecords, err := d.service.GetDictList(nil, searchParam["page"], searchParam["size"])
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("分页查询失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(response.NewPageData(totalRecords, totalPages, list)))
}

// handleSave 处理保存
func (d *DictController) handleSave(w http.ResponseWriter, r *http.Request) {
	var dictType *domain.DictType
	if err := json.NewDecoder(r.Body).Decode(&dictType); err != nil {
		log.Printf("字典类型新增参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("新增参数解析失败"))
		return
	}

	utils.CloseBodyError("字典类型新增失败", w, r)

	saveDictType, err := d.service.SaveDictType(dictType)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkDataResp(saveDictType))
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
}

func (d *DictController) handleUpdate(w http.ResponseWriter, r *http.Request) {
	var dictType *domain.DictType
	if err := json.NewDecoder(r.Body).Decode(&dictType); err != nil {
		log.Printf("字典类型修改参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("修改参数解析失败"))
		return
	}

	utils.CloseBodyError("字典类修改失败", w, r)

	// 如果ID类型设置为*Int64这种就可以使用nil判断是否为空
	if dictType.Id == nil {
		response.WriteJson(w, response.FailMessageResp("ID不能为空"))
		return
	}

	// 如果类型类型设置为*string这种就可以使用nil判断是否为空
	if dictType.DictType == "" {
		response.WriteJson(w, response.FailMessageResp("类型不能为空"))
		return
	}

	isExist := d.service.CheckTypeIsExist(dictType.Id, dictType.DictType)
	if isExist {
		response.WriteJson(w, response.FailMessageResp("类型已存在"))
		return
	}

	result, err := d.service.UpdateDictType(dictType)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("修改成功", result))
}

// handleDelete 处理查询
func (d *DictController) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		response.WriteJson(w, response.FailMessageResp("ID不能为空"))
		return
	}

	id := utils.ConvertToInt64(idStr)
	if id <= 10 {
		response.WriteJson(w, response.FailMessageResp("系统字典不允许删除"))
		return
	}

	result := d.service.DeleteDictType(id)
	if !result {
		response.WriteJson(w, response.FailMessageResp("删除失败"))
		return
	}
	
	response.WriteJson(w, response.OkMessageResp("删除成功"))
}
