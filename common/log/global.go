package log

import (
	"context"
	"fmt"
	"github.com/hoang-hs/base/common"
)

var globalLogger *logger

func InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Info(addCtxValue(ctx, msg), args...)
}

func Info(msg string, args ...interface{}) {
	globalLogger.Info(msg, args...)
}

func DebugCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Debug(addCtxValue(ctx, msg), args...)
}

func Debug(msg string, args ...interface{}) {
	globalLogger.Debug(msg, args...)
}

func WarnCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Warn(addCtxValue(ctx, msg), args...)
}

func Warn(msg string, args ...interface{}) {
	globalLogger.Warn(msg, args...)
}

func ErrorCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Error(addCtxValue(ctx, msg), args...)
}

func Error(msg string, args ...interface{}) {
	globalLogger.Error(msg, args...)
}

func FatalCtx(ctx context.Context, msg string, args ...interface{}) {
	globalLogger.Fatal(addCtxValue(ctx, msg), args...)
}

func Fatal(msg string, args ...interface{}) {
	globalLogger.Fatal(msg, args...)
}

func GetLogger() *logger {
	return globalLogger
}

func addCtxValue(ctx context.Context, msg string) string {
	traceId := common.GetTraceId(ctx)
	return fmt.Sprintf("%s, trace_id:[%s]", msg, traceId)
}
