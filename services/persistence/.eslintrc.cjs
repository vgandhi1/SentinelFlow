module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: { project: './tsconfig.json' },
  env: { node: true, es2022: true },
  extends: ['eslint:recommended'],
  overrides: [
    {
      files: ['*.ts'],
      extends: ['eslint:recommended'],
      rules: { 'no-unused-vars': 'off' },
    },
  ],
  ignorePatterns: ['dist', 'node_modules'],
};
