<template>
  <div class="py-4">
    <div class="justify-center flex">
      <n-heatmap class=" " color-theme="blue" :data="yearData" :loading-data="yearData"
        :first-day-of-week="firstDayOfWeek" :loading="loading" size="large" :fill-calendar-leading="internalValue === 'recent'"
        :show-color-indicator="false" :tooltip="{ placement: 'bottom', delay: 500 }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
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

function convertToHeatmapData(rawData: Record<string, number>, year?: string): HeatmapData {
  return Object.entries(rawData).map(([dateStr, value]) => {
    let fullDateStr = dateStr
    // heatmaps 内层的 key 格式为 "MM-DDT00:00:00+08:00"，需要拼接年份
    if (year && /^\d{2}-\d{2}T/.test(dateStr)) {
      fullDateStr = `${year}-${dateStr}`
    }
    return {
      timestamp: new Date(fullDateStr).getTime(),
      value
    }
  })
}

const yearData = computed(() => {
  if (!props.heatmapData) return []

  if (internalValue.value === 'recent') {
    if (props.heatmapData.past_year_heatmap) {
      return convertToHeatmapData(props.heatmapData.past_year_heatmap)
    }
    return []
  }

  const yearKey = String(internalValue.value)
  if (props.heatmapData.heatmaps?.[yearKey]) {
    return convertToHeatmapData(props.heatmapData.heatmaps[yearKey], yearKey)
  }
  return []
})

setInterval(() => {
  console.log(yearData.value)
}, 2000)
</script>
