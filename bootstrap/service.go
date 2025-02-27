package bootstrap

import (
	"github.com/hoang-hs/base/core/service"
	"go.uber.org/fx"
)

func BuildServicesModules() fx.Option {
	return fx.Options(
		fx.Provide(service.NewBaseService),
	)
}
