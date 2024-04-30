package datasource

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"log"
)

var ProviderSet = wire.NewSet(InitMysqlDB)

type MysqlDB struct {
	db *sql.DB
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
	log.Println("mysql connect success")
	return &MysqlDB{db: db}
}

// GetDb 返回数据库
func (s *MysqlDB) GetDb() (*sql.DB, error) {
	return s.db, nil
}
