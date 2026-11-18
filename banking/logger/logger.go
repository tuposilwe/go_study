package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	// Use zap.NewDevelopment() for colored, human-readable console output
	var err error

	config := zap.NewProductionConfig()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.StacktraceKey = ""
	config.EncoderConfig = encoderConfig

	log, err = config.Build(zap.AddCallerSkip(1))
	defer log.Sync() // Flushes buffer, if any

	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zap.Field) {
	log.Info(message, fields...)
}
func Debug(message string, fields ...zap.Field) {
	log.Debug(message, fields...)
}
func Error(message string, fields ...zap.Field) {
	log.Error(message, fields...)
}
