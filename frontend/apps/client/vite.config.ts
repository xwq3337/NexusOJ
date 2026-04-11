import { defineConfig, loadEnv } from 'vite'
import { fileURLToPath, URL } from 'node:url'
import { dirname, resolve } from 'node:path'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// monorepo 根目录（.env 所在目录）
const monorepoRoot = resolve(dirname(fileURLToPath(import.meta.url)), '../../')

export default defineConfig(({ mode, command }) => {
  const env = loadEnv(mode, monorepoRoot, '')
  const isDev = mode === 'development'
  const isBuild = command === 'build'
  return {
    server: {
      port: Number(env.CLIENT_DEV_PORT),
      host: env.VITE_HOST,
      proxy: {
        // 使用通配符匹配多个后端路径
        '^/(sse|service)': {
          target: env.BACKEND_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/(sse|service)/, ''),
          onError: (err: any, req: any, res: any) => {
            console.log({
              msg: '后端代理错误，检查服务是否启动在8080端口',
              err,
              req,
              res
            })
            // 可以在这里添加降级逻辑
          }
        },
        '/agent': {
          target: env.AGENT_URL,
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/agent/, ''),
          onError: (err: any, req: any, res: any) => {
            console.log('代理错误，检查服务是否启动在5557端口')
            // 可以在这里添加降级逻辑
          }
        },
        '/ws': {
          target: env.WEBSOCKET_URL,
          secure: false,
          changeOrigin: true
        }
      }
    },
    plugins: [vue(), vueDevTools()],
    define: {
      // 'process.env.API_KEY': JSON.stringify(env.GEMINI_API_KEY)
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
