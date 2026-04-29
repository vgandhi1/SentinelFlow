// SentinelFlow Stream Processor: consume raw telemetry from Kafka,
// validate, run anomaly detection, produce to validated topic or DLQ.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"sentinelflow/stream-processor/internal/anomaly"
	"sentinelflow/stream-processor/internal/config"
	"sentinelflow/stream-processor/internal/dlq"
	"sentinelflow/stream-processor/internal/kafka"
	"sentinelflow/stream-processor/internal/validation"
)

func main() {
	cfg := config.Load()
	brokers := strings.Split(cfg.KafkaBrokers, ",")

	consumer := kafka.NewConsumer(brokers, cfg.ConsumeTopic, cfg.ConsumerGroup)
	defer consumer.Close()

	producer := kafka.NewProducer(brokers, cfg.ProduceTopic)
	defer producer.Close()

	dlqPub := dlq.NewKafkaPublisher(brokers, cfg.DLQTopicValidation, cfg.DLQTopicAnomaly)
	defer dlqPub.Close()

	detector := &anomaly.Detector{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			msg, payload, err := consumer.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Printf("consumer read error: %v", err)
				if len(msg.Value) > 0 {
					_ = dlqPub.PublishValidation(ctx, msg.Value, "parse_error: "+err.Error(), msg.Offset, msg.Partition)
				}
				continue
			}
			rawPayload := msg.Value

			rec := &validation.Record{
				IdempotencyKey:  payload.IdempotencyKey,
				TimestampUnixNs: payload.TimestampUnixNs,
				SensorID:        payload.SensorID,
				MetricName:      payload.MetricName,
				Value:           payload.Value,
			}
			if err := validation.Validate(rec); err != nil {
				_ = dlqPub.PublishValidation(ctx, rawPayload, err.Error(), msg.Offset, msg.Partition)
				continue
			}

			res := detector.Check(payload.MetricName, payload.Value)
			if res.Anomalous {
				_ = dlqPub.PublishAnomaly(ctx, rawPayload, "anomaly_detected", res.Score, msg.Offset, msg.Partition)
			}

			if err := producer.Produce(ctx, payload.IdempotencyKey, payload); err != nil {
				log.Printf("produce error: %v", err)
				continue
			}
		}
	}()

	log.Printf("stream-processor running (consume=%s produce=%s dlq_validation=%s dlq_anomaly=%s)",
		cfg.ConsumeTopic, cfg.ProduceTopic, cfg.DLQTopicValidation, cfg.DLQTopicAnomaly)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	log.Print("shutdown complete")
}
