<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NTag, NButton, NSpin, NEmpty, NPagination } from 'naive-ui'
import { RefreshCw, Sparkles, BookOpen } from 'lucide-vue-next'
import { userApi } from '@nexusoj/server'
import { difficultyMap } from '@/constants'
import type { RecommendedProblem } from '@nexusoj/type'

const router = useRouter()
const problems = ref<RecommendedProblem[]>([])
const total = ref(0)
const loading = ref(false)
const currentPage = ref(1)
const PAGE_SIZE = 8

const reasonMap: Record<string, { label: string; color: string }> = {
  difficulty_match: { label: '难度匹配', color: '#38bdf8' },
  tag_practice: { label: '标签巩固', color: '#a78bfa' },
  similar_users: { label: '相似用户', color: '#fbbf24' },
  contextual: { label: '上下文推荐', color: '#fb923c' },
  fresh: { label: '新题推荐', color: '#4ade80' },
}

const fetchRecommendations = async (refresh = false) => {
  loading.value = true
  try {
    const { code, info } = await userApi.getRecommendations(
      currentPage.value,
      PAGE_SIZE,
      refresh,
    )
    if (code === 200 && info) {
      problems.value = info.problems || []
      total.value = info.total || 0
    }
  } catch (e) {
    console.error('Failed to fetch recommendations:', e)
  } finally {
    loading.value = false
  }
}

const handleRefresh = () => {
  currentPage.value = 1
  fetchRecommendations(true)
}

const goToProblem = (problemId: number) => {
  router.push({ name: 'ProblemDetail', params: { id: problemId } })
}

onMounted(() => {
  fetchRecommendations()
})
</script>

<template>
  <div class="rounded-xl p-5 cyber-glow-card" :style="{
    backgroundColor: 'var(--surface-primary)',
    border: '1px solid var(--border-color)',
  }">
    <!-- Header -->
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-sm font-semibold flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
        <Sparkles :size="16" :style="{ color: 'var(--accent-color)' }" />
        为你推荐
      </h3>
      <!-- TODO: 实现换一批功能 -->
      <!-- <NButton size="tiny" quaternary :loading="loading" @click="handleRefresh">
        <template #icon>
          <RefreshCw :size="14" />
        </template>
        换一批
      </NButton> -->
    </div>

    <!-- Problem List -->
    <NSpin :show="loading">
      <div v-if="problems.length > 0" class="space-y-3">
        <div
          v-for="problem in problems"
          :key="problem.problem_id"
          class="rounded-lg p-3.5 cursor-pointer transition-all duration-200 hover:scale-[1.01]"
          :style="{
            backgroundColor: 'var(--surface-secondary)',
            border: '1px solid var(--border-color)',
          }"
          @click="goToProblem(problem.problem_id)"
        >
          <!-- Title & Difficulty -->
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2 min-w-0">
              <BookOpen :size="14" class="shrink-0" :style="{ color: 'var(--accent-color)' }" />
              <span
                class="text-sm font-medium truncate"
                :style="{ color: 'var(--text-primary)' }"
              >
                {{ problem.title }}
              </span>
            </div>
            <NTag
              size="tiny"
              :type="difficultyMap[Math.min(Math.round(problem.difficulty) - 1, 4)]?.type as any"
              :bordered="false"
            >
              {{ difficultyMap[Math.min(Math.round(problem.difficulty) - 1, 4)]?.text }}
            </NTag>
          </div>

          <!-- Tags -->
          <div class="flex flex-wrap gap-1.5 mb-2">
            <NTag
              v-for="tag in problem.tags?.slice(0, 3)"
              :key="tag"
              size="tiny"
              :bordered="false"
              :style="{
                backgroundColor: 'rgba(14, 165, 233, 0.08)',
                color: 'var(--accent-color)',
              }"
            >
              {{ tag }}
            </NTag>
          </div>

          <!-- Reason & Score -->
          <div class="flex items-center justify-between">
            <span
              class="text-xs flex items-center gap-1"
              :style="{ color: reasonMap[problem.reason]?.color || 'var(--text-tertiary)' }"
            >
              <Sparkles :size="10" />
              {{ reasonMap[problem.reason]?.label || problem.reason }}
            </span>
            <div class="flex items-center gap-2">
              <div class="w-16 h-1.5 rounded-full overflow-hidden" :style="{ backgroundColor: 'var(--surface-tertiary)' }">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :style="{
                    width: `${Math.round(problem.score * 100)}%`,
                    backgroundColor: reasonMap[problem.reason]?.color || 'var(--accent-color)',
                  }"
                />
              </div>
              <span class="text-xs font-terminal" :style="{ color: 'var(--text-tertiary)' }">
                {{ Math.round(problem.score * 100) }}%
              </span>
            </div>
          </div>
        </div>
      </div>

      <NEmpty v-else-if="!loading" description="暂无推荐题目" />
    </NSpin>

    <!-- Pagination -->
    <div v-if="total > PAGE_SIZE" class="flex justify-center pt-4">
      <NPagination
        v-model:page="currentPage"
        :page-size="PAGE_SIZE"
        :item-count="total"
        size="small"
        @update:page="() => fetchRecommendations()"
      />
    </div>
  </div>
</template>
