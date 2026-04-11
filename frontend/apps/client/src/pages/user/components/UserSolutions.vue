<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex items-center gap-3 mb-4 overflow-x-auto">
      <n-select
        v-model:value="filters.difficulty"
        :options="difficultyOptions"
        placeholder="难度"
        clearable
        size="small"
        style="min-width: 120px; width: 120px;"
      />
      <n-select
        v-model:value="filters.tag"
        :options="tagOptions"
        placeholder="标签"
        clearable
        size="small"
        style="min-width: 120px; width: 120px;"
      />
      <n-input
        v-model:value="filters.search"
        placeholder="搜索题解..."
        size="small"
        clearable
        class="flex-1"
        style="min-width: 200px;"
      >
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
    </div>

    <!-- Solutions List -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <n-spin size="large" />
    </div>

    <div v-else-if="solutions.length === 0" class="text-center py-12">
      <FileCode :size="48" :style="{ color: 'var(--text-tertiary)' }" class="mx-auto mb-3" />
      <p :style="{ color: 'var(--text-tertiary)' }">暂无题解</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="solution in solutions"
        :key="solution.id"
        class="p-5 rounded-xl transition-all duration-200 hover:scale-[1.02] hover:shadow-lg cursor-pointer"
        :style="{
          backgroundColor: 'var(--surface-secondary)',
          border: '1px solid var(--border-color)'
        }"
        @click="handleViewSolution(solution)"
      >
        <!-- Problem Badge -->
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2 text-xs" :style="{ color: 'var(--text-tertiary)' }">
            <Eye :size="14" />
            <span>{{ solution.view }} 阅读</span>
          </div>
        </div>

        <!-- Title -->
        <h4 class="font-semibold mb-2 line-clamp-2" :style="{ color: 'var(--text-primary)' }">
          {{ solution.title }}
        </h4>

        <!-- Problem Name -->
        <p class="text-sm mb-3 truncate" :style="{ color: 'var(--text-secondary)' }">
          {{ solution.problem_title }}
        </p>

        <!-- Tags -->
        <div class="flex flex-wrap gap-2 mb-3">
          <span
            v-for="tag in solution.tags.slice(0, 3)"
            :key="tag"
            class="px-2 py-0.5 rounded text-xs"
            :style="{
              backgroundColor: 'var(--surface-primary)',
              color: 'var(--text-tertiary)'
            }"
          >
            {{ tag }}
          </span>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between text-xs" :style="{ color: 'var(--text-tertiary)' }">
          <span class="flex items-center gap-1">
            <ThumbsUp :size="12" />
            {{ solution.like }}
          </span>
          <span>{{ formatRelativeTime(solution.created_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-6">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-size="pageSize"
        show-quick-jumper
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NSelect,
  NInput,
  NSpin,
  NPagination,
  type SelectOption
} from 'naive-ui'
import {
  Search,
  FileCode,
  Eye,
  ThumbsUp
} from 'lucide-vue-next'
import { formatRelativeTime } from '@/utils/format'
import { difficultyMap } from '@/constants'
import { solutionApi } from '@nexusoj/server'

const props = defineProps<{
  userId?: Number
}>()

const router = useRouter()

const solutions = ref<any[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(12)
const totalPages = ref(0)

const filters = ref({
  difficulty: null as number | null,
  tag: null as string | null,
  search: ''
})

const difficultyOptions: SelectOption[] = [
  { label: '简单', value: 0 },
  { label: '中等', value: 1 },
  { label: '困难', value: 2 },
  { label: '极难', value: 3 }
]

const tagOptions: SelectOption[] = [
  { label: '数组', value: '数组' },
  { label: '链表', value: '链表' },
  { label: '树', value: '树' },
  { label: '动态规划', value: '动态规划' },
  { label: '贪心', value: '贪心' },
  { label: '二分查找', value: '二分查找' },
  { label: '图论', value: '图论' },
  { label: '数学', value: '数学' }
]

const getDifficultyColor = (difficulty: number) => {
  return difficultyMap[difficulty]?.color || '#6b7280'
}

const getDifficultyBgColor = (difficulty: number) => {
  const color = getDifficultyColor(difficulty)
  return color.replace(')', ', 0.1)').replace('rgb', 'rgba')
}

const getDifficultyLabel = (difficulty: number) => {
  return difficultyMap[difficulty]?.text || '未知'
}

const fetchSolutions = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      user_id: props.userId ? Number(props.userId) : undefined,
      tag: filters.value.tag || undefined,
      keyword: filters.value.search || undefined
    }
    const response = await solutionApi.getSolutions(params)
    solutions.value = response.info?.solutions || []
    totalPages.value = Math.ceil((response.info?.total || 0) / pageSize.value)
  } catch (error) {
    console.error('Failed to fetch solutions:', error)
  } finally {
    loading.value = false
  }
}

const handleViewSolution = (solution: any) => {
  router.push({
    name: 'SolutionDetail',
    params: { id: solution.id }
  })
}

watch([currentPage, filters], () => {
  fetchSolutions()
}, { deep: true })

onMounted(() => {
  fetchSolutions()
})
</script>
