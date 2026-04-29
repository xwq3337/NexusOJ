<template>
  <div class="flex flex-col md:flex-row gap-6 mb-8">
    <!-- Left: Numbers -->
    <div class="flex flex-col gap-4 md:w-72 shrink-0">
      <!-- Problems + Contests -->
      <div class="p-5 rounded-xl transition-all duration-300 cyber-glow-card"
        :style="{ backgroundColor: 'var(--surface-primary)', border: '1px solid var(--border-color)' }">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2.5 rounded-lg" :style="{ backgroundColor: 'rgba(14, 165, 233, 0.1)', color: 'var(--accent-color)' }">
            <Code2 :size="20" />
          </div>
          <span class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Problems / Contests</span>
        </div>
        <div class="flex items-end gap-3">
          <h3 class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
            {{ problemCount ?? '—' }}
          </h3>
          <span class="text-lg mb-0.5" :style="{ color: 'var(--border-color)' }">/</span>
          <h3 class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
            {{ contestCount ?? '—' }}
          </h3>
        </div>
      </div>

      <!-- Blogs -->
      <div class="p-5 rounded-xl transition-all duration-300 cyber-glow-card"
        :style="{ backgroundColor: 'var(--surface-primary)', border: '1px solid var(--border-color)' }">
        <div class="flex items-center gap-3 mb-2">
          <div class="p-2.5 rounded-lg" :style="{ backgroundColor: 'rgba(20, 184, 166, 0.1)', color: 'var(--secondary)' }">
            <BookOpen :size="20" />
          </div>
          <span class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Blogs</span>
        </div>
        <h3 class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
          {{ blogCount ?? '—' }}
        </h3>
      </div>
    </div>

    <!-- Right: Avg. Judge Time -->
    <div class="flex-1 p-6 rounded-xl transition-all duration-300 cyber-glow-card"
      :style="{ backgroundColor: 'var(--surface-primary)', border: '1px solid var(--border-color)' }">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center gap-3">
          <div class="p-2.5 rounded-lg" :style="{ backgroundColor: 'rgba(74, 222, 128, 0.1)', color: 'var(--success-color)' }">
            <Zap :size="20" />
          </div>
          <span class="text-sm" :style="{ color: 'var(--text-tertiary)' }">Avg. Judge Time</span>
        </div>
        <div class="flex items-end gap-1">
          <span class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
            {{ judgeStats ? judgeStats.global.avg_max_time_ms.toFixed(1) : '—' }}
          </span>
          <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">ms</span>
        </div>
      </div>

      <div v-if="judgeStats" class="space-y-2">
        <div v-for="{ lang, name } in sortedLanguages" :key="name" class="flex items-center gap-3">
          <span class="text-xs w-20 shrink-0 font-terminal" :style="{ color: 'var(--text-secondary)' }">
            {{ name }}
          </span>
          <div class="flex-1 h-2 rounded-full overflow-hidden" :style="{ backgroundColor: 'var(--surface-tertiary)' }">
            <div class="h-full rounded-full transition-all duration-500"
              :style="{
                width: barWidth(lang.avg_max_time_ms) + '%',
                backgroundColor: barColor(lang.avg_max_time_ms),
              }"></div>
          </div>
          <span class="text-xs w-14 text-right shrink-0" :style="{ color: 'var(--text-tertiary)' }">
            {{ lang.avg_max_time_ms.toFixed(1) }}ms
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Code2, BookOpen, Zap } from 'lucide-vue-next'
import { problemApi, blogApi, contestApi, judgeApi, type JudgeStatsResponse } from '@nexusoj/server'

const problemCount = ref<number | null>(null)
const contestCount = ref<number | null>(null)
const blogCount = ref<number | null>(null)
const judgeStats = ref<JudgeStatsResponse | null>(null)

const sortedLanguages = computed(() => {
  if (!judgeStats.value) return []
  return Object.entries(judgeStats.value.per_language)
    .map(([name, lang]) => ({ name, lang }))
    .sort((a, b) => a.lang.avg_max_time_ms - b.lang.avg_max_time_ms)
})

const maxAvgTime = () => {
  if (!judgeStats.value) return 1
  const times = Object.values(judgeStats.value.per_language).map(l => l.avg_max_time_ms)
  return Math.max(...times, 1)
}

const barWidth = (ms: number) => Math.min((ms / maxAvgTime()) * 100, 100)

const barColor = (ms: number) => {
  const pct = barWidth(ms)
  if (pct <= 10) return 'rgb(14, 165, 233)'   // blue
  if (pct <= 30) return '#18a058'               // green
  if (pct <= 60) return 'rgb(240, 160, 32)'    // orange
  return '#ef4444'                               // red
}

onMounted(async () => {
  const results = await Promise.allSettled([
    problemApi.Count(),
    contestApi.getContestList(1, 1),
    blogApi.Count(),
    judgeApi.GetStats(),
  ])

  if (results[0].status === 'fulfilled') {
    problemCount.value = (results[0].value as any).info ?? results[0].value
  }
  if (results[1].status === 'fulfilled') {
    const info = (results[1].value as any).info
    contestCount.value = info?.total ?? null
  }
  if (results[2].status === 'fulfilled') {
    blogCount.value = (results[2].value as any).info ?? results[2].value
  }
  if (results[3].status === 'fulfilled') {
    judgeStats.value = (results[3].value as any).info ?? results[3].value
  }
})
</script>
