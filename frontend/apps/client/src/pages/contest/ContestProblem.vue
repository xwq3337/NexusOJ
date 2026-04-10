<script setup lang="ts">
import { defineAsyncComponent, onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
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
  NTag,
  NIcon,
  NSpace
} from 'naive-ui'
import { useMessage } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import type { RemovableRef } from '@vueuse/core'
import {
  LANGUAGE_CONFIG,
  EDITOR_THEME_OPTIONS,
  type EDITOR_THEHE,
  type LanguageValue,
  difficultyMap,
  STATUS_COLORS
} from '@/constants'
import { Play, RotateCcw, Settings, Target, Tag, Clock, Cpu, ArrowLeft } from 'lucide-vue-next'
import { contestApi, ideApi } from '@nexusoj/server'
import { useUserStore } from '@/stores/useUserStore'
import indexedDBService from '@/services/indexedDB'
import { formatAcceptance } from '@/utils/format'

const codeEditor = defineAsyncComponent(() => import('@/components/AceEditor/AceEditor.vue'))
const message = useMessage()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const contestId = route.params.id as string
const problemLabel = route.params.label as string

// 编辑器设置
const fontSize = useLocalStorage('editor_font_size', 14)
const Language = useLocalStorage<LanguageValue>('language', 'cpp')
const Theme = useLocalStorage('editor_theme', 'chrome') as RemovableRef<EDITOR_THEHE>

const languageToApi = (lang: LanguageValue): string => {
  return LANGUAGE_CONFIG[lang].apiValue
}

// 数据
const contest = ref<any>(null)
const contestProblem = ref<any>(null) // 比赛题目信息（label, score, my_status）
const problem = ref<any>({
  title: '',
  context: '',
  difficulty: 0,
  input_description: '',
  output_description: '',
  tags: [],
  accept: 0,
  submission: 0,
  judge_config: { time_limit: 1000, memory_limit: 256 },
  judge_sample: [],
  tips: ''
})
const ProblemContext = ref('')

// 编辑器
const code = ref('')
const test_case = ref({ input: '', expected: '', output: '' })
const result = ref<string | null>(null)
const isRunning = ref(false)
const isTesting = ref(false)
const hoverBgColor2 = ref('transparent')

// 倒计时
const now = ref(Date.now())
let countdownInterval: ReturnType<typeof setInterval> | null = null

const countdown = computed(() => {
  if (!contest.value?.end_at) return '--:--:--'
  const endAt = new Date(contest.value.end_at).getTime()
  const beginAt = new Date(contest.value.begin_at).getTime()
  const current = now.value

  let diff: number
  if (current >= beginAt && current < endAt) {
    diff = endAt - current
  } else if (current < beginAt) {
    diff = beginAt - current
  } else {
    return '00:00:00'
  }

  if (diff <= 0) return '00:00:00'
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})

const contestStatus = computed(() => {
  if (!contest.value) return ''
  return contest.value.status
})

const isContestLive = computed(() => contestStatus.value === 'Live')

// IndexedDB key: 避免与普通题目冲突
const dbKey = `contest_${contestId}_${problemLabel}`

const loadCodeForLanguage = async (language: LanguageValue) => {
  try {
    const codeRecord = await indexedDBService.getCode(dbKey, language)
    if (codeRecord) {
      code.value = codeRecord.code
    } else {
      code.value = indexedDBService.getDefaultCode(language)
    }
  } catch {
    code.value = indexedDBService.getDefaultCode(language)
  }
}

const saveCodeToDB = async (language: LanguageValue) => {
  try {
    await indexedDBService.saveCode(dbKey, language, code.value)
  } catch (e) {
    console.error('Failed to save code:', e)
  }
}

const saveTestCaseToDB = async () => {
  try {
    await indexedDBService.saveTestCase(dbKey, test_case.value.input, test_case.value.expected)
  } catch (e) {
    console.error('Failed to save test case:', e)
  }
}

watch(code, () => saveCodeToDB(Language.value))
watch(() => [test_case.value.input, test_case.value.expected], () => saveTestCaseToDB())
watch(Language, async (newLang, oldLang) => {
  if (oldLang) await saveCodeToDB(oldLang)
  await loadCodeForLanguage(newLang)
})

// 数据获取
const fetchData = async () => {
  // 获取完整题目详情
  try {
    const {code, info} = await contestApi.getContestProblemDetail(contestId, problemLabel)
    if (code === 200 && info) {
      problem.value = info.problem
      contest.value = info.contest
      buildMarkdownContent()
      // 预填样例
      if (info.problem.judge_sample?.length > 0) {
        test_case.value.input = info.problem.judge_sample[0].input
        test_case.value.expected = info.problem.judge_sample[0].expected
      }
    }
  } catch (e) {
    console.error(e)
  }
}

