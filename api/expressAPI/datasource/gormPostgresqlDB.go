package datasource

import (
	"apiProject/api/expressAPI/config"
	cfg "apiProject/api/expressAPI/types/config"
	"apiProject/api/expressAPI/types/domain"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

// customLogger 自定义日志记录器
type customLogger struct {
	logger.Interface
	out io.Writer
}

// LogMode 实现日志输出的 Format 方法
func (c *customLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *c
	return &newLogger
}

func (c *customLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	c.log("INFO", msg, data...)
}

func (c *customLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	c.log("WARN", msg, data...)
}

func (c *customLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	c.log("ERROR", msg, data...)
}

func (c *customLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 获取 SQL 查询语句和执行的行数
	sql, rows := fc()
	// 获取文件名和行号
	_, file, line, _ := runtime.Caller(2)
	// 计算查询的执行时间
	duration := time.Since(begin)
	// 格式化日志
	c.log("TRACE", fmt.Sprintf("%s:%d\n[%v] [rows:%d]\nSQL: %s\r\n[elapsed: %v]", file, line, duration, rows, sql, duration), nil)
}

// 自定义日志输出格式
func (c *customLogger) log(level string, msg string, data ...interface{}) {
	fmt.Fprintln(c.out, fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), level, fmt.Sprintf(msg, data...)))
}

type GormPostgresSqlDb struct {
	PsDB *gorm.DB
}

// InitGormPostgresSql gorm初始化postgresql连接
func InitGormPostgresSql() *GormPostgresSqlDb {
	var err error
	viperConfig := config.ReadConfig("api/expressApi/config", "application", "yml")
	if viperConfig == nil {
		log.Println("无Postgresql配置信息")
	}
	var postgresqlConfig cfg.PostgresqlConfig
	if err := viperConfig.Unmarshal(&postgresqlConfig); err != nil {
		log.Printf("获取Postgresql配置错误:%+v", err)
	}

	postgresqlItem := postgresqlConfig.Postgresql
	log.Printf("获取postgresqlItem:%v", postgresqlItem)
	// 创建日志文件
	file, err := os.OpenFile("gorm_sql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建 MultiWriter，输出到控制台和日志文件
	multiWriter := io.MultiWriter(os.Stdout, file)

	// 创建自定义的 logger 实例
	customGormLogger := &customLogger{out: multiWriter}

	// 自定义日志输出函数
	//log.SetFlags(0) // 禁用默认的时间格式
	//log.SetPrefix("") // 去掉日志前缀

	// 自定义日志输出
	//gormLogger := logger.New(
	//	// 使用自定义时间格式
	//	log.New(multiWriter, "", 0), // 设置日志输出
	//	logger.Config{
	//		LogLevel:                  logger.Info, // 设置日志级别为 Info
	//		SlowThreshold:             time.Second, // 慢查询时间阈值
	//		IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
	//		Colorful:                  true,        // 允许控制台输出颜色
	//	},
	//)

	// 构建连接字符串
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		postgresqlItem.Host, postgresqlItem.Port, postgresqlItem.User, postgresqlItem.Password, postgresqlItem.DbName)
	// 连接到数据库
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tb_", // 可以设置表前缀
			SingularTable: true,  // 禁用复数表名
		},
		//Logger: logger.Default.LogMode(logger.Info),
		Logger: customGormLogger, // 使用自定义的日志器
	})

	if err != nil {
		log.Printf("Failed to connect to Postgresql database:%v", err)
	}
	log.Printf("Gorm Create Postgresql Database connection successful")

	sqlDb, err := db.DB()
	if err != nil {
		log.Printf("Failed to Get the Postgresql Db==%v", err)
	}

	err = sqlDb.Ping()
	if err != nil {
		log.Printf("Failed to Ping the Postgresql Db==%v", err)
	}
	log.Printf("Postgresql Database 状态：%v", sqlDb.Stats())

	setConnPoolProps(sqlDb)

	//var storeAdmin domain.StoreAdmin
	//if err := gormDB.PsDB.First(&storeAdmin, 122).Error; err != nil {
	//	fmt.Printf("User not found:%v", err)
	//} else {
	//	marshal, err := json.Marshal(storeAdmin)
	//	if err != nil {
	//		log.Printf("Error marshaling data: %v", err)
	//	}
	//	fmt.Printf("User found:%s", marshal)
	//}
	var storeAdminList []*domain.StoreAdmin
	err = db.Table("tb_store_admin").
		Select("id, user_name, mobile, real_name, status_value, store_name, merchant_id, merchant_name, " +
			"TO_CHAR(create_time, 'YYYY-MM-DD HH24:MI:SS') AS create_time, " +
			"TO_CHAR(update_time, 'YYYY-MM-DD HH24:MI:SS') AS update_time").
		//Where("id = ?", 122).
		Order("id DESC").
		Limit(2).
		Scan(&storeAdminList).Error

	marshal, err := json.MarshalIndent(&storeAdminList, "", "    ")
	if err != nil {
		log.Printf("店铺管理员转换json错误: %v", err)
	}
	log.Printf("店铺管理员:\r\n%s", marshal)

	return &GormPostgresSqlDb{
		PsDB: db,
	}
}

// 自定义日志输出函数
func logWithTimestamp(format string, v ...interface{}) {
	// 获取当前时间并格式化为毫秒
	currentTime := time.Now().Format("2006-01-02 15:04:05.000")
	log.Printf("[%s] %s", currentTime, fmt.Sprintf(format, v...))
}

func setConnPoolProps(db *sql.DB) {
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	log.Println("Postgresql Database Connection Pool Setting success")
}
