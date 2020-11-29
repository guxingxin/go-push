package gateway

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
)

var customLogger ZapLogger

func InitLogger() (err error) {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.

	fmt.Println(G_config.LogPath)
	file := &lumberjack.Logger{
		Filename:   filepath.Join(G_config.LogPath, "gateway.log"), // 日志文件路径
		MaxSize:    1024,                                           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10,                                             // 日志文件最多保存多少个备份
		MaxAge:     30,                                             // 文件最多保存多少天
		Compress:   true,                                           // 是否压缩
	}
	cf := zap.NewProductionEncoderConfig()
	cf.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cf),                         // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(file)), // 打印到控制台和文件
		zap.DebugLevel, // 日志级别
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	if logger != nil {
		customLogger.Sugar = logger.Sugar()
	}

	DebugW("printf config", "config", G_config)

	return err
}

//Sync flushes buffered logs (if any).
func Sync() {
	if customLogger.Sugar != nil {
		customLogger.Sugar.Sync()
	}
}

func DebugW(msg string, keysAndValues ...interface{}) {
	if customLogger.Sugar != nil {
		customLogger.Sugar.Debugw(msg, keysAndValues...)
	}
}

type ZapLogger struct {
	Sugar *zap.SugaredLogger
}
