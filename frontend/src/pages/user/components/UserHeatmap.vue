<template>
  <div class="py-4">
    <div class="justify-center flex">
      <n-radio-group class="justify-center flex " v-model:value="value" name="year">
        <n-radio-button v-for="range in dateRanges" :key="range.value" :value="range.value" :label="range.label" />
      </n-radio-group>
    </div>
    <div class="justify-center flex">
      <n-heatmap class=" " color-theme="blue" :data="yearData" :loading-data="yearData"
        :first-day-of-week="firstDayOfWeek" :loading="loading" size="large" :fill-calendar-leading="value === 'recent'"
        :show-color-indicator="false" :tooltip="{ placement: 'bottom', delay: 500 }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { NHeatmap, NRadioButton, NRadioGroup } from 'naive-ui'
import type { HeatmapFirstDayOfWeek } from 'naive-ui'
import { heatmapMockData } from 'naive-ui'

const value = ref<'recent' | number>('recent')
const dateRanges = [
  {
    value: 'recent',
    label: '最近一年'
  },
  {
    value: 2025,
    label: 2025
  },
  {
    value: 2024,
    label: 2024
  },
  {
    value: 2023,
    label: 2023
  },
  {
    value: 2022,
    label: 2022
  }
].map((r) => {
  return {
    ...r,
    label: r.label.toString()
  }
}) as {
  value: 'recent' | number
  label: string
}[]

const yearData = computed(() => {
  return heatmapMockData(value.value)
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

const loading = ref(false)
const firstDayOfWeek = ref<HeatmapFirstDayOfWeek>(0)


</script>
