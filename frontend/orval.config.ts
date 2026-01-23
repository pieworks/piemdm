import { defineConfig } from 'orval';

export default defineConfig({
  api: {
    input: '../shared/api/swagger.json',
    output: {
      target: '../packages/web/api-client/index.ts',
      client: 'axios',
      mode: 'single',

    },
  },
});
