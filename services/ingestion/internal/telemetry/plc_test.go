package telemetry

import (
	"testing"
)

func TestParsePLCJSON(t *testing.T) {
	t.Parallel()
	p, err := ParsePLCJSON([]byte(`{"station_id":"5","torque":1.5,"timestamp":1700000000}`))
	if err != nil {
		t.Fatal(err)
	}
	if p.StationID != "5" || p.Torque != 1.5 || p.Timestamp != 1700000000 {
		t.Fatalf("unexpected: %+v", p)
	}
}

func TestParsePLCJSON_Invalid(t *testing.T) {
	t.Parallel()
	_, err := ParsePLCJSON([]byte(`{}`))
	if err == nil {
		t.Fatal("expected error")
	}
}
