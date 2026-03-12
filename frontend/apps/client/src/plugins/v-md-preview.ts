import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
// Katex
import createKatexPlugin from '@kangc/v-md-editor/lib/plugins/katex/cdn'
// 代码行号
import createLineNumbertPlugin from '@kangc/v-md-editor/lib/plugins/line-number/index'
// Tip
import createTipPlugin from '@kangc/v-md-editor/lib/plugins/tip/index';
import '@kangc/v-md-editor/lib/plugins/tip/tip.css';
import { type App } from 'vue'
import hljs from 'highlight.js'
// Github主题
VMdPreview.use(githubTheme, {
  Hljs: hljs
})
VMdPreview.use(createKatexPlugin())
VMdPreview.use(createLineNumbertPlugin())
VMdPreview.use(createTipPlugin())
export const VMdPreviewPlugin = {
  install(app: App) {
    app.use(VMdPreview)
  }
}
