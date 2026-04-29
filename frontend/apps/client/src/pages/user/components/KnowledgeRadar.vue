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
    <div v-if="!tagProgress || Object.keys(tagProgress).length === 0" class="h-60 flex items-center justify-center">
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">暂无做题记录</p>
    </div>

    <!-- 力导向图 -->
    <div v-else class="h-80">
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
import { GraphChart } from 'echarts/charts'
import { TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { use } from 'echarts/core'
import { Target } from 'lucide-vue-next'

use([CanvasRenderer, GraphChart, TooltipComponent])

interface Props {
  tagProgress: Record<string, number>
  tagTotal?: Record<string, number>
  strongestTags?: string[]
  weakestTags?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  tagTotal: () => ({}),
  strongestTags: () => [],
  weakestTags: () => [],
})

const getProgressColor = (progress: number) => {
  if (progress >= 0.6) return '#4ade80'
  if (progress >= 0.3) return '#fbbf24'
  return '#f87171'
}

const chartOption = computed(() => {
  const tags = Object.keys(props.tagProgress)
  if (tags.length === 0) return {}

  const nodes = tags.map(tag => {
    const progress = props.tagProgress[tag] ?? 0
    const total = props.tagTotal[tag] ?? 1
    const size = Math.max(36, Math.min(70, Math.sqrt(total) * 8))
    const color = getProgressColor(progress)

    return {
      name: tag,
      symbolSize: size,
      symbol: 'roundRect' as const,
      itemStyle: {
        color: '#1e293b',
        borderColor: color,
        borderWidth: 2,
        borderRadius: 10,
      },
      label: {
        show: true,
        position: 'inside' as const,
        formatter: `{name|${tag}}\n{progress|${Math.round(progress * 100)}%}`,
        rich: {
          name: {
            fontSize: 10,
            color: '#e2e8f0',
            lineHeight: 16,
            align: 'center' as const,
          },
          progress: {
            fontSize: 11,
            fontWeight: 'bold' as const,
            color: color,
            lineHeight: 16,
            align: 'center' as const,
          },
        },
      },
      _progress: progress,
      _total: total,
    }
  })

  return {
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      borderColor: 'rgba(14, 165, 233, 0.3)',
      textStyle: { color: '#e2e8f0', fontSize: 12 },
      formatter: (params: any) => {
        const tag = params.name
        const progress = props.tagProgress[tag] ?? 0
        const total = props.tagTotal[tag] ?? 0
        const solved = Math.round(progress * total)
        return `<b>${tag}</b><br/>进度: ${Math.round(progress * 100)}%<br/>已解决: ${solved}/${total}`
      },
    },
    series: [
      {
        type: 'graph',
        layout: 'force',
        roam: false,
        draggable: true,
        force: {
          repulsion: 150,
          gravity: 0.15,
          edgeLength: [60, 120],
          layoutAnimation: true,
          friction: 0.6,
        },
        data: nodes,
        links: [],
        emphasis: {
          focus: 'adjacency',
          itemStyle: {
            shadowBlur: 12,
            shadowColor: 'rgba(14, 165, 233, 0.4)',
          },
        },
        animationDuration: 800,
        animationEasingUpdate: 'quinticInOut',
      },
    ],
  }
})
</script>
