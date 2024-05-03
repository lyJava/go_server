package router

import (
	"apiProject/api/expressAPI/controller"
	"apiProject/api/expressAPI/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr    string                          // 服务启动端口
	express service.ExpressServiceInterface // 快递接口
	user    service.UserServiceInterface    // 用户服务接口
}

func NewAPIServer(add string, express service.ExpressServiceInterface, user service.UserServiceInterface) *APIServer {
	return &APIServer{
		addr:    add,
		express: express,
		user:    user,
	}
}

func (s *APIServer) Serve() {
	router := mux.NewRouter()
	//childRouter := router.PathPrefix("/dev-api").Subrouter()
	//childRouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	// 快递管理控制器
	expressController := controller.ExpressControllerInit(s.express, s.user)
	expressController.RegisterRoutes(router)

	// 用户管理控制器
	userController := controller.UserControllerInit(s.user)
	userController.RegisterRoutes(router)

	// 验证码控制器
	captchaController := controller.CaptchaControllerInit()
	captchaController.RegisterRoutes(router)

	// 文件上传控制器
	uploadController := controller.UploadControllerInit()
	uploadController.RegisterRoutes(router)

	// 文件下载控制器
	downloadController := controller.DownloadControllerInit()
	downloadController.RegisterRoutes(router)

	log.Println("api server starting at====", s.addr)
	log.Fatalln(http.ListenAndServe(s.addr, router))
}
