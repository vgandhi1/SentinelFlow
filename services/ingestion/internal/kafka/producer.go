package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/segmentio/kafka-go"
)

// Producer abstracts Kafka produce. Implement with segmentio/kafka-go or confluent-kafka-go.
type Producer interface {
	Send(ctx context.Context, key string, value []byte) error
	SendBatch(ctx context.Context, key string, values [][]byte) error
	Close() error
}

// NewWriterProducer builds a Kafka producer from broker list and topic. Brokers is a comma-separated list.
func NewWriterProducer(brokersCSV, topic string) (*WriterProducer, error) {
	brokers := splitBrokers(brokersCSV)
	if len(brokers) == 0 || topic == "" {
		return nil, fmt.Errorf("kafka: brokers and topic required")
	}
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
	return &WriterProducer{w: w}, nil
}

func splitBrokers(csv string) []string {
	var out []string
	for _, p := range strings.Split(csv, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// WriterProducer is a segmentio/kafka-go based implementation of Producer.
type WriterProducer struct {
	w *kafka.Writer
}

func (p *WriterProducer) Send(ctx context.Context, key string, value []byte) error {
	if p == nil || p.w == nil {
		return fmt.Errorf("kafka writer not initialized")
	}
	return p.w.WriteMessages(ctx, kafka.Message{Key: []byte(key), Value: value})
}

func (p *WriterProducer) SendBatch(ctx context.Context, key string, values [][]byte) error {
	if p == nil || p.w == nil {
		return fmt.Errorf("kafka writer not initialized")
	}
	msgs := make([]kafka.Message, 0, len(values))
	k := []byte(key)
	for _, v := range values {
		msgs = append(msgs, kafka.Message{Key: k, Value: v})
	}
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *WriterProducer) Close() error {
	if p == nil || p.w == nil {
		return nil
	}
	return p.w.Close()
}

// NoopProducer is a stub for tests or when Kafka is unavailable (traffic goes to spool).
type NoopProducer struct{}

func (NoopProducer) Send(ctx context.Context, key string, value []byte) error   { return nil }
func (NoopProducer) SendBatch(ctx context.Context, key string, values [][]byte) error { return nil }
func (NoopProducer) Close() error                                               { return nil }

// ErrProducerAlwaysFails is used to simulate Kafka outage in tests.
type ErrProducerAlwaysFails struct{}

func (ErrProducerAlwaysFails) Send(ctx context.Context, key string, value []byte) error {
	return fmt.Errorf("kafka unavailable")
}
func (ErrProducerAlwaysFails) SendBatch(ctx context.Context, key string, values [][]byte) error {
	return fmt.Errorf("kafka unavailable")
}
func (ErrProducerAlwaysFails) Close() error { return nil }
