import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vitest/config'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  server: {
    port: 5173,
    proxy: {
      // 开发环境将 API 请求代理到后端（避免 CORS，PRD §4.4）
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/healthz': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },

  build: {
    chunkSizeWarningLimit: 600,
    rollupOptions: {
      output: {
        // Rollup 4 仅支持函数形式 manualChunks：Vue 生态单独分包，其余三方库归 vendor
        manualChunks(id: string) {
          if (id.includes('node_modules')) {
            if (id.includes('pinia') || id.includes('@vue') || id.includes('/vue')) {
              return 'vue'
            }
            return 'vendor'
          }
        },
      },
    },
  },

  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
})
