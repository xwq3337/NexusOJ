<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    <!-- 总进度 - 仪表盘 -->
    <div
      class="p-6 rounded-xl flex flex-col items-center justify-center"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div class="text-center mb-4">
        <h3 class="text-sm font-medium mb-1" :style="{ color: 'var(--text-secondary)' }">总进度</h3>
        <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">已完成题目</p>
      </div>
      <n-progress
        type="dashboard"
        :gap-offset-degree="180"
        :percentage="progressPercentage"
        :stroke-width="10"
        :color="progressColor"
        :rail-color="railColor"
        :show-indicator="false"
      />
      <div class="text-center mt-2">
        <span class="text-3xl font-bold" :style="{ color: 'var(--text-primary)' }">
          {{ solved }}
        </span>
        <p class="text-xs" :style="{ color: 'var(--text-secondary)' }">
          / {{ total }} 题
        </p>
      </div>
    </div>

    <!-- 总通过率 - 环形进度 -->
    <div
      class="p-6 rounded-xl flex flex-col items-center justify-center"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div class="text-center mb-4">
        <h3 class="text-sm font-medium mb-1" :style="{ color: 'var(--text-secondary)' }">总通过率</h3>
        <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">提交成功比例</p>
      </div>
      <n-progress
        type="circle"
        :percentage="acceptanceRate"
        :stroke-width="10"
        :color="acceptanceColor"
        :rail-color="railColor"
      >

      </n-progress>
    </div>

    <!-- 总提交数 -->
    <div
      class="p-6 rounded-xl flex flex-col items-center justify-center"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div class="text-center mb-4">
        <h3 class="text-sm font-medium mb-1" :style="{ color: 'var(--text-secondary)' }">总提交数</h3>
        <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">累计提交次数</p>
      </div>
      <div class="text-center">
        <span class="text-4xl font-bold text-blue-400">{{ submissionCount }}</span>
        <p class="text-xs mt-2" :style="{ color: 'var(--text-secondary)' }">
          平均 {{ avgPerDay }} 次/天
        </p>
      </div>
    </div>

    <!-- Rating -->
    <div
      class="p-6 rounded-xl flex flex-col items-center justify-center"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }"
    >
      <div class="text-center mb-4">
        <h3 class="text-sm font-medium mb-1" :style="{ color: 'var(--text-secondary)' }">Rating</h3>
        <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">当前排名</p>
      </div>
      <div class="text-center">
        <span
          class="text-4xl font-bold"
          :style="{ color: getRatingColor(rating) }"
        >
          {{ rating }}
        </span>
        <div
          class="inline-block px-2 py-0.5 rounded text-xs mt-2"
          :style="{
            backgroundColor: getRatingBgColor(rating),
            color: getRatingColor(rating)
          }"
        >
          {{ getRatingTitle(rating) }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { NProgress } from 'naive-ui'

const props = defineProps<{
  userId?: string
}>()

// 用户数据
const solved = ref(156)
const total = ref(320)
const submissionCount = ref(892)
const rating = ref(2150)

// 计算进度百分比
const progressPercentage = computed(() => {
  return total.value > 0 ? Math.round((solved.value / total.value) * 100) : 0
})

// 计算通过率
const acceptanceRate = computed(() => {
  return submissionCount.value > 0 ? Number(((solved.value / submissionCount.value) * 100).toFixed(1)) : 0
})

// 计算平均每天提交数
const avgPerDay = computed(() => {
  return (submissionCount.value / 365).toFixed(1)
})

// 进度颜色（根据百分比）
const progressColor = computed(() => {
  const percentage = progressPercentage.value
  if (percentage >= 80) return '#10b981'
  if (percentage >= 60) return '#3b82f6'
  if (percentage >= 40) return '#f59e0b'
  return '#ef4444'
})

// 通过率颜色
const acceptanceColor = computed(() => {
  const rate = acceptanceRate.value
  if (rate >= 80) return '#10b981'
  if (rate >= 60) return '#3b82f6'
  if (rate >= 40) return '#f59e0b'
  return '#ef4444'
})

// 轨道颜色
const railColor = 'var(--surface-secondary)'

// Rating 颜色
const getRatingColor = (rating: number) => {
  if (rating >= 2400) return '#ef4444' // Red
  if (rating >= 2100) return '#f97316' // Orange
  if (rating >= 1900) return '#eab308' // Yellow
  if (rating >= 1700) return '#8b5cf6' // Purple
  if (rating >= 1500) return '#3b82f6' // Blue
  if (rating >= 1200) return '#10b981' // Green
  return '#6b7280' // Gray
}

// Rating 背景颜色
const getRatingBgColor = (rating: number) => {
  if (rating >= 2400) return 'rgba(239, 68, 68, 0.1)'
  if (rating >= 2100) return 'rgba(249, 115, 22, 0.1)'
  if (rating >= 1900) return 'rgba(234, 179, 8, 0.1)'
  if (rating >= 1700) return 'rgba(139, 92, 246, 0.1)'
  if (rating >= 1500) return 'rgba(59, 130, 246, 0.1)'
  if (rating >= 1200) return 'rgba(16, 185, 129, 0.1)'
  return 'rgba(107, 114, 128, 0.1)'
}

// Rating 称号
const getRatingTitle = (rating: number) => {
  if (rating >= 2400) return 'Grandmaster'
  if (rating >= 2100) return 'Master'
  if (rating >= 1900) return 'Candidate Master'
  if (rating >= 1700) return 'Expert'
  if (rating >= 1500) return 'Specialist'
  if (rating >= 1200) return 'Pupil'
  return 'Newbie'
}

// 从 API 获取用户统计数据
const fetchUserStats = async () => {
  try {
    // TODO: 从 API 获取实际数据
    // const response = await Request.get(`/api/user/${props.userId}/stats`)
    // solved.value = response.solved
    // total.value = response.total
    // submissionCount.value = response.submissionCount
    // rating.value = response.rating
  } catch (error) {
    console.error('Failed to fetch user stats:', error)
  }
}

onMounted(() => {
  fetchUserStats()
})
</script>
