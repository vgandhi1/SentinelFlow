# Archived repository

**Status:** Archived on 2026-06-28  
**Reason:** Near-duplicate IIoT telemetry reference stack; consolidated under edge-telemetry-plane + aegis-foresight.

## Successors

| Need | Use instead |
|------|-------------|
| Edge → cloud telemetry ingest (MQTT, NATS, Kafka, TimescaleDB) | [**edge-telemetry-plane**](https://github.com/vgandhi1/edge-telemetry-plane) |
| Manufacturing correlation + streaming ML | [**aegis**](https://github.com/vgandhi1/aegis) (local folder: `aegis-foresight`) |
| Factory demo plant (digital twin + OEE + vision) | [**factory-system-AI**](https://github.com/vgandhi1/factory-system-AI) |

## What SentinelFlow covered

SentinelFlow demonstrated MQTT/OPC UA ingestion, Kafka stream processing, GraphQL API, and TimescaleDB storage for IIoT dashboards. **edge-telemetry-plane** covers the same architectural story with a clearer edge/cloud boundary (Rust edge + DETCP protobuf contract) and is the portfolio's preferred ingest-layer reference.

## Local clone

```bash
git clone https://github.com/vgandhi1/edge-telemetry-plane.git
```

This repository is read-only on GitHub.
