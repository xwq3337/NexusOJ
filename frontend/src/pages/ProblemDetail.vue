<template>
  <div class="flex h-screen overflow-hidden">
    <n-split direction="horizontal" :max="0.75" :min="0.25">
      <template #1>
        <div
          class="p-4 border-b"
          :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <h1 class="text-2xl font-bold mb-2" :style="{ color: 'var(--text-primary)' }">
            {{ problem.title }}
          </h1>
          <div class="flex gap-4 text-sm">
            <span
              class="font-medium"
              :class="difficultyMap[Number(problem.difficulty) - 1]?.color"
              >{{ difficultyMap[Number(problem.difficulty) - 1]?.text }}</span
            >
            <span class="text-black">通过率: {{ problem.acceptance }}</span>
            <span class="text-black">标签: {{ problem.tags.join(', ') }}</span>
            <span class="text-black"
              >判题配置: {{ problem.judge_config.time_limit }}s
              {{ problem.judge_config.memory_limit }}MB</span
            >
          </div>
        </div>
        <MarkdownPreview
          :text="ProblemContext"
          :style="{
            padding: '5px',
            height: 'calc(100vh - 128px)',
            overflowY: 'auto',
            backgroundColor: 'var(--surface-secondary)',
            color: 'var(--text-primary)'
          }"
        />
      </template>
      <template #2>
        <div
          class="h-12 border-b flex items-center justify-between px-4"
          :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <div class="flex items-center gap-2">
            <!-- TODO: 主题切换-背景 -->
            <n-select
              v-model:value="Language"
              :options="languageOptions"
              :style="{ width: '120px' }"
              :dropdown-props="{ style: { maxHeight: '200px', overflowY: 'auto' } }"
              placeholder="选择语言"
            >
            </n-select>
            <n-popover trigger="click" placement="bottom">
              <template #trigger>
                <n-button
                  :style="{ color: 'var(--text-primary)' }"
                  class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded transition-colors"
                >
                  <Settings :size="14" /> 设置
                </n-button>
              </template>
              <n-grid x-gap="12" :cols="2" :style="{ width: '15rem' }">
                <n-gi>
                  <n-select
                    v-model:value="Theme"
                    :options="themeOptions"
                    :dropdown-props="{
                      style: { maxHeight: '200px', overflowY: 'auto' }
                    }"
                    placeholder="选择主题"
                  >
                  </n-select>
                </n-gi>
                <n-gi>
                  <n-input-number
                    v-model:value="fontSize"
                    :update-value-on-input="false"
                    placeholder=""
                    :min="10"
                    :max="30"
                  />
                </n-gi>
              </n-grid>
            </n-popover>
          </div>
          <div class="flex items-center gap-2">
            <button
              class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded transition-colors"
              :style="{
                color: 'var(--text-primary)',
                backgroundColor: hoverBgColor2
              }"
              @mouseenter="() => (hoverBgColor2 = 'var(--surface-tertiary)')"
              @mouseleave="() => (hoverBgColor2 = 'transparent')"
            >
              <RotateCcw :size="14" /> 重置
            </button>
            <button
              @click="handleTest"
              :disabled="isTesting"
              class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium bg-green-600 hover:bg-green-700 rounded transition-colors"
              :style="{ color: 'var(--text-primary)' }"
              :class="isTesting ? 'opacity-70 cursor-wait' : ''"
            >
              <Play :size="14" /> {{ isTesting ? '测试中...' : '测试代码' }}
            </button>
            <button
              @click="handleRun"
              :disabled="isRunning"
              class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium bg-green-600 hover:bg-green-700 rounded transition-colors"
              :style="{ color: 'var(--text-primary)' }"
              :class="isRunning ? 'opacity-70 cursor-wait' : ''"
            >
              <Play :size="14" /> {{ isRunning ? '运行中...' : '运行代码' }}
            </button>
          </div>
        </div>

        <div class="flex-1 relative font-mono" style="height: 66%">
          <codeEditor
            @change="handleEditorChange"
            :value="code"
            @update:languageOptions="
              ($event) =>
                (languageOptions = $event.map((i) => {
                  return { value: i, label: i }
                }))
            "
            @update:themeOptions="
              ($event) =>
                (themeOptions = $event.map((i) => {
                  return { value: i, label: i }
                }))
            "
            :theme="Theme"
            :language="Language"
          />
        </div>

        <div
          class="h-40 border-t px-2 py-0 font-mono text-sm overflow-y-auto"
          :style="{
            backgroundColor: 'var(--surface-secondary)',
            borderColor: 'var(--border-color)',
            borderTopWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <n-tabs type="line" animated>
            <n-tab-pane name="Case" tab="测试用例">
              <n-grid x-gap="12" :cols="2" :style="{ height: '100%' }">
                <n-gi>
                  <n-input
                    autosize
                    type="textarea"
                    v-model:value="test_case.input"
                    placeholder="输入"
                  ></n-input>
                </n-gi>
                <n-gi>
                  <n-input
                    autosize
                    disabled
                    type="textarea"
                    v-model:value="test_case.output"
                    placeholder="输出"
                  ></n-input>
                </n-gi>
              </n-grid>
            </n-tab-pane>
            <n-tab-pane name="console" tab="控制台"> {{ result }}(fake) </n-tab-pane>
          </n-tabs>
        </div>
      </template>
    </n-split>
    <AiAssistant />
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  NSplit,
  NPopover,
  NButton,
  NSelect,
  NTabPane,
  NTabs,
  NInput,
  NCol,
  NGrid,
  NGi,
  NInputNumber,
  NTag
} from 'naive-ui'
import { Play, RotateCcw, Settings } from 'lucide-vue-next'
import AiAssistant from '@/components/AiAssistant.vue'
import { useLocalStorage } from '@vueuse/core'
import { RemovableRef } from '@vueuse/core'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import Request from '@/services/api'
import { useMessage } from 'naive-ui'
import { TagColor } from 'naive-ui/es/tag/src/common-props'
const message = useMessage()
const route = useRoute()
const codeEditor = defineAsyncComponent(() => import('@/components/AceEditor/AceEditor.vue'))
const languageOptions = ref([])
const themeOptions = ref([])
const ProblemContext = ref('')
const fontSize = useLocalStorage('editor_font_size', 14)
const problem = reactive({
  id: route.params.id,
  title: '',
  context: '',
  difficulty: '',
  acceptance: '',
  input_description: '',
  output_description: '',
  judge_config: {} as { time_limit: number; memory_limit: number },
  judge_sample: [],
  tip: '',
  tags: []
})
onMounted(async () => {
  await Request.get('/problem/' + route.params.id)
    .then((res) => {
      problem.id = res.info.id
      problem.title = res.info.title
      problem.difficulty = res.info.difficulty
      problem.acceptance = res.info.accept
      problem.context = res.info.context
      problem.input_description = res.info.input_description
      problem.output_description = res.info.output_description
      problem.judge_config = res.info.judge_config
      problem.judge_sample = res.info.judge_sample
      problem.tip = res.info.tip
      problem.tags = res.info.tags
    })
    .finally(() => {
      var res = ''
      problem.judge_sample.forEach((item, index) => {
        res += `### 样例输入${index + 1}\n\n\`\`\`\n${item.input} \n\`\`\`\n`
        res += `### 样例输出${index + 1}\n\n\`\`\`\n${item.expected}\n \`\`\`\n`
        if (index != problem.judge_sample.length - 1) {
          res += `--- \n\n`
        }
        ProblemContext.value =
          `## 题目描述\n${problem.context} \n\n` +
          `## 输入格式\n${problem.input_description}\n\n` +
          `## 输出格式\n${problem.output_description}\n\n ` +
          `## 样例\n${res}\n\n` +
          `## 提示\n ${problem.tip} \n\n`
      })
    })
})
const Language = useLocalStorage('language', 'c_cpp') as RemovableRef<
  'c_cpp' | 'javascript' | 'sql' | 'text' | 'python' | 'java' | 'csharp'
