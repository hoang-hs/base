package consumer

import (
	"context"
	"github.com/hoang-hs/base/common/log"
	"github.com/hoang-hs/base/configs"
	"github.com/hoang-hs/base/pkg"
	ikafka "github.com/hoang-hs/base/pkg/kafka"
	"go.uber.org/fx"
)

func NewTestKafkaConsumer(lc fx.Lifecycle, recoveryInterceptor *pkg.RecoveryInterceptor, handleMessage ikafka.HandleFunc) {
	quit := make(chan bool)
	var consumer *ikafka.Consumer
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("OnStart CDC consumer")
			consumer = ikafka.NewConsumer(configs.Get().Kafka, quit, recoveryInterceptor.RecoveryConsumer)
			go consumer.Run(ctx, handleMessage)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("OnStop CDC consumer")
			quit <- true
			return consumer.Close()
		},
	})
}
