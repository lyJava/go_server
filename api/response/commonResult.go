package response

import (
	"encoding/json"
	"net/http"
)

// PageData 分页数据结构体
type PageData struct {
	TotalRecords int64       `json:"totalRecords"` // 总条数
	TotalPages   int64       `json:"totalPages"`   // 总页数
	Content      interface{} `json:"content"`      // 数据集合
}

// NewPageData 创建分页数据结构体
func NewPageData(totalRecords, totalPages int64, data any) PageData {
	return PageData{
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		Content:      data,
	}
}

// Response 返回结果
type Response struct {
	Code    int    `json:"code"`              // 响应码
	Message string `json:"message,omitempty"` // 相应信息
	Data    any    `json:"data,omitempty"`    // 响应数据（json:"data,omitempty" 表示为空忽略）
}

// OkCodeResp 返回响应码的成功信息
//
// 参数:
//
//	code (int): 响应码，表示失败的原因。
//
// 返回:
//
//	Response: 包含指定响应码和消息的 Response 结构
func OkCodeResp(code int) Response {
	return Response{
		Code:    code,
		Message: "请求成功",
	}
}

// OkDataResp 返回响应数据的成功信息
//
// 参数
//
//	data (any): 响应数据
//
// 返回
//
//	Response: 包含指定响应码、消与响应数据的的结构
func OkDataResp(data any) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func OkCodeMessageResp(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func OkMessageResp(message string) Response {
	return Response{
		Code:    200,
		Message: message,
	}
}

func OkCodeMessageData(message string, data any) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// OkFullResp 返回完整的成功响应
//
// 参数:
//
//	code (int): 响应码
//	message (string): 描述响应消息
//	data interface{}: 响应数据体
//
// 返回:
//
//	Response: 包含指定响应码和消息的 Response 结构
func OkFullResp(code int, message string, data any) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func FailCodeResp(code int) Response {
	return Response{
		Code:    code,
		Message: "请求失败",
	}
}

// FailCodeMessageResp 返回响应码和消息
//
// 参数:
//
//	code (int): 响应码，表示失败的原因。
//	message (string): 描述响应的消息。
//
// 返回:
//
//	Response: 包含指定响应码和消息的 Response 结构
func FailCodeMessageResp(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// FailMessageResp 返回响应码和消息
//
// 参数:
//
//	message (string): 描述响应的消息。
//
// 返回:
//
//	Response: 包含指定响应码和消息的 Response 结构
func FailMessageResp(message string) Response {
	return Response{
		Code:    500,
		Message: message,
	}
}

func FailFullResp(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func ExceptionResp(code int, message string, data interface{}) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// WriteJson 将数据转换json写入到http.ResponseWriter
func WriteJson(w http.ResponseWriter, data Response) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(data.Code)
	json.NewEncoder(w).Encode(data)
}
