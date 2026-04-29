package telemetry

import (
	"encoding/json"
	"errors"
	"fmt"
)

// PLC matches the JSON contract consumed by Aegis correlation-worker (aegis.telemetry.raw).
type PLC struct {
	StationID string  `json:"station_id"`
	Torque    float64 `json:"torque"`
	Timestamp int64   `json:"timestamp"`
}

var (
	errEmptyStation = errors.New("station_id required")
)

// ParsePLCJSON decodes MQTT payload bytes into PLC. Returns an error if JSON is invalid
// or required fields are missing (caller logs generically; payload is not logged).
func ParsePLCJSON(b []byte) (PLC, error) {
	var p PLC
	if err := json.Unmarshal(b, &p); err != nil {
		return PLC{}, fmt.Errorf("json: %w", err)
	}
	if p.StationID == "" {
		return PLC{}, errEmptyStation
	}
	return p, nil
}

// MarshalJSON returns canonical JSON for JetStream / Kafka normalized topics.
func MarshalJSON(p PLC) ([]byte, error) {
	return json.Marshal(p)
}
