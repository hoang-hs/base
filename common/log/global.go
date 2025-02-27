package log

import (
	"context"
	"fmt"
	"github.com/hoang-hs/base/common"
	"go.uber.org/zap"
)

var globalLogger *logger

// log format
func InfofCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Infof(addCtxValue(ctx, msg), args...)
}

func Infof(msg string, args ...interface{}) {
	globalLogger.Infof(msg, args...)
}

func DebugfCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debugf(addCtxValue(ctx, msg), args...)
}

func Debugf(msg string, args ...interface{}) {
	globalLogger.Debugf(msg, args...)
}

func WarnfCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warnf(addCtxValue(ctx, msg), args...)
}

func Warnf(msg string, args ...interface{}) {
	globalLogger.Warnf(msg, args...)
}

func ErrorfCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Errorf(addCtxValue(ctx, msg), args...)
}

func Errorf(msg string, args ...interface{}) {
	globalLogger.Errorf(msg, args...)
}

func FatalfCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Fatalf(addCtxValue(ctx, msg), args...)
}

func Fatalf(msg string, args ...interface{}) {
	globalLogger.Fatalf(msg, args...)
}

// log fields
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Info(addCtxValue(ctx, msg), fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Debug(addCtxValue(ctx, msg), fields...)
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Warn(addCtxValue(ctx, msg), fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Error(addCtxValue(ctx, msg), fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	globalLogger.Fatal(addCtxValue(ctx, msg), fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func GetLogger() *logger {
	return globalLogger
}

func addCtxValue(ctx context.Context, msg string) string {
	traceId := common.GetTraceId(ctx)
	return fmt.Sprintf("%s, trace_id:[%s]", msg, traceId)
}