const buildMarkdownContent = () => {
  let sampleContent = ''
  problem.value.judge_sample?.forEach((item: any, index: number) => {
    sampleContent += `### 样例输入${index + 1}\n\n\`\`\`\n${item.input}\n\`\`\`\n`
    sampleContent += `### 样例输出${index + 1}\n\n\`\`\`\n${item.expected}\n\`\`\`\n`
    if (index !== problem.value.judge_sample.length - 1) {
      sampleContent += `---\n\n`
    }
  })
  ProblemContext.value =
    `## 题目描述\n${problem.value.context}\n\n` +
    `## 输入格式\n${problem.value.input_description}\n\n` +
    `## 输出格式\n${problem.value.output_description}\n\n` +
    `## 样例\n${sampleContent}\n\n` +
    `## 提示\n${problem.value.tips}\n\n`
}

// 自测运行
const handleTest = async () => {
  if (!code.value.trim()) {
    message.warning('请输入代码')
    return
  }
  isTesting.value = true
  result.value = null
  test_case.value.output = ''
  try {
    const response = await ideApi.RunCode({
      submission_id: Date.now(),
      code: code.value,
      language: languageToApi(Language.value),
      test_cases: [{ case_id: 1, stdin: test_case.value.input, expected: test_case.value.expected }],
      message: '',
      seccomp_profile: ''
    })
    if (response.info) {
      if (response.info.verdict === 'WrongAnswer' || response.info.verdict === 'Accepted') {
        test_case.value.output = response.info.result[0]?.stdout ?? response.info.result[0]?.stderr ?? ''
      } else {
        test_case.value.output = response.info.result[0]?.stderr ?? ''
      }
      result.value = response.info.verdict
    }
  } catch (e) {
    console.error(e)
  }
  isTesting.value = false
}

// 提交代码（比赛专属 API）
const handleRun = async () => {
  if (!userStore.id) {
    message.error('请先登录')
    return
  }
  if (!code.value.trim()) {
    message.error('请输入代码')
    return
  }
  if (!isContestLive.value) {
    message.error('比赛未在进行中，无法提交')
    return
  }

  isRunning.value = true
  result.value = null
  try {
    const res = await contestApi.submitContestProblem(contestId, {
      problem_label: problemLabel,
      code: code.value,
      language: languageToApi(Language.value)
    })
    if (res.code === 200 && res.info) {
      result.value = res.info.verdict
      message.success('提交成功')
    }
  } catch (e: any) {
    message.error('提交失败: ' + (e?.message || '未知错误'))
  }
  isRunning.value = false
}

const handleEditorChange = (newCode: string) => {
  code.value = newCode
}

const goBack = () => {
  router.push({ name: 'ContestProblems', params: { id: contestId } })
}

onMounted(async () => {
  document.body.style.overflow = 'hidden'
  document.documentElement.style.overflow = 'hidden'

  await indexedDBService.init()
  await fetchData()
  await loadCodeForLanguage(Language.value)

  // 加载保存的测试用例
  const testCaseRecord = await indexedDBService.getTestCase(dbKey)
  if (testCaseRecord) {
    test_case.value.input = testCaseRecord.input
    test_case.value.expected = testCaseRecord.expected
  }

  countdownInterval = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  document.body.style.overflow = ''
  document.documentElement.style.overflow = ''
  if (countdownInterval) clearInterval(countdownInterval)
})
</script>

