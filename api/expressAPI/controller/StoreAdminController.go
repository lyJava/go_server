package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

// StoreAdminController 店铺管理员控制器
type StoreAdminController struct {
	storeAdminService service.StoreAdminServiceInterface // 用户服务接口
}

// StoreAdminControllerInit 店铺管理员控制器初始化
func StoreAdminControllerInit(s service.StoreAdminServiceInterface) *StoreAdminController {
	return &StoreAdminController{
		storeAdminService: s,
	}
}

// RegisterRoutes 注册快递服务请求路由
func (e *StoreAdminController) RegisterRoutes(r *mux.Router) {
	// 新增
	r.HandleFunc("/dev-api/storeAdmin", e.handlerCrete).Methods("POST")
	// 查询详情
	r.HandleFunc("/dev-api/storeAdmin/detail", e.handlerDetail).Methods("GET")
	// 查询
	r.HandleFunc("/dev-api/storeAdmin/{dataId}", e.handlerGet).Methods("GET")
	// 删除
	r.HandleFunc("/dev-api/storeAdmin/{dataId}", e.handlerDelete).Methods("DELETE")
	// 多条件查询分页
	r.HandleFunc("/dev-api/storeAdmin/list", e.handlerSelectPageParam).Methods("POST")
	// 批量删除
	r.HandleFunc("/dev-api/storeAdmin/batchDelete", e.handlerBatchDelete).Methods("POST")
	// 批量删除
	r.HandleFunc("/dev-api/storeAdmin/batchSave", e.handlerBatchSave).Methods("POST")
}

func (e *StoreAdminController) handlerCrete(w http.ResponseWriter, r *http.Request) {
	var storeAdmin *domain.StoreAdmin

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("新增参数读取错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}
	defer r.Body.Close()

	err = sonic.Unmarshal(body, &storeAdmin)
	if err != nil {
		log.Printf("新增参数解析错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}

	t, err := e.storeAdminService.Save(storeAdmin)
	if err != nil {
		log.Printf("新增操作错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

func (e *StoreAdminController) handlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}
	t, err := e.storeAdminService.SelectById(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("获取数据失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

func (e *StoreAdminController) handlerDetail(w http.ResponseWriter, r *http.Request) {
	var dataId = r.URL.Query().Get("dataId")
	log.Printf("获取的参数：%s", dataId)
	if dataId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}
	id := utils.ConvertToInt64(dataId)
	t, err := e.storeAdminService.SelectById(id)
	if err != nil {
		log.Printf("获取详情数据错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("获取详情数据失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(t))
}

func (e *StoreAdminController) handlerSelectPageParam(w http.ResponseWriter, r *http.Request) {
	var searchParam param.StoreAdminSearchParam
	// 单个新增数据量比较小使用io.ReadAll
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("获取分页错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("分页参数解析失败"))
		return
	}

	defer utils.CloseBodyError("快递分页请求", w, r)

	if err := sonic.Unmarshal(body, &searchParam); err != nil {
		log.Printf("sonic.Unmarshal错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("获取分页数据失败"))
		return
	}

	column, order := utils.HandlerColumnOrder(searchParam.Column, searchParam.Order)
	searchParam.Column = column
	searchParam.Order = order

	log.Println("快递分页查询参数===", searchParam)

	expressPage, totalRecords, totalPages, err := e.storeAdminService.PageList(&searchParam)
	if err != nil {
		log.Printf("分页查询错误===%s+v", err)
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "查询分页数据失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(response.NewPageData(totalRecords, totalPages, expressPage)))
}

func (e *StoreAdminController) handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	ids := []any{queryId}
	t, err := e.storeAdminService.BatchDelete(ids)
	if err != nil {
		log.Printf("删除操作错误===%s+v", err)
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "删除失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

func (e *StoreAdminController) handlerBatchDelete(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("批量删除参数读取错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("批量删除参数获取失败"))
		return
	}
	defer utils.CloseBodyError("快递批量删除请求", w, r)

	var ids []any
	err = json.Unmarshal(body, &ids)
	if err != nil {
		log.Printf("批量删除参数解析错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("批量删除参数解析失败"))
		return
	}

	if len(ids) == 0 {
		response.WriteJson(w, response.FailMessageResp("请至少选择一条数据"))
		return
	}

	if len(ids) > 20 {
		response.WriteJson(w, response.FailMessageResp("批量删除单次操作不能超过20条"))
		return
	}
	t, err := e.storeAdminService.BatchDelete(ids)
	if err != nil {
		log.Printf("批量删除参数操作错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("批量删除失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}

// handlerBatchSave 批量新增
func (e *StoreAdminController) handlerBatchSave(w http.ResponseWriter, r *http.Request) {
	var list []*domain.StoreAdmin
	// 大批量情况下使用json.NewDecoder与Decode
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Printf("批量新增参数解析错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("批量新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("快递批量新增请求", w, r)

	if len(list) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量新增参数不能为空"))
		return
	}

	if len(list) > 1000 {
		response.WriteJson(w, response.FailMessageResp("批量新增单次最多不能超过1000条"))
		return
	}

	t, err := e.storeAdminService.BatchSave(list)
	if err != nil {
		log.Printf("批量新增操作错误===%s+v", err)
		response.WriteJson(w, response.FailMessageResp("批量新增失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
}
