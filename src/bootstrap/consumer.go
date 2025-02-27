package bootstrap

import (
	"github.com/hoang-hs/base/src/present/consumer"
	"go.uber.org/fx"
)

func BuildConsumerModule() fx.Option {
	return fx.Options(
		fx.Invoke(consumer.NewTestKafkaConsumer),
	)
}
