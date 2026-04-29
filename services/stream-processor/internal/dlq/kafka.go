package dlq

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// KafkaPublisher publishes DLQ messages to Kafka topics.
type KafkaPublisher struct {
	validationWriter *kafka.Writer
	anomalyWriter   *kafka.Writer
}

// NewKafkaPublisher creates a DLQ publisher for the given brokers and topic names.
func NewKafkaPublisher(brokers []string, validationTopic, anomalyTopic string) *KafkaPublisher {
	return &KafkaPublisher{
		validationWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    validationTopic,
			Balancer: &kafka.LeastBytes{},
		},
		anomalyWriter: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    anomalyTopic,
			Balancer: &kafka.LeastBytes{},
		},
	}
}

// PublishValidation sends a validation failure to the DLQ.
func (p *KafkaPublisher) PublishValidation(ctx context.Context, rawPayload []byte, reason string, offset int64, partition int32) error {
	msg := DLQMessage{
		Reason:     reason,
		Topic:      "",
		Partition:  partition,
		Offset:     offset,
		Timestamp:  time.Now().UTC(),
		RawPayload: rawPayload,
	}
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.validationWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(offset, 10)),
		Value: val,
	})
}

// PublishAnomaly sends an anomaly detection result to the DLQ.
func (p *KafkaPublisher) PublishAnomaly(ctx context.Context, rawPayload []byte, reason string, score float64, offset int64, partition int32) error {
	msg := DLQMessage{
		Reason:     reason,
		Topic:      "",
		Partition:  partition,
		Offset:     offset,
		Timestamp:  time.Now().UTC(),
		RawPayload: rawPayload,
		Extra:      map[string]any{"anomaly_score": score},
	}
	val, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.anomalyWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(offset, 10)),
		Value: val,
	})
}

// Close closes both writers.
func (p *KafkaPublisher) Close() error {
	if err := p.validationWriter.Close(); err != nil {
		return err
	}
	return p.anomalyWriter.Close()
}
