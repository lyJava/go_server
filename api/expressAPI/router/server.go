package router

import (
	"apiProject/api/expressAPI/controller"
	"apiProject/api/expressAPI/service"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	addr    string                   // 地址
	express service.ExpressInterface //快递接口
	user    service.UserServiceInterface
}

func NewAPIServer(add string, express service.ExpressInterface, user service.UserServiceInterface) *APIServer {
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

	expressTest := controller.NewExpressService(s.express, s.user)
	expressTest.RegisterRoutes(router)

	userTest := controller.NewUserServiceInterfaceTest(s.user)
	userTest.RegisterRoutes(router)

	log.Println("api server starting at====", s.addr)
	log.Fatalln(http.ListenAndServe(s.addr, router))
}
