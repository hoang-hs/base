package kafka

import (
	"context"
	"errors"
	"github.com/hoang-hs/base/common/log"
	"io"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	timeoutKafka = 2 * time.Second
)

type HandleFunc func(ctx context.Context, m *kafka.Message) error

func NewConsumer(
	consumer *kafka.Consumer,
	onStop func(),
	onRecover func(ctx context.Context, m *kafka.Message),
) *Consumer {
	return &Consumer{
		consumer:  consumer,
		onStop:    onStop,
		onRecover: onRecover,
	}
}

type Consumer struct {
	consumer  *kafka.Consumer
	onStop    func()
	onRecover func(ctx context.Context, m *kafka.Message)
}

func (c *Consumer) Start(ctx context.Context, handleMessage HandleFunc) {
	defer c.onStop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			func() {
				m, err := c.consumer.ReadMessage(timeoutKafka)
				switch {
				case errors.Is(err, io.EOF):
					return
				case err != nil:
					log.Error("Failed to read message", log.Err(err))
					return
				}
				// process message
				log.Info("New message received", log.String("topic", *m.TopicPartition.Topic),
					log.String("offset", m.TopicPartition.Offset.String()), log.Int32("partition", m.TopicPartition.Partition))
				err = handleMessage(ctx, m)
				if err != nil {
					log.Error("Failed to handle message", log.Err(err))
				}
			}()
		}
	}
}
