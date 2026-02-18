<template>
  <div class="flex h-screen overflow-hidden">
    <n-split direction="horizontal" :max="0.75" :min="0.25">
      <template #1>
        <div class="h-full overflow-y-auto markdown-container" :style="{
          backgroundColor: 'var(--surface-secondary)',
          color: 'var(--text-primary)'
        }">
          <div class="p-4 border-b" :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }">
            <h1 class="text-2xl font-bold mb-2" :style="{ color: 'var(--text-primary)' }">
              {{ problem.title }}
            </h1>
            <div class="flex gap-4 text-sm">
              <span class="font-medium" :class="difficultyMap[Number(problem.difficulty) - 1]?.color">{{
                difficultyMap[Number(problem.difficulty) - 1]?.text }}</span>
              <span class="text-black">通过率: {{ formatAcceptance(problem.accept, problem.submission) }}</span>
              <span class="text-black">标签: {{ problem.tags.join(', ') }}</span>
              <span class="text-black">判题配置: {{ problem.judge_config.time_limit }}s
                {{ problem.judge_config.memory_limit }}MB</span>
            </div>
          </div>
          <v-md-preview :text="ProblemContext" style="height: 50rem;" :style="{
            padding: '5px',
            backgroundColor: 'transparent'
          }" />
        </div>
      </template>
      <template #2>
        <div class="h-12 border-b flex items-center justify-between px-4" :style="{
          backgroundColor: 'var(--surface-primary)',
          borderColor: 'var(--border-color)',
          borderWidth: '1px',
          borderStyle: 'solid'
        }">
          <div class="flex items-center gap-2">
            <!-- TODO: 题目页 主题切换 -->
            <n-select v-model:value="Language"
              :options="Object.values(LANGUAGE_CONFIG).map(config => ({ value: config.value, label: config.label }))"
              :style="{ width: '140px' }" :dropdown-props="{ style: { maxHeight: '200px', overflowY: 'auto' } }"
              placeholder="选择语言">
            </n-select>
            <n-popover trigger="click" placement="bottom">
              <template #trigger>
                <n-button :style="{ color: 'var(--text-primary)' }"
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded transition-colors">
                  <Settings :size="14" /> 设置
                </n-button>
              </template>
              <n-grid x-gap="12" :cols="2" :style="{ width: '15rem' }">
                <n-gi>
                  <n-select v-model:value="Theme" :options="EDITOR_THEME_OPTIONS.map(
                    (i) => {
                      return { label: i, value: i }
                    }
                  )" :dropdown-props="{
                    style: { maxHeight: '200px', overflowY: 'auto' }
                  }" placeholder="选择主题">
                  </n-select>
                </n-gi>
                <n-gi>
                  <n-input-number v-model:value="fontSize" :update-value-on-input="false" placeholder="" :min="10"
                    :max="30" />
                </n-gi>
              </n-grid>
            </n-popover>
          </div>
          <div class="flex items-center gap-2">
            <n-button class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded transition-colors" :style="{
              color: 'var(--text-primary)',
              backgroundColor: hoverBgColor2
            }" @mouseenter="() => (hoverBgColor2 = 'var(--surface-tertiary)')"
              @mouseleave="() => (hoverBgColor2 = 'transparent')">
              <RotateCcw :size="14" /> 重置
            </n-button>
            <n-button @click="handleTest" :disabled="isTesting"
              class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium rounded transition-colors"
              :style="{ color: 'var(--text-primary)' }" :class="isTesting ? 'opacity-70 cursor-wait' : ''">
              <Play :size="14" /> {{ isTesting ? '测试中...' : '自测运行' }}
            </n-button>
            <n-button @click="handleRun" :disabled="isRunning"
              type="success"
              class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium rounded transition-colors"
              :style="{ color: 'var(--text-primary)' }" :class="isRunning ? 'opacity-70 cursor-wait' : ''">
               {{ isRunning ? '提交中...' : '提交代码' }}
            </n-button>
          </div>
        </div>

        <div class="flex-1 relative font-mono" style="height: 66%">
          <codeEditor @change="handleEditorChange" :value="code" " :theme="Theme"
            :language="LANGUAGE_CONFIG[Language].aceMode" />
        </div>

        <div class="h-full border-t px-2 py-0 font-mono text-sm overflow-y-auto" :style="{
          backgroundColor: 'var(--surface-secondary)',
          borderColor: 'var(--border-color)',
          borderTopWidth: '1px',
          borderStyle: 'solid'
        }">
          <n-tabs type="line" animated>
            <n-tab-pane name="Case" tab="测试用例">
              <n-grid x-gap="12" :cols="2" :style="{ height: '100%' }">
                <n-gi>
                  <n-input autosize type="textarea" v-model:value="test_case.input" placeholder="输入"></n-input>
                </n-gi>
                <n-gi>
                  <n-input autosize disabled type="textarea" v-model:value="test_case.output"
                    placeholder="输出"></n-input>
                </n-gi>
              </n-grid>
            </n-tab-pane>
            <n-tab-pane name="console" tab="控制台"> 模拟 {{ result }} </n-tab-pane>
          </n-tabs>
        </div>
      </template>
    </n-split>
    <AiAssistant />
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import {
  NSplit,
  NPopover,
  NButton,
  NSelect,
  NTabPane,
  NTabs,
  NInput,
  NGrid,
  NGi,
  NInputNumber,
} from 'naive-ui'
import { LANGUAGE_CONFIG, EDITOR_THEME_OPTIONS, type EDITOR_THEHE, type LanguageValue, LANGUAGE_OPTIONS } from '@/constants'
import { Play, RotateCcw, Settings } from 'lucide-vue-next'
import AiAssistant from '@/components/AiAssistant.vue'
import { useLocalStorage } from '@vueuse/core'
import { RemovableRef } from '@vueuse/core'
import Request from '@/services/api'
import { useMessage } from 'naive-ui'
import indexedDBService from '@/services/indexedDB'
import { useUserStore } from '@/stores/useUserStore'
const message = useMessage()
const route = useRoute()
const userStore = useUserStore()
const codeEditor = defineAsyncComponent(() => import('@/components/AceEditor/AceEditor.vue'))
// const codeEditor = defineAsyncComponent(() => import('@/components/CodeMirror/CodeMirror.vue'))
const ProblemContext = ref('')
const fontSize = useLocalStorage('editor_font_size', 14)
const problem = ref<Problem>({
  id: 0,
  user_id: '',
  title: '',
  context: '',
  difficulty: 0,
  input_description: '',
  output_description: '',
  tags: [],
  accept: 0,
  judge_config: {
    time_limit: 1000,
    memory_limit: 65536
  },
  judge_case: [],
  judge_sample: [],
  tips: ''
})

