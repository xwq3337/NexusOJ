/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-12-16 13:21:42
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-14 11:29:38
 * @FilePath: /frontend/vite.config.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import vue from '@vitejs/plugin-vue'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, '.', '')
  return {
    server: {
      port: 3000,
      host: '0.0.0.0',
      proxy: {
        '/api': {
          // target: 'http://127.0.0.1:8080/',
          target: 'http://47.109.57.7:8080',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
          onError: (err, req, res) => {
            console.log('代理错误，检查后端服务是否启动在8080端口')
            // 可以在这里添加降级逻辑
            
          }
        },
        '/ws': {
          // target: 'ws://127.0.0.1:8080/',
          target: 'http://47.109.57.7:8080',
          secure: false,
          changeOrigin: true
        }
      }
    },
    plugins: [vue()],
    define: {
      'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY),
      'process.env.GEMINI_API_KEY': JSON.stringify(env.GEMINI_API_KEY)
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    }
  }
})
