package config

import "os"

// Config holds stream-processor configuration.
type Config struct {
	KafkaBrokers    string
	ConsumeTopic    string
	ProduceTopic    string
	ConsumerGroup   string
	DLQTopicValidation string
	DLQTopicAnomaly    string
}

// Load reads config from environment.
func Load() *Config {
	return &Config{
		KafkaBrokers:        getEnv("KAFKA_BROKERS", "localhost:9092"),
		ConsumeTopic:        getEnv("KAFKA_TOPIC_RAW", "sentinelflow.telemetry.raw"),
		ProduceTopic:        getEnv("KAFKA_TOPIC_VALIDATED", "sentinelflow.telemetry.validated"),
		ConsumerGroup:       getEnv("KAFKA_CONSUMER_GROUP", "stream-processor"),
		DLQTopicValidation:  getEnv("KAFKA_DLQ_TOPIC_VALIDATION", "sentinelflow.dlq.validation"),
		DLQTopicAnomaly:     getEnv("KAFKA_DLQ_TOPIC_ANOMALY", "sentinelflow.dlq.anomaly"),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
