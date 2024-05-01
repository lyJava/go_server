package controller

import (
	"apiProject/api/expressAPI/interceptor"
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"github.com/bytedance/sonic"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type ExpressInterfaceTest struct {
	expressInter service.ExpressInterface     //快递服务接口
	userInter    service.UserServiceInterface // 用户服务接口
}

// NewExpressService 创建新的请求快递服务
func NewExpressService(e service.ExpressInterface, u service.UserServiceInterface) *ExpressInterfaceTest {
	return &ExpressInterfaceTest{
		expressInter: e,
		userInter:    u,
	}
}

// RegisterRoutes 注册快递服务请求路由
func (e *ExpressInterfaceTest) RegisterRoutes(r *mux.Router) {
	// 新增
	r.HandleFunc("/express", e.handlerCrete).Methods("POST")
	// 查询详情
	r.HandleFunc("/express/detail", e.handlerDetail).Methods("GET")
	// 分页查询
	r.HandleFunc("/expressPage/list", e.handlerSelectPage).Methods("GET")
	// 查询
	r.HandleFunc("/express/{dataId}", e.handlerGet).Methods("GET")
	// 删除
	r.HandleFunc("/express/{dataId}", e.handlerDelete).Methods("DELETE")
	// 多条件查询分页
	r.HandleFunc("/expressPage/list", e.handlerSelectPageParam).Methods("POST")
	// 修改
	r.HandleFunc("/express/update", e.handlerUpdate).Methods("PUT")
	// 批量新增
	r.HandleFunc("/express/batchDelete", e.handlerBatchDelete).Methods("POST")
	// 批量新增
	r.HandleFunc("/express/batchAdd", interceptor.WithJWTAuthorization(e.handlerBatchInsert, e.userInter)).Methods("POST")
}

func (e *ExpressInterfaceTest) handlerCrete(w http.ResponseWriter, r *http.Request) {
	var express *types.Express

	/*if err := json.NewDecoder(r.Body).Decode(&express); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}*/

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}
	defer r.Body.Close()

	//err = json.Unmarshal(body, &express)
	err = sonic.Unmarshal(body, &express)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("新增失败"))
		return
	}

	t, err := e.expressInter.CreateExpress(express)
	if err != nil {
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

func (e *ExpressInterfaceTest) handlerGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}
	t, err := e.expressInter.GetExpress(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("获取数据失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

func (e *ExpressInterfaceTest) handlerDetail(w http.ResponseWriter, r *http.Request) {
	var dataId = r.URL.Query().Get("dataId")
	log.Printf("获取的参数：%s", dataId)
	if dataId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}
	id := utils.ConvertToInt64(dataId)
	t, err := e.expressInter.GetExpress(id)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("获取详情数据失败"))
		return
	}

	response.WriteJson(w, response.OkDataResp(t))
	return
}

func (e *ExpressInterfaceTest) handlerSelectPage(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	expressName := url.Query().Get("expressName")
	page := url.Query().Get("page")
	size := url.Query().Get("size")

	expressPage, totalRecords, totalPages, err := e.expressInter.SelectExpressPage(expressName, page, size)
	if err != nil {
		return
	}
	/*pageDataMap := make(map[string]interface{})
	pageDataMap["content"] = expressPage
	pageDataMap["totalRecords"] = totalRecords
	pageDataMap["totalPages"] = totalPages
	response.WriteJson(w, response.OkDataResp(pageDataMap))*/
	pageData := response.PageData{
		Content:      expressPage,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
	}

	response.WriteJson(w, response.OkDataResp(pageData))
	return
}

func (e *ExpressInterfaceTest) handlerSelectPageParam(w http.ResponseWriter, r *http.Request) {
	var param types.ExpressSearchParam
	/*if err := json.NewDecoder(r.Body).Decode(&searchParam); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("查询失败"))
		return
	}*/
	// 单个新增数据量比较小使用io.ReadAll
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("分页参数解析失败"))
		return
	}

	defer r.Body.Close()

	if err := sonic.Unmarshal(body, &param); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("查询分页数据失败"))
		return
	}

	expressPage, totalRecords, totalPages, err := e.expressInter.SelectExpressPageByParam(&param)
	if err != nil {
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "查询分页数据失败"))
		return
	}

	pageData := response.PageData{
		Content:      expressPage,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
	}
	response.WriteJson(w, response.OkDataResp(pageData))
	return
}

func (e *ExpressInterfaceTest) handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var queryId = vars["dataId"]
	if queryId == "" {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}
	t, err := e.expressInter.DeleteById(utils.ConvertToInt64(queryId))
	if err != nil {
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "删除失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

func (e *ExpressInterfaceTest) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "获取参数失败"))
		return
	}
	defer r.Body.Close()

	var express *types.Express
	err = json.Unmarshal(body, &express)
	if err != nil {
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "解析参数失败"))
		return
	}

	t, err := e.expressInter.UpdateExpress(express)
	if err != nil {
		response.WriteJson(w, response.FailCodeMessageResp(http.StatusInternalServerError, "更新失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

func (e *ExpressInterfaceTest) handlerBatchDelete(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("批量删除参数获取失败"))
		return
	}
	defer r.Body.Close()

	var ids []string
	err = json.Unmarshal(body, &ids)
	if err != nil {
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
	t, err := e.expressInter.BatchDeleteByIds(ids)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("批量删除失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}

// handlerBatchInsert 批量新增
func (e *ExpressInterfaceTest) handlerBatchInsert(w http.ResponseWriter, r *http.Request) {
	var list []*types.Express
	// 大批量情况下使用json.NewDecoder与Decode
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		log.Println(err)
		response.WriteJson(w, response.FailMessageResp("批量新增参数解析失败"))
		return
	}

	defer r.Body.Close()

	if len(list) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量新增参数不能为空"))
		return
	}

	if len(list) > 1000 {
		response.WriteJson(w, response.FailMessageResp("批量新增单次最多不能超过1000条"))
		return
	}

	t, err := e.expressInter.BatchCreateExpress(list)
	if err != nil {
		response.WriteJson(w, response.FailMessageResp("批量新增失败"))
		return
	}
	response.WriteJson(w, response.OkDataResp(t))
	return
}
