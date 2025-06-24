package main

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/controller"
	"apiProject/api/expressAPI/router"
	"github.com/google/wire"
)

//	func InitializeConfig() (*types.MysqlConfig, error) {
//		wire.Build(config.InitConfig)
//		return &types.MysqlConfig{}, nil
//	}
func InitializeApplication() (*Application, error) {
	wire.Build(
		router.APIServer{},
		config.InitConfig(),
		controller.ExpressController{},
		controller.UserController{},
		NewApplication,
	)
	return &Application{}, nil
}
