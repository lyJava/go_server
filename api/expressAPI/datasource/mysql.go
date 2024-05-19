package datasource

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.uber.org/zap"
)

type MysqlDB struct {
	Db *sql.DB
}

// InitMysqlDB 初始化数据库
func InitMysqlDB(cfg mysql.Config) *MysqlDB {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	//db.SetMaxIdleConns(20)
	//db.SetMaxOpenConns(10)
	log.Println("Mysql connection initialized successfully")
	return &MysqlDB{Db: db}
}

// GetDb 返回数据库
func (s *MysqlDB) GetDb() (*sql.DB, error) {
	db := s.Db
	row := db.QueryRow("SELECT VERSION()")
	var version string

	err := row.Scan(&version)
	if err != nil {
		log.Printf("查询Mysql数据库版本失败==%+v", err)
		zap.L().Sugar().Errorf("查询Mysql数据库版本失败==%+v", err)
		return nil, err
	}
	//zap.L().Info("mysql连接信息", zap.String("Mysql current version:", version))
	zap.L().Sugar().Infoln("Mysql current version:", version)
	return db, nil
}

func InitializeMysqlDB(cfg mysql.Config) (*MysqlDB, error) {
	wire.Build(InitMysqlDB)
	return &MysqlDB{}, nil
}
