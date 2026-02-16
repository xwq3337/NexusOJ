<template>
  <div
    class="rounded-xl p-6"
    :style="{
      backgroundColor: 'var(--surface-primary)',
      borderColor: 'var(--border-color)',
      borderWidth: '1px',
      borderStyle: 'solid'
    }"
  >
    <div class="flex items-center justify-between mb-4">
      <div>
        <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">提交热力图</h3>
        <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">过去365天的提交活动</p>
      </div>
      <div class="flex items-center gap-2 text-xs" :style="{ color: 'var(--text-secondary)' }">
        <span>Less</span>
        <div class="flex items-center gap-1">
          <span
            v-for="level in 5"
            :key="level"
            class="w-3 h-3 rounded-sm transition-colors"
            :style="{ backgroundColor: getHeatmapColor(level - 1) }"
          ></span>
        </div>
        <span>More</span>
      </div>
    </div>

    <div class="overflow-x-auto">
      <div class="flex gap-1 min-w-max">
        <!-- Day labels -->
        <div
          class="flex flex-col gap-1 mr-2 text-xs pr-2"
          :style="{ color: 'var(--text-tertiary)', borderRight: '1px solid var(--border-color)' }"
        >
          <span class="h-3"></span>
          <span class="h-3"></span>
          <span>Mon</span>
          <span class="h-3"></span>
          <span>Wed</span>
          <span class="h-3"></span>
          <span>Fri</span>
          <span class="h-3"></span>
          <span class="h-3"></span>
        </div>

        <!-- Heatmap grid -->
        <div class="flex gap-1">
          <div v-for="(week, weekIndex) in heatmapGrid" :key="weekIndex" class="flex flex-col gap-1">
            <div
              v-for="(day, dayIndex) in week"
              :key="`${weekIndex}-${dayIndex}`"
              class="w-3 h-3 rounded-sm transition-all hover:scale-125 cursor-pointer relative"
              :style="{
                backgroundColor: day.count > 0 ? getHeatmapColor(day.level) : 'transparent',
                border: day.count === 0 ? '1px solid var(--border-light)' : 'none',
                opacity: day.count === 0 ? '0.3' : '1'
              }"
              :title="getTooltip(day)"
            ></div>
          </div>
        </div>
      </div>

      <!-- Month labels -->
      <div class="flex mt-2 pl-8 text-xs" :style="{ color: 'var(--text-tertiary)' }">
        <span
          v-for="(month, index) in monthLabels"
          :key="index"
          class="flex-shrink-0"
          :style="{ width: `${month.width}px` }"
        >
          {{ month.label }}
        </span>
      </div>
    </div>

    <div class="mt-4 flex items-center justify-center gap-6 text-sm" :style="{ color: 'var(--text-secondary)' }">
      <span>
        总提交: <strong class="text-blue-400">{{ totalSubmissions }}</strong>
      </span>
      <span>
        最长连续: <strong class="text-green-400">{{ longestStreak }}</strong> 天
      </span>
      <span>
        当前连续: <strong class="text-yellow-400">{{ currentStreak }}</strong> 天
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface HeatmapDay {
  date: string
  count: number
  level: number
}

const props = defineProps<{
  userId?: string
}>()

const heatmapData = ref<HeatmapDay[]>([])
const heatmapGrid = ref<HeatmapDay[][]>([])
const totalSubmissions = ref(0)
const longestStreak = ref(0)
const currentStreak = ref(0)
const monthLabels = ref<Array<{ label: string; width: number }>>([])

// 热力图颜色（绿色渐变，类似 GitHub）
const getHeatmapColor = (level: number) => {
  const colors = [
    'rgba(35, 134, 54, 0.2)', // Level 0 - 浅绿
    'rgba(35, 134, 54, 0.4)', // Level 1
    'rgba(35, 134, 54, 0.6)', // Level 2
    'rgba(35, 134, 54, 0.8)', // Level 3
    'rgba(35, 134, 54, 1)' // Level 4 - 深绿
  ]
  return colors[level] || colors[0]
}

const getTooltip = (day: HeatmapDay) => {
  if (day.count === 0) return `${day.date}: 无提交`
  return `${day.date}: ${day.count} 次提交`
}

// 生成模拟热力图数据（实际应从 API 获取）
const generateHeatmapData = () => {
  const data: HeatmapDay[] = []
  const today = new Date()

  for (let i = 364; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    const dateStr = date.toISOString().split('T')[0]

    // 模拟数据：30% 概率没有提交，其余随机 1-15 次提交
    const hasSubmission = Math.random() > 0.3
    const count = hasSubmission ? Math.floor(Math.random() * 15) + 1 : 0

    data.push({
      date: dateStr,
      count,
      level: count > 0 ? Math.min(Math.ceil(count / 4), 4) : 0
    })
  }

  return data
}

// 计算统计数据
const calculateStats = (data: HeatmapDay[]) => {
  // 总提交数
  totalSubmissions.value = data.reduce((sum, day) => sum + day.count, 0)

  // 计算最长连续天数
  let maxStreak = 0
  let currentStreakCount = 0

  for (const day of data) {
    if (day.count > 0) {
      currentStreakCount++
      maxStreak = Math.max(maxStreak, currentStreakCount)
    } else {
      currentStreakCount = 0
    }
  }
  longestStreak.value = maxStreak

  // 计算当前连续天数（从今天往前数）
  currentStreak.value = 0
  for (let i = data.length - 1; i >= 0; i--) {
    if (data[i].count > 0) {
      currentStreak.value++
    } else {
      break
    }
  }
}

// 将数据转换为每周网格（7天 x 53周）
const transformToGrid = (data: HeatmapDay[]) => {
  const weeks: HeatmapDay[][] = []
  const weekSize = 7

  for (let i = 0; i < data.length; i += weekSize) {
    const week = data.slice(i, i + weekSize)
    weeks.push(week)
  }

  return weeks
}

// 生成月份标签
const generateMonthLabels = (data: HeatmapDay[]) => {
  const months: Array<{ label: string; width: number }> = []
  const weekWidth = 16 // 每周的宽度（包含gap）

  let currentMonth = ''
  let weekCount = 0
  const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']

  for (let i = 0; i < data.length; i += 7) {
    const date = new Date(data[i].date)
    const month = monthNames[date.getMonth()]

    if (month !== currentMonth) {
      if (currentMonth !== '') {
        months.push({ label: currentMonth, width: weekCount * weekWidth })
      }
      currentMonth = month
      weekCount = 1
    } else {
      weekCount++
    }
  }

  // 添加最后一个月
  if (currentMonth !== '') {
    months.push({ label: currentMonth, width: weekCount * weekWidth })
  }

  monthLabels.value = months
}

// 从 API 获取用户热力图数据
const fetchHeatmapData = async () => {
  try {
    // TODO: 从 API 获取实际数据
    // const response = await Request.get(`/api/user/${props.userId}/heatmap`)
    // heatmapData.value = response.data

    // 使用模拟数据
    const mockData = generateHeatmapData()
    heatmapData.value = mockData
    calculateStats(mockData)
    heatmapGrid.value = transformToGrid(mockData)
    generateMonthLabels(mockData)
  } catch (error) {
    console.error('Failed to fetch heatmap data:', error)
  }
}

onMounted(() => {
  fetchHeatmapData()
})
</script>
