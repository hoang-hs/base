package bootstrap

import (
	"github.com/hoang-hs/base/src/pkg/observe"
	"go.uber.org/fx"
)

func BuildObserveServicesModules() fx.Option {
	return fx.Options(
		fx.Invoke(observe.NewMetric),
		fx.Invoke(observe.InitTrace),
	)
}
