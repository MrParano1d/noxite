/// <reference types="vitest" />

import { resolve } from 'path'
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import dts from 'vite-plugin-dts'

export default defineConfig({
  plugins: [dts({ insertTypesEntry: true })],
  build: {
    lib: {
      entry: resolve(__dirname, 'src/index.ts'),
      name: 'NoxiteWebAdapters',
      fileName: 'noxite-web-adapters',
    },
    rollupOptions: {
      external: ['@easy/core-essentials'],
      output: {
        globals: {
          '@easy/core-essentials': 'CoreEssentials',
        },
      },
    },
  },
  define: {
    'process.env': {},
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    include: ['**/*.test.ts'],
    coverage: {
      provider: 'istanbul',
      all: true,
      include: './src/**/*.ts',
    },
  },
})
