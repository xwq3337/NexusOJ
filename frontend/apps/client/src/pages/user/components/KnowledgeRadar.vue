<template>
  <div class="rounded-xl p-5 cyber-glow-card" :style="{
    backgroundColor: 'var(--surface-primary)',
    border: '1px solid var(--border-color)'
  }">
    <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
      <Target :size="16" :style="{ color: 'var(--accent-color)' }" />
      知识掌握度
    </h3>

    <!-- 无数据状态 -->
    <div v-if="!tagScores || Object.keys(tagScores).length === 0" class="h-60 flex items-center justify-center">
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">暂无做题记录</p>
    </div>  
    <!-- TODO: 统计各个标签的总数 Redis -->
    <!-- 雷达图 -->
    <div v-else class="h-72">
      <VChart :option="chartOption" autoresize class="w-full h-full" />
    </div>
    <!-- 标签概览 -->
    <div v-if="strongestTags.length > 0 || weakestTags.length > 0" class="mt-4 space-y-2">
      <div v-if="strongestTags.length > 0" class="flex items-start gap-2">
        <span class="text-xs font-medium shrink-0 mt-0.5" :style="{ color: '#4ade80' }">强项</span>
        <div class="flex flex-wrap gap-1">
          <span v-for="tag in strongestTags" :key="tag" class="px-2 py-0.5 rounded text-xs" :style="{
            backgroundColor: 'rgba(74, 222, 128, 0.1)',
            color: '#4ade80',
            border: '1px solid rgba(74, 222, 128, 0.2)'
          }">{{ tag }}</span>
        </div>
      </div>
      <div v-if="weakestTags.length > 0" class="flex items-start gap-2">
        <span class="text-xs font-medium shrink-0 mt-0.5" :style="{ color: '#f87171' }">待加强</span>
        <div class="flex flex-wrap gap-1">
          <span v-for="tag in weakestTags" :key="tag" class="px-2 py-0.5 rounded text-xs" :style="{
            backgroundColor: 'rgba(248, 113, 113, 0.1)',
            color: '#f87171',
            border: '1px solid rgba(248, 113, 113, 0.2)'
          }">{{ tag }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import VChart from 'vue-echarts'
import * as echarts from 'echarts/core'
import { RadarChart } from 'echarts/charts'
import { TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { use } from 'echarts/core'
import { Target } from 'lucide-vue-next'

use([CanvasRenderer, RadarChart, TooltipComponent, LegendComponent])

interface Props {
  tagScores: Record<string, number>
  strongestTags?: string[]
  weakestTags?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  strongestTags: () => [],
  weakestTags: () => [],
})

const chartOption = computed(() => {
  const tags = Object.keys(props.tagScores)
  if (tags.length === 0) return {}

  // 雷达图最多显示 8 个维度，按掌握度排序取最有代表性的
  const sorted = tags.sort((a, b) => props.tagScores[b] - props.tagScores[a])
  const displayTags = sorted.length > 8
    ? [...sorted.slice(0, 4), ...sorted.slice(-4)]
    : sorted

  const indicators = displayTags.map(tag => ({
    name: tag,
    max: 1.0,
  }))

  const values = displayTags.map(tag => props.tagScores[tag])

  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      borderColor: 'rgba(14, 165, 233, 0.3)',
      textStyle: { color: '#e2e8f0', fontSize: 12 },
      formatter: (params: any) => {
        if (!params.value) return ''
        const points = params.value
        return displayTags
          .map((tag, i) => `${tag}: ${(points[i] * 100).toFixed(0)}%`)
          .join('<br/>')
      },
    },
    radar: {
      indicator: indicators,
      shape: 'polygon',
      radius: '65%',
      axisName: {
        color: 'var(--text-secondary)',
        fontSize: 11,
      },
      splitArea: {
        areaStyle: {
          color: ['rgba(14, 165, 233, 0.02)', 'rgba(14, 165, 233, 0.05)'],
        },
      },
      splitLine: {
        lineStyle: { color: 'rgba(14, 165, 233, 0.15)' },
      },
      axisLine: {
        lineStyle: { color: 'rgba(14, 165, 233, 0.15)' },
      },
    },
    series: [
      {
        type: 'radar',
        data: [
          {
            value: values,
            name: '掌握度',
            areaStyle: {
              color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                { offset: 0, color: 'rgba(14, 165, 233, 0.35)' },
                { offset: 1, color: 'rgba(14, 165, 233, 0.05)' },
              ]),
            },
            lineStyle: {
              color: '#0ea5e9',
              width: 2,
              shadowColor: 'rgba(14, 165, 233, 0.4)',
              shadowBlur: 6,
            },
            itemStyle: {
              color: '#0ea5e9',
              borderColor: '#0ea5e9',
              borderWidth: 2,
            },
            symbol: 'circle',
            symbolSize: 6,
          },
        ],
      },
    ],
  }
})
</script>