const Language = useLocalStorage<LanguageValue>('language', 'cpp')
const Theme = useLocalStorage('editor_theme', 'chrome') as RemovableRef<EDITOR_THEHE>

// 将内部语言值转换为API需要的值
const languageToApi = (lang: LanguageValue): string => {
  return LANGUAGE_CONFIG[lang].apiValue
}

import {difficultyMap} from '@/constants'
import { formatAcceptance } from '@/utils/format'
import { Problem } from '@/types/problem'
import { problemApi } from '@/services/problem'
// 使用 IndexedDB 存储代码
const code = ref('')
const testCaseInput = ref('')
const testCaseExpected = ref('')

const result = ref<string | null>(null)
const isRunning = ref(false)
const isTesting = ref(false)
const hoverBgColor2 = ref('transparent')

const test_case = ref({
  input: '1 2',
  expected: '3',
  output: ''
})

// 加载指定语言的代码
const loadCodeForLanguage = async (language: LanguageValue) => {
  try {
    const codeRecord = await indexedDBService.getCode(String(route.params.id), language)
    if (codeRecord) {
      code.value = codeRecord.code
    } else {
      // 如果没有保存的代码，使用默认模板
      code.value = indexedDBService.getDefaultCode(language)
    }
  } catch (error) {
    console.error(`Failed to load code for language ${language}:`, error)
    // 出错时使用默认模板

    code.value = indexedDBService.getDefaultCode(language)
  }
}

// 初始化 IndexedDB 并加载数据
const initIndexedDB = async () => {
  try {
    await indexedDBService.init()

    // 加载当前语言的代码
    await loadCodeForLanguage(Language.value)

    // 加载保存的测试用例
    const testCaseRecord = await indexedDBService.getTestCase(String(route.params.id))
    if (testCaseRecord) {
      test_case.value.input = testCaseRecord.input
      test_case.value.expected = testCaseRecord.expected
      testCaseInput.value = testCaseRecord.input
      testCaseExpected.value = testCaseRecord.expected
    }
  } catch (error) {
    console.error('Failed to load data from IndexedDB:', error)
  }
}

// 保存代码到 IndexedDB (按语言分别保存)
const saveCodeToDB = async (language: LanguageValue) => {
  try {
    await indexedDBService.saveCode(String(route.params.id), language, code.value)
  } catch (error) {
    console.error('Failed to save code to IndexedDB:', error)
  }
}

