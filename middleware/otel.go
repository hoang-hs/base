package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func Otel() gin.HandlerFunc {
	return otelgin.Middleware(base.Get().Name)
}
