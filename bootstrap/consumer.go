package bootstrap

import (
	"github.com/hoang-hs/base/present/consumer"
	"go.uber.org/fx"
)

func BuildConsumerModule() fx.Option {
	return fx.Options(
		fx.Invoke(consumer.NewTestKafkaConsumer),
	)
}
