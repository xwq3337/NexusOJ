<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex items-center gap-3 mb-4 overflow-x-auto">
      <n-select v-model:value="filters.status" :options="statusOptions" placeholder="状态" clearable size="small"
        style="min-width: 120px; width: 120px;" />
      <n-select v-model:value="filters.language" :options="languageOptions" placeholder="语言" clearable size="small"
        style="min-width: 120px; width: 120px;" />
      <n-input v-model:value="filters.problem" placeholder="搜索题目..." size="small" clearable class="flex-1"
        style="min-width: 200px;">
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
    </div>

    <!-- Initial Loading -->
    <div v-if="loading && submissions.length === 0" class="flex justify-center items-center py-12">
      <n-spin size="large" />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && submissions.length === 0" class="text-center py-12">
      <Inbox :size="48" :style="{ color: 'var(--text-tertiary)' }" class="mx-auto mb-3" />
      <p :style="{ color: 'var(--text-tertiary)' }">暂无提交记录</p>
    </div>

    <!-- Submissions List -->
    <div v-else class="space-y-3">
      <div v-for="submission in submissions" :key="submission.id"
        class="p-4 rounded-xl transition-all duration-200 hover:scale-[1.01] cursor-pointer" :style="{
          backgroundColor: 'var(--surface-secondary)',
          border: '1px solid var(--border-color)'
        }" @click="handleViewSubmission(submission)">
        <div class="flex items-start justify-between gap-4">
          <!-- Problem Info -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-2">
              <span class="px-2 py-0.5 rounded text-xs font-medium" :style="{
                backgroundColor: STATUS_COLORS[submission.verdict]?.color || 'var(--surface-primary)',
                color: STATUS_COLORS[submission.verdict]?.textColor || 'var(--text-secondary)'
              }">
                {{ submission.verdict }}
              </span>
              <span class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                {{ submission.language }}
              </span>
            </div>
            <h4 class="font-medium mb-1 truncate" :style="{ color: 'var(--text-primary)' }">
              {{ submission.problem_title }}
            </h4>
            <div class="flex items-center gap-3 text-xs" :style="{ color: 'var(--text-tertiary)' }">
              <span class="flex items-center gap-1">
                <Clock :size="12" />
                {{ formatTime(submission.max_time) }}
              </span>
              <span class="flex items-center gap-1">
                <Cpu :size="12" />
                {{ formatMemory(submission.max_memory) }}
              </span>
              <span class="flex items-center gap-1">
                <Calendar :size="12" />
                {{ formatRelativeTime(submission.created_at) }}
              </span>
            </div>
          </div>

        </div>
      </div>

      <!-- Load More Trigger -->
      <div ref="loadMoreTrigger" class="flex justify-center items-center py-6">
        <div v-if="loadingMore" class="flex items-center gap-2" :style="{ color: 'var(--text-tertiary)' }">
          <n-spin size="small" />
          <span class="text-sm">加载更多...</span>
        </div>
        <div v-else-if="!hasMore" class="text-sm" :style="{ color: 'var(--text-tertiary)' }">
          已加载全部数据
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  NSelect,
  NInput,
  NSpin,
  type SelectOption
} from 'naive-ui'
import {
  Search,
  Inbox,
  Clock,
  Cpu,
  Calendar
} from 'lucide-vue-next'
import { useIntersectionObserver } from '@vueuse/core'
import { useDebounceFn } from '@vueuse/core'
import { userApi } from '@nexusoj/server'
import { formatMemory, formatRelativeTime, formatTime } from '@/utils/format'
import { STATUS_COLORS } from '@/constants'
import { GetRecordListResponse, VerdictType } from '@nexusoj/type'
// 根据props获取用户 id
const props = defineProps<{ userId: string }>()

const router = useRouter()

const submissions = ref<GetRecordListResponse[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const hasMore = ref(true)

const filters = ref({
  status: null as string | null,
  language: null as string | null,
  problem: ''
})
// TODO使用constant.ts里面的常量

const statusOptions: SelectOption[] = [
  { label: '通过', value: 'Accepted' },
  { label: '答案错误', value: 'WrongAnswer' },
  { label: '超时', value: 'TimeLimitExceeded' },
  { label: '内存超限', value: 'MemoryLimitExceeded' },
  { label: '运行错误', value: 'RuntimeError' },
  { label: '编译错误', value: 'CompilationError' }
]
// TODO使用constant.ts里面的常量
const languageOptions: SelectOption[] = [
  { label: 'C', value: 'c' },
  { label: 'C++', value: 'cpp' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'Go', value: 'go' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'Rust', value: 'rust' }
]

// Load more trigger element for intersection observer
const loadMoreTrigger = ref<HTMLElement | null>(null)

// Fetch submissions - append mode for infinite scroll
const fetchSubmissions = async (isLoadMore = false) => {
  if (isLoadMore) {
    if (loadingMore.value || !hasMore.value) return
    loadingMore.value = true
  } else {
    loading.value = true
    currentPage.value = 1
    submissions.value = []
  }

  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      verdict: filters.value.status as VerdictType,
      language: filters.value.language || undefined
    }
    const response = await userApi.getUserRecordList(props.userId, params)
    const newSubmissions = response.info || []

    if (isLoadMore) {
      submissions.value = [...submissions.value, ...newSubmissions]
    } else {
      submissions.value = newSubmissions
    }

    // Determine if there's more data based on returned items count
    hasMore.value = newSubmissions.length === pageSize.value
  } catch (error) {
    console.error('Failed to fetch submissions:', error)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

// Load more function for intersection observer
const loadMore = useDebounceFn(() => {
  if (!loadingMore.value && hasMore.value && !loading.value) {
    currentPage.value++
    fetchSubmissions(true)
  }
}, 200)

// Set up intersection observer for infinite scroll
const { stop } = useIntersectionObserver(
  loadMoreTrigger,
  ([{ isIntersecting }]) => {
    if (isIntersecting) {
      loadMore()
    }
  },
  {
    threshold: 0.1,
    rootMargin: '100px'
  }
)

const handleViewSubmission = (submission: any) => {
  router.push({
    name: 'RecordDetail',
    params: { id: submission.id }
  })
}

// Watch filters changes - reset pagination
watch(filters, () => {
  currentPage.value = 1
  hasMore.value = true
  fetchSubmissions(false)
}, { deep: true })

onMounted(() => {
  fetchSubmissions(false)
})

onUnmounted(() => {
  stop()
})
</script>
