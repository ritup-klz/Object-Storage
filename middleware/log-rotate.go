package middleware

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewZapLogger() *zap.Logger {
	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./logs/server.log", // log file path
		MaxSize:    10,                  // megabytes
		MaxBackups: 5,                   // number of old files to retain
		MaxAge:     30,                  // days
		Compress:   true,                // gzip old logs
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // structured logs
		logWriter,
		zap.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger
}
