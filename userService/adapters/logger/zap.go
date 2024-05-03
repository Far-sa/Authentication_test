package logger

import (
	"os"
	"user-svc/ports"

	// Example config library (replace if needed)
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	config ports.Config
	logger *zap.Logger
}

func NewZapLogger(config ports.Config) (*zapLogger, error) {
	logConfig := config.GetLoggerConfig()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultEncoder := zapcore.NewJSONEncoder(encoderCfg)

	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.Filename,
		LocalTime:  logConfig.LocalTime,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	})

	stdOutWriter := zapcore.AddSync(os.Stdout)
	defaultLogLevel := zapcore.InfoLevel
	core := zapcore.NewTee(
		zapcore.NewCore(defaultEncoder, writer, defaultLogLevel),
		zapcore.NewCore(defaultEncoder, stdOutWriter, zap.InfoLevel),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &zapLogger{config: config, logger: logger}, nil

}

// Info logs an informational message.
func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Error logs an error message.
func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}
