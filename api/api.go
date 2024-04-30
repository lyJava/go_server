package main

import (
	"apiProject/api/response"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// APIServer API服务结构
type APIServer struct {
	// 地址
	addr string
}

// NewAPIServer 创建API服务
//
// 参数
//
//	addr (string):地址
//
// 返回
//
//	APIServer
func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run() error {

	router := http.NewServeMux()

	// 这里默认是GET请求，可以这样写router.HandleFunc("GET /users/{userId}"
	router.HandleFunc("/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		log.Printf("user id: %s", userID)
		//goland:noinspection GoUnhandledErrorResult
		//w.Write([]byte("user id: " + userID))
		//json.NewEncoder(w).Encode(response.OkFullResp(200, "请求成功", userID))
		response.WriteJson(w, response.OkFullResp(200, "请求成功", userID))
	})

	middleWareChain := MiddlewareChain(RequestLogMiddleware, RequestAuthMiddleware)

	childRouter := http.NewServeMux()
	childRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	server := http.Server{
		Addr:    s.addr,
		Handler: middleWareChain(childRouter),
	}
	log.Printf("server has started %s", s.addr)
	return server.ListenAndServe()
}

// RequestLogMiddleware 请求日志中间件
func RequestLogMiddleware(next http.Handler) http.HandlerFunc {
	return func(resp http.ResponseWriter, request *http.Request) {
		log.Printf("请求方法：%s，请求URL：%s", request.Method, request.URL.Path)
		// 调用下一个中间件或处理器
		next.ServeHTTP(resp, request)
	}
}

// RequestAuthMiddleware 验证权限中间件
//
// 参数
//
//	next http.Handler: 下一个处理器
//
// 返回
//
//	http.HandlerFunc: 请求处理器函数
func RequestAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(resp http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		// pathParts := strings.Split(path, "/")
		// log.Println(pathParts, pathParts[len(pathParts)-1])
		userID := path[strings.LastIndex(path, "/")+1:]
		token := request.Header.Get("Authorization")
		log.Printf("请求头用户ID: %s， 请求头token: %s", userID, token)
		if token == "" {
			//http.Error(resp, "请求头凭证为空", http.StatusUnauthorized)
			json.NewEncoder(resp).Encode(response.FailCodeMessageResp(http.StatusUnauthorized, "用户【"+userID+"】请求头凭证为空"))
			return
		}

		if token != "Bearer token" {
			log.Printf("无权限")
			json.NewEncoder(resp).Encode(response.FailCodeMessageResp(http.StatusPaymentRequired, "用户【"+userID+"】权限验证失败"))
			return
		}
		next.ServeHTTP(resp, request)
	}
}

// ContentTypeMiddleware 设置中间件的Content-Type
//
// 参数
//
//	next http.Handler: 下一个中间件或处理器
func ContentTypeMiddleware(next http.Handler) http.HandlerFunc {
	return func(resp http.ResponseWriter, request *http.Request) {
		// 设置Content-Type为application/json，这样其他的地方不需要重复设置content-type
		resp.Header().Set("Content-Type", "application/json;charset=utf-8")
		// 调用下一个中间件或处理器
		next.ServeHTTP(resp, request)
	}
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return ContentTypeMiddleware(next)
	}
}
