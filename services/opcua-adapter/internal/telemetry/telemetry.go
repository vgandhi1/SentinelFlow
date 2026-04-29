package telemetry

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SensorTelemetry matches the JSON shape expected by Stream Processor and Persistence.
type SensorTelemetry struct {
	IdempotencyKey  string            `json:"idempotency_key"`
	TimestampUnixNs int64             `json:"timestamp_unix_ns"`
	SensorID        string            `json:"sensor_id"`
	AssetID         string            `json:"asset_id"`
	MetricName      string            `json:"metric_name"`
	Value           float64           `json:"value"`
	Tags            map[string]string  `json:"tags,omitempty"`
	Critical        bool              `json:"critical"`
}

// NewFromOPCUA builds a SensorTelemetry from OPC UA value and mapping.
func NewFromOPCUA(sensorID, assetID, metricName string, value float64, critical bool) *SensorTelemetry {
	now := time.Now().UTC()
	key := uuid.New().String() + "-" + now.Format("20060102150405.000000000")
	return &SensorTelemetry{
		IdempotencyKey:  key,
		TimestampUnixNs: now.UnixNano(),
		SensorID:        sensorID,
		AssetID:         assetID,
		MetricName:      metricName,
		Value:           value,
		Tags:            map[string]string{"source": "opcua"},
		Critical:        critical,
	}
}

// ToJSON returns JSON bytes for Kafka.
func (s *SensorTelemetry) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}
