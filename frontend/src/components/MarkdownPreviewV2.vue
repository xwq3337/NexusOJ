<script setup lang="ts">
import { computed } from 'vue'
import MarkdownIt from 'markdown-it'
import markdownItHighlightjs from 'markdown-it-highlightjs'
import markdownItKatex from 'markdown-it-katex'
import 'highlight.js/styles/github.css'
import 'katex/dist/katex.min.css'

const props = defineProps({
  text: {
    type: String,
    required: true
  }
})

// 初始化 markdown-it 实例
const md = new MarkdownIt({
  html: true, // 允许 HTML 标签
  linkify: true, // 自动转换 URL 为链接
  typographer: true, // 启用排版优化
  breaks: true // 转换换行符为 <br>
})


// 使用插件
md.use(markdownItKatex, {
  throwOnError: false,
  errorColor: '#cc0000'
})
md.use(markdownItHighlightjs, {
  auto: true,
  code: true
})

// 渲染 markdown 为 HTML
const renderedHtml = computed(() => {
  return md.render(props.text || '')
})
</script>

<template>
  <div class="markdown-preview prose max-w-none" v-html="renderedHtml"></div>
</template>

<style scoped>
.markdown-preview {
  padding: 1rem;
  line-height: 1.6;
  color: #24292f;
}

.markdown-preview :deep(h1),
.markdown-preview :deep(h2),
.markdown-preview :deep(h3),
.markdown-preview :deep(h4),
.markdown-preview :deep(h5),
.markdown-preview :deep(h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

.markdown-preview :deep(h1) {
  font-size: 2em;
  border-bottom: 1px solid #d0d7de;
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h2) {
  font-size: 1.5em;
  border-bottom: 1px solid #d0d7de;
  padding-bottom: 0.3em;
}

.markdown-preview :deep(h3) {
  font-size: 1.25em;
}

.markdown-preview :deep(p) {
  margin-top: 0;
  margin-bottom: 16px;
}

.markdown-preview :deep(code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(175, 184, 193, 0.2);
  border-radius: 6px;
  font-family: ui-monospace, SFMono-Regular, SF Mono, Menlo, Consolas, monospace;
}

.markdown-preview :deep(pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f6f8fa;
  border-radius: 6px;
  margin-bottom: 16px;
}

.markdown-preview :deep(pre code) {
  padding: 0;
  background-color: transparent;
  border-radius: 0;
}

.markdown-preview :deep(blockquote) {
  padding: 0 1em;
  color: #57606a;
  border-left: 0.25em solid #d0d7de;
  margin: 0 0 16px 0;
}

.markdown-preview :deep(ul),
.markdown-preview :deep(ol) {
  padding-left: 2em;
  margin-bottom: 16px;
}

.markdown-preview :deep(li) {
  margin-top: 0.25em;
}

.markdown-preview :deep(table) {
  border-spacing: 0;
  border-collapse: collapse;
  margin-bottom: 16px;
  width: 100%;
}

.markdown-preview :deep(table th),
.markdown-preview :deep(table td) {
  padding: 6px 13px;
  border: 1px solid #d0d7de;
}

.markdown-preview :deep(table th) {
  font-weight: 600;
  background-color: #f6f8fa;
}

.markdown-preview :deep(table tr) {
  background-color: #ffffff;
  border-top: 1px solid #d0d7de;
}

.markdown-preview :deep(table tr:nth-child(2n)) {
  background-color: #f6f8fa;
}

.markdown-preview :deep(img) {
  max-width: 100%;
  box-sizing: content-box;
  background-color: #ffffff;
}

.markdown-preview :deep(a) {
  color: #0969da;
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.markdown-preview :deep(hr) {
  height: 0.25em;
  padding: 0;
  margin: 24px 0;
  background-color: #d0d7de;
  border: 0;
}

.markdown-preview :deep(.katex-display) {
  margin: 1em 0;
  overflow-x: auto;
  overflow-y: hidden;
}

.markdown-preview :deep(.katex) {
  font-size: 1.1em;
}

.markdown-preview :deep(.katex sup) {
  vertical-align: super !important;
}

.markdown-preview :deep(.katex sub) {
  vertical-align: sub !important;
}
</style>
