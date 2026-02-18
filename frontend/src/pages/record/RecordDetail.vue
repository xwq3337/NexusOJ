<template>
  <div :style="{ backgroundColor: 'var(--bg-primary)' }" class="min-h-screen p-6">
    <div class="max-w-7xl mx-auto">
      <!-- 评测概览区域 -->
      <div class="mb-6 p-6 rounded-xl" :style="{ backgroundColor: 'var(--surface-primary)' }">
        <div class="flex flex-wrap justify-between items-center mb-4">
          <div>
            <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">评测详情</h1>
            <div class="flex items-center gap-4 mt-2">
              <div class="flex items-center gap-2">
                <span class="text-gray-500">ID:</span>
                <span class="font-mono" :style="{ color: 'var(--text-primary)' }">{{
                  record.id
                  }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-gray-500">题目:</span>
                <span class="font-medium" :style="{ color: 'var(--text-primary)' }">{{
                  record.problem_title
                  }}</span>
              </div>
              <div class="flex items-center gap-2">
                <span class="text-gray-500">用户:</span>
                <span class="font-medium" :style="{ color: 'var(--text-primary)' }">{{
                  record.username
                  }}</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-4">
            <n-tag :color="STATUS_COLORS[record.verdict]" size="large">
              {{ record.verdict }}
            </n-tag>
            <n-tag type="info" size="large">
              {{ LANGUAGE_CONFIG[record.language]?.label || 'Unknown' }}
            </n-tag>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t"
          :style="{ borderColor: 'var(--border-color)' }">
          <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <div class="text-2xl font-bold text-blue-400">{{ formatTime(record.max_time) }}</div>
            <div class="text-sm text-gray-500">最大时间</div>
          </div>
          <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <div class="text-2xl font-bold text-green-400">{{ formatMemory(record.max_memory) }}</div>
            <div class="text-sm text-gray-500">最大内存</div>
          </div>
          <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <div class="text-2xl font-bold text-purple-400">{{ testCases.length }}</div>
            <div class="text-sm text-gray-500">测试用例</div>
          </div>
          <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
            <div class="text-2xl font-bold text-yellow-400">
              {{ formatDate(record.created_at) }}
            </div>
            <div class="text-sm text-gray-500">提交时间</div>
          </div>
        </div>
      </div>

      <!-- 代码展示区域 -->
      <div class="mb-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-bold" :style="{ color: 'var(--text-primary)' }">代码</h2>
          <n-button @click="copyCode" size="small" type="primary">
            <template #icon>
              <Copy :size="14" />
            </template>
            复制代码
          </n-button>
        </div>
        <div class="rounded-xl overflow-hidden border" :style="{ borderColor: 'var(--border-color)' }">
          <pre class="p-4 text-sm overflow-x-auto" :style="{
            backgroundColor: 'var(--surface-tertiary)',
            color: 'var(--text-primary)',
            fontFamily: 'monospace',
            lineHeight: '1.5'
          }" v-html="highlightedCode"></pre>

        </div>
      </div>

      <!-- 评测结果区域 -->
      <div>
        <h2 class="text-xl font-bold mb-4" :style="{ color: 'var(--text-primary)' }">评测结果</h2>

        <div class="space-y-4">
          <div v-for="(testCase, index) in testCases" :key="testCase.case_id" class="p-4 rounded-xl border" :style="{
            borderColor: STATUS_COLORS[testCase.status].borderColor,
            backgroundColor: STATUS_COLORS[testCase.status].color,
          }">
            <div class="flex justify-between items-center mb-3">
              <div class="flex items-center gap-3">
                <n-tag :color="STATUS_COLORS[testCase.status]">{{ testCase.status }}</n-tag>
                <span class="font-mono text-sm">#{{ testCase.case_id }}</span>
              </div>
              <div class="flex items-center gap-4 text-sm">
                <div class="flex items-center gap-1">
                  <Clock :size="14" />
                  <span>{{ formatTime(testCase.time) }}</span>
                </div>
                <div class="flex items-center gap-1">
                  <Database :size="14" />
                  <span>{{ formatMemory(testCase.memory) }}</span>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">输入</h4>
                <pre class="p-3 rounded text-sm overflow-x-auto" :style="{
                  backgroundColor: 'var(--surface-tertiary)',
                  color: 'var(--text-secondary)'
                }">{{ testCase.stdin }}</pre>
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">输出</h4>
                <pre class="p-3 rounded text-sm overflow-x-auto" :style="{
                  backgroundColor: 'var(--surface-tertiary)',
                  color: 'var(--text-secondary)'
                }">{{ testCase.stdout }}</pre>
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">期望输出</h4>
                <pre class="p-3 rounded text-sm overflow-x-auto" :style="{
                  backgroundColor: 'var(--surface-tertiary)',
                  color: 'var(--text-secondary)'
                }">{{ testCase.expected }}</pre>
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">错误输出</h4>
                <pre v-if="testCase.stderr" class="p-3 rounded text-sm overflow-x-auto" :style="{
                  backgroundColor: 'var(--surface-tertiary)',
                  color: 'var(--text-secondary)'
                }">{{ testCase.stderr }}</pre>
                <div v-else class="p-3 rounded text-sm text-gray-500"
                  :style="{ backgroundColor: 'var(--surface-tertiary)' }">
                  无错误输出
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NTag, NButton, useMessage } from 'naive-ui'
import { Clock, Database, Copy } from 'lucide-vue-next'
import { useClipboard } from '@vueuse/core'
import { LANGUAGE_CONFIG, STATUS_COLORS } from '@/constants'
import hljs from 'highlight.js/lib/common'
import javascript from 'highlight.js/lib/languages/javascript'
import cpp from 'highlight.js/lib/languages/cpp'
import python from 'highlight.js/lib/languages/python'
import java from 'highlight.js/lib/languages/java'
import rust from 'highlight.js/lib/languages/rust'
import go from 'highlight.js/lib/languages/go'
import { useTheme } from '@/composables/useTheme'
import { recordApi } from '@/services/record'
import { GetRecordDetailResponse, JudgeTestCaseResult } from '@/types/record'
const { theme } = useTheme()
import { formatTime, formatMemory, formatDate } from '@/utils/format'
// 动态加载 highlight.js 主题 CSS
let highlightThemeLoaded: string | null = null

