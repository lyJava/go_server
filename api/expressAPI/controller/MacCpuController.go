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

type MacCpuController struct {
	service service.MacCpuService
}

func MacCpuControllerInit(cpuService service.MacCpuService) *MacCpuController {
	return &MacCpuController{
		service: cpuService,
	}
}

func (cpu *MacCpuController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/mac/cpu/save", cpu.handlerSave).Methods("POST")
	router.HandleFunc("/mac/cpu/{id}", cpu.handlerSelect).Methods("GET")
	router.HandleFunc("/mac/cpu/update", cpu.handlerUpdate).Methods("POST")
	router.HandleFunc("/mac/cpu/batch/save", cpu.handlerBatchSave).Methods("POST")
	router.HandleFunc("/mac/cpu/delete/{id}", cpu.handlerDelete).Methods("GET")
	router.HandleFunc("/mac/cpu/batch/delete", cpu.handleBatchDelete).Methods("POST")
}

func (cpu *MacCpuController) handlerSave(w http.ResponseWriter, r *http.Request) {
	var macCpu *domain.MacCpu
	if err := json.NewDecoder(r.Body).Decode(&macCpu); err != nil {
		zap.L().Sugar().Errorf("苹果处理器新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果处理器新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果处理器新增请求", w, r)

	result, err := cpu.service.Save(macCpu)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("新增成功", result))
}

func (cpu *MacCpuController) handlerSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" || cast.ToInt64(id) == 0 {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := cpu.service.SelectById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器通过ID查询错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("查询成功", result))
}

func (cpu *MacCpuController) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	var macCpu *domain.MacCpu
	if err := json.NewDecoder(r.Body).Decode(&macCpu); err != nil {
		zap.L().Sugar().Errorf("苹果处理器修改参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果处理器修改参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果处理器修改请求", w, r)

	result, err := cpu.service.Update(macCpu)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器修改错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("修改成功", result))
}

func (cpu *MacCpuController) handlerBatchSave(w http.ResponseWriter, r *http.Request) {
	var list []*domain.MacCpu
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果处理器新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果处理器批量新增请求", w, r)

	result, err := cpu.service.BatchSave(list)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("批量新增成功", result))
}

func (cpu *MacCpuController) handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := cpu.service.DeleteById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器通过ID删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}

func (cpu *MacCpuController) handleBatchDelete(w http.ResponseWriter, r *http.Request) {
	ids := []any{nil}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		zap.L().Sugar().Errorf("苹果处理器批量删除参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果处理器删除参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果处理器批量删除请求", w, r)

	result, err := cpu.service.BatchDeleteByIds(ids)
	if err != nil {
		zap.L().Sugar().Errorf("苹果处理器通过ID删除删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}
