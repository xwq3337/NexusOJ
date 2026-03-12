import { API_BASE_URL } from '@nexusoj/config';
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
        // 使用通配符匹配多个后端路径
        '^/(sse|service)': {
          target: API_BASE_URL.development,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/(sse|service)/, ''),
          onError: (err, req, res) => {
            console.log('后端代理错误，检查服务是否启动在8080端口')
            // 可以在这里添加降级逻辑
          }
        },
        '/agent': {
          target: API_BASE_URL.agent,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/agent/, ''),
          onError: (err, req, res) => {
            console.log('代理错误，检查服务是否启动在5557端口')
            // 可以在这里添加降级逻辑
          }
        },
        '/ws': {
          target: API_BASE_URL.ws,
          secure: false,
          changeOrigin: true
        },
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
          // drop_console: true, // 移除 console 语句
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
