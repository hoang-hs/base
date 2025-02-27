package bootstrap

import (
	"github.com/hoang-hs/base/src/present/controller"
	"github.com/hoang-hs/base/src/present/validator"
	"go.uber.org/fx"
)

func BuildControllerModule() fx.Option {
	return fx.Options(
		fx.Provide(controller.NewBaseController),
	)
}

func BuildValidator() fx.Option {
	return fx.Options(
		fx.Provide(validator.NewValidator),
		fx.Invoke(validator.RegisterValidations),
	)
}
