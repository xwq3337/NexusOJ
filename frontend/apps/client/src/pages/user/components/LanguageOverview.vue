<template>
  <div class="rounded-xl p-5 cyber-glow-card" :style="{
    backgroundColor: 'var(--surface-primary)',
    border: '1px solid var(--border-color)'
  }">
    <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
      <Code2 :size="16" :style="{ color: 'var(--accent-color)' }" />
      语言概览
    </h3>

    <div v-if="!languages || Object.keys(languages).length === 0" class="h-30 flex items-center justify-center">
      <p class="text-sm" :style="{ color: 'var(--text-tertiary)' }">暂无数据</p>
    </div>

    <div v-else class="space-y-3">
      <div v-for="lang in displayLanguages" :key="lang.key" class="space-y-1">
        <div class="flex items-center justify-between">
          <span class="text-xs truncate max-w-40" :style="{ color: lang.color }">{{ lang.label }}</span>
          <span class="text-xs font-terminal" :style="{ color: lang.color }">{{ lang.count }} 次</span>
        </div>
        <div class="h-1.5 rounded-full overflow-hidden" :style="{
          backgroundColor: 'var(--surface-secondary)'
        }">
          <div class="h-full rounded-full transition-all duration-500" :style="{
            width: `${lang.percent}%`,
            backgroundColor: lang.color,
          }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Code2 } from 'lucide-vue-next'
import { LANGUAGE_CONFIG } from '@/constants'

interface Props {
  languages: Record<string, number>
}

const props = defineProps<Props>()

const displayLanguages = computed(() => {
  const entries = Object.entries(props.languages)
  if (entries.length === 0) return []

  const maxCount = Math.max(...entries.map(([, count]) => count))

  return entries
    .map(([key, count]) => {
      const config = LANGUAGE_CONFIG[key as keyof typeof LANGUAGE_CONFIG]
      return {
        key,
        label: config?.label ?? key,
        color: config?.color.textColor ?? 'var(--text-secondary)',
        count,
        percent: maxCount > 0 ? (count / maxCount) * 100 : 0,
      }
    })
    .sort((a, b) => b.count - a.count)
    .slice(0, 8)
})
</script>
