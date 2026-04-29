package backpressure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolicy_Decide(t *testing.T) {
	policy := &Policy{Cap: 100, Threshold80: 0.80, Threshold95: 0.95}

	assert.Equal(t, Accept, policy.Decide(50, false))
	assert.Equal(t, Accept, policy.Decide(79, false))
	assert.Equal(t, Drop, policy.Decide(80, false))
	assert.Equal(t, Accept, policy.Decide(80, true))
	assert.Equal(t, NACK, policy.Decide(95, true))
	assert.Equal(t, NACK, policy.Decide(95, false))
}
