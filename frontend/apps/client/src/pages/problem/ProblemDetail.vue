<template>
  <div class="flex h-screen overflow-hidden" :style="{ backgroundColor: 'var(--bg-primary)' }">
    <n-split direction="horizontal" :max="0.75" :min="0.25">
      <template #1>
        <div ref="problemPanelRef" class="h-full overflow-y-auto markdown-container" :style="{
          backgroundColor: 'var(--surface-secondary)',
          color: 'var(--text-primary)'
        }">
          <div class="p-4 border-b" :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }">
            <h1 class="text-2xl font-bold mb-2 flex items-center gap-3" :style="{ color: 'var(--text-primary)' }">
              {{ problem.title }}
              <n-button size="tiny" type="primary" ghost @click="openSolutionDrawer">
                <template #icon>
                  <BookOpen :size="14" />
                </template>
                题解
              </n-button>
              <n-button size="tiny" type="info" ghost @click="openRecordDrawer">
                <template #icon>
                  <History :size="14" />
                </template>
                提交记录
              </n-button>
            </h1>
            <div class="flex flex-wrap items-center gap-x-4 gap-y-2 text-sm">
              <!-- 难度标签 -->
              <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md font-medium text-xs"
                :class="difficultyMap[Number(problem.difficulty) - 1]?.color" :style="{
                  backgroundColor: `var(--surface-tertiary)`
                }">
                {{ difficultyMap[Number(problem.difficulty) - 1]?.text }}
              </span>

              <!-- 通过率 -->
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Target :size="14" />
                <span>{{ formatAcceptance(problem.accept, problem.submission) }}</span>
              </span>

              <!-- 标签 -->
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Tag :size="14" />
                <span>{{ problem.tags.join(', ') || '暂无标签' }}</span>
              </span>

              <!-- 时间限制 -->
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Clock :size="14" />
                <span>{{ problem.judge_config.time_limit }}s</span>
              </span>

              <!-- 内存限制 -->
              <span class="inline-flex items-center gap-1.5" :style="{ color: 'var(--text-secondary)' }">
                <Cpu :size="14" />
                <span>{{ problem.judge_config.memory_limit }}MB</span>
              </span>
            </div>
          </div>
          <v-md-preview :text="ProblemContext" style="height: 50rem;" :style="{
            padding: '5px',
            backgroundColor: 'transparent',
          }" />
        </div>
      </template>
      <template #2>
        <!-- TODO移动端的工具栏，应采用flex布局 -->
        <div class="@container flex flex-col h-full">
          <div class="h-12 border-b flex items-center justify-between px-4 shrink-0" :style="{
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
                    <Settings :size="14" /> <span class="@md:inline hidden">设置</span>
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
            <!-- 操作按钮(重置，运行，提交) -->
            <div class="flex items-center gap-2">
              <n-button class="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded transition-colors"
                :style="{
                  color: 'var(--text-primary)',
                  backgroundColor: hoverBgColor2
                }" @mouseenter="() => (hoverBgColor2 = 'var(--surface-tertiary)')"
                @mouseleave="() => (hoverBgColor2 = 'transparent')"
                @click="code = indexedDBService.getDefaultCode(Language)">
                <RotateCcw :size="14" /> <span class="@md:inline hidden">重置</span>
              </n-button>
              <n-button @click="handleTest" :disabled="isTesting"
                class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium rounded transition-colors"
                :style="{ color: 'var(--text-primary)' }" :class="isTesting ? 'opacity-70 cursor-wait' : ''">
                <Play :size="14" /> <span class="@md:inline hidden">{{ isTesting ? '测试中...' : '自测运行' }}</span>
              </n-button>
              <n-button @click="handleRun" :disabled="isRunning" type="success"
                class="flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium rounded transition-colors"
                :style="{ color: 'var(--text-primary)' }" :class="isRunning ? 'opacity-70 cursor-wait' : ''">
                <Send :size="14" /> <span class="@md:inline hidden">{{ isRunning ? '提交中...' : '提交代码' }}</span>
              </n-button>
            </div>
          </div>

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
                <n-tabs type="line" animated>
                  <n-tab-pane name="Case" tab="测试用例">
                    <n-grid x-gap="12" :cols="2">
                      <n-gi>
                        <n-input type="textarea" autosize v-model:value="test_case.input" placeholder="输入" />
                      </n-gi>
                      <n-gi class="overflow-y-scroll">
                        <n-input type="textarea" autosize disabled v-model:value="test_case.output" placeholder="输出" />
                      </n-gi>
                    </n-grid>
                  </n-tab-pane>
                  <n-tab-pane name="console" tab="控制台"> 模拟 {{ result }} </n-tab-pane>
                </n-tabs>
              </div>
            </template>
          </n-split>
        </div>
      </template>
    </n-split>
    <AiAssistant />

    <!-- 题解抽屉 -->
    <NDrawer v-model:show="showSolutionDrawer" :width="drawerWidth" placement="left">
      <NDrawerContent :title="currentSolution ? currentSolution.title : '题解列表'" closable>
        <!-- 题解详情 -->
        <template v-if="currentSolution">
          <div class="flex items-center gap-2 mb-3">
            <n-button size="small" quaternary @click="backToList" class="mb-3">
              <template #icon>
                <ArrowLeft :size="14" />
              </template>
              返回列表
            </n-button>
            <NAvatar v-if="currentSolution.avatar" round :size="28" :src="currentSolution.avatar" />
            <span class="text-sm" :style="{ color: 'var(--text-secondary)' }">
              {{ currentSolution.username }}
            </span>
            <span class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
              · {{ currentSolution.view }} 阅读 · {{ currentSolution.like }} 点赞
            </span>
          </div>
          <NDivider style="margin: 8px 0" />
          <v-md-preview v-if="currentSolution.context" :text="currentSolution.context" />
        </template>

        <!-- 题解列表 -->
        <template v-else>
          <div class="flex justify-end mb-3">
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button size="small" type="primary" :disabled="!hasAccepted" @click="goCreateSolution">
                  <template #icon>
                    <Plus :size="14" />
                  </template>
                  发布题解
                </n-button>
              </template>
              {{ hasAccepted ? '发布你的题解' : '通过本题后才能发布题解' }}
            </n-tooltip>
          </div>
          <NSpin :show="solutionLoading">
            <NEmpty v-if="!solutionLoading && solutionList.length === 0" description="暂无题解" />
            <NList v-else>
              <NListItem v-for="item in solutionList" :key="item.id"
                class="cursor-pointer hover:bg-(--surface-tertiary) rounded-lg p-2!"
                @click="viewSolutionDetail(item.id)">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium text-xl truncate flex-1" :style="{ color: 'var(--text-primary)' }">
                      {{ item.title }}
                    </span>
                    <n-tag v-if="item.username === 'NexusOJ'" size="tiny" :bordered="false"
                      style="background: linear-gradient(135deg, #f59e0b, #d97706); color: #fff; font-weight: 600; flex-shrink: 0;">
                      官方
                    </n-tag>
                  </div>
                  <div class="flex items-center gap-2 text-xs" :style="{ color: 'var(--text-primary)' }">
                    <span>{{ item.username }}</span>
                    <NDivider vertical />
                    <span>{{ item.view }} 阅读</span>
                    <NDivider vertical />
                    <span>{{ item.like }} 赞</span>
                  </div>
                  <div v-if="item.tags?.length" class="flex gap-1 mt-1">
                    <NTag v-for="tag in item.tags" :key="tag" size="tiny" type="info" :bordered="false">
                      {{ tag }}
                    </NTag>
                  </div>
                </div>
              </NListItem>
            </NList>
          </NSpin>
        </template>
      </NDrawerContent>
    </NDrawer>

    <!-- 提交记录抽屉 -->
    <NDrawer v-model:show="showRecordDrawer" :width="drawerWidth" placement="left">
      <NDrawerContent :title="currentRecord ? '提交详情' : '提交记录'" closable :style="{ color: 'var(--text-primary)' }">
        <!-- 记录详情 -->
        <template v-if="currentRecord">
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <n-button size="small" quaternary @click="backToRecordList" class="mb-3">
                <template #icon>
                  <ArrowLeft :size="14" />
                </template>
                返回列表
              </n-button>
              <n-tag size="small" :bordered="false" :style="convertToCss(STATUS_COLORS[currentRecord.verdict])">
                {{ currentRecord.verdict }}
              </n-tag>
              <n-tag size="small" :bordered="false"
                :style="convertToCss(LANGUAGE_CONFIG[currentRecord.language as keyof typeof LANGUAGE_CONFIG]?.color)">
                {{ LANGUAGE_CONFIG[currentRecord.language as keyof typeof LANGUAGE_CONFIG]?.label ||
                  currentRecord.language }}
              </n-tag>
            </div>
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div>
                <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">用户</div>
                <div class="font-medium" :style="{ color: 'var(--text-primary)' }">{{ currentRecord.username }}</div>
              </div>
              <div>
                <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">题目</div>
                <div class="font-medium" :style="{ color: 'var(--text-primary)' }">{{ currentRecord.problem_title }}
                </div>
              </div>
              <div>
                <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">耗时</div>
                <div class="font-medium" :style="{ color: 'var(--text-primary)' }">{{ formatTime(currentRecord.max_time)
                  }}</div>
              </div>
              <div>
                <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">内存</div>
                <div class="font-medium" :style="{ color: 'var(--text-primary)' }">{{
                  formatMemory(currentRecord.max_memory) }}</div>
              </div>
            </div>
            <NDivider style="margin: 8px 0" />
            <div>
              <div class="text-xs mb-1" :style="{ color: 'var(--text-tertiary)' }">提交时间</div>
              <div class="text-sm" :style="{ color: 'var(--text-secondary)' }">{{ formatDate(currentRecord.created_at)
                }}</div>
            </div>
            <div v-if="currentRecord.code">
              <div class="text-xs mb-1" :style="{ color: 'var(--text-tertiary)' }">代码</div>
              <n-card content-style="padding: 0;" :bordered="false">
                <pre class="p-3 text-sm overflow-x-auto font-mono"
                  :style="{ backgroundColor: 'var(--surface-secondary)', color: 'var(--text-primary)' }">{{ currentRecord.code
                  }}</pre>
              </n-card>
            </div>
          </div>
        </template>
        <template v-else>

          <!-- 筛选栏 -->
          <div class="flex flex-wrap gap-2 mb-4">
            <NInput v-model:value="recordUserSearch" placeholder="搜索用户" size="small" clearable style="width: 140px"
              @keydown.enter="recordPage = 1; fetchRecords()" />
            <NSelect v-model:value="recordStatusFilter" :options="STATUS_OPTIONS" placeholder="状态" size="small"
              clearable style="width: 130px" @update:value="recordPage = 1; fetchRecords()" />
            <NSelect v-model:value="recordLangFilter" :options="languageFilterOptions" placeholder="语言" size="small"
              clearable style="width: 140px" @update:value="recordPage = 1; fetchRecords()" />
            <NButton size="small" @click="resetRecordFilters">重置</NButton>
          </div>

          <!-- 记录列表 -->
          <NSpin :show="recordLoading">
            <NEmpty v-if="!recordLoading && recordList.length === 0" description="暂无提交记录" />
            <NList v-else>
              <NListItem v-for="record in recordList" :key="record.id"
                class="cursor-pointer hover:bg-(--surface-tertiary) rounded-lg p-2!"
                @click="viewRecordDetail(record.id)">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between mb-1">
                    <span class="text-xs font-medium" :style="{ color: 'var(--accent-color)' }">#{{ record.id }}</span>
                    <span class="text-xs" :style="{ color: 'var(--text-secondary)' }">{{ formatDate(record.created_at)
                      }}</span>
                  </div>
                  <div class="flex items-center gap-2 mb-1">
                    <NTag size="small" :bordered="false" :style="convertToCss(STATUS_COLORS[record.verdict])">
                      {{ record.verdict }}
                    </NTag>
                    <NTag size="small" :bordered="false"
                      :style="convertToCss(LANGUAGE_CONFIG[record.language as keyof typeof LANGUAGE_CONFIG]?.color)">
                      {{ LANGUAGE_CONFIG[record.language as keyof typeof LANGUAGE_CONFIG]?.label || record.language }}
                    </NTag>
                  </div>
                  <div class="flex items-center gap-3 text-xs" :style="{ color: 'var(--text-primary)' }">
                    <span class="font-medium">{{ record.username }}</span>
                    <span :style="{ color: 'var(--text-secondary)' }">{{ formatTime(record.max_time) }}</span>
                    <span :style="{ color: 'var(--text-secondary)' }">{{ formatMemory(record.max_memory) }}</span>
                  </div>
                </div>
              </NListItem>
            </NList>
          </NSpin>

          <!-- 分页 -->
          <div v-if="recordTotal > 0" class="flex justify-center mt-4">
            <NPagination v-model:page="recordPage" :page-size="recordPageSize" :item-count="recordTotal" size="small"
              @update:page="fetchRecords" />
          </div>
        </template>
      </NDrawerContent>
    </NDrawer>
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, onMounted, onUnmounted, reactive, ref, watch } from 'vue'
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
  NDrawer,
  NDrawerContent,
  NList,
  NListItem,
  NTag,
  NAvatar,
  NSpin,
  NEmpty,
  NDivider,
  NPagination,
  NTooltip,
} from 'naive-ui'
import { LANGUAGE_CONFIG, EDITOR_THEME_OPTIONS, type EDITOR_THEHE, type LanguageValue, LANGUAGE_OPTIONS, convertToCss, STATUS_MESSAGE } from '@/constants'
import { Play, ArrowLeft, RotateCcw, Settings, Target, Tag, Clock, Cpu, BookOpen, History, Plus, Send } from 'lucide-vue-next'
import AiAssistant from '@/components/AiAssistant.vue'
import MarkdownPreviewV2 from '@/components/MarkdownPreviewV2.vue'
import { useLocalStorage } from '@vueuse/core'
import { RemovableRef } from '@vueuse/core'
import { useMessage } from 'naive-ui'
import indexedDBService from '@/services/indexedDB'
import { useUserStore } from '@/stores/useUserStore'
const message = useMessage()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const codeEditor = defineAsyncComponent(() => import('@/components/AceEditor/AceEditor.vue'))
// const codeEditor = defineAsyncComponent(() => import('@/components/CodeMirror/CodeMirror.vue'))
const ProblemContext = ref('')
const fontSize = useLocalStorage('editor_font_size', 14)
const problem = ref<Problem>({
  id: 0,
  user_id: 0,
  title: '加载中...',
  context: '',
  difficulty: 1,
  input_description: '',
  output_description: '',
  tags: [],
  accept: 0,
  submission: 0,
  judge_config: {
    time_limit: 1,
    memory_limit: 64
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

import { difficultyMap } from '@/constants'
import { formatAcceptance } from '@/utils/format'
import { Problem } from '@nexusoj/type'
import { ideApi, problemApi } from '@nexusoj/server'
import { solutionApi } from '@nexusoj/server'
import { recordApi } from '@nexusoj/server'
import type { SolutionWithAuthor, GetRecordListResponse, GetRecordDetailResponse, JudgeVerdictType } from '@nexusoj/type'
import { STATUS_OPTIONS, STATUS_COLORS } from '@/constants'
import { formatMemory, formatDate, formatTime } from '@/utils/format'
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
  await problemApi.getProblemDetail(route.params.id as string)
    .then((res) => {
      const { info, code, msg } = res
      if (code !== 200 || !info) {
        message.error(msg || '获取题目详情失败')
        return
      }
      const { problem: pb, my_status } = info
      problem.value = pb
      hasAccepted.value = my_status === 'accepted'
      // 如果有样例，使用第一个样例作为默认测试用例
      if (pb.judge_sample && pb.judge_sample.length > 0 && !testCaseInput.value) {
        test_case.value.input = pb.judge_sample[0].input
        test_case.value.expected = pb.judge_sample[0].expected
        test_case.value = JSON.parse(JSON.stringify(test_case.value))
        testCaseInput.value = pb.judge_sample[0].input
        testCaseExpected.value = pb.judge_sample[0].expected
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
  await ideApi.RunCode({
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
  })
    .then((response) => {
      const { info, code } = response
      if (code !== 200 || !info) {
        message.error('自测运行失败')
        return
      }
      if (info.verdict == 'WrongAnswer' || info.verdict == 'Accepted') {
        test_case.value.output = info.result[0].stdout ?? info.result[0].stderr ?? ''
      } else {
        test_case.value.output = info.result[0].stderr ?? ''
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
    await problemApi.SubmitCode({
      problem_id: String(route.params.id),
      user_id: userStore.id,
      code: code.value,
      language: languageToApi(Language.value)
    })
      .then((res) => {
        const { code, info, msg } = res
        if (code === 200 && info) {
          if (info.verdict === 'Accepted') {
            hasAccepted.value = true
            message.success(STATUS_MESSAGE[info.verdict as JudgeVerdictType])
          } else {
            message.error(STATUS_MESSAGE[info.verdict as JudgeVerdictType])
          }
        }
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
onMounted(() => {
  // 禁止 body 和 html 滚动
  document.body.style.overflow = 'hidden'
  document.documentElement.style.overflow = 'hidden'

})
onUnmounted(() => {
  // 恢复 body 和 html 滚动
  document.body.style.overflow = ''
  document.documentElement.style.overflow = ''
})

// 题解抽屉
const problemPanelRef = ref<HTMLElement | null>(null)
const drawerWidth = ref(500)
const showSolutionDrawer = ref(false)
const solutionList = ref<SolutionWithAuthor[]>([])
const solutionLoading = ref(false)
const currentSolution = ref<SolutionWithAuthor | null>(null)

const fetchSolutions = async () => {
  solutionLoading.value = true
  try {
    const { code, info } = await solutionApi.getSolutions({
      problem_id: Number(route.params.id),
      page: 1,
      page_size: 50,
    })
    if (code === 200 && info) {
      solutionList.value = info.solutions || []
    }
  } catch (e) {
    console.error(e)
  } finally {
    solutionLoading.value = false
  }
}

const openSolutionDrawer = () => {
  if (problemPanelRef.value) {
    drawerWidth.value = problemPanelRef.value.offsetWidth
  }
  showSolutionDrawer.value = true
  currentSolution.value = null
  fetchSolutions()
}

const viewSolutionDetail = async (id: number) => {
  try {
    const { code, info } = await solutionApi.getSolutionDetail(id)
    if (code === 200 && info) {
      currentSolution.value = info
    }
  } catch (e) {
    console.error(e)
  }
}

const backToList = () => {
  currentSolution.value = null
}

// 检查用户是否已通过本题
const hasAccepted = ref(false)

const goCreateSolution = () => {
  router.push({
    name: 'SolutionCreate',
    query: {
      problem_id: route.params.id as string,
      problem_title: problem.value.title,
    },
  })
}

// 提交记录抽屉
const showRecordDrawer = ref(false)
const recordList = ref<GetRecordListResponse[]>([])
const recordLoading = ref(false)
const recordTotal = ref(0)
const recordPage = ref(1)
const recordPageSize = ref(15)
const recordStatusFilter = ref<string | null>(null)
const recordLangFilter = ref<string | null>(null)
const recordUserSearch = ref('')
const currentRecord = ref<GetRecordDetailResponse | null>(null)
const recordDetailLoading = ref(false)

const languageFilterOptions = Object.keys(LANGUAGE_CONFIG).map((lang) => ({
  label: LANGUAGE_CONFIG[lang as keyof typeof LANGUAGE_CONFIG].label,
  value: lang,
}))

const fetchRecords = async () => {
  recordLoading.value = true
  try {
    const params: Record<string, any> = {
      problem_id: route.params.id,
      page: recordPage.value,
      page_size: recordPageSize.value,
    }
    if (recordStatusFilter.value) params.verdict = recordStatusFilter.value
    if (recordLangFilter.value) params.language = recordLangFilter.value
    if (recordUserSearch.value) params.search = recordUserSearch.value

    const { code, info } = await recordApi.getRecordList(params)
    if (code === 200 && info) {
      recordList.value = info.data || []
      recordTotal.value = info.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    recordLoading.value = false
  }
}

const openRecordDrawer = () => {
  if (problemPanelRef.value) {
    drawerWidth.value = problemPanelRef.value.offsetWidth
  }
  recordPage.value = 1
  recordStatusFilter.value = null
  recordLangFilter.value = null
  recordUserSearch.value = ''
  showRecordDrawer.value = true
  currentRecord.value = null
  fetchRecords()
}

const viewRecordDetail = async (id: number) => {
  recordDetailLoading.value = true
  try {
    const { code, info } = await recordApi.getRecordDetail(String(id))
    if (code === 200 && info) {
      currentRecord.value = info
    }
  } catch (e) {
    console.error(e)
  } finally {
    recordDetailLoading.value = false
  }
}

const backToRecordList = () => {
  currentRecord.value = null
}

const resetRecordFilters = () => {
  recordStatusFilter.value = null
  recordLangFilter.value = null
  recordUserSearch.value = ''
  recordPage.value = 1
  fetchRecords()
}
</script>
<style scoped>
.markdown-container {
  overflow-y: auto;
  overflow-x: hidden;
  border-radius: 0.5rem;
  flex-shrink: 0;
}
</style>
