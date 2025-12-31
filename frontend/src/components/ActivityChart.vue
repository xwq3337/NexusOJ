<template>
  <div
    class="rounded-xl p-6 h-full"
    :style="{
      backgroundColor: 'var(--surface-primary)',
      borderColor: 'var(--border-color)',
      borderWidth: '1px',
      borderStyle: 'solid'
    }"
  >
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-lg font-bold text-white">Global Activity</h3>
        <p class="text-sm text-gray-400">本周所有的提交数量</p>
      </div>
      <div class="flex items-center gap-2">
        <span class="w-2 h-2 rounded-full bg-blue-500"></span>
        <span class="text-xs text-gray-400">提交数量</span>
      </div>
    </div>

    <div class="h-[300px] w-full">
      <VChart :option="option" autoresize class="w-full h-full" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import VChart from 'vue-echarts'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { use } from 'echarts/core'
import { ACTIVITY_DATA } from '../../constants'

use([CanvasRenderer, LineChart, GridComponent, TooltipComponent])

const data = ACTIVITY_DATA

const option = ref({
  tooltip: {
    trigger: 'axis',
    backgroundColor: '#0f172a',
    borderColor: '#334155',
    textStyle: { color: '#e2e8f0' }
  },
  xAxis: {
    type: 'category',
    data: data.map((d) => d.name),
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: '#94a3b8' }
  },
  yAxis: {
    type: 'value',
    axisLine: { show: false },
    axisTick: { show: false },
    axisLabel: { color: '#94a3b8' },
    splitLine: { lineStyle: { color: '#334155' } }
  },
  grid: { left: '3%', right: '3%', bottom: '10%', top: '10%' },
  series: [
    {
      data: data.map((d) => d.submissions),
      type: 'line',
      smooth: true,
      symbol: 'none',
      lineStyle: { color: '#3b82f6', width: 3 },
      areaStyle: {
        color: new (echarts as any).graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(59,130,246,0.3)' },
          { offset: 1, color: 'rgba(59,130,246,0)' }
        ])
      }
    }
  ]
})
</script>
