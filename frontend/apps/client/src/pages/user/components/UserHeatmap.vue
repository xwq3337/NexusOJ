<template>
  <div class="py-4">
    <div class="justify-center flex">
      <n-heatmap color-theme="blue" :data="yearData" :loading-data="yearData"
        :first-day-of-week="firstDayOfWeek" :loading="loading" size="small" :fill-calendar-leading="internalValue === 'recent'"
        :show-color-indicator="false" :tooltip="{ placement: 'bottom', delay: 500 }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import {
  eachDayOfInterval,
  startOfDay,
  subYears,
  startOfYear,
  endOfYear,
} from 'date-fns'
import { NHeatmap } from 'naive-ui'
import type { HeatmapFirstDayOfWeek, HeatmapData } from 'naive-ui'

interface HeatmapRawData {
  heatmaps?: Record<string, Record<string, number>>
  past_year_heatmap?: Record<string, number>
}

const props = defineProps<{
  userId?: Number
  period?: 'recent' | '2026' | '2025'
  heatmapData?: HeatmapRawData
}>()

const internalValue = ref<'recent' | '2026' | '2025'>(props.period || 'recent')
const loading = ref(false)
const firstDayOfWeek = ref<HeatmapFirstDayOfWeek>(0)

watch(() => props.period, (newValue) => {
  if (newValue) {
    internalValue.value = newValue
  }
})

function generateFullRange(start: Date, end: Date, rawData: Record<string, number>, year?: string): HeatmapData {
  // 生成完整日期范围内每一天的初始数据 (value: 0)
  const allDays = eachDayOfInterval({ start, end })
  const resultMap = new Map<number, number>()

  for (const day of allDays) {
    resultMap.set(startOfDay(day).getTime(), 0)
  }

  // 用真实数据覆盖
  for (const [dateStr, value] of Object.entries(rawData)) {
    let fullDateStr = dateStr
    if (year && /^\d{2}-\d{2}T/.test(dateStr)) {
      fullDateStr = `${year}-${dateStr}`
    }
    const ts = startOfDay(new Date(fullDateStr)).getTime()
    if (!isNaN(ts)) {
      resultMap.set(ts, value)
    }
  }

  return Array.from(resultMap.entries()).map(([timestamp, value]) => ({
    timestamp,
    value,
  }))
}

const yearData = computed(() => {
  if (!props.heatmapData) return []

  if (internalValue.value === 'recent') {
    if (props.heatmapData.past_year_heatmap) {
      const end = new Date()
      const start = subYears(end, 1)
      return generateFullRange(start, end, props.heatmapData.past_year_heatmap)
    }
    return []
  }

  const yearKey = String(internalValue.value)
  if (props.heatmapData.heatmaps?.[yearKey]) {
    const yearNum = Number(yearKey)
    const start = startOfYear(new Date(yearNum, 0, 1))
    const end = endOfYear(new Date(yearNum, 11, 31))
    return generateFullRange(start, end, props.heatmapData.heatmaps[yearKey], yearKey)
  }
  return []
})
</script>
