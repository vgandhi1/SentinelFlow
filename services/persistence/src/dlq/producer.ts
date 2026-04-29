import { Kafka } from 'kafkajs';
import { config } from '../config';
import type { PersistenceDLQEnvelope } from './envelope';

const kafka = new Kafka({
  clientId: 'sentinelflow-persistence-dlq',
  brokers: config.kafkaBrokers.split(','),
});

const producer = kafka.producer();
let connected = false;

export async function getDLQProducer() {
  if (!connected) {
    await producer.connect();
    connected = true;
  }
  return producer;
}

export async function publishToDLQ(envelope: PersistenceDLQEnvelope): Promise<void> {
  const p = await getDLQProducer();
  const raw = envelope.raw_payload as Record<string, unknown> | null;
  const key = raw && typeof raw.idempotency_key === 'string' ? raw.idempotency_key : undefined;
  const value = Buffer.from(JSON.stringify(envelope), 'utf8');
  await p.send({
    topic: config.dlqTopicPersistence,
    messages: [{ key, value }],
  });
}

export async function closeDLQProducer(): Promise<void> {
  if (connected) {
    await producer.disconnect();
    connected = false;
  }
}
