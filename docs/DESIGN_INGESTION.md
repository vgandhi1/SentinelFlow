# Technical Design: High-Concurrency Ingestion Service

## 1. Problem Statement

The factory floor generates bursts of telemetry (vibration, torque, heat) that can overwhelm standard RESTful endpoints. We need a system that can absorb these spikes without crashing or losing data.

## 2. Architecture: The Worker Pool Model

To handle 100k+ events/sec, we utilize a non-blocking fan-in pattern.

### Concurrency Flow

1. **Listeners:** Multiple Goroutines listen on MQTT topics.
2. **Internal Channel:** Messages are pushed into a **Buffered Channel**. This acts as our primary in-memory shock absorber.
3. **Worker Pool:** A fixed set of workers (based on `runtime.NumCPU()`) pulls from the channel, performs Protobuf unmarshaling, and batches the data for Kafka.

## 3. Backpressure Management

When the Kafka broker experiences high latency (e.g. during a partition rebalance), the system triggers reactive backpressure:

- **Threshold 1 (80% Buffer):** Service begins "Load Shedding" by dropping non-critical diagnostic telemetry.
- **Threshold 2 (95% Buffer):** Service issues `NACK` (Negative Acknowledgment) to MQTT clients, forcing sensors to retry/buffer at the edge.

## 4. Data Integrity

We utilize **Idempotency Keys** (UUID + Nanosecond Timestamp) generated at the sensor level. Even if the network retries a packet, the persistence layer uses a "UPSERT" logic to ensure no duplicate records exist in the time-series database.
