package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	/* 无日志切割 */
	//buffer := &bytes.Buffer{}                                                            // bytes.Buffer
	//console := os.Stdout                                                                 // 控制台
	//file, err := os.OpenFile("go_comm_mqtt.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0755) // 文件
	//if err != nil {
	//	log.Fatalf("create file go_modbus_mqtt.log failed: %v", err)
	//}
	//
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors:     true,
	//	FullTimestamp:   true,
	//	TimestampFormat: time.RFC3339Nano,
	//})
	//logrus.SetOutput(io.MultiWriter(buffer, console, file))

	/* 添加hook方式，使用之后开启颜色，只有控制台有颜色，输出到文件没有颜色了 */
	//rotateHook := newLfsHook("./log/rotate", "debug", 5)
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors:     true,
	//	FullTimestamp:   true,
	//	TimestampFormat: time.RFC3339Nano,
	//})
	//logrus.AddHook(rotateHook)

	/* 说是lumberjack包不更新维护了，但用起来一切正常 */
	logger := &lumberjack.Logger{
		LocalTime:  true,
		Filename:   "./log/go_comm_mqtt.log",
		MaxSize:    20, // 一个文件最大为nM
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

func newLfsHook(logName string, logLevel string, maxRemainCnt uint) logrus.Hook {
	writer, err := rotatelogs.New(
		logName+"-%Y%m%d%H%M.log",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Second*60),

		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		// rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	// 使用了lfshook软件包创建了一个新的日志钩子，该钩子将日志记录到指定的日志文件中。
	// lfshook.WriterMap指定了每个日志级别所使用的写入器（writer）。
	// 在这个函数中，所有的日志级别都使用同一个写入器writer。
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	return lfsHook
}
