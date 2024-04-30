package main

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/datasource"
	"apiProject/api/expressAPI/router"
	"apiProject/api/expressAPI/service/impl"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

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

	// 设置时区为亚洲/上海
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalln(err)
	}
	cfg := mysql.Config{
		User:                 config.EnvConfig.DbUser,
		Passwd:               config.EnvConfig.DbPass,
		Addr:                 config.EnvConfig.DbAddress,
		DBName:               config.EnvConfig.DbName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            false,
		Collation:            "utf8mb4_general_ci",
		Loc:                  loc,
	}
	expressSQL := datasource.NewMysqlDB(cfg)

	db, err := expressSQL.Init()
	if err != nil {
		log.Printf("初始化失败===%v", err)
	}

	/*dbWire, err := InitializeWire(expressSQL)
	log.Printf("获取数据库信息===%v", dbWire)
	if err != nil {
		fmt.Printf("failed to create db: %s\n", err)
		os.Exit(2)
	}*/

	express := impl.NewExpressDB(db)
	user := impl.NewUserDB(db)
	api := router.NewAPIServer(":3000", express, user)
	api.Serve()

}
