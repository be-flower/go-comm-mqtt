package logger

import (
	"io"
	"os"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/sirupsen/logrus"
)

func init() {

	/* 说是lumberjack包不更新维护了，但用起来一切正常 */
	logger := &lumberjack.Logger{
		LocalTime:  true,
		Filename:   "./log/go_comm_mqtt.log",
		MaxSize:    50, // 一个文件最大为nM
		MaxBackups: 5,  // 最多同时保存n份文件(加上正在使用的文件共n+1份)
		MaxAge:     30, // 一个文件最多同时存在30天
		Compress:   false,
	}
	writers := []io.Writer{
		logger,
		os.Stdout, // 控制台
	}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logrus.SetOutput(fileAndStdoutWriter)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})
}
