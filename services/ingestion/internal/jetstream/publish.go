package jetstream

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	streamName  = "AEGIS"
	subjectPLC  = "aegis.telemetry.raw"
	subjectMask = "aegis.>"
)

// Publisher publishes normalized PLC JSON to JetStream for the Aegis enrich path.
type Publisher struct {
	nc *nats.Conn
	js nats.JetStreamContext
}

// Close drains the NATS connection.
func (p *Publisher) Close() {
	if p == nil || p.nc == nil {
		return
	}
	_ = p.nc.Drain()
}

// ConnectPublisher connects to NATS and returns a publisher after ensuring stream AEGIS exists.
func ConnectPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}
	js, err := nc.JetStream()
	if err != nil {
		_ = nc.Drain()
		return nil, fmt.Errorf("jetstream: %w", err)
	}
	if err := ensureStream(js); err != nil {
		_ = nc.Drain()
		return nil, err
	}
	return &Publisher{nc: nc, js: js}, nil
}

func ensureStream(js nats.JetStreamContext) error {
	_, err := js.StreamInfo(streamName)
	if err == nil {
		return nil
	}
	if err != nats.ErrStreamNotFound {
		return err
	}
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subjectMask},
		Storage:  nats.FileStorage,
		MaxAge:   24 * time.Hour,
	})
	return err
}

// PublishPLC sends one JSON payload to aegis.telemetry.raw.
func (p *Publisher) PublishPLC(data []byte) error {
	if p == nil || p.js == nil {
		return fmt.Errorf("publisher not initialized")
	}
	_, err := p.js.Publish(subjectPLC, data)
	return err
}
