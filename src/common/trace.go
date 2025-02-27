package common

import "context"

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ""
	if ctx.Value(TraceIdName) != nil {
		traceId = ctx.Value(TraceIdName).(string)
	}
	return traceId
}
