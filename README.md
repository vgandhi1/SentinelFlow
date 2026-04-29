# SentinelFlow: High-Throughput IIoT Ingestion Engine

SentinelFlow is a distributed, industrial-grade backend platform designed to ingest and process high-frequency telemetry data from thousands of factory-floor sensors. It serves as the data backbone for real-time monitoring and "Digital Twin" visualizations in manufacturing environments.

## System Overview

The platform handles the **High Velocity** requirement of modern automotive factories, ensuring that data from PLCs (Programmable Logic Controllers) is captured, validated, and stored with sub-second latency.

### Core Modules

| Module | Language | Role |
|--------|----------|------|
| **Ingestion** | Go | High-concurrency edge service for MQTT/Protobuf termination → Kafka |
| **Stream Processor** | Go | Real-time validation and anomaly detection |
| **Persistence** | Node.js/TS | Distributed consumer optimized for Time-Series (TimescaleDB) |

### Tech Stack

- **Languages:** Golang 1.21+, TypeScript (Node.js 20)
- **Communication:** gRPC, Protobuf, Kafka, MQTT
- **Databases:** TimescaleDB (PostgreSQL), Redis (cache), BadgerDB (local spool)
- **Infrastructure:** Kubernetes, Docker, Helm, Prometheus/Grafana

## Quick Start

```bash
docker-compose up -d   # Start core dependencies + observability stack
cd services/ingestion && go run main.go
```

See [docs/PLAN.md](docs/PLAN.md) for full architecture and [docs/DESIGN_INGESTION.md](docs/DESIGN_INGESTION.md) for ingestion design.

## Observability

Prometheus/Grafana is included in local Compose for out-of-the-box monitoring:

- **Prometheus:** [http://localhost:9090](http://localhost:9090)
- **Grafana:** [http://localhost:3000](http://localhost:3000) (default login: `admin` / `admin`)
- **Exporters:** Kafka (`9308`), Redis (`9121`), TimescaleDB (`9187`), node exporter (`9100`), cAdvisor (`8080`)

Provisioned config lives under `deploy/observability/`:

- Prometheus scrape/rules config: `deploy/observability/prometheus.yml`, `deploy/observability/rules/`
- Grafana provisioning + starter dashboard: `deploy/observability/grafana/`

Optional app targets (`host.docker.internal:2112-2116`) are pre-listed in Prometheus config. Expose `/metrics` on those ports in each service to add application-level telemetry.

## Repository Layout

- `api/proto/` — Protobuf definitions for telemetry
- `services/ingestion/` — Go MQTT/Protobuf → Kafka edge service
- `services/stream-processor/` — Go validation, anomaly detection, DLQ
- `services/persistence/` — Node.js/TS Kafka → TimescaleDB + Redis + DLQ
- `services/opcua-adapter/` — Go OPC UA → Kafka raw topic
- `services/graphql-api/` — GraphQL API over TimescaleDB (Digital Twin / dashboards)
- `docs/` — Design, reliability, DLQ, OPC UA, contributing

## Development

- **Lint:** Go: `golangci-lint run`; TS: `npm run lint`
- **Tests:** Unit + integration for Kafka/Postgres; see [CONTRIBUTING](docs/CONTRIBUTING.md)

## Reliability

Designed for graceful degradation: local spooling on Kafka outage, Redis fallback on DB lag, adaptive sampling and circuit breakers. See [docs/RELIABILITY.md](docs/RELIABILITY.md).

## Roadmap

Implemented: **DLQ** (validation, anomaly, persistence topics), **OPC UA adapter** (poll nodes → Kafka raw), **GraphQL API** (telemetry, sensors, assets). See [docs/DLQ.md](docs/DLQ.md), [docs/DESIGN_OPCUA.md](docs/DESIGN_OPCUA.md), [docs/ROADMAP.md](docs/ROADMAP.md).
