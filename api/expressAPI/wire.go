package main

import (
	"apiProject/api/expressAPI/datasource"
	"github.com/google/wire"
)

func InitializeWire(db *datasource.MysqlDB) (*datasource.MysqlDB, error) {
	wire.Build(datasource.ProviderSet)
	return db, nil
}
