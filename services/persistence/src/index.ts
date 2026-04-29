import { createConsumer, parsePayload } from './consumer/kafka';
import { getPool, initSchema, upsertTelemetry, TelemetryRow } from './db/timescale';
import { setFallback } from './cache/redis';
import { buildDLQEnvelope } from './dlq/envelope';
import { publishToDLQ, closeDLQProducer } from './dlq/producer';
import { config } from './config';

async function main(): Promise<void> {
  const consumer = await createConsumer();
  const pool = getPool();
  const client = await pool.connect();
  await initSchema(client);

  await consumer.run({
    eachMessage: async ({ message }) => {
      const rawValue = message.value;
      const payload = rawValue ? parsePayload(rawValue) : null;
      if (!payload) return;
      const row: TelemetryRow = {
        idempotency_key: payload.idempotency_key,
        time: new Date(Number(payload.timestamp_unix_ns) / 1e6),
        sensor_id: payload.sensor_id,
        asset_id: payload.asset_id,
        metric_name: payload.metric_name,
        value: payload.value,
        tags: payload.tags,
        critical: payload.critical ?? false,
      };
      try {
        await upsertTelemetry(client, row);
        await setFallback(payload.idempotency_key, JSON.stringify(row));
      } catch (err) {
        const envelope = buildDLQEnvelope(
          err instanceof Error ? err.message : String(err),
          payload,
          err instanceof Error ? err.name : undefined
        );
        try {
          await publishToDLQ(envelope);
        } catch (dlqErr) {
          console.error('DLQ publish failed:', dlqErr);
        }
        console.error('persistence error:', err);
      }
    },
  });

  console.log(`persistence consuming ${config.kafkaTopic} (DLQ: ${config.dlqTopicPersistence})`);
}

process.on('SIGTERM', async () => {
  await closeDLQProducer();
  process.exit(0);
});

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
