import Redis from 'ioredis';
import { config } from '../config';

let redis: Redis | null = null;

export function getRedis(): Redis {
  if (!redis) {
    redis = new Redis(config.redisUrl);
  }
  return redis;
}

const FALLBACK_PREFIX = 'sentinelflow:fallback:';
const TTL = config.redisFallbackTtlSeconds;

export async function setFallback(key: string, payload: string): Promise<void> {
  const r = getRedis();
  await r.setex(FALLBACK_PREFIX + key, TTL, payload);
}

export async function getFallback(key: string): Promise<string | null> {
  const r = getRedis();
  return r.get(FALLBACK_PREFIX + key);
}
