package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/configs"
)

func RegisterGinEngineDefault(cf *configs.Config) *gin.Engine {
	engine := gin.New()
	engine.Use(Recovery())
	engine.Use(Otel(cf.Server.Name))
	engine.Use(Tracer())
	engine.Use(Log())
	engine.Use(Cors())
	group := engine.Group(cf.Server.Http.Prefix)
	group.GET("/ping", HealthCheckEndpoint)
	return engine
}
