package validation

// Record is a minimal telemetry record for validation (replace with proto when generated).
type Record struct {
	IdempotencyKey  string
	TimestampUnixNs int64
	SensorID        string
	MetricName      string
	Value           float64
}

// Validate checks basic invariants. Returns nil if valid.
func Validate(r *Record) error {
	if r.IdempotencyKey == "" {
		return ErrMissingIdempotencyKey
	}
	if r.TimestampUnixNs <= 0 {
		return ErrInvalidTimestamp
	}
	if r.SensorID == "" || r.MetricName == "" {
		return ErrMissingFields
	}
	return nil
}

// Validation errors
var (
	ErrMissingIdempotencyKey = &ValidationError{Msg: "missing idempotency_key"}
	ErrInvalidTimestamp     = &ValidationError{Msg: "invalid timestamp_unix_ns"}
	ErrMissingFields         = &ValidationError{Msg: "missing sensor_id or metric_name"}
)

// ValidationError is a validation failure.
type ValidationError struct {
	Msg string
}

func (e *ValidationError) Error() string { return e.Msg }
