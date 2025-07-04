// vite.config.js
import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'
import {defineConfig} from 'vite'


export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
    extensions: ['.js', '.vue', '.json'] // 确保支持 .js 文件导入
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    // 允许的主机名列表
    allowedHosts: [
      'c6f4-14-145-46-241.ngrok-free.app'
    ],
    proxy: {
      '/api': {
        target: 'http://192.168.235.129:8080',
        changeOrigin: true,
        rewrite: (path) => {
          // 将 /api/ 替换为 '' 并将剩余部分作为查询参数
          return path; // 直接返回
        }
      }
    }
  }
})
