package controller

import (
	"apiProject/api/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type SseController struct {
}

func SseControllerInit() *SseController {
	return &SseController{}
}

// RegisterRoutes 注册SSE控制器请求路由
func (sse *SseController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/dev-api/sse", sse.handlerSSE).Methods(utils.GET, utils.POST)
}

func (sse *SseController) handlerSSE(w http.ResponseWriter, r *http.Request) {
	// 记录请求头
	fmt.Println("Request Headers:", utils.ToJsonFormat(r.Header))
	// 设置 CORS 请求头, 允许所有域访问
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// 设置响应头，告知浏览器这是一个 SSE 流
	w.Header().Set("Content-Type", "text/event-stream;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// 打印日志，确认请求到达
	fmt.Println("Client connected")

	// 使用 goroutine 异步推送数据，避免阻塞主请求
	sse.pushSSEData(w)
}

func (sse *SseController) pushSSEData(w http.ResponseWriter) {
	// 保持连接，不断推送数据
	//for {
	// 模拟数据
	data := fmt.Sprintf("data: 当前时间：%s\n\n", time.Now().Format(time.DateTime))

	// 向客户端写入数据
	_, err := w.Write([]byte(data))
	if err != nil {
		// 如果发生错误，说明客户端可能已断开连接，退出循环
		fmt.Println("Error writing SSE data:", err)
		return
	}

	// 每隔 5 秒推送一次数据
	//time.Sleep(5 * time.Second)
	//}
}
