package ports

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	// Other logging methods...
}

// type Logger interface {
//     Info(ctx context.Context, msg string, fields map[string]interface{})
//     Warn(ctx context.Context, msg string, fields map[string]interface{})
//     Error(ctx context.Context, msg string, fields map[string]interface{})
//     Debug(ctx context.Context, msg string, fields map[string]interface{})
// }