const loadHighlightTheme = async (isDark: boolean) => {
  const themeName = isDark ? 'dark' : 'light'

  // 如果已经加载了相同的主题，跳过
  if (highlightThemeLoaded === themeName) {
    return
  }

  // 移除旧的样式表
  const oldStyle = document.getElementById('highlightjs-theme')
  if (oldStyle) {
    oldStyle.remove()
  }

  try {
    // 动态导入 CSS 文本内容
    const cssContent = isDark
      ? (await import('highlight.js/styles/github-dark.css?raw')).default
      : (await import('highlight.js/styles/github.css?raw')).default

    // 创建 style 标签并插入
    const style = document.createElement('style')
    style.id = 'highlightjs-theme'
    style.textContent = cssContent
    document.head.appendChild(style)

    highlightThemeLoaded = themeName
  } catch (error) {
    console.error('加载 highlight.js 主题失败:', error)
  }
}

// 注册语言
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('python', python)
hljs.registerLanguage('java', java)
hljs.registerLanguage('rust', rust)
hljs.registerLanguage('go', go)

// 初始化时加载主题
onMounted(() => {
  loadHighlightTheme(theme.value === 'dark')
})

// 监听主题变化
watch(theme, (newTheme) => {
  loadHighlightTheme(newTheme === 'dark')
})

onMounted(async () => {
  if (!route.params.id) return
  await recordApi.getRecordDetail(route.params.id as string).then((res) => {
    const { info } = res
    record.value = info
    testCases.value = info.judge_result || []
  })
})

// 获取当前路由参数
const route = useRoute()

// 记录数据
const record = ref({} as GetRecordDetailResponse)

// 解析测试用例
const testCases = ref<JudgeTestCaseResult[]>([])

// 代码高亮
const highlightedCode = computed(() => {
  const langMap: Record<string, string> = {
    cpp: 'cpp',
    python: 'python',
    java: 'java',
    javascript: 'javascript',
    csharp: 'csharp',
    c: 'c',
    text: 'plaintext',
    rust: 'rust',
    go: 'go'
  }

  const language = langMap[record.value.language] || 'plaintext'
  const code = record.value.code || ''
  try {
    const result = hljs.highlight(code, { language })
    return result.value || ''
  } catch (error) {
    console.error('代码高亮失败:', error)
    return code
  }
})


const { copy } = useClipboard()
const message = useMessage()

// 复制代码
const copyCode = async () => {
  try {
    await copy(record.value.code)
    message.success('代码已复制到剪贴板')
  } catch (err) {
    message.error('复制代码失败')
  }
}
</script>
