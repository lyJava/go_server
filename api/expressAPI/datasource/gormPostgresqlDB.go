package datasource

import (
	"apiProject/api/expressAPI/config"
	cfg "apiProject/api/expressAPI/types/config"
	"context"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
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
	consoleOut                io.Writer // 控制台输出
	fileOut                   io.Writer // 文件输出
	ignoreRecordNotFoundError bool      // 是否忽略记录未找到错误
	colorful                  bool      // 是否启用彩色输出
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
	execSql, rows := fc()
	// 获取文件名和行号
	_, file, line, ok := runtime.Caller(2)
	if ok {
		// 计算查询的执行时间
		duration := time.Since(begin)
		c.log("TRACE", fmt.Sprintf("%s:%d\n%v\n[rows:%d]\nSQL: %s\r\n[elapsed: %v]", file, line, duration, rows, execSql, duration))
	} else {
		c.log("TRACE", fmt.Sprintf("%s:%d\ngorm获取执行信息失败\n[rows:%d]\nSQL: %s\r\n", file, line, rows, execSql))
	}
}

// 自定义日志输出格式
/*func (c *customLogger) log(level string, msg string, data ...interface{}) {
	fmt.Fprintln(c.out, fmt.Sprintf("%s [%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), level, fmt.Sprintf(msg, data...)))
}*/

func (cusLog *customLogger) log(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	// 根据日志级别设置颜色
	var logMessage string
	// 控制台输出
	if cusLog.consoleOut != nil {
		if cusLog.colorful {
			switch level {
			case "INFO":
				logMessage = color.GreenString("[%s] [%s] %s\n", timestamp, level, message)
			case "WARN":
				logMessage = color.YellowString("[%s] [%s] %s\n", timestamp, level, message)
			case "ERROR":
				logMessage = color.RedString("[%s] [%s] %s\n", timestamp, level, message)
			case "TRACE":
				logMessage = color.HiCyanString("[%s] [%s] %s\n", timestamp, level, message)
			default:
				logMessage = fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
			}
		} else {
			// 控制台输出，不添加颜色
			logMessage = fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
		}
		cusLog.consoleOut.Write([]byte(logMessage))
	}

	if cusLog.fileOut != nil {
		// 输出到文件不添加颜色
		cusLog.fileOut.Write([]byte(fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)))
	}

	//c.out.Write([]byte(fmt.Sprintf("%s %s %s\n", timestamp, level, message)))
	/*if f, ok := c.out.(*os.File); ok {
		f.Sync() // 刷新文件缓冲区
	}*/
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
	// TODO 这里的文件关闭会导致每次只能记录最后一条sql
	//defer file.Close()

	// 创建 MultiWriter，输出到控制台和日志文件
	// multiWriter := io.MultiWriter(os.Stdout, file)

	// 创建自定义的 logger 实例
	customGormLogger := &customLogger{
		//out:                       multiWriter,
		consoleOut:                os.Stdout,
		fileOut:                   file,
		ignoreRecordNotFoundError: true,
		colorful:                  true,
	}

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

	return &GormPostgresSqlDb{
		PsDB: db,
	}
}

func setConnPoolProps(db *sql.DB) {
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)
	log.Println("Postgresql Database Connection Pool Setting success")
}
