package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

type LoggerError struct {
	msg string
}

func (e LoggerError) Error() string {
	return e.msg
}

func newZapLogger() error {
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	enccoderConfig.StacktraceKey = "" // disable stacktrace
	config.EncoderConfig = enccoderConfig

	var err error
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return LoggerError{
			msg: "unable to create zap logger, " + err.Error(),
		}
	}
	return nil
}

func init() {
	err := newZapLogger()
	if err != nil {
		panic(err.Error())
	}
}

func Info(message string, fields ...zap.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	zapLog.Fatal(message, fields...)
}
