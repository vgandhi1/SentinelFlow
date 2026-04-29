package config

import (
	"os"
	"runtime"
	"strconv"
)

// Config holds ingestion service configuration (env or defaults).
type Config struct {
	MQTTBroker   string
	MQTTTopic    string
	KafkaBrokers string // empty disables Kafka (NoopProducer); set for history / replay path
	KafkaTopic   string
	NATSURL      string // empty disables JetStream publish to Aegis enrich path
	ChannelCap   int
	NumWorkers   int
	BufferPct80  float64
	BufferPct95  float64
	SpoolDir     string
}

// Load reads config from environment with defaults.
func Load() *Config {
	channelCap := 100000
	if v := os.Getenv("INGEST_CHANNEL_CAP"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			channelCap = n
		}
	}
	workers := runtime.NumCPU() * 2
	if v := os.Getenv("INGEST_NUM_WORKERS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			workers = n
		}
	}
	return &Config{
		MQTTBroker:   getEnv("MQTT_BROKER", "tcp://localhost:1883"),
		MQTTTopic:    getEnv("MQTT_TOPIC", "factory/telemetry/#"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", ""),
		KafkaTopic:   getEnv("KAFKA_TOPIC_RAW", "sentinelflow.telemetry.raw"),
		NATSURL:      getEnv("NATS_URL", ""),
		ChannelCap:   channelCap,
		NumWorkers:   workers,
		BufferPct80:  0.80,
		BufferPct95:  0.95,
		SpoolDir:     getEnv("SPOOL_DIR", "./spool"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
