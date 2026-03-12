import './assets/base.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
const pinia = createPinia()
pinia.use(piniaPluginPersistedstate)
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'

import App from './App.vue'
import router from './router'
import { usePerformanceStore } from './stores/performance'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
const app = createApp(App)
import VueMarkdownEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'
import VMdPreview from '@kangc/v-md-editor/lib/preview';
import createKatexPlugin from '@kangc/v-md-editor/lib/plugins/katex/cdn';
import '@kangc/v-md-editor/lib/style/preview.css';


VueMarkdownEditor.use(vuepressTheme, {
  Prism,
})
VueMarkdownEditor.use(createKatexPlugin)

VMdPreview.use(vuepressTheme, {
  Prism,
})
VMdPreview.use(createKatexPlugin)

app.use(VueMarkdownEditor)
app.use(VMdPreview)
app.use(pinia)
app.use(router)
app.use(ElementPlus, {
  locale: zhCn,
})
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
// 初始化性能监控
const performanceStore = usePerformanceStore(pinia)
performanceStore.initializeMonitoring()

app.mount('#app')
