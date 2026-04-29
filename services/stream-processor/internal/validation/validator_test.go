package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	valid := &Record{
		IdempotencyKey:  "uuid-ts",
		TimestampUnixNs: 1e9,
		SensorID:        "s1",
		MetricName:      "temperature",
		Value:           42.0,
	}
	assert.NoError(t, Validate(valid))

	assert.Error(t, Validate(&Record{IdempotencyKey: ""}))
	assert.Error(t, Validate(&Record{IdempotencyKey: "x", TimestampUnixNs: 0}))
	assert.Error(t, Validate(&Record{IdempotencyKey: "x", TimestampUnixNs: 1, SensorID: ""}))
}