// 保存测试用例到 IndexedDB
const saveTestCaseToDB = async () => {
  try {
    await indexedDBService.saveTestCase(
      String(route.params.id),
      test_case.value.input,
      test_case.value.expected
    )
  } catch (error) {
    console.error('Failed to save test case to IndexedDB:', error)
  }
}

// 监听代码变化，自动保存
watch(code, () => {
  saveCodeToDB(Language.value)
})

// 监听测试用例变化，自动保存
watch(
  () => [test_case.value.input, test_case.value.expected],
  () => {
    saveTestCaseToDB()
  }
)

// 监听语言变化，切换到对应语言的代码
watch(Language, async (newLanguage, oldLanguage) => {
  // 保存旧语言的代码
  if (oldLanguage) {
    await saveCodeToDB(oldLanguage)
  }
  // 加载新语言的代码
  await loadCodeForLanguage(newLanguage)
})

onMounted(async () => {
  // 初始化 IndexedDB
  await initIndexedDB()
  if(route.params.id == undefined) return
  await problemApi.getProblemDetail(route.params.id as string)
    .then((res) => {
      const {info} = res
      problem.value = info
      // 如果有样例，使用第一个样例作为默认测试用例
      if (info.judge_sample && info.judge_sample.length > 0 && !testCaseInput.value) {
        test_case.value.input = info.judge_sample[0].input
        test_case.value.expected = info.judge_sample[0].expected
        test_case.value = JSON.parse(JSON.stringify(test_case.value))
        testCaseInput.value = info.judge_sample[0].input
        testCaseExpected.value = info.judge_sample[0].expected
      }
    })
    .finally(() => {
      var res = ''
      problem.value.judge_sample.forEach((item, index) => {
        res += `### 样例输入${index + 1}\n\n\`\`\`\n${item.input} \n\`\`\`\n`
        res += `### 样例输出${index + 1}\n\n\`\`\`\n${item.expected}\n \`\`\`\n`
        if (index != problem.value.judge_sample.length - 1) {
          res += `--- \n\n`
        }
        ProblemContext.value =
          `## 题目描述\n${problem.value.context} \n\n` +
          `## 输入格式\n${problem.value.input_description}\n\n` +
          `## 输出格式\n${problem.value.output_description}\n\n ` +
          `## 样例\n${res}\n\n` +
          `## 提示\n ${problem.value.tips} \n\n`
      })
    })
})

const handleTest = async () => {
  isTesting.value = true
  result.value = null
  test_case.value.output = ''
  await Request.post('/ide/submit', {
    submission_id: 123456789,
    code: code.value,
    language: languageToApi(Language.value),
    test_cases: [
      {
        case_id: 1,
        stdin: test_case.value.input,
        expected: test_case.value.expected
      }
    ],
    message: '',
    seccomp_profile: '',
    resources_limits: {
      cpu_time: 100000,
      memory_bytes: 67108864,
      stack_bytes: 10485760,
      output_bytes: 10485760
    },
    timeout: 10 * 1000
  })
    .then((res) => {
      if (res.info.verdict == 'CompilationError') {
        test_case.value.output = res.info.result[0].stderr ?? ''
      } else if (res.info.verdict == 'WrongAnswer' || res.info.verdict == 'Accepted') {
        test_case.value.output = res.info.result[0].stdout
      }
    })
    .catch((err) => {
      console.log(err)
    })
    .finally(() => (isTesting.value = false))
}

const handleRun = async () => {
  if (!userStore.id) {
    message.error('请先登录')
    return
  }

  if (!code.value.trim()) {
    message.error('请输入代码')
    return
  }

  isRunning.value = true
  result.value = null

  try {
    await Request.post('/problem/submit', {
      problem_id: String(route.params.id),
      user_id: userStore.id,
      code: code.value,
      language: languageToApi(Language.value)
    })
      .then((res) => {
        message.success('代码提交成功')
        // 可以根据返回结果更新UI
        console.log('Run result:', res)
      })
      .catch((err) => {
        console.error('Run failed:', err)
        message.error('运行失败: ' + (err.response?.data?.message || err.message))
      })
  } finally {
    isRunning.value = false
  }
}

const handleEditorChange = (newCode: string) => {
  code.value = newCode
}
onMounted(()=>{
      // 禁止 body 和 html 滚动
    document.body.style.overflow = 'hidden'
    document.documentElement.style.overflow = 'hidden'

})
onUnmounted(() => {
    // 恢复 body 和 html 滚动
    document.body.style.overflow = ''
    document.documentElement.style.overflow = ''
})
</script>
<style scoped>
.markdown-container {
  overflow-y: auto;
  overflow-x: hidden;
  border-radius: 0.5rem;
  flex-shrink: 0;
}
</style>
