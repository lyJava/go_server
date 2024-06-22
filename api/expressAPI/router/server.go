package router

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/controller"
	"apiProject/api/expressAPI/service"
	"apiProject/api/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
)

type APIServer struct {
	addr         string                          // 服务启动端口
	express      service.ExpressServiceInterface // 快递服务接口
	user         service.UserServiceInterface    // 用户服务接口
	rabbitmqConn *amqp.Connection                // rabbitmq连接
	dict         service.DictTypeService         // 字典服务接口
	testExpress  service.TestExpressService      // 测试快递服务接口
	macCpu       service.MacCpuService           // 苹果处理器接口
	macMemory    service.MacMemoryService        // 苹果内存接口
}

// NewAPIServer 创建API服务
func NewAPIServer(add string, express service.ExpressServiceInterface, user service.UserServiceInterface,
	conn *amqp.Connection, d service.DictTypeService, te service.TestExpressService, cpu service.MacCpuService, memory service.MacMemoryService) *APIServer {
	return &APIServer{
		addr:         add,
		express:      express,
		user:         user,
		rabbitmqConn: conn,
		dict:         d,
		testExpress:  te,
		macCpu:       cpu,
		macMemory:    memory,
	}
}

var queueName = "create_order_queue"
var exchangeName = "create_order_exchange"
var routingKey = "create_order_routing_key"

// Serve 启动API服务
func (s *APIServer) Serve() {
	router := mux.NewRouter()

	requestLog := config.EnvConfig.RequestLog
	if requestLog {
		router.Use(RequestLogMiddleware)
	}

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

	// 订单控制器(测试rabbitmq)
	orderController := controller.OrderControllerInit(s.rabbitmqConn)
	orderController.RegisterRoutes(router)

	// 字典类型控制器
	dictController := controller.DictControllerInit(s.dict)
	dictController.RegisterRoutes(router)

	// 测试快递控制器
	testExpressController := controller.TestExpressControllerInit(s.testExpress)
	testExpressController.RegisterRoutes(router)

	// 大模型控制器
	llmsController := controller.LlmsControllerInit(utils.CreateModel("llama3"))
	llmsController.RegisterRoutes(router)

	// 苹果处理器控制器
	macCpuController := controller.MacCpuControllerInit(s.macCpu)
	macCpuController.RegisterRoutes(router)

	// 苹果内存控制器
	macMemoryController := controller.MacMemoryControllerInit(s.macMemory)
	macMemoryController.RegisterRoutes(router)

	// 测试请求日志
	router.HandleFunc("/user/test/{id}", UserHandler).Methods("GET")

	log.Println("api server starting at port====", s.addr)
	log.Fatalln(http.ListenAndServe(s.addr, router))
}
