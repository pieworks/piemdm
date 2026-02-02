/// <reference types="vitest" />
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

import { defineConfig, loadEnv } from 'vite';
import { prismjsPlugin } from 'vite-plugin-prismjs';
import vueDevTools from 'vite-plugin-vue-devtools';

// env must start with this prefix. in .env config
const envPrefix = 'VITE_';
// process.env = { ...process.env, ...loadEnv(mode, process.cwd(), envPrefix) };

export default ({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd(), envPrefix) };

  return defineConfig({
    base: process.env.VITE_BASE,
    define: {
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'false',
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
        views: resolve(__dirname, 'src/views'),
        components: resolve(__dirname, 'src/components'),
        'workflow-vue': '@pieteams/workflow-vue',
        vue$: 'vue/dist/vue.runtime.esm-bundler.js',
      },
    },
    plugins: [
      vue(),
      prismjsPlugin({
        languages: ['json', 'css', 'javascript', 'typescript', 'vue'],
        plugins: ['line-numbers'],
        theme: 'tomorrow',
        css: true,
      }),

      vueDevTools(),
    ],
    envPrefix: envPrefix,
    server: {
      // if your frontend not in the localhost, please uncomment the https config meanwhile
      host: process.env.VITE_HOST,
      port: Number(process.env.VITE_PORT),
    },
    // @ts-expect-error vitest config
    test: {
      environment: 'happy-dom',
      globals: true,
      setupFiles: ['src/test/setup.js'],
      coverage: {
        provider: 'v8',
        reporter: ['text', 'json', 'html'],
      },
      testTimeout: 10000,
    },
  });
};
