import pg from 'pg';
import { config } from '../config';

const { Pool } = pg;

export type TelemetryRow = {
  idempotency_key: string;
  time: Date;
  sensor_id: string;
  asset_id: string;
  metric_name: string;
  value: number;
  tags?: Record<string, string>;
  critical: boolean;
};

let pool: pg.Pool | null = null;

export function getPool(): pg.Pool {
  if (!pool) {
    pool = new Pool({ connectionString: config.pgConnectionString });
  }
  return pool;
}

export async function upsertTelemetry(client: pg.PoolClient, row: TelemetryRow): Promise<void> {
  const timeNs = row.time.getTime() * 1e6;
  await client.query(
    `INSERT INTO telemetry (idempotency_key, time, sensor_id, asset_id, metric_name, value, tags, critical)
     VALUES ($1, to_timestamp($2::double precision / 1e9), $3, $4, $5, $6, $7::jsonb, $8)
     ON CONFLICT (idempotency_key) DO UPDATE SET
       value = EXCLUDED.value,
       tags = EXCLUDED.tags`,
    [
      row.idempotency_key,
      timeNs / 1e9,
      row.sensor_id,
      row.asset_id,
      row.metric_name,
      row.value,
      row.tags ? JSON.stringify(row.tags) : null,
      row.critical,
    ]
  );
}

export async function initSchema(client: pg.PoolClient): Promise<void> {
  await client.query(`
    CREATE TABLE IF NOT EXISTS telemetry (
      idempotency_key TEXT PRIMARY KEY,
      time TIMESTAMPTZ NOT NULL,
      sensor_id TEXT NOT NULL,
      asset_id TEXT NOT NULL,
      metric_name TEXT NOT NULL,
      value DOUBLE PRECISION NOT NULL,
      tags JSONB,
      critical BOOLEAN NOT NULL DEFAULT false
    );
  `);
  // For TimescaleDB: run SELECT create_hypertable('telemetry', 'time', if_not_exists => true); once
}
