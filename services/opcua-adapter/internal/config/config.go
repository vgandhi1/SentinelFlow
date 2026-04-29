package config

import (
	"fmt"
	"os"
	"strings"
)

// NodeMapping maps one OPC UA node to telemetry fields.
type NodeMapping struct {
	NodeID     string // e.g. "ns=2;s=Machine1.Temperature"
	SensorID   string
	AssetID    string
	MetricName string
}

// Config holds OPC UA adapter configuration.
type Config struct {
	Endpoint        string
	KafkaBrokers    string
	KafkaTopicRaw   string
	PublishInterval int   // ms
	SamplingInterval int  // ms
	NodeMappings    []NodeMapping
}

// Load reads config from environment (node mappings can be extended from file later).
func Load() *Config {
	return &Config{
		Endpoint:         getEnv("OPCUA_ENDPOINT", "opc.tcp://localhost:4840"),
		KafkaBrokers:     getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaTopicRaw:    getEnv("KAFKA_TOPIC_RAW", "sentinelflow.telemetry.raw"),
		PublishInterval:   getEnvInt("OPCUA_PUBLISH_INTERVAL_MS", 500),
		SamplingInterval: getEnvInt("OPCUA_SAMPLING_INTERVAL_MS", 500),
		NodeMappings:     loadNodeMappings(), // from env OPCUA_NODES or default
	}
}

func loadNodeMappings() []NodeMapping {
	// Format: OPCUA_NODES="nodeId1,sensor1,asset1,metric1;nodeId2,sensor2,asset2,metric2"
	s := getEnv("OPCUA_NODES", "")
	if s == "" {
		return nil
	}
	var out []NodeMapping
	for _, part := range strings.Split(s, ";") {
		fields := strings.SplitN(part, ",", 4)
		if len(fields) != 4 {
			continue
		}
		out = append(out, NodeMapping{
			NodeID:     strings.TrimSpace(fields[0]),
			SensorID:   strings.TrimSpace(fields[1]),
			AssetID:    strings.TrimSpace(fields[2]),
			MetricName: strings.TrimSpace(fields[3]),
		})
	}
	return out
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return defaultVal
}
