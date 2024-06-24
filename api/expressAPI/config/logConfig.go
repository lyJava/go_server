package config

import (
	"log"
	"os"
	"time"
)


func init() {
    log.SetFlags(0)
    log.SetOutput(new(logWriter))
}

type logWriter struct{}

// Write 自定义的logWriter设置日志时间格式yyyy-MM-dd HH:mm:ss.SSS
func (writer logWriter) Write(bytes []byte) (int, error) {
    return os.Stdout.Write(append([]byte(time.Now().Format("2006-01-02 15:04:05.000 ")), bytes...))
}

