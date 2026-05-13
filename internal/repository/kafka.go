package repository

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// KafkaProducer producer wrapper.
type KafkaProducer struct {
	w *kafka.Writer
}

// NewKafkaProducer creates new producer.
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		w: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// Send sends message to Kafka.
func (k *KafkaProducer) Send(ctx context.Context, key, value []byte) error {
	return k.w.WriteMessages(ctx, kafka.Message{Key: key, Value: value})
}

// Close closes writer.
func (k *KafkaProducer) Close() error { return k.w.Close() }
