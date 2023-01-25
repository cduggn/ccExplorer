package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func init() {
	_, err := New()
	if err != nil {
		panic(err.Error())
	}
}

type LoggerError struct {
	msg string
}

func (e LoggerError) Error() string {
	return e.msg
}

func New() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	enccoderConfig := zap.NewProductionEncoderConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	enccoderConfig.StacktraceKey = "" // disable stacktrace
	config.EncoderConfig = enccoderConfig

	var err error
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, LoggerError{
			msg: "unable to create zap logger, " + err.Error(),
		}
	}
	return zapLog, nil
}

//func setLogLevel(level string) error {

//atomicLevel := zap.NewAtomicLevel()
//atomicLevel.SetLevel(zapcore.InfoLevel)
//zapLog = zapLog.WithOptions(zap.IncreaseLevel(atomicLevel))
//
//
//zapLog.Info("setting log level to " + level)
//switch level {
//case "debug":
//	zapLog.Sugar().SetLevel(zap.DebugLevel)
//case "info":
//	zapLog.SetLevel(zap.InfoLevel)
//case "warning":
//	zapLog.SetLevel(zap.WarnLevel)
//case "error":
//	zapLog.SetLevel(zap.ErrorLevel)
//case "critical":
//	zapLog.SetLevel(zap.FatalLevel)
//default:
//	return LoggerError{
//		msg: "invalid log level: " + level,
//	}
//}
//}

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
