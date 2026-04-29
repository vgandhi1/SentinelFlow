# Dead Letter Queues (DLQ)

Failed Kafka messages are published to dedicated DLQ topics so they can be inspected, corrected, and replayed instead of being dropped.

## Topics

| Topic | Source | When |
|-------|--------|------|
| `sentinelflow.dlq.validation` | Stream Processor | Parse error or validation failure (missing idempotency_key, invalid timestamp, etc.) |
| `sentinelflow.dlq.anomaly` | Stream Processor | Record marked anomalous by anomaly detector |
| `sentinelflow.dlq.persistence` | Persistence | UPSERT or Redis write error |

Override via env: `KAFKA_DLQ_TOPIC_VALIDATION`, `KAFKA_DLQ_TOPIC_ANOMALY`, `KAFKA_DLQ_TOPIC_PERSISTENCE`.

## Message formats

- **Validation / Anomaly (Stream Processor):** JSON envelope with `reason`, `topic`, `partition`, `offset`, `timestamp`, `raw_payload` (original message), and optional `extra` (e.g. `anomaly_score`).
- **Persistence:** JSON envelope with `reason`, `error_code`, `timestamp`, `raw_payload` (original telemetry payload).

## Retention

Configure retention per DLQ topic so old failures don’t fill the cluster:

```properties
# Example: 7 days, 1GB max
retention.ms=604800000
retention.bytes=1073741824
```

Create topics with these settings or use Kafka broker defaults and monitor size.

## Replay procedure

1. **Consume** from the DLQ topic (e.g. with a one-off script or Kafka Console Consumer).
2. **Inspect** `reason` and `raw_payload`; fix data or environment (e.g. DB connectivity).
3. **Re-publish** `raw_payload` to the original input topic:
   - Validation DLQ → `sentinelflow.telemetry.raw`
   - Anomaly DLQ → optionally back to raw or to a separate “reviewed” topic
   - Persistence DLQ → `sentinelflow.telemetry.validated`
4. **Delete** or skip the DLQ message after successful replay to avoid double processing.

## Alerting

Consume DLQ topics in a small service or pipeline that increments metrics (e.g. Prometheus counter per topic/reason) and optionally sends alerts when rate or count exceeds a threshold.
