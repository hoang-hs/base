package bootstrap

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	log2 "github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/configs"
	"github.com/hoang-hs/base/src/present/middleware"
	"github.com/hoang-hs/base/src/present/router"
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
					log2.Fatal("Cannot start application", log2.Err(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log2.Info("Stopping HTTP server")
			return nil
		},
	})
}
