export const typeDefs = `#graphql
  type TelemetryPoint {
    time: String!
    sensorId: String!
    assetId: String!
    metricName: String!
    value: Float!
    idempotencyKey: String!
  }

  type Sensor {
    sensorId: String!
    assetId: String!
    metricName: String!
  }

  type Asset {
    assetId: String!
  }

  type Query {
    telemetry(sensorId: String, assetId: String, metricName: String, from: String, to: String, limit: Int): [TelemetryPoint!]!
    sensors(assetId: String): [Sensor!]!
    assets: [Asset!]!
  }
`;
