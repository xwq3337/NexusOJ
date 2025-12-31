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
            <n-tag :type="getStatusType(record.verdict)" size="large">
              {{ record.verdict }}
            </n-tag>
            <n-tag type="info" size="large">
              {{ getLanguageLabel(record.language) }}
            </n-tag>
          </div>
        </div>

        <div
          class="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t"
          :style="{ borderColor: 'var(--border-color)' }"
        >
          <div
            class="text-center p-3 rounded-lg"
            :style="{ backgroundColor: 'var(--surface-secondary)' }"
          >
            <div class="text-2xl font-bold text-blue-400">{{ record.max_time }}ms</div>
            <div class="text-sm text-gray-500">最大时间</div>
          </div>
          <div
            class="text-center p-3 rounded-lg"
            :style="{ backgroundColor: 'var(--surface-secondary)' }"
          >
            <div class="text-2xl font-bold text-green-400">{{ record.max_memory }}MB</div>
            <div class="text-sm text-gray-500">最大内存</div>
          </div>
          <div
            class="text-center p-3 rounded-lg"
            :style="{ backgroundColor: 'var(--surface-secondary)' }"
          >
            <div class="text-2xl font-bold text-purple-400">{{ testCases.length }}</div>
            <div class="text-sm text-gray-500">测试用例</div>
          </div>
          <div
            class="text-center p-3 rounded-lg"
            :style="{ backgroundColor: 'var(--surface-secondary)' }"
          >
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
        <div
          class="rounded-xl overflow-hidden border"
          :style="{ borderColor: 'var(--border-color)' }"
        >
          <pre
            class="p-4 text-sm overflow-x-auto"
            :style="{
              backgroundColor: 'var(--surface-tertiary)',
              color: 'var(--text-primary)',
              fontFamily: 'monospace',
              lineHeight: '1.5'
            }"
            v-html="highlightedCode"
          ></pre>
        </div>
      </div>

      <!-- 评测结果区域 -->
      <div>
        <h2 class="text-xl font-bold mb-4" :style="{ color: 'var(--text-primary)' }">评测结果</h2>

        <div class="space-y-4">
          <div
            v-for="(testCase, index) in testCases"
            :key="testCase.case_id"
            class="p-4 rounded-xl border"
            :style="{
              borderColor: getTestCaseBorderColor(testCase.status),
              backgroundColor: 'var(--surface-secondary)'
            }"
          >
            <div class="flex justify-between items-center mb-3">
              <div class="flex items-center gap-3">
                <n-tag :type="getTestCaseStatusType(testCase.status)">{{ testCase.status }}</n-tag>
                <span class="font-mono text-sm">#{{ testCase.case_id }}</span>
              </div>
              <div class="flex items-center gap-4 text-sm">
                <div class="flex items-center gap-1">
                  <Clock :size="14" />
                  <span>{{ testCase.time }}ms</span>
                </div>
                <div class="flex items-center gap-1">
                  <Database :size="14" />
                  <span>{{ testCase.memory }}MB</span>
                </div>
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">输入</h4>
                <pre
                  class="p-3 rounded text-sm overflow-x-auto"
                  :style="{
                    backgroundColor: 'var(--surface-tertiary)',
                    color: 'var(--text-secondary)'
                  }"
                  >{{ testCase.stdin }}</pre
                >
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">输出</h4>
                <pre
                  class="p-3 rounded text-sm overflow-x-auto"
                  :style="{
                    backgroundColor: 'var(--surface-tertiary)',
                    color: 'var(--text-secondary)'
                  }"
                  >{{ testCase.stdout }}</pre
                >
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">期望输出</h4>
                <pre
                  class="p-3 rounded text-sm overflow-x-auto"
                  :style="{
                    backgroundColor: 'var(--surface-tertiary)',
                    color: 'var(--text-secondary)'
                  }"
                  >{{ testCase.expected }}</pre
                >
              </div>
              <div>
                <h4 class="font-medium mb-2 text-sm text-gray-500">错误输出</h4>
                <pre
                  v-if="testCase.stderr"
                  class="p-3 rounded text-sm overflow-x-auto"
                  :style="{
                    backgroundColor: 'var(--surface-tertiary)',
                    color: 'var(--text-secondary)'
                  }"
                  >{{ testCase.stderr }}</pre
                >
                <div
                  v-else
                  class="p-3 rounded text-sm text-gray-500"
                  :style="{ backgroundColor: 'var(--surface-tertiary)' }"
                >
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
import { ref, onMounted, computed } from 'vue'
import hljs from 'highlight.js'
import { useRoute } from 'vue-router'
import { NTag, NButton } from 'naive-ui'
import { Clock, Database, Copy } from 'lucide-vue-next'
import Request from '@/services/api'
import { safeJsonParse } from '@/utils/safeJsonParse'

