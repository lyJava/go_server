package datasource

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"log"
)

var ProviderSet = wire.NewSet(NewMysqlDB)

type MysqlDB struct {
	db *sql.DB
}

func NewMysqlDB(cfg mysql.Config) *MysqlDB {
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

func (s *MysqlDB) Init() (*sql.DB, error) {
	return s.db, nil
}
