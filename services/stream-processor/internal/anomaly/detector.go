package anomaly

// Detector marks telemetry as normal or anomalous. Stub for real logic (e.g. threshold, ML).
type Detector struct{}

// Result is the outcome of anomaly check.
type Result struct {
	Anomalous bool
	Score     float64
}

// Check returns anomaly result for the given value and metric. Stub: always normal.
func (d *Detector) Check(metricName string, value float64) Result {
	return Result{Anomalous: false, Score: 0}
}
