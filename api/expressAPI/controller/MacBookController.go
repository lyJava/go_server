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


type MacBookController struct {
	service service.MacBookService
}

func MacBookControllerInit(bookService service.MacBookService) *MacBookController {
	return &MacBookController{
		service: bookService,
	}
}

func (macBook *MacBookController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/mac/book/save", macBook.handlerSave).Methods("POST")
	router.HandleFunc("/mac/book/{id}", macBook.handlerSelect).Methods("GET")
	router.HandleFunc("/mac/book/update", macBook.handlerUpdate).Methods("POST")
	router.HandleFunc("/mac/book/batch/save", macBook.handlerBatchSave).Methods("POST")
	router.HandleFunc("/mac/book/delete/{id}", macBook.handlerDelete).Methods("GET")
	router.HandleFunc("/mac/book/batch/delete", macBook.handleBatchDelete).Methods("POST")
}

func (macBook *MacBookController) handlerSave(w http.ResponseWriter, r *http.Request) {
	var book *domain.MacBook
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果笔记本新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果笔记本新增请求", w, r)

	result, err := macBook.service.Save(book)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("新增成功", result))
}

func (macBook *MacBookController) handlerSelect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" || cast.ToInt64(id) == 0 {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := macBook.service.SelectById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本通过ID查询错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("查询成功", result))
}

func (macBook *MacBookController) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	var book *domain.MacBook
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本修改参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果笔记本修改参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果笔记本修改请求", w, r)

	result, err := macBook.service.Update(book)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本修改错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("修改成功", result))
}

func (macBook *MacBookController) handlerBatchSave(w http.ResponseWriter, r *http.Request) {
	var list []*domain.MacBook
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量新增参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果笔记本批量新增参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果笔记本批量新增请求", w, r)

	if len(list) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量新增参数不能为空"))
		return
	}
	if len(list) > 10 {
		response.WriteJson(w, response.FailMessageResp("批量新增单次操作不能超过10条"))
		return
	}

	result, err := macBook.service.BatchSave(list)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量新增错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("批量新增成功", result))
}

func (macBook *MacBookController) handlerDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id = vars["id"]
	if id == "" || cast.ToInt64(id) == 0 {
		response.WriteJson(w, response.FailMessageResp("ID参数不能为空"))
		return
	}

	result, err := macBook.service.DeleteById(cast.ToInt64(id))
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本通过ID删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}

func (macBook *MacBookController) handleBatchDelete(w http.ResponseWriter, r *http.Request) {
	ids := []any{nil}
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		zap.L().Sugar().Errorf("苹果笔记本批量删除参数解析错误===%+v", err)
		response.WriteJson(w, response.FailMessageResp("苹果笔记本删除参数解析失败"))
		return
	}

	defer utils.CloseBodyError("苹果笔记本批量删除请求", w, r)

	if len(ids) == 0 {
		response.WriteJson(w, response.FailMessageResp("批量删除参数不能为空"))
		return
	}
	if len(ids) > 100 {
		response.WriteJson(w, response.FailMessageResp("批量删除单次操作不能超过100条"))
		return
	}

	result, err := macBook.service.BatchDeleteByIds(ids)
	if err != nil {
		zap.L().Sugar().Errorf("苹果笔记本通过ID删除删除错误:%+v", err)
		response.WriteJson(w, response.FailMessageResp(err.Error()))
		return
	}

	response.WriteJson(w, response.OkCodeMessageData("删除成功", result))
}
