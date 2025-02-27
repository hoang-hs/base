package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/src/configs"
	cors "github.com/rs/cors/wrapper/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type RoutersIn struct {
	fx.In
	cf     *configs.Config
	engine *gin.Engine
}

func RegisterGinRouters(in RoutersIn) {
	in.engine.Use(cors.AllowAll())

	group := in.engine.Group(in.cf.Server.Http.Prefix)
	// http swagger serve
	if in.cf.Swagger.Enabled {
		group.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	registerPublicRouters(group, in)

	adminGroup := group.Group("/")
	adminGroup.Use()
	{
		registerAdminRouters(adminGroup, in)
	}
}

func registerPublicRouters(r *gin.RouterGroup, in RoutersIn) {

}

func registerAdminRouters(r *gin.RouterGroup, in RoutersIn) {

}
