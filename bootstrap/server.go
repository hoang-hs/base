package bootstrap

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/common/log"
	"github.com/hoang-hs/base/configs"
	"github.com/hoang-hs/base/present/middleware"
	"github.com/hoang-hs/base/present/router"
	"go.uber.org/fx"
)

func BuildHTTPServerModule() fx.Option {
	return fx.Options(
		fx.Provide(gin.New),
		fx.Invoke(middleware.RegisterGinEngineDefault),
		fx.Invoke(router.RegisterGinRouters),
		fx.Invoke(NewHttpServer),
	)
}

func NewHttpServer(lc fx.Lifecycle, engine *gin.Engine) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := engine.Run(fmt.Sprintf(":%s", configs.Get().Server.Http.Address)); err != nil {
					log.Fatal("Cannot start application", log.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			return nil
		},
	})
}
