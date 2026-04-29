import * as db from './db';

export const resolvers = {
  Query: {
    telemetry: async (
      _: unknown,
      args: {
        sensorId?: string;
        assetId?: string;
        metricName?: string;
        from?: string;
        to?: string;
        limit?: number;
      }
    ) => {
      const rows = await db.queryTelemetry({
        sensorId: args.sensorId ?? undefined,
        assetId: args.assetId ?? undefined,
        metricName: args.metricName ?? undefined,
        from: args.from ?? undefined,
        to: args.to ?? undefined,
        limit: args.limit ?? 1000,
      });
      return rows.map((r) => ({
        time: r.time.toISOString(),
        sensorId: r.sensor_id,
        assetId: r.asset_id,
        metricName: r.metric_name,
        value: r.value,
        idempotencyKey: r.idempotency_key,
      }));
    },
    sensors: async (_: unknown, args: { assetId?: string }) => {
      const rows = await db.querySensors(args.assetId ?? undefined);
      return rows.map((r) => ({
        sensorId: r.sensor_id,
        assetId: r.asset_id,
        metricName: r.metric_name,
      }));
    },
    assets: async () => {
      const rows = await db.queryAssets();
      return rows.map((r) => ({ assetId: r.asset_id }));
    },
  },
};
