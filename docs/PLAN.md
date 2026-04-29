# SentinelFlow: System Plan

This document maps the prompt specifications to the repository layout and implementation priorities.

## 1. Goals (from prompts)

- **High-throughput IIoT ingestion** for factory-floor telemetry (100k+ events/sec).
- **Sub-second latency** from sensor to store; backbone for Digital Twin and real-time monitoring.
- **Graceful degradation**: local spooling, Redis fallback, adaptive sampling, circuit breakers.
- **Data integrity**: idempotency keys (UUID + nanosecond timestamp), UPSERT in persistence.

## 2. Architecture Summary

```
[Sensors/PLCs] --MQTT--> [Ingestion (Go)] --Kafka--> [Stream Processor (Go)] --Kafka--> [Persistence (Node/TS)] --> TimescaleDB
                              |                              |                                    |
                         BadgerDB (spool)              Validation + Anomaly                  Redis (cache)
                         Circuit breaker                     detection
```

- **Ingestion:** MQTT/Protobuf termination, worker pool, buffered channel, batch to Kafka. Backpressure: 80% load shedding, 95% NACK. Circuit breaker → local spool.
- **Stream Processor:** Consumes from Kafka, validates, runs anomaly detection, produces to downstream topic(s).
- **Persistence:** Consumes from Kafka, UPSERTs to TimescaleDB, serves last 5 min from Redis on DB lag.

## 3. Repository Layout

```
SentinelFlow/
├── README.md
├── docker-compose.yml
├── .gitignore
├── docs/
│   ├── PLAN.md                 # This file
│   ├── DESIGN_INGESTION.md     # Ingestion technical design
│   ├── RELIABILITY.md          # Failure modes, circuit breakers, SLAs
│   ├── CONTRIBUTING.md         # Fork/branch, lint, tests, PR policy
│   ├── DLQ.md                  # Dead letter queues: retention, replay
│   └── DESIGN_OPCUA.md         # OPC UA adapter design
├── api/
│   └── proto/
│       ├── telemetry.proto     # Sensor event schema
│       └── README.md           # How to generate stubs
├── services/
│   ├── ingestion/              # Go – edge MQTT/Protobuf → Kafka
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── internal/
│   │   │   ├── mqtt/
│   │   │   ├── worker/
│   │   │   ├── kafka/
│   │   │   ├── backpressure/
│   │   │   ├── spool/         # BadgerDB local replay
│   │   │   └── config/
│   │   └── Dockerfile
│   ├── stream-processor/       # Go – validation + anomaly
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── internal/
│   │   │   ├── consumer/
│   │   │   ├── validation/
│   │   │   ├── anomaly/
│   │   │   └── config/
│   │   └── Dockerfile
│   ├── persistence/           # Node.js/TS – Kafka → TimescaleDB + Redis
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── src/
│   │   │   ├── index.ts
│   │   │   ├── consumer/
│   │   │   ├── db/
│   │   │   ├── cache/
│   │   │   ├── config/
│   │   │   └── dlq/
│   │   └── Dockerfile
│   ├── opcua-adapter/         # Go – OPC UA → Kafka raw topic
│   │   ├── go.mod
│   │   ├── main.go
│   │   ├── internal/
│   │   │   ├── config/
│   │   │   ├── kafka/
│   │   │   ├── opcua/
│   │   │   └── telemetry/
│   │   └── Dockerfile
│   └── graphql-api/          # Node.js/TS – GraphQL over TimescaleDB
│       ├── package.json
│       ├── tsconfig.json
│       ├── src/
│       │   ├── index.ts
│       │   ├── schema.ts
│       │   ├── resolvers.ts
│       │   ├── db.ts
│       │   └── config.ts
│       └── Dockerfile
├── .github/
│   └── workflows/
│       └── ci.yml             # Lint + tests (Go + TS)
└── prompts/                   # Source prompts (readme, design, contributing, reliability)
```

## 4. Implementation Priorities

| Priority | Component            | Deliverable |
|----------|----------------------|-------------|
| P0       | Plan + docs          | PLAN.md, DESIGN_INGESTION.md, RELIABILITY.md, CONTRIBUTING.md |
| P0       | Root + infra         | README, docker-compose, .gitignore |
| P0       | Proto contract       | telemetry.proto + generation notes |
| P1       | Ingestion            | Worker pool, MQTT listener, Kafka batch producer, backpressure thresholds, idempotency in payload |
| P1       | Persistence          | Kafka consumer, TimescaleDB writer (UPSERT), Redis fallback wiring |
| P2       | Stream processor     | Kafka consumer, validation, anomaly detection stub |
| P2       | Reliability          | BadgerDB spool + replay, circuit breaker (sony/gobreaker), MQTT NACK |
| P3       | Observability        | Prometheus metrics, Grafana placeholders |
| P3       | K8s/Helm             | Charts and deployment manifests |

## 5. Tech Stack (from prompts)

- **Languages:** Go 1.21+, Node.js 20, TypeScript.
- **Comms:** gRPC (optional for control plane), Protobuf (telemetry), Kafka, MQTT.
- **Storage:** TimescaleDB (PostgreSQL), Redis, BadgerDB (local spool).
- **Infra:** Docker, docker-compose; later Kubernetes, Helm, Prometheus/Grafana.

## 6. Quality Gates (from contributing prompt)

- Go: `golangci-lint run`; tests with `testify/assert`; benchmarks for ingestion worker changes.
- TS: `npm run lint`; unit tests for business logic.
- Integration tests for any service touching Kafka or Postgres.
- PR: one approval; docs in `/docs` on API contract changes; CI 100% pass.

---

*Generated from prompts in `prompts/*.md`. Update this plan as the system evolves.*
