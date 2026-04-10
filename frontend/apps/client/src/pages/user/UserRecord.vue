<script setup lang="ts">
import { ref, watch, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { NDataTable, NTag, NEmpty, NSpin, NCard, NPagination, NButton, NSelect } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { useUserStore } from '@/stores/useUserStore'
import { formatMemory, formatDate, formatTime } from '@/utils/format'
import { LANGUAGE_CONFIG, STATUS_COLORS } from '@/constants'
import router from '@/router'
import { userApi } from '@nexusoj/server'
import type { GetRecordListResponse } from '@nexusoj/type'

const userStore = useUserStore()
const userId = String(userStore.id)

const PAGE_SIZE = 20
const currentPage = ref(1)
const total = ref(0)
const records = ref<GetRecordListResponse[]>([])
const loading = ref(false)
const verdictFilter = ref<string | null>(null)
const languageFilter = ref<string | null>(null)

const fetchRecords = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: PAGE_SIZE,
    }
    if (verdictFilter.value) params.verdict = verdictFilter.value
    if (languageFilter.value) params.language = languageFilter.value
    const { code, info } = await userApi.getUserRecordList(userId, params)
    if (code === 200 && info) {
      records.value = info.records || []
      total.value = info.total || 0
    }
  } catch (e) {
    console.error('Failed to fetch records:', e)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  verdictFilter.value = null
  languageFilter.value = null
  currentPage.value = 1
  fetchRecords()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  fetchRecords()
}

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 150,
    render(row) {
      return h('span', { style: { color: 'var(--text-link)', cursor: 'pointer' }, onClick: () => router.push({ name: 'RecordDetail', params: { id: row.id } }) }, row.id)
    }
  },
  {
    title: '题目',
    key: 'problem_title',
    width: 150,
    render(row) {
      return h('div',{
        style: { color: 'var(--accent-color)' },
      }, [
        h('div', { style: { fontWeight: 'bold', color: 'var(--text-primary)', cursor: 'pointer' }, onClick: () => router.push({ name: 'ProblemDetail', params: { id: row.problem_id } }) }, row.problem_title),
        h('div', { style: { fontSize: '12px', color: 'var(--text-secondary)' } }, `ID: ${row.problem_id}`)
      ])
    }
  },
  {
    title: '状态',
    key: 'verdict',
    width: 120,
    render(row) {
      return h(NTag, { size: 'small', color: STATUS_COLORS[row.verdict], bordered: false }, { default: () => row.verdict })
    }
  },
  {
    title: '语言',
    key: 'language',
    width: 120,
    render(row) {
      const lang = LANGUAGE_CONFIG[row.language as keyof typeof LANGUAGE_CONFIG]
      return h(NTag, { size: 'small', color: lang?.color }, { default: () => lang?.label || row.language })
    }
  },
  {
    title: '内存',
    key: 'max_memory',
    width: 100,
    render(row) {
      return formatMemory(row.max_memory)
    }
  },
  {
    title: '时间',
    key: 'max_time',
    width: 100,
    render(row) {
      return formatTime(row.max_time)
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return formatDate(row.created_at)
    }
  }
]

watch([verdictFilter, languageFilter], () => {
  currentPage.value = 1
  fetchRecords()
})

onMounted(() => {
  fetchRecords()
})
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="mb-6">
      <h1 class="text-2xl font-semibold mb-2" :style="{ color: 'var(--text-primary)' }">提交记录</h1>
      <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">共 {{ total }} 条记录</p>
    </div>

    <!-- 筛选 -->
    <NCard class="mb-4" :bordered="false" content-style="padding: 12px 16px;">
      <div class="flex gap-4 items-center flex-nowrap overflow-x-auto">
        <NSelect
          v-model:value="verdictFilter"
          :options="Object.entries(STATUS_COLORS).map(([value]) => ({ label: value, value }))"
          placeholder="状态筛选"
          clearable
          style="min-width: 140px"
        />
        <NSelect
          v-model:value="languageFilter"
          :options="Object.entries(LANGUAGE_CONFIG).map(([value, cfg]) => ({ label: cfg.label, value }))"
          placeholder="语言筛选"
          clearable
          style="min-width: 120px"
        />
        <NButton size="small" @click="resetFilters">重置</NButton>
      </div>
    </NCard>

    <!-- 表格 -->
    <NCard :bordered="false" content-style="padding: 0;">
      <NSpin :show="loading">
        <n-data-table :columns="columns" :data="records" :row-key="(row: any) => row.id" size="small" :loading="loading" />
        <div class="flex justify-end p-4" v-if="total > 0">
          <NPagination v-model:page="currentPage" :page-size="PAGE_SIZE" :item-count="total" @update:page="handlePageChange" />
        </div>
      </NSpin>
    </NCard>
  </div>
</template>
