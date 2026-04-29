# SentinelFlow: High-Throughput IIoT Ingestion Engine

SentinelFlow is a distributed, industrial-grade backend platform designed to ingest and process high-frequency telemetry data from thousands of factory-floor sensors. It serves as the data backbone for real-time monitoring and "Digital Twin" visualizations in manufacturing environments.



## 🏗 System Overview
The platform is designed to handle the "High Velocity" requirement of modern automotive factories, ensuring that data from PLCs (Programmable Logic Controllers) is captured, validated, and stored with sub-second latency.

### Core Modules
* **Ingestion (Go):** High-concurrency edge service for MQTT/Protobuf termination.
* **Stream Processor (Go):** Real-time validation and anomaly detection.
* **Persistence (Node.js/TS):** Distributed consumer optimized for Time-Series (TimescaleDB).

## 🛠 Tech Stack
* **Languages:** Golang 1.21+, TypeScript (Node.js 20).
* **Communication:** gRPC, Protobuf, Kafka, MQTT.
* **Databases:** TimescaleDB (PostgreSQL), Redis (Cache).
* **Infrastructure:** Kubernetes, Docker, Helm, Prometheus/Grafana.

## 🚀 Quick Start
```bash
docker-compose up -d  # Start Kafka & TimescaleDB
cd services/ingestion && go run main.go