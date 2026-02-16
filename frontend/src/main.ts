import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './assets/main.css'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import {VMdEditorPlugin,VMdPreviewPlugin} from './plugins'
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
type ThemeType = 'light' | 'dark'
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const theme = savedTheme || (prefersDark ? 'dark' : 'light')

  // 设置初始主题
  updateTheme(theme as ThemeType)

  // 监听 localStorage 变化（多标签页同步）
  window.addEventListener('storage', (e) => {
    if (e.key === 'theme' && e.newValue) {
      updateTheme(e.newValue as ThemeType)
    }
  })
}

function updateTheme(theme: ThemeType) {
  document.body.classList.remove('light', 'dark')
  document.documentElement.classList.remove('light', 'dark')

  document.body.classList.add(theme)
  document.documentElement.classList.add(theme === 'dark' ? 'dark' : 'light')
}

initTheme()

const app = createApp(App)
app.use(router)
app.use(pinia)
app.use(VMdEditorPlugin)
app.use(VMdPreviewPlugin)
app.mount('#app')
