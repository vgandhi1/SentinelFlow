function getEnv(key: string, defaultValue: string): string {
  return process.env[key] ?? defaultValue;
}

export const config = {
  kafkaBrokers: getEnv('KAFKA_BROKERS', 'localhost:9092'),
  kafkaTopic: getEnv('KAFKA_TOPIC_VALIDATED', 'sentinelflow.telemetry.validated'),
  consumerGroup: getEnv('KAFKA_CONSUMER_GROUP', 'persistence'),
  dlqTopicPersistence: getEnv('KAFKA_DLQ_TOPIC_PERSISTENCE', 'sentinelflow.dlq.persistence'),
  pgConnectionString: getEnv(
    'PG_URL',
    'postgres://sentinel:sentinel@localhost:5432/sentinelflow'
  ),
  redisUrl: getEnv('REDIS_URL', 'redis://localhost:6379'),
  redisFallbackTtlSeconds: 300,
};
