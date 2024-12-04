package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.SugaredLogger

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}

	zapLog = logger.Sugar()
}

func Info(args ...any) {
	zapLog.Info(args)
}

func Infof(format string, args ...any) {
	zapLog.Infof(format, args...)
}

func Infow(msg string, keysAndValues ...any) {
	zapLog.Infow(msg, keysAndValues...)
}

func Error(args ...any) {
	zapLog.Error(args)
}

func Fatal(args ...any) {
	zapLog.Fatal(args)
}

func Debug(args ...any) {
	zapLog.Debug(args)
}
