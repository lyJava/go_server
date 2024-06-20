package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
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
}

func (cpu *MacCpuController) handlerSave(w http.ResponseWriter, r *http.Request) {
	var macCpu *domain.MacCpu
	if err := json.NewDecoder(r.Body).Decode(&macCpu); err != nil {
		log.Printf("苹果处理器新增参数解析错误===%+v", err)
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
