package controller

import (
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/types/domain"
	"apiProject/api/expressAPI/types/param"
	"apiProject/api/response"
	"apiProject/api/utils"
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
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
	r.HandleFunc("/storeAdmin", e.handlerCrete).Methods("POST")
	// 修改
	r.HandleFunc("/storeAdmin", e.handlerUpdate).Methods("PUT")
	// 查询详情
	r.HandleFunc("/storeAdmin/detail", e.handlerDetail).Methods("GET")
	// 查询
	r.HandleFunc("/storeAdmin/{dataId}", e.handlerGet).Methods("GET")
	// 删除
	r.HandleFunc("/storeAdmin/{dataId}", e.handlerDelete).Methods("DELETE")
	// 多条件查询分页
	r.HandleFunc("/storeAdmin/list", e.handlerSelectPageParam).Methods("POST")
	// 批量删除
	r.HandleFunc("/storeAdmin/batchDelete", e.handlerBatchDelete).Methods("POST")
	// 批量删除
	r.HandleFunc("/storeAdmin/batchSave", e.handlerBatchSave).Methods("POST")
	// pdf水印
	r.HandleFunc("/storeAdmin/pdf/water", e.handlerPdfWater).Methods(utils.GET)
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

func (e *StoreAdminController) handlerUpdate(w http.ResponseWriter, r *http.Request) {
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

	count := int64(0)

	if storeAdmin.Id == 0 {
		response.WriteJson(w, response.FailMessageResp("数据ID不能为空"))
		return
	} else {
		count, err = e.storeAdminService.SelectCountById(storeAdmin.Id)
		if err != nil {
			response.WriteJson(w, response.FailMessageResp("新增失败"))
			return
		}
	}

	if count > 0 {
		t, err := e.storeAdminService.Update(storeAdmin)
		if err != nil {
			log.Printf("新增操作错误===%s+v", err)
			response.WriteJson(w, response.FailMessageResp("新增失败"))
			return
		}
		response.WriteJson(w, response.OkDataResp(t))
	}
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

func (e *StoreAdminController) handlerPdfWater(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handlerPdfWater")
	// 读取一个 PDF 文件
	filePath := "/Users/yangge/GolandProjects/apiProject/upload/example.pdf" // 这里指定你的 PDF 文件路径
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", err)
		return
	}
	defer file.Close()

	// 创建一个 multipart.Reader
	//body := &bytes.Buffer{}
	//_, err = io.Copy(body, file)
	//if err != nil {
	//	fmt.Printf("Failed to copy file content: %s\n", err)
	//	return
	//}
	//multipartReader := multipart.NewReader(body, "boundary")
	//
	//// 获取文件部分 (通常文件位于 "file" 字段)
	//part, err := multipartReader.NextPart()
	//if err == io.EOF {
	//	fmt.Printf("No more parts available===%s\n", err)
	//	return
	//} else if err != nil {
	//	fmt.Printf("Failed to get part from multipart: %s\n", err)
	//	return
	//}

	// 传递文件部分到 AddTextWaterForPdf
	watermarkedPDF, err := utils.AddTextWaterForPdfFile(file, "公开作品，禁止商用")
	if err != nil {
		fmt.Printf("Failed to add watermark: %s\n", err)
		return
	}

	//texts := []string{"机密文件", "严格保密", "仅限内部使用"}
	//multipleWatermarkedPDF, err := utils.AddMultipleWatermarksToPdf(filePath, texts)
	//if err != nil {
	//	fmt.Printf("Failed to add watermark: %s\n", err)
	//	return
	//}

	// 保存带水印的 PDF
	//outputFile := "output_watermarked.pdf"
	//err = os.WriteFile(outputFile, watermarkedPDF, 0644)
	//if err != nil {
	//	fmt.Printf("Failed to write output file: %s\n", err)
	//	return
	//}
	//
	//// 保存带水印的 PDF
	//outputFile2 := "multiple_output_watermarked.pdf"
	//err = os.WriteFile(outputFile2, multipleWatermarkedPDF, 0644)
	//if err != nil {
	//	fmt.Printf("Failed to write output file: %s\n", err)
	//	return
	//}

	// 输出文件路径信息
	//fmt.Printf("PDF with watermark saved to: %s\n", outputFile2)

	// 设置正确的 HTTP 响应头
	w.Header().Set("Content-Type", "application/pdf")
	// 这样会直接下载
	//w.Header().Set("Content-Disposition", "attachment; filename=\"output_watermarked.pdf\"")

	// 不直接下载
	w.Header().Set("Content-Disposition", "inline; filename=\""+uuid.NewString()+".pdf\"")

	// 将带水印的 PDF 写入响应
	_, err = w.Write(watermarkedPDF)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write output file: %s", err), http.StatusInternalServerError)
		return
	}
}
