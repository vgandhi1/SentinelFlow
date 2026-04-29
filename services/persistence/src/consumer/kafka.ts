import { Kafka } from 'kafkajs';
import { config } from '../config';

const kafka = new Kafka({
  clientId: 'sentinelflow-persistence',
  brokers: config.kafkaBrokers.split(','),
});

export async function createConsumer() {
  const consumer = kafka.consumer({ groupId: config.consumerGroup });
  await consumer.connect();
  await consumer.subscribe({ topic: config.kafkaTopic, fromBeginning: true });
  return consumer;
}

export type TelemetryPayload = {
  idempotency_key: string;
  timestamp_unix_ns: number;
  sensor_id: string;
  asset_id: string;
  metric_name: string;
  value: number;
  tags?: Record<string, string>;
  critical?: boolean;
};

export function parsePayload(raw: Buffer): TelemetryPayload | null {
  try {
    return JSON.parse(raw.toString()) as TelemetryPayload;
  } catch {
    return null;
  }
}
