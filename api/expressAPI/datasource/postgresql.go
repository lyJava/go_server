package datasource

import (
	"apiProject/api/expressAPI/config"
	"apiProject/api/expressAPI/types"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresqlDB struct {
	Db *sql.DB
}

func InitPostgresql() *PostgresqlDB {
	viperConfig := config.ReadConfig("api/expressApi/config", "application", "yml")
	var postgresqlConfig types.PostgresqlConfig
	viperConfig.Unmarshal(&postgresqlConfig)

	postgresqlItem := postgresqlConfig.Postgresql
	// 构建连接字符串
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresqlItem.Host, postgresqlItem.Port, postgresqlItem.User, postgresqlItem.Password, postgresqlItem.DbName)
	// 连接到数据库
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Postgresql连接失败:%+v", err)
	}
	return &PostgresqlDB{
		Db: db,
	}
}

func (d *PostgresqlDB) GetPostgresqlDB() (*sql.DB, error) {
	postgresqlDb := d.Db
	row := postgresqlDb.QueryRow("SELECT VERSION()")
	var version string

	err := row.Scan(&version)
	if err != nil {
		log.Printf("查询Postgresql数据库版本失败==%+v", err)
		zap.L().Sugar().Errorf("查询Postgresql数据库版本失败==%+v", err)
		return nil, err
	}
	//zap.L().Info("Postgresql连接信息", zap.String("Postgresql current version:",version))
	zap.L().Sugar().Infoln("Postgresql current version:", version)
	return postgresqlDb, err
}
