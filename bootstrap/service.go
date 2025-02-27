package bootstrap

import (
	"github.com/hoang-hs/base/core/service"
	"github.com/hoang-hs/base/pkg"
	"github.com/hoang-hs/base/pkg/alert"
	"go.uber.org/fx"
)

func BuildServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(service.NewBaseService),
		fx.Provide(alert.NewTelegram),
		fx.Provide(pkg.NewRecoveryInterceptor),
	)
}
