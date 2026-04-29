import pg from 'pg';
import { config } from './config';

const pool = new pg.Pool({ connectionString: config.pgConnectionString });

export type TelemetryRow = {
  time: Date;
  sensor_id: string;
  asset_id: string;
  metric_name: string;
  value: number;
  idempotency_key: string;
};

export async function queryTelemetry(params: {
  sensorId?: string;
  assetId?: string;
  metricName?: string;
  from?: string;
  to?: string;
  limit?: number;
}): Promise<TelemetryRow[]> {
  const conditions: string[] = [];
  const values: unknown[] = [];
  let i = 1;
  if (params.sensorId) {
    conditions.push(`sensor_id = $${i++}`);
    values.push(params.sensorId);
  }
  if (params.assetId) {
    conditions.push(`asset_id = $${i++}`);
    values.push(params.assetId);
  }
  if (params.metricName) {
    conditions.push(`metric_name = $${i++}`);
    values.push(params.metricName);
  }
  if (params.from) {
    conditions.push(`time >= $${i++}::timestamptz`);
    values.push(params.from);
  }
  if (params.to) {
    conditions.push(`time <= $${i++}::timestamptz`);
    values.push(params.to);
  }
  const where = conditions.length ? `WHERE ${conditions.join(' AND ')}` : '';
  const limit = Math.min(params.limit ?? 1000, 10000);
  values.push(limit);
  const res = await pool.query(
    `SELECT time, sensor_id, asset_id, metric_name, value, idempotency_key FROM telemetry ${where} ORDER BY time DESC LIMIT $${i}`,
    values
  );
  return res.rows as TelemetryRow[];
}

export async function querySensors(assetId?: string): Promise<{ sensor_id: string; asset_id: string; metric_name: string }[]> {
  if (assetId) {
    const res = await pool.query(
      'SELECT DISTINCT sensor_id, asset_id, metric_name FROM telemetry WHERE asset_id = $1',
      [assetId]
    );
    return res.rows;
  }
  const res = await pool.query('SELECT DISTINCT sensor_id, asset_id, metric_name FROM telemetry');
  return res.rows;
}

export async function queryAssets(): Promise<{ asset_id: string }[]> {
  const res = await pool.query('SELECT DISTINCT asset_id FROM telemetry');
  return res.rows;
}
