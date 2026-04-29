# Reliability and Failure Modes

In a manufacturing environment, software failure is not an option. We design for "Graceful Degradation."

## 1. Failure Scenarios

| Scenario | Impact | Mitigation Strategy |
| :--- | :--- | :--- |
| **Kafka Outage** | Data cannot be streamed. | **Local Spooling:** Workers write data to a local BadgerDB (KV store) and replay it once Kafka is online. |
| **Database Lag** | Real-time dashboards freeze. | **Redis Fallback:** The API serves the last 5 minutes of data directly from Redis cache. |
| **Sensor 'Noise'** | Network congestion. | **Adaptive Sampling:** The Ingestion service instructs sensors to reduce frequency via an MQTT Control Topic. |

## 2. Circuit Breakers
We use the `sony/gobreaker` library for downstream dependencies. If the Kafka producer fails more than 5 times in 10 seconds, the circuit opens, and we divert traffic to the Local Spooler immediately to prevent memory exhaustion.

## 3. Monitoring (SLAs)
* **Availability:** 99.9%
* **P99 Latency:** < 200ms (Sensor to Kafka)