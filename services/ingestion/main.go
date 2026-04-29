// SentinelFlow Ingestion: high-concurrency MQTT edge service → Kafka (raw history) and/or NATS JetStream (Aegis enrich).
// MQTT payloads for the enrich path must be JSON: {"station_id","torque","timestamp"} matching Aegis PLCMessage.
// See docs/DESIGN_INGESTION.md for architecture.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sentinelflow/ingestion/internal/backpressure"
	"sentinelflow/ingestion/internal/config"
	"sentinelflow/ingestion/internal/jetstream"
	"sentinelflow/ingestion/internal/kafka"
	"sentinelflow/ingestion/internal/mqtt"
	"sentinelflow/ingestion/internal/spool"
	"sentinelflow/ingestion/internal/telemetry"
	"sentinelflow/ingestion/internal/worker"
)

func main() {
	cfg := config.Load()
	ch := make(chan worker.RawMessage, cfg.ChannelCap)
	policy := &backpressure.Policy{
		Cap:         cfg.ChannelCap,
		Threshold80: cfg.BufferPct80,
		Threshold95: cfg.BufferPct95,
	}
	_ = policy

	var prod kafka.Producer = kafka.NoopProducer{}
	if cfg.KafkaBrokers != "" {
		p, err := kafka.NewWriterProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
		if err != nil {
			log.Fatalf("kafka: %v", err)
		}
		prod = p
		defer func() { _ = p.Close() }()
	}

	var jsPub *jetstream.Publisher
	if cfg.NATSURL != "" {
		p, err := jetstream.ConnectPublisher(cfg.NATSURL)
		if err != nil {
			log.Fatalf("nats jetstream: %v", err)
		}
		jsPub = p
		defer jsPub.Close()
	}

	spl, err := spool.NewSpooler(cfg.SpoolDir)
	if err != nil {
		log.Fatalf("spool: %v", err)
	}
	defer spl.Close()

	handler := newBridgeHandler(prod, jsPub, spl)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker.Pool(ctx, ch, handler, cfg.NumWorkers, 100)

	byteChan := make(chan []byte, cfg.ChannelCap)
	go func() {
		for b := range byteChan {
			ch <- worker.RawMessage{Payload: b}
		}
	}()
	listener, err := mqtt.NewListener(cfg.MQTTBroker, cfg.MQTTTopic, byteChan, func() {
		log.Print("backpressure: NACK sent to MQTT client")
	})
	if err != nil {
		log.Fatalf("mqtt: %v", err)
	}

	go func() {
		if err := listener.Run(ctx); err != nil {
			log.Printf("mqtt listener: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	log.Print("shutdown complete")
}

type bridgeHandler struct {
	prod  kafka.Producer
	nats  *jetstream.Publisher
	spool *spool.Spooler
}

func newBridgeHandler(prod kafka.Producer, nats *jetstream.Publisher, spool *spool.Spooler) *bridgeHandler {
	return &bridgeHandler{prod: prod, nats: nats, spool: spool}
}

func (h *bridgeHandler) HandleBatch(ctx context.Context, batch []worker.RawMessage) error {
	if len(batch) == 0 {
		return nil
	}
	values := make([][]byte, len(batch))
	var skippedEnrich int
	for i := range batch {
		payload := batch[i].Payload
		values[i] = payload
		if h.nats != nil {
			plc, err := telemetry.ParsePLCJSON(payload)
			if err != nil {
				skippedEnrich++
				continue
			}
			j, mErr := telemetry.MarshalJSON(plc)
			if mErr != nil {
				skippedEnrich++
				continue
			}
			if err := h.nats.PublishPLC(j); err != nil {
				log.Printf("jetstream publish failed: %v", err)
			}
		}
	}
	if skippedEnrich > 0 {
		log.Printf("ingestion: skipped %d MQTT message(s) for enrich path (not valid PLC JSON)", skippedEnrich)
	}
	err := h.prod.SendBatch(ctx, "ingestion", values)
	if err != nil {
		for i := range batch {
			_ = h.spool.Write(ctx, string(batch[i].Payload), batch[i].Payload)
		}
		return err
	}
	return nil
}
