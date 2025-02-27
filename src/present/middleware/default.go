package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/src/configs"
)

func RegisterGinEngineDefault(cf *configs.Config, engine *gin.Engine) *gin.Engine {
	engine.Use(Recovery())
	engine.Use(Otel(cf.Server.Name))
	engine.Use(Tracer())
	engine.Use(Log())
	engine.Use(Cors())
	group := engine.Group(cf.Server.Http.Prefix)
	group.GET("/ping", HealthCheckEndpoint)
	return engine
}
