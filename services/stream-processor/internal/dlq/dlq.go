package dlq

import (
	"context"
	"encoding/json"
	"time"
)

// DLQMessage is the envelope sent to DLQ topics (validation or anomaly).
type DLQMessage struct {
	Reason      string          `json:"reason"`
	Topic       string          `json:"topic"`
	Partition   int32           `json:"partition"`
	Offset      int64           `json:"offset"`
	Timestamp   time.Time       `json:"timestamp"`
	RawPayload  json.RawMessage  `json:"raw_payload"`
	Extra       map[string]any  `json:"extra,omitempty"`
}

// Publisher sends failed messages to DLQ topics.
type Publisher interface {
	PublishValidation(ctx context.Context, rawPayload []byte, reason string, offset int64, partition int32) error
	PublishAnomaly(ctx context.Context, rawPayload []byte, reason string, score float64, offset int64, partition int32) error
}
