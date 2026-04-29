function getEnv(key: string, defaultValue: string): string {
  return process.env[key] ?? defaultValue;
}

export const config = {
  port: parseInt(getEnv('PORT', '4000'), 10),
  pgConnectionString: getEnv(
    'PG_URL',
    'postgres://sentinel:sentinel@localhost:5432/sentinelflow'
  ),
};
