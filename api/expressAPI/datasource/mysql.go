package datasource

import (
	"apiProject/api/expressAPI/config"
	cfg "apiProject/api/expressAPI/types/config"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"go.uber.org/zap"
)

type MysqlDB struct {
	Db *sql.DB
}

// InitMysqlDB 初始化mysql
func InitMysqlDB() *MysqlDB {
	// 从配置文件中获取mysql的配置信息
	viperConfig := config.ReadConfig("api/expressAPI/config", "application", "yml")
	if viperConfig == nil {
		zap.L().Sugar().Errorf("无mysql配置信息")
		return nil
	}

	var sqlConfig cfg.SqlConfig
	if err := viperConfig.Unmarshal(&sqlConfig); err != nil {
		zap.L().Sugar().Errorf("获取mysql配置错误:%+v", err)
		return nil
	}

	mysqlItem := sqlConfig.Mysql

	// 构建mysql连接相关配置
	buildConfig := mysql.Config{
		User:                 mysqlItem.Username,
		Passwd:               mysqlItem.Password,
		Addr:                 mysqlItem.Url,
		DBName:               mysqlItem.Database,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            false,
		Collation:            "utf8mb4_general_ci",
		Loc:                  config.EnvConfig.Loc,
	}

	db, err := sql.Open("mysql", buildConfig.FormatDSN())
	if err != nil {
		zap.L().Sugar().Errorf("mysql连接错误==%+v", err)
		//log.Fatalln(err) 这里会直接退出程序
		return nil
	}
	err = db.Ping()
	if err != nil {
		zap.L().Sugar().Errorf("mysql-ping错误==%+v", err)
		return nil
	}
	//db.SetMaxIdleConns(20)
	//db.SetMaxOpenConns(10)
	zap.L().Sugar().Info("Mysql connection initialized successfully")
	return &MysqlDB{Db: db}
}

// GetDb 返回数据库
func (s *MysqlDB) GetDb() (*sql.DB, error) {
	db := s.Db
	row := db.QueryRow("SELECT VERSION()")
	var version string

	err := row.Scan(&version)
	if err != nil {
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
