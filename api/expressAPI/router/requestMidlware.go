package router

import (
	"apiProject/api/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// RequestLogMiddleware 记录请求和响应相关日志
func RequestLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// 读取请求体
		var requestBodyBytes []byte
		if r.Body != nil {
			requestBodyBytes, _ = io.ReadAll(r.Body)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))

		// 创建一个记录响应数据的 ResponseWriter
		rec := &responseRecorder{ResponseWriter: w, body: &bytes.Buffer{}}

		defer r.Body.Close()

		// 调用下一个处理程序
		next.ServeHTTP(rec, r)

		// 计算请求持续时间
		duration := time.Since(startTime)

		zap.L().Sugar().Infof("请求响应描述==%s,请求响应码===%d", rec.status, rec.statusCode)
		queryParam := r.URL.Query().Encode()
		queryParam = queryParam + "\njson格式:\n" + utils.QueryParamToJson(queryParam)

		go logRequest(r, queryParam, requestBodyBytes, utils.ForcedToJsonFormat(rec.body.Bytes()), time.Duration(duration))
	})
}

// logRequest 异步记录请求日志
func logRequest(r *http.Request, queryParam string, requestBody []byte, responseBody string, duration time.Duration) {
	headerParams := utils.ToJsonFormat(r.Header)
	zap.L().Sugar().Infof("\n请求地址: %s,\n请求方法: %s,\n请求头参数: %s,\n请求URL参数: %s,\n请求URL路径参数: %s,\n请求body: %s,\n返回body: %s,\n请求耗时: %s",
		r.URL.String(),
		r.Method,
		headerParams,
		queryParam,
		formatMap(mux.Vars(r)),
		string(requestBody),
		responseBody,
		duration,
	)
}

// responseRecorder 是一个自定义的http.ResponseWriter用于捕获响应数据
type responseRecorder struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	status     string
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func (rec *responseRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.status = http.StatusText(statusCode)
}

// UserHandler 是一个处理程序，用于处理 RESTFUL 路径参数
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	response := map[string]string{"message": fmt.Sprintf("User ID: %s", userID)}
	json.NewEncoder(w).Encode(response)
}

func formatMap(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("{")
	for k, v := range m {
		sb.WriteString(fmt.Sprintf("%s: %s, ", k, v))
	}
	result := sb.String()
	if len(result) > 1 {
		result = result[:len(result)-2]
	}

	result += "}"
	return result
}
