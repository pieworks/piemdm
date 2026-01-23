import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: '../shared/api/swagger.json',
    output: {
      target: '../packages/web/api-client/index.ts',
      client: 'axios',
      mode: 'single',
      override: {
        mutator: {
          path: './src/api/axios.ts', // Optional: Custom axios instance if needed, or remove for default
          name: 'customInstance',
        },
      },
    },
  },
});
