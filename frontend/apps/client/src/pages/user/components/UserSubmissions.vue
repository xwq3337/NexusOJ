<script setup lang="ts">
import { ref, watch, onMounted, h } from 'vue'
import { useRouter } from 'vue-router'
import { NDataTable,NButton, NTag, NEmpty, NSpin, NPagination, NSelect, NInput } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { userApi } from '@nexusoj/server'
import { formatMemory, formatRelativeTime, formatTime } from '@/utils/format'
import { STATUS_COLORS, LANGUAGE_CONFIG } from '@/constants'
import type { GetRecordListResponse } from '@nexusoj/type'

const props = defineProps<{ userId: string }>()
const router = useRouter()

const PAGE_SIZE = 20
const currentPage = ref(1)
const total = ref(0)
const submissions = ref<GetRecordListResponse[]>([])
const loading = ref(false)
const statusFilter = ref<string | null>(null)
const languageFilter = ref<string | null>(null)

const fetchSubmissions = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: PAGE_SIZE,
    }
    if (statusFilter.value) params.verdict = statusFilter.value
    if (languageFilter.value) params.language = languageFilter.value
    const { code, info } = await userApi.getUserRecordList(props.userId, params)
    if (code === 200 && info) {
      submissions.value = info.records || []
      total.value = info.total || 0
    }
  } catch (e) {
    console.error('Failed to fetch submissions:', e)
  } finally {
    loading.value = false
  }
}

const resetFilters = () => {
  statusFilter.value = null
  languageFilter.value = null
  currentPage.value = 1
  fetchSubmissions()
}

const handleViewSubmission = (row: any) => {
  router.push({ name: 'RecordDetail', params: { id: row.id } })
}

const columns: DataTableColumns<any> = [
  {
    title: 'ID',
    key: 'id',
    width: 200,
  },
  {
    title: '题目',
    key: 'problem_title',
    width: 150,
    render(row) {
      return row.problem_title
    }
  },
  {
    title: '状态',
    key: 'verdict',
    width: 130,
    render(row) {
      return h(NTag, { size: 'small', color: STATUS_COLORS[row.verdict] }, { default: () => row.verdict })
    }
  },
  {
    title: '语言',
    key: 'language',
    width: 100,
    render(row) {
      const lang = LANGUAGE_CONFIG[row.language as keyof typeof LANGUAGE_CONFIG]
      return lang ? lang.label : row.language
    }
  },
  {
    title: '时间',
    key: 'max_time',
    width: 80,
    render(row) {
      return formatTime(row.max_time)
    }
  },
  {
    title: '内存',
    key: 'max_memory',
    width: 80,
    render(row) {
      return formatMemory(row.max_memory)
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 140,
    render(row) {
      return formatRelativeTime(row.created_at)
    }
  }
]

watch([statusFilter, languageFilter], () => {
  currentPage.value = 1
  fetchSubmissions()
})

onMounted(() => {
  fetchSubmissions()
})
</script>

<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex gap-4 items-center flex-nowrap overflow-x-auto">
      <NSelect
        v-model:value="statusFilter"
        :options="Object.entries(STATUS_COLORS).map(([value]) => ({ label: value, value }))"
        placeholder="状态筛选"
        clearable
        size="small"
        style="min-width: 140px"
      />
      <NSelect
        v-model:value="languageFilter"
        :options="Object.entries(LANGUAGE_CONFIG).map(([value, cfg]) => ({ label: cfg.label, value }))"
        placeholder="语言筛选"
        clearable
        size="small"
        style="min-width: 120px"
      />
      <NButton size="small" @click="resetFilters">重置</NButton>
    </div>

    <!-- Table -->
    <NSpin :show="loading">
      <n-data-table
        :columns="columns"
        :data="submissions"
        :row-key="(row: any) => row.id"
        size="small"
        :loading="loading"
        :pagination="false"
        :row-class-name="() => 'cursor-pointer'"
        @click-row="handleViewSubmission"
      />
      <div class="flex justify-end py-4" v-if="total > 0">
        <NPagination
          v-model:page="currentPage"
          :page-size="PAGE_SIZE"
          :item-count="total"
          @update:page="fetchSubmissions"
        />
      </div>
    </NSpin>
  </div>
</template>
