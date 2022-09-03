package gflog

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var logger *zap.Logger

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	enc := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	logger = zap.New(zapcore.NewTee(
		zapcore.NewCore(enc,
			zapcore.Lock(os.Stdout),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level < zapcore.ErrorLevel
			})),
		zapcore.NewCore(enc,
			zapcore.Lock(os.Stderr),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= zapcore.ErrorLevel
			})),
	)).Named("gflog").WithOptions(zap.AddCallerSkip(1))
	// redirect output from fmt package to logger
	zap.RedirectStdLog(logger)
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Debug(fillTrace(ctx, msg), fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Info(fillTrace(ctx, msg), fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Warn(fillTrace(ctx, msg), fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Error(fillTrace(ctx, msg), fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger.Fatal(fillTrace(ctx, msg), fields...)
}

func fillTrace(ctx context.Context, str string) string {
	if ctx == nil {
		return str
	}
	var sb strings.Builder
	sb.WriteString("[TRACE: ")
	sb.WriteString(GetTraceID(ctx))
	sb.WriteString("] ")
	sb.WriteString(str)
	return sb.String()
}
