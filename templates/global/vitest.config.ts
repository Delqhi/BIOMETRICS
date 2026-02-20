import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    globals: true,
    environment: 'node',
    include: ['src/**/*.test.ts', 'src/**/*.spec.ts'],
    exclude: ['node_modules', 'dist', 'coverage'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules',
        'dist',
        'coverage',
        '**/*.test.ts',
        '**/*.spec.ts',
        '**/*.config.ts',
        '**/vite.config.ts'
      ]
    },
    testTimeout: 10000,
    hookTimeout: 10000,
    env: {
      NODE_ENV: 'test',
      VITEST: 'true'
    }
  }
});
