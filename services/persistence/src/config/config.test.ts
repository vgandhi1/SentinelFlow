import { config } from './index';

describe('config', () => {
  it('exposes kafka and db settings with defaults', () => {
    expect(config.kafkaBrokers).toBeDefined();
    expect(config.kafkaTopic).toBeDefined();
    expect(config.pgConnectionString).toBeDefined();
    expect(config.redisUrl).toBeDefined();
  });
});
