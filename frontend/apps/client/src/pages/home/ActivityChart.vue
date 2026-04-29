<template>
  <div
    class="rounded-xl p-6 h-full"
    :style="{
      backgroundColor: 'var(--surface-primary)',
      border: '1px solid var(--border-color)'
    }"
  >
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">总活跃量</h3>
        <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">最近 7 天的提交数量</p>
      </div>
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 rounded-full" :style="{ backgroundColor: 'var(--accent-color)' }"></span>
        <span class="text-xs" :style="{ color: 'var(--text-tertiary)' }">提交数量</span>
      </div>
    </div>

    <div class="h-75 w-full">
      <VChart :option="chartOption" autoresize class="w-full h-full" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import VChart from 'vue-echarts'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { use } from 'echarts/core'
import { recordApi } from '@nexusoj/server'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const chartOption = ref(buildOption([], []))

function buildOption(dates: string[], counts: number[]) {
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#0f172a',
      borderColor: '#1e3a5f',
      textStyle: { color: '#e2e8f0' }
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#94a3b8' }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#94a3b8' },
      splitLine: { lineStyle: { color: '#1e3a5f' } }
    },
    grid: { left: '3%', right: '3%', bottom: '10%', top: '10%' },
    series: [
      {
        data: counts,
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: {
          color: '#0ea5e9',
          width: 3,
          shadowColor: 'rgba(14, 165, 233, 0.4)',
          shadowBlur: 10
        },
        areaStyle: {
          color: new (echarts as any).graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(14, 165, 233, 0.25)' },
            { offset: 1, color: 'rgba(14, 165, 233, 0)' }
          ])
        }
      }
    ]
  }
}

onMounted(async () => {
  try {
    const res: any = await recordApi.getDailyActivity()
    const data: { date: string; count: number }[] = res.info ?? res
    if (data?.length) {
      const dates = data.map(d => d.date.slice(5)) // "04/22"
      const counts = data.map(d => d.count)
      chartOption.value = buildOption(dates, counts)
    }
  } catch {
    // keep empty chart
  }
})
</script>
