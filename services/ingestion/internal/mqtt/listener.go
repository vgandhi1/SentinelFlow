package mqtt

import (
	"context"
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Listener subscribes to MQTT topics and pushes raw payloads onto a channel.
type Listener struct {
	client mqtt.Client
	topic  string
	out    chan<- []byte
	onNACK func()
}

// NewListener creates an MQTT listener that forwards messages to out.
func NewListener(broker, topic string, out chan<- []byte, onNACK func()) (*Listener, error) {
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("sentinelflow-ingestion")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	l := &Listener{client: client, topic: topic, out: out, onNACK: onNACK}
	return l, nil
}

// Run subscribes and forwards messages until ctx is cancelled.
func (l *Listener) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	token := l.client.Subscribe(l.topic, 1, func(_ mqtt.Client, msg mqtt.Message) {
		payload := msg.Payload()
		select {
		case l.out <- payload:
			msg.Ack()
		case <-ctx.Done():
			msg.Ack()
			return
		default:
			if l.onNACK != nil {
				l.onNACK()
			}
			// paho v1.4 Message has no Nack(); Ack releases the broker slot while this layer drops under pressure.
			msg.Ack()
		}
	})
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("subscribe: %w", token.Error())
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		l.client.Disconnect(250)
	}()
	wg.Wait()
	return nil
}
