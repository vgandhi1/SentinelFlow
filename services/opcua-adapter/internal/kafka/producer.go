package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer writes telemetry to the raw topic.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a writer for the raw topic.
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

// Produce sends a message with key (e.g. idempotency_key) and value (JSON).
func (p *Producer) Produce(ctx context.Context, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: value,
	})
}

// Close closes the writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
