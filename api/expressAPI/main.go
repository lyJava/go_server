package main

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/controller"
	"apiProject/api/expressAPI/datasource"
	"apiProject/api/expressAPI/rabbitmq"
	"apiProject/api/expressAPI/router"
	"apiProject/api/expressAPI/service"
	"apiProject/api/expressAPI/service/impl"
	"apiProject/api/expressAPI/types"
	"apiProject/api/utils"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
)

type Application struct {
	port              string                        //端口号
	cfg               *types.MysqlConfig            //配置
	expressService    *service.UserServiceInterface //快递接口
	userService       *service.UserServiceInterface //用户接口
	expressController *controller.ExpressController //快递控制器
	UserController    *controller.UserController    // 用户控制器
}

func NewApplication(
	port string,
	cfg *types.MysqlConfig,
	expressService *service.UserServiceInterface,
	userService *service.UserServiceInterface,
	expressController *controller.ExpressController,
	serController *controller.UserController) *Application {
	return &Application{
		port:              port,
		cfg:               cfg,
		expressService:    expressService,
		userService:       userService,
		expressController: expressController,
		UserController:    serController,
	}
}

func main() {

	//publicKeyName := "public.pem"
	//privateKeyName := "private.pem"
	//bits := 2048
	//err := utils.GenerateRSAKey2File(bits, publicKeyName, privateKeyName)
	//if err != nil {
	//	panic(err)
	//}

	/*publicKeyName := "public.pem"
	privateKeyName := "private.pem"

	publicKey, err := os.ReadFile(publicKeyName)
	if err != nil {
		panic(err)
	}
	privateKey, err := os.ReadFile(privateKeyName)
	if err != nil {
		panic(err)
	}

	password := []byte("123456")
	log.Printf("原始密码===%s", password)
	cipherdata, err := utils.EncryptOAEP(publicKey, password)
	if err != nil {
		panic(err)
	}
	ciphertext := base64.StdEncoding.EncodeToString(cipherdata)
	log.Printf("加密后的密码===%s", ciphertext)

	cipherdata, _ = base64.StdEncoding.DecodeString(ciphertext)
	password, err = utils.DecryptOAEP(privateKey, cipherdata)
	if err != nil {
		log.Printf("解密错误===%v", err)
		panic(err)
	}
	log.Printf("解密后的密码===%s", string(password))*/

	cfg := mysql.Config{
		User:                 config.EnvConfig.DbUser,
		Passwd:               config.EnvConfig.DbPass,
		Addr:                 config.EnvConfig.DbAddress,
		DBName:               config.EnvConfig.DbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            false,
		Collation:            "utf8mb4_general_ci",
		Loc:                  config.EnvConfig.Loc,
	}

	rabbitmqCfg := config.EnvConfig.Rabbitmq
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		rabbitmqCfg.Username,
		rabbitmqCfg.Password,
		rabbitmqCfg.Host,
		rabbitmqCfg.Port,
		rabbitmqCfg.VirtualHost,
	)
	log.Println(connString)

	//rabbitmq.Consume(ch, queue.Name, false)
	/*expressSQL := datasource.InitMysqlDB(cfg)

	db, err := expressSQL.GetDb()
	if err != nil {
		log.Printf("获取数据库信息===%v", err)
	}

	express := impl.NewExpressDB(db)
	user := impl.NewUserDB(db)
	api := router.NewAPIServer(":3000", express, user)
	api.Serve()*/

	//initializeConfig, _ := InitializeConfig()
	//
	//cfg := mysql.Config{
	//	User:                 initializeConfig.DbUser,
	//	Passwd:               initializeConfig.DbPass,
	//	Addr:                 initializeConfig.DbAddress,
	//	DBName:               initializeConfig.DbName,
	//	Net:                  "tcp",
	//	AllowNativePasswords: true,
	//	ParseTime:            false,
	//	Collation:            "utf8mb4_general_ci",
	//	Loc:                  initializeConfig.Loc,
	//}

	//application, _ := InitializeApplication()
	//log.Printf("1===%v", application.expressController)

	//cache.Set("test123", "我只是测试数据", 30)

	// 生成验证码图片
	//img := utils.GenerateImage(code)

	/*fileName := "captcha" + code + ".png"
	// 创建文件
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	//将图片保存为 PNG 格式
	err = png.Encode(file, img)
	if err != nil {
		fmt.Println("Failed to save image:", err)
		return
	}

	fmt.Println("验证码图片已保存为 ", fileName)*/

	//utils.MergeImage()
	//utils.MergeImages2("./image/图片.png", "./image/图片1.png")
	//utils.MergeImages2("./image/test15.png", "./image/1714422153891.jpg", 1)
	expressMySQL := datasource.InitMysqlDB(cfg)
	//expressSQL := mysqlDB.InitMysqlDB(cfg)

	postgresql := datasource.InitPostgresql()
	postgresqlDb, err := postgresql.GetPostgresqlDB()
	if err != nil {
		log.Printf("获取Postgresql数据库信息失败===%v", err)
	}

	mySqlDb, err := expressMySQL.GetDb()
	if err != nil {
		log.Printf("获取Mysql数据库信息失败===%v", err)
	}

	express := impl.NewExpressDB(mySqlDb)
	user := impl.NewUserDB(mySqlDb)
	dict := impl.NewDictTypeDb(postgresqlDb)
	testExpressDB := impl.NewTestDictTypeDb(postgresqlDb)

	serverPort := ":" + utils.ConvertIntToStr(config.EnvConfig.ServerPort)
	api := router.NewAPIServer(serverPort, express, user, rabbitmq.ConnectRabbitmq(connString), dict, testExpressDB)
	api.Serve()

	//log.Printf("获取数据库信息===%v", dbWire)
	//mysqlConfig, _ := InitializeEnglishGreeter()

	////go:build wireinject
	//// +build wireinject
}