>
const Theme = useLocalStorage('editor_theme', 'chrome') as RemovableRef<
  'chrome' | 'github' | 'monokai' | 'xcode' | 'dracula' | 'clouds' | 'terminal'
>

const difficultyMap = [
  { text: '简单', color: 'text-green-400', type: 'success' },
  { text: '容易', color: 'text-yellow-400', type: 'warning' },
  { text: '中等', color: 'text-orange-400', type: 'info' },
  { text: '困难', color: 'text-red-400', type: 'error' },
  { text: '极难', color: 'text-purple-400', type: 'error' }
]
const code = useLocalStorage(
  'code' + route.params.id,
  `#include <iostream>
using namespace std;

int main() {
    // Your code here
    return 0;
}`
)

const result = ref<string | null>(null)
const isRunning = ref(false)
const isTesting = ref(false)
const hoverBgColor2 = ref('transparent')

const test_case = ref({
  input: '1 2',
  expected: '3',
  output: ''
})
const language_translator = (s: string) => {
  const map = {
    c_cpp: 'cpp',
    javascript: 'javascript',
    sql: 'sql',
    text: 'plaintext',
    python: 'python',
    java: 'java',
    csharp: 'csharp'
  }
  return map[s]
}
const handleTest = async () => {
  isTesting.value = true
  result.value = null
  isTesting.value = true
  test_case.value.output = ''
  await Request.post('/ide/submit', {
    submission_id: 123456789,
    code: code.value,
    language: language_translator(Language.value),
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

const handleRun = () => {
  isRunning.value = true
  result.value = null
  setTimeout(() => {
    isRunning.value = false
    const outcomes = ['Accepted', 'Wrong Answer', 'Time Limit Exceeded']
    const randomOutcome = outcomes[Math.floor(Math.random() * outcomes.length)]
    result.value = randomOutcome
  }, 1500)
}

const handleEditorChange = (newCode: string) => {
  code.value = newCode
}
</script>
