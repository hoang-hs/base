package middleware

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/common"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"math/rand"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanContext := trace.SpanContextFromContext(c.Request.Context())
		span := trace.SpanFromContext(c.Request.Context())
		span.SetName(fmt.Sprintf("[%s] %s", c.Request.Method, c.FullPath()))
		var traceId string
		if spanContext.TraceID().IsValid() {
			traceId = spanContext.TraceID().String()
		} else {
			traceIdByte := make([]byte, 16)
			rand.Read(traceIdByte)
			traceId = hex.EncodeToString(traceIdByte[:])
		}
		c.Set(common.TraceIdName, traceId)
		traceContext := context.WithValue(c.Request.Context(), common.TraceIdName, traceId)
		ctxMetaData := metadata.AppendToOutgoingContext(traceContext, []string{common.TraceIdName, traceId}...)
		c.Request = c.Request.WithContext(ctxMetaData)
		c.Next()
	}
}
