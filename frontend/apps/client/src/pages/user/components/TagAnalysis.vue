<template>
  <div class="rounded-xl p-5 cyber-glow-card" :style="{
    backgroundColor: 'var(--surface-primary)',
    border: '1px solid var(--border-color)'
  }">
    <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
      <BarChart3 :size="16" :style="{ color: 'var(--accent-color)' }" />
      标签掌握度
    </h3>

    <div v-if="!tagScores || Object.keys(tagScores).length === 0" class="h-30 flex items-center justify-center">
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">暂无数据</p>
    </div>

    <div v-else class="space-y-3">
      <div v-for="tag in displayTags" :key="tag" class="space-y-1">
        <div class="flex items-center justify-between">
          <span class="text-xs truncate max-w-40" :style="{ color: 'var(--text-secondary)' }">{{ tag }}</span>
          <span class="text-xs font-terminal" :style="{
            color: getScoreColor(tagScores[tag])
          }">{{ (tagScores[tag] * 100).toFixed(0) }}%</span>
        </div>
        <div class="h-1.5 rounded-full overflow-hidden" :style="{
          backgroundColor: 'var(--surface-secondary)'
        }">
          <div class="h-full rounded-full transition-all duration-500" :style="{
            width: `${tagScores[tag] * 100}%`,
            backgroundColor: getScoreColor(tagScores[tag]),
          }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { BarChart3 } from 'lucide-vue-next'

interface Props {
  tagScores: Record<string, number>
}

const props = defineProps<Props>()

// 按掌握度从低到高排序，最多显示 8 个
const displayTags = computed(() => {
  const tags = Object.keys(props.tagScores)
  return tags
    .sort((a, b) => props.tagScores[a] - props.tagScores[b])
    .slice(0, 8)
})

const getScoreColor = (score: number) => {
  if (score >= 0.7) return '#4ade80'
  if (score >= 0.4) return '#fbbf24'
  return '#f87171'
}
</script>
