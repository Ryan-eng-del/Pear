package logs

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var LG *zap.Logger

type LogConfig struct {
	DebugFileName string `json:"debugFileName"`
	InfoFileName string `json:"infoFileName"`
	WarningFileName string `json:"warningFileName"`
	MaxSize int `json:"maxSize"`
	MaxAge int `json:"maxAge"`
	MaxBackups int `json:"maxBackups"`
}

func InitLogger (cfg *LogConfig) (err error) {
	maxSize := cfg.MaxSize
	maxAge := cfg.MaxAge
	maxBackups := cfg.MaxBackups
	writeSyncerDebug := getLogWriter(cfg.DebugFileName, maxSize, maxBackups, maxAge)
	writeSyncerInfo := getLogWriter(cfg.InfoFileName, maxSize, maxBackups, maxAge)
	writeSyncerWarn := getLogWriter(cfg.WarningFileName, maxSize, maxBackups,maxAge)

	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl <= zapcore.DebugLevel
	})

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.WarnLevel
	})

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	// 文件输出， writer 使用第三方实现的日志分割 writer
	debugCore := zapcore.NewCore(getEncoder(), writeSyncerDebug,debugPriority)
	infoCore := zapcore.NewCore(getEncoder(), writeSyncerInfo, infoPriority)
	warnCore := zapcore.NewCore(getEncoder(), writeSyncerWarn, highPriority)

	// 标准输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	stdCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	// core 抽象成 
	// 1. encoder 也就是这条日志的打印编码或打印形式，
	// 2. writer 往哪里去写，怎么去写，哪里去写指的是文件还是标准输出，怎么去写指的是日志分割
	// 3. level 这个级别的日志，怎么打印，以及往那里输出
	// tee 将多个 core 组合在一起，实现不同级别日志的不同入口（std or file）
	cores := zapcore.NewTee(debugCore, infoCore, warnCore, stdCore)
	LG = zap.New(cores, zap.AddCaller())
	zap.ReplaceGlobals(LG)
	return	
}


func getEncoder () zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}


func getLogWriter (filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename: filename,
		MaxSize: maxSize,
		MaxBackups: maxBackup,
		MaxAge:  maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GinLogger () gin.HandlerFunc {
	return func (c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		LG.Info(path, 
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip",c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors",c.Errors.ByType(gin.ErrorTypePrivate).String()),
		  zap.Duration("duration", cost))
	}
}