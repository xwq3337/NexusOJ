<script setup lang="ts">
import { ref, computed, onMounted, inject } from 'vue'
import { h } from 'vue'
import {
  NDataTable,
  NTag,
  NPagination,
  NCard,
  NSelect,
  NButton,
  NResult
} from 'naive-ui'
import { STATUS_OPTIONS, STATUS_COLORS, LANGUAGE_CONFIG } from '@/constants'
import { contestApi } from '@nexusoj/server'
import { formatMemory, formatDate, formatTime } from '@/utils/format'
import type { ContestRecordItem } from '@nexusoj/type'
import { useRouter } from 'vue-router'

const router = useRouter()
const { contestId, problems } = inject<any>('contestData')

const loading = ref(false)
const records = ref<ContestRecordItem[]>([])
const totalRecords = ref(0)
const currentPage = ref(1)
const pageSize = ref(15)
const statusFilter = ref<string | null>(null)
const languageFilter = ref<string | null>(null)
const problemFilter = ref<string | null>(null)

const problemOptions = computed(() =>
  (problems.value || []).map((p: any) => ({
    label: `${p.label}. ${p.title}`,
    value: p.label
  }))
)

const languageOptions = Object.keys(LANGUAGE_CONFIG).map((lang) => ({
  label: LANGUAGE_CONFIG[lang as keyof typeof LANGUAGE_CONFIG].label,
  value: lang
}))

const fetchRecords = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (statusFilter.value) params.verdict = statusFilter.value
    if (languageFilter.value) params.language = languageFilter.value
    if (problemFilter.value) params.problem_label = problemFilter.value
    const { code, info } = await contestApi.getContestSubmissions(contestId, params)
    if (code === 200 && info) {
      records.value = info.list || []
      totalRecords.value = info.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const columns = [
  {
    title: '提交ID',
    key: 'id',
    width: 180,
    render: (row: ContestRecordItem) =>
      h(
        'span',
        {
          class: 'cursor-pointer font-medium',
          style: { color: 'var(--accent-color)' },
          onClick: () => router.push({ name: 'ContestSubmissionDetail', params: { id: contestId.value, rid: row.id } })
        },
        row.id
      )
  },
  {
    title: '状态',
    key: 'verdict',
    width: 100,
    render: (row: ContestRecordItem) => {
      const colors = STATUS_COLORS[row.verdict as keyof typeof STATUS_COLORS]
      return h(
        NTag,
        {
          size: 'small',
          bordered: false,
          style: colors
            ? {
              backgroundColor: colors.color,
              color: colors.textColor,
              border: `1px solid ${colors.borderColor}`
            }
            : {}
        },
        { default: () => row.verdict }
      )
    }
  },
  {
    title: '题目',
    key: 'problem_title',
    ellipsis: true,
    render: (row: ContestRecordItem) =>
      h('div', { class: 'flex items-center gap-2' }, [
        h('span', {
          class: 'shrink-0 inline-flex items-center justify-center w-6 h-6 rounded-md text-xs font-bold',
          style: { background: 'var(--accent-color)', color: '#fff' }
        }, row.problem_label),
        h('span', {
          style: { color: 'var(--text-primary)' },
          class: 'truncate'
        }, row.problem_title)
      ])
  },
  {
    title: '用户',
    key: 'username',
    ellipsis: true,
    render: (row: ContestRecordItem) =>
      h('div', { 
        class: 'flex items-center gap-2' ,
        onClick : () => router.push({ name: 'UserHomePage', params: { id: row.user_id } }),
      }, row.username)
  },
  {
    title: '语言',
    key: 'language',
    width: 130,
    render: (row: ContestRecordItem) => {
      const config = LANGUAGE_CONFIG[row.language as keyof typeof LANGUAGE_CONFIG]
      return h(
        NTag,
        {
          size: 'small',
          bordered: false,
          type: config ? undefined : 'default',
          style: config
            ? {
              backgroundColor: config.color.color,
              color: config.color.textColor,
              border: `1px solid ${config.color.borderColor}`
            }
            : {}
        },
        { default: () => config?.label || row.language }
      )
    }
  },
  {
    title: '耗时',
    key: 'max_time',
    width: 90,
    render: (row: ContestRecordItem) => h('span', { style: 'color: var(--text-secondary)' }, formatTime(row.max_time))
  },
  {
    title: '内存',
    key: 'max_memory',
    width: 90,
    render: (row: ContestRecordItem) => h('span', { style: 'color: var(--text-secondary)' }, formatMemory(row.max_memory))
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 220,
    render: (row: ContestRecordItem) => h('span', { style: 'color: var(--text-tertiary)' }, formatDate(row.created_at))
  }
]

const resetFilters = () => {
  statusFilter.value = null
  languageFilter.value = null
  problemFilter.value = null
  currentPage.value = 1
  fetchRecords()
}

onMounted(() => {
  fetchRecords()
})
</script>

<template>
  <div>
    <!-- Filters -->
    <div class="flex flex-row gap-3 mb-4 items-center">
      <NSelect v-model:value="statusFilter" :options="STATUS_OPTIONS" placeholder="状态筛选" clearable
        @update:value="currentPage = 1; fetchRecords()" />
      <NSelect v-model:value="languageFilter" :options="languageOptions" placeholder="语言筛选" clearable
        @update:value="currentPage = 1; fetchRecords()" />
      <NSelect v-model:value="problemFilter" :options="problemOptions" placeholder="题目筛选" clearable
        @update:value="currentPage = 1; fetchRecords()" />
      <NButton @click="resetFilters">重置</NButton>
    </div>
    <!-- Table -->
    <NCard content-style="padding: 0;">
      <NDataTable :columns="columns" :data="records" :loading="loading" size="small"
        :row-key="(row: ContestRecordItem) => row.id" />
      <div class="flex justify-end p-4">
        <NPagination v-model:page="currentPage" v-model:page-size="pageSize" :item-count="totalRecords"
          :page-sizes="[10, 15, 20, 50]" show-size-picker @update:page="fetchRecords"
          @update:page-size="currentPage = 1; fetchRecords()" />
      </div>
    </NCard>
  </div>
</template>
