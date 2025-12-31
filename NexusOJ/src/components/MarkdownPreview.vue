<script setup lang="ts">
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
import createKatexPlugin from '@kangc/v-md-editor/lib/plugins/katex/cdn'
import createLineNumbertPlugin from '@kangc/v-md-editor/lib/plugins/line-number/index'

import hljs from 'highlight.js'

VMdPreview.use(githubTheme, {
  Hljs: hljs
})
VMdPreview.use(createKatexPlugin())
VMdPreview.use(createLineNumbertPlugin())

const props = defineProps({
  text: {
    type: String,
    required: true
  }
})
</script>

<template>
  <v-md-preview :text="text"></v-md-preview>
</template>
