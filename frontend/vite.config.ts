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
        '/service': {
          target: 'http://127.0.0.1:8080/',
          // target: 'http://47.109.57.7:8080',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/service/, ''),
          onError: (err, req, res) => {
            console.log('代理错误，检查后端服务是否启动在8080端口')
            // 可以在这里添加降级逻辑
          }
        },
        '/agent': {
          // target : "http://127.0.0.1:5557/",
          target: 'http://47.109.57.7:5557/',
          changeOrigin: true,
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
        },
        '/test': {
          target: 'http://localhost:8000',
          rewrite: (path) => path.replace(/^\/test/, ''),
          onError: (err, req, res) => {
            console.log('代理错误，检查服务是否启动在8000端口')
            // 可以在这里添加降级逻辑
          }
        }
      }
    },
    plugins: [vue()],
    define: {
      'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY)
    },
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    build: {
      outDir: 'dist',
      assetsDir: 'assets',
      sourcemap: true, // 生成 source map 以便调试
      minify: 'terser', // 使用 terser 进行压缩
      terserOptions: {
        compress: {
          drop_console: true, // 移除 console 语句
          drop_debugger: true // 移除 debugger 语句
        }
      },
      cssMinify: true, // 压缩 CSS
      cssCodeSplit: true, // 启用 CSS 代码分割
      reportCompressedSize: true, // 显示压缩后的文件大小
      rollupOptions: {
        output: {
          chunkFileNames: 'assets/js/[hash].js',
          entryFileNames: 'assets/js/[hash].js',
          assetFileNames: 'assets/[ext]/[hash].[ext]',
          manualChunks: {
            // Ace 571.43kB
            vendor: ['vue', 'vue-router', 'pinia'], // 将 Vue 和 Axios 分离到一个单独的 chunk 中
            utils: ['lodash', 'dayjs', 'axios'], // 将 Lodash 分离到一个单独的 chunk 中
            editor: ['codemirror', '@kangc/v-md-editor'], // 将 CodeMirror 分离到一个单独的 chunk 中
            // index.js 1498.91kB
            // ui.js 1386.27kB
            ui: ['naive-ui', 'echarts', 'vue-echarts'] // 将 Naive UI 分离到一个单独的 chunk 中
          }
        }
      }
    }
  }
})
