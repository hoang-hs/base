package middleware

import (
	"github.com/gin-gonic/gin"
)

func RegisterHandlerDefault(engine *gin.Engine) {
	engine.Use(Recovery())
	engine.Use(Otel())
	engine.Use(Tracer())
	engine.Use(Log())
	engine.Use(Cors())
}
