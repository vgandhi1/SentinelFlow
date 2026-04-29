# SentinelFlow Roadmap

Planned enhancements to support reliability, broader industrial protocols, and frontend consumption.

---

## [x] Dead Letter Queues (DLQ) for failed Kafka messages

**Goal:** Capture messages that fail validation, processing, or persistence so they can be inspected, corrected, and replayed instead of being dropped.

**Scope:**
- **Stream Processor:** Publish invalid or anomalous records (or raw payloads) to a dedicated DLQ topic (e.g. `sentinelflow.dlq.validation`, `sentinelflow.dlq.anomaly`) with metadata (reason, offset, timestamp).
- **Persistence:** On UPSERT or Redis errors, write the message to a DLQ topic (e.g. `sentinelflow.dlq.persistence`) with error code and payload.
- **Operational:** Optional DLQ consumer for alerting, dashboards, and manual replay; document retention and replay procedure in `/docs`.

**References:** `services/stream-processor` (validation, anomaly), `services/persistence` (consumer, db).

---

## [x] OPC UA protocol translation at the edge

**Goal:** Ingest telemetry from OPC UA servers (common in factory automation) alongside MQTT, and normalize to the same Protobuf/Kafka pipeline.

**Scope:**
- **Edge component:** New adapter or sub-service that connects to OPC UA endpoints (subscriptions or polling), maps nodes/variables to `SensorTelemetry` fields, and publishes to the same internal channel or directly to Kafka (or via existing Ingestion API).
- **Configuration:** Node ID mapping, sampling intervals, and security (certificates, username/password) configurable per OPC UA server.
- **Placement:** “At the edge” may mean a separate binary (e.g. `services/opcua-adapter`) or an optional module within the Ingestion service; document in `docs/DESIGN_INGESTION.md` or a new `docs/DESIGN_OPCUA.md`.

**References:** `api/proto/telemetry.proto`, `services/ingestion`, new `services/opcua-adapter` (or equivalent).

---

## [x] GraphQL layer for complex frontend queries

**Goal:** Expose a GraphQL API so frontends (Digital Twin, dashboards) can request flexible, nested queries over telemetry and metadata without many REST endpoints.

**Scope:**
- **New service or layer:** GraphQL API (Node/TS or Go) that resolves queries against TimescaleDB (and Redis for hot path); support time ranges, aggregation, filtering by sensor/asset/metric, and optional relations (assets, sites).
- **Schema:** Define types (e.g. `TelemetryPoint`, `Sensor`, `Asset`, `TimeSeries`) and queries (e.g. `telemetry(sensorId, from, to)`, `sensors(assetId)`); document in `/docs` and keep schema in version control.
- **Performance:** Use DataLoader-style batching and consider query complexity limits; cache frequently used shapes in Redis where it fits the existing fallback strategy.

**References:** `services/persistence` (DB schema, Redis), new `services/graphql-api` or `services/api`; `docs/PLAN.md` (add to layout).

---

All three items above are implemented. See `services/stream-processor` (DLQ validation/anomaly), `services/persistence` (DLQ persistence), `docs/DLQ.md`; `services/opcua-adapter` and `docs/DESIGN_OPCUA.md`; `services/graphql-api` (GraphQL over TimescaleDB).
