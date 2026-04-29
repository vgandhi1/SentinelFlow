package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// Producer writes to Kafka.
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a writer for the given brokers and topic.
func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{writer: w}
}

// Produce sends a JSON payload with optional key.
func (p *Producer) Produce(ctx context.Context, key string, payload interface{}) error {
	val, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: val})
}

// Close closes the writer.
func (p *Producer) Close() error {
	return p.writer.Close()
}
