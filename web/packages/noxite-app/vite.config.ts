/// <reference types="vitest" />

import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'
import { resolve } from 'path'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag !== 'ez-icon' && tag.startsWith('ez-'),
        },
      },
    }),
    VueI18nPlugin({ include: resolve(__dirname, './src/locales/**') }),
  ],
  define: {
    '__APP_VERSION__': JSON.stringify(process.env.npm_package_version),
    'process.env': {},
  },
  build: {
    lib: {
      entry: resolve(__dirname, 'src/main.ts'),
      name: 'NoxiteWebapp',
      fileName: 'noxite-webapp',
    },
    rollupOptions: {
      external: ['@easy/ui', 'vue', '@easy/core-essentials'],
      output: {
        globals: {
          '@easy/ui': 'easyUI',
          'vue': 'Vue3',
          '@easy/core-essentials': 'CoreEssentials',
        },
      },
    },
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    coverage: {
      provider: 'istanbul',
      all: true,
      include: ['src/**/*.{ts,vue}'],
      exclude: ['src/main.ts', 'src/App.vue', 'src/testutils'],
    },
  },
  root: '.',
})
