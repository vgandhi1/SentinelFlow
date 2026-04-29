package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

// TelemetryPayload matches the JSON payload from ingestion.
type TelemetryPayload struct {
	IdempotencyKey  string            `json:"idempotency_key"`
	TimestampUnixNs int64             `json:"timestamp_unix_ns"`
	SensorID        string            `json:"sensor_id"`
	AssetID         string            `json:"asset_id"`
	MetricName      string            `json:"metric_name"`
	Value           float64           `json:"value"`
	Tags            map[string]string  `json:"tags,omitempty"`
	Critical        bool              `json:"critical"`
}

// Consumer wraps a Kafka reader for the raw telemetry topic.
type Consumer struct {
	reader *kafka.Reader
}

// NewConsumer creates a consumer for the given brokers and topic.
func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: r}
}

// ReadMessage reads the next message (blocking).
func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, *TelemetryPayload, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return msg, nil, err
	}
	var payload TelemetryPayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return msg, nil, err
	}
	return msg, &payload, nil
}

// Close closes the reader.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