<template>
  <div class="flex h-screen overflow-hidden">
    <n-split direction="horizontal" :max="0.75" :min="0.25">
      <!-- 左侧：题目描述 -->
      <template #1>
        <div class="h-full overflow-y-auto markdown-container" :style="{
          backgroundColor: 'var(--surface-secondary)',
          color: 'var(--text-primary)'
        }">
          <!-- 题目头部 -->
          <div class="p-4 border-b" :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }">
            <div class="flex items-center gap-3 mb-2">
              <NButton quaternary circle size="small" @click="goBack">
                <template #icon><NIcon :size="18"><ArrowLeft /></NIcon></template>
              </NButton>
              <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">
                {{ problemLabel }}. {{ problem.title || '加载中...' }}
              </h1>
              <NTag v-if="contestStatus" :type="isContestLive ? 'success' : contestStatus === 'Upcoming' ? 'warning' : 'default'" size="small" round>
                {{ isContestLive ? '比赛中' : contestStatus === 'Upcoming' ? '未开始' : '已结束' }}
              </NTag>
              <span v-if="isContestLive" class="font-mono text-sm" :style="{ color: 'var(--contest-timer-value)' }">
                {{ countdown }}
              </span>
            </div>
            <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm">
              <span
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md font-medium text-xs"
                :class="difficultyMap[Math.max(0, Number(problem.difficulty) - 1)]?.color"
                :style="{ backgroundColor: 'var(--surface-tertiary)' }"
              >
                {{ difficultyMap[Math.max(0, Number(problem.difficulty) - 1)]?.text }}
              </span>
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Target :size="14" />
                <span>{{ formatAcceptance(problem.accept, problem.submission) }}</span>
              </span>
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Tag :size="14" />
                <span>{{ problem.tags?.join(', ') || '暂无标签' }}</span>
              </span>
              <span v-if="problem.judge_config?.time_limit" class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Clock :size="14" />
                <span>{{ problem.judge_config.time_limit }}s</span>
              </span>
              <span v-if="problem.judge_config?.memory_limit" class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Cpu :size="14" />
                <span>{{ problem.judge_config.memory_limit }}MB</span>
              </span>
            </div>
          </div>
          <!-- Markdown 内容 -->
          <v-md-preview :text="ProblemContext" style="min-height: 50rem" :style="{
            padding: '5px',
            backgroundColor: 'transparent'
          }" />
        </div>
      </template>

      <!-- 右侧：编辑器 -->
      <template #2>
        <div class="flex flex-col h-full">
          <!-- 工具栏 -->
          <div class="h-12 border-b flex items-center justify-between px-4 shrink-0" :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }">
            <NSpace :size="8" align="center">
              <NSelect
                v-model:value="Language"
                :options="Object.values(LANGUAGE_CONFIG).map(c => ({ value: c.value, label: c.label }))"
                :style="{ width: '140px' }"
                :dropdown-props="{ style: { maxHeight: '200px', overflowY: 'auto' } }"
                placeholder="选择语言"
              />
              <NPopover trigger="click" placement="bottom">
                <template #trigger>
                  <NButton :style="{ color: 'var(--text-primary)' }" class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded">
                    <NIcon><Settings :size="14" /></NIcon> 设置
                  </NButton>
                </template>
                <NGrid x-gap="12" :cols="2" :style="{ width: '15rem' }">
                  <NGi>
                    <NSelect
                      v-model:value="Theme"
                      :options="EDITOR_THEME_OPTIONS.map(i => ({ label: i, value: i }))"
                      :dropdown-props="{ style: { maxHeight: '200px', overflowY: 'auto' } }"
                      placeholder="选择主题"
                    />
                  </NGi>
                  <NGi>
                    <NInputNumber v-model:value="fontSize" :update-value-on-input="false" :min="10" :max="30" />
                  </NGi>
                </NGrid>
              </NPopover>
            </NSpace>

            <NSpace :size="8" align="center">
              <NButton
                :style="{ color: 'var(--text-primary)', backgroundColor: hoverBgColor2 }"
                @mouseenter="hoverBgColor2 = 'var(--surface-tertiary)'"
                @mouseleave="hoverBgColor2 = 'transparent'"
                @click="code = indexedDBService.getDefaultCode(Language)"
              >
                <template #icon><NIcon><RotateCcw :size="14" /></NIcon></template>
                重置
              </NButton>
              <NButton :disabled="isTesting" @click="handleTest">
                <template #icon><NIcon><Play :size="14" /></NIcon></template>
                {{ isTesting ? '测试中...' : '自测运行' }}
              </NButton>
              <NButton type="success" :disabled="isRunning || !isContestLive" @click="handleRun">
                {{ isRunning ? '提交中...' : '提交代码' }}
              </NButton>
            </NSpace>
          </div>

          <!-- 编辑器 + 控制台 -->
          <n-split direction="vertical" :default-size="0.7" :min="0.3" :max="0.8">
            <template #1>
              <div class="h-full font-mono">
                <codeEditor @change="handleEditorChange" :value="code" :theme="Theme" :language="Language" />
              </div>
            </template>
            <template #2>
              <div class="border-t px-2 py-0 font-mono text-sm overflow-y-auto h-full" :style="{
                backgroundColor: 'var(--surface-secondary)',
                borderColor: 'var(--border-color)',
                borderTopWidth: '1px',
                borderStyle: 'solid'
              }">
                <NTabs type="line" animated>
                  <NTabPane name="Case" tab="测试用例">
                    <NGrid x-gap="12" :cols="2">
                      <NGi>
                        <NInput type="textarea" autosize v-model:value="test_case.input" placeholder="输入" />
                      </NGi>
                      <NGi class="overflow-y-scroll">
                        <NInput type="textarea" autosize disabled v-model:value="test_case.output" placeholder="输出" />
                      </NGi>
                    </NGrid>
                  </NTabPane>
                  <NTabPane name="console" tab="控制台">
                    <div v-if="result" class="p-2 rounded text-sm font-medium" :style="{
                      color: (STATUS_COLORS as any)[result]?.textColor,
                      backgroundColor: (STATUS_COLORS as any)[result]?.color
                    }">
                      {{ result }}
                    </div>
                    <div v-else :style="{ color: 'var(--text-tertiary)' }">等待运行...</div>
                  </NTabPane>
                </NTabs>
              </div>
            </template>
          </n-split>
        </div>
      </template>
    </n-split>
  </div>
</template>

<style scoped>
.markdown-container {
  overflow-y: auto;
  overflow-x: hidden;
  border-radius: 0.5rem;
  flex-shrink: 0;
}
</style>
