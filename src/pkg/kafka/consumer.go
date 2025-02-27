package kafka

import (
	"context"
	"errors"
	"github.com/hoang-hs/base/src/common"
	log2 "github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/configs"
	"io"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	timeoutKafka = 2 * time.Second
)

type HandleFunc func(ctx context.Context, m *kafka.Message) error

func NewConsumer(
	cf *configs.Kafka,
	onStop chan bool,
	onRecover func(ctx context.Context, m *kafka.Message),
) *Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  cf.Host,
		"group.id":           cf.Consumer.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log2.Fatal("Failed to create new consumer", log2.Err(err))
	}
	err = c.SubscribeTopics([]string{cf.Consumer.Topic}, nil)
	if err != nil {
		log2.Fatal("Failed to subscribe topic")
	}

	return &Consumer{
		consumer:  c,
		onStop:    onStop,
		onRecover: onRecover,
	}
}

type Consumer struct {
	consumer  *kafka.Consumer
	onStop    chan bool
	onRecover func(ctx context.Context, m *kafka.Message)
}

func (c *Consumer) Run(ctx context.Context, handleMessage HandleFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.onStop:
			return
		default:
			func() {
				ctxMessage := common.CreateNewCtx()
				m, err := c.consumer.ReadMessage(timeoutKafka)
				defer func(m *kafka.Message) {
					if r := recover(); r != nil {
						log2.Error("Recovered from panic", log2.Any("panic", r))
						if c.onRecover != nil {
							c.onRecover(context.WithoutCancel(ctx), m)
						}
					}
				}(m)
				switch {
				case errors.Is(err, io.EOF):
					return
				case err != nil:
					log2.Error("Failed to read message", log2.Err(err))
					return
				}
				// process message
				log2.Info("New message received", log2.String("topic", *m.TopicPartition.Topic),
					log2.String("offset", m.TopicPartition.Offset.String()), log2.Int32("partition", m.TopicPartition.Partition))
				err = handleMessage(ctxMessage, m)
				if err != nil {
					log2.Error("Failed to handle message", log2.Err(err))
				}
			}()
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}
