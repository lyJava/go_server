package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"net/http"
)

type MacMemoryController struct {
	service service.MacMemoryService
}

func MacMemoryControllerInit(memoryService service.MacMemoryService) *MacMemoryController {
	return &MacMemoryController{
		service: memoryService,
	}
}

func (memory *MacMemoryController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/mac/memory/save", memory.handlerSave).Methods("POST")
	router.HandleFunc("/mac/memory/{id}", memory.handlerSelect).Methods("GET")
	router.HandleFunc("/mac/memory/update", memory.handlerUpdate).Methods("POST")
	router.HandleFunc("/mac/memory/batch/save", memory.handlerBatchSave).Methods("POST")
	router.HandleFunc("/mac/memory/delete/{id}", memory.handlerDelete).Methods("GET")
	router.HandleFunc("/mac/memory/batch/delete", memory.handleBatchDelete).Methods("POST")
}

func (memory *MacMemoryController) handlerSave(w http.ResponseWriter, r *http.Request) {
	var macMemory *domain.MacMemory
	if err := json.NewDecoder(r.Body).Decode(&macMemory); err != nil {
		zap.L().Sugar().Errorf("苹果内存新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果内存新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果内存新增请求", w, r)

	result, err := memory.service.Save(macMemory)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("新增成功", result))
}

func (memory *MacMemoryController) handlerSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" || cast.ToInt64(id) == 0 {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := memory.service.SelectById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存通过ID查询错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("查询成功", result))
}

func (memory *MacMemoryController) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	var macMemory *domain.MacMemory
	if err := json.NewDecoder(r.Body).Decode(&macMemory); err != nil {
		zap.L().Sugar().Errorf("苹果内存修改参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果内存修改参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果内存修改请求", w, r)

	result, err := memory.service.Update(macMemory)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存修改错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("修改成功", result))
}

func (memory *MacMemoryController) handlerBatchSave(w http.ResponseWriter, r *http.Request) {
	var list []*domain.MacMemory
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		zap.L().Sugar().Errorf("苹果内存批量新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果内存批量新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果内存批量新增请求", w, r)

	if len(list) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量新增参数不能为空"))
		return
	}
	if len(list) > 10 {
		response.WriteJson(w, response.FailMessageResp("批量新增单次操作不能超过10条"))
		return
	}

	result, err := memory.service.BatchSave(list)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存批量新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("批量新增成功", result))
}

func (memory *MacMemoryController) handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" || cast.ToInt64(id) == 0 {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := memory.service.DeleteById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存通过ID删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}

func (memory *MacMemoryController) handleBatchDelete(w http.ResponseWriter, r *http.Request) {
	ids := []any{nil}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		zap.L().Sugar().Errorf("苹果内存批量删除参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果内存删除参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果内存批量删除请求", w, r)

	if len(ids) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量删除参数不能为空"))
		return
	}
	if len(ids) > 100 {
		response.WriteJson(w, response.FailMessageResp("批量删除单次操作不能超过100条"))
		return
	}

	result, err := memory.service.BatchDeleteByIds(ids)
	if err != nil {
		zap.L().Sugar().Errorf("苹果内存通过ID删除删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}
