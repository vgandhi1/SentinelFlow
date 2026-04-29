package backpressure

// Policy decides whether to accept, drop (load shed), or NACK based on buffer usage.
type Policy struct {
	Cap         int
	Threshold80 float64
	Threshold95 float64
}

// Decision is the action to take for an incoming message.
type Decision int

const (
	Accept Decision = iota
	Drop
	NACK
)

// Decide returns Accept, Drop, or NACK given current channel length and whether the event is critical.
func (p *Policy) Decide(currentLen int, critical bool) Decision {
	if p.Cap <= 0 {
		return Accept
	}
	usage := float64(currentLen) / float64(p.Cap)
	if usage >= p.Threshold95 {
		return NACK
	}
	if usage >= p.Threshold80 && !critical {
		return Drop
	}
	return Accept
}
