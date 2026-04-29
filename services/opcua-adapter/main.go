package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"sentinelflow/opcua-adapter/internal/config"
	"sentinelflow/opcua-adapter/internal/kafka"
	"sentinelflow/opcua-adapter/internal/opcua"
	"sentinelflow/opcua-adapter/internal/telemetry"
)

func main() {
	cfg := config.Load()
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	producer := kafka.NewProducer(brokers, cfg.KafkaTopicRaw)
	defer producer.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	valCh := make(chan opcua.Value, 64)
	reader := opcua.NewReader(cfg)
	go func() {
		if err := reader.Run(ctx, valCh); err != nil {
			log.Printf("opcua reader: %v", err)
		}
	}()
	go func() {
		for v := range valCh {
			tel := telemetry.NewFromOPCUA(v.SensorID, v.AssetID, v.MetricName, v.Value, v.Critical)
			jsonBytes, err := tel.ToJSON()
			if err != nil {
				log.Printf("telemetry marshal: %v", err)
				continue
			}
			if err := producer.Produce(ctx, tel.IdempotencyKey, jsonBytes); err != nil {
				log.Printf("kafka produce: %v", err)
			}
		}
	}()
	log.Printf("opcua-adapter running (endpoint=%s topic=%s)", cfg.Endpoint, cfg.KafkaTopicRaw)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	cancel()
	log.Print("shutdown complete")
}
