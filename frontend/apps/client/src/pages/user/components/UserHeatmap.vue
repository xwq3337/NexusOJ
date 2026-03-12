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
import type { HeatmapFirstDayOfWeek } from 'naive-ui'
import { heatmapMockData } from 'naive-ui'

const props = defineProps<{
  userId?: string
  period?: 'recent' | '2025' | '2024'
}>()

const internalValue = ref<'recent' | number>(props.period || 'recent')
const loading = ref(false)
const firstDayOfWeek = ref<HeatmapFirstDayOfWeek>(0)

watch(() => props.period, (newValue) => {
  if (newValue) {
    internalValue.value = newValue
  }
})

const yearData = computed(() => {
  return heatmapMockData(internalValue.value)
})

const dataStats = computed(() => {
  const data = yearData.value
  const total = data.length
  const zeros = data.filter(d => d.value === 0).length
  const maxValue = Math.max(...data.map(d => d.value ?? 0))
  const avgValue
    = Math.round(
      (data.reduce((sum, d) => sum + (d.value ?? 0), 0) / total) * 100
    ) / 100

  return {
    total,
    zeros,
    maxValue,
    avgValue,
    zeroPercent: Math.round((zeros / total) * 100)
  }
})
</script>
