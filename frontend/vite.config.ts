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
        '/agent' : {
          // target : "http://127.0.0.1:5557/",
          target : "http://47.109.57.7:5557/",
          changeOrigin : true,
          rewrite: (path) => path.replace(/^\/agent/, ''),
          onError: (err, req, res) => {
            console.log('代理错误，检查服务是否启动在5557端口')
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