onMounted(async () => {
  await Request.get('/record/' + route.params.id).then((res) => {
    record.value = res.info
  })
})

// 评测记录数据结构
interface Record {
  id: number
  code: string
  created_at: string
  deleted_at: string | null
  judge_result: string // JSON字符串
  language: string
  max_memory: number
  max_time: number
  problem_id: string
  problem_title: string
  updated_at: string
  user_id: string
  username: string
  verdict: string
}

// 测试用例数据结构
interface TestCase {
  time: number
  stdin: string
  memory: number
  status: string
  stderr: string
  stdout: string
  case_id: string
  expected: string
}

// 获取当前路由参数
const route = useRoute()

// 记录数据
const record = ref<Record>({} as Record)

// 解析测试用例
const testCases = computed<TestCase[]>(() => {
  const result = safeJsonParse(record.value.judge_result)
  return result.isVaild() ? result.value : []
})

// 代码高亮
const highlightedCode = computed(() => {
  const langMap = {
    cpp: 'cpp',
    python: 'python',
    java: 'java',
    javascript: 'javascript',
    csharp: 'csharp',
    c: 'c',
    text: 'plaintext'
  }

  const language = langMap[record.value.language] || 'plaintext'
  const highlighted = hljs.highlight(record.value.code ?? '', { language })
  return highlighted.value
})

// 获取状态类型
const getStatusType = (verdict: string) => {
  switch (verdict) {
    case 'Accepted':
      return 'success'
    case 'WrongAnswer':
    case 'RuntimeError':
    case 'TimeLimitExceeded':
    case 'MemoryLimitExceeded':
      return 'error'
    case 'CompilationError':
      return 'warning'
    default:
      return 'default'
  }
}

// 获取测试用例状态类型
const getTestCaseStatusType = (status: string) => {
  switch (status) {
    case 'Accepted':
      return 'success'
    case 'WrongAnswer':
    case 'RuntimeError':
    case 'TimeLimitExceeded':
    case 'MemoryLimitExceeded':
      return 'error'
    case 'CompilationError':
      return 'warning'
    default:
      return 'default'
  }
}

// 获取测试用例边框颜色
const getTestCaseBorderColor = (status: string) => {
  switch (status) {
    case 'Accepted':
      return 'var(--success-color)'
    case 'WrongAnswer':
    case 'RuntimeError':
    case 'TimeLimitExceeded':
    case 'MemoryLimitExceeded':
      return 'var(--error-color)'
    case 'CompilationError':
      return 'var(--warning-color)'
    default:
      return 'var(--border-color)'
  }
}

// 获取语言标签
const getLanguageLabel = (lang: string) => {
  switch (lang) {
    case 'cpp':
      return 'C++'
    case 'python':
      return 'Python'
    case 'java':
      return 'Java'
    case 'javascript':
      return 'JavaScript'
    case 'csharp':
      return 'C#'
    default:
      return lang
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 复制代码
const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(record.value.code)
    // 这里可以添加一个提示，但需要从app中获取message实例
    console.log('代码已复制到剪贴板')
  } catch (err) {
    console.error('复制代码失败:', err)
  }
}
</script>
