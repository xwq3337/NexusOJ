<template>
  <div class="animate-fade-in">
    <!-- Hero Section with Cyberpunk Scanline -->
    <div class="relative overflow-hidden rounded-2xl p-8 md:p-12 mb-12 cyber-scanline-overlay" :style="{
      background: `linear-gradient(135deg, var(--hero-bg-from), var(--hero-bg-to))`,
      border: '1px solid var(--hero-border)',
      boxShadow: '0 0 40px rgba(14, 165, 233, 0.08)'
    }">
      <!-- Grid pattern background -->
      <div class="absolute inset-0 cyber-grid-bg opacity-40"></div>

      <div class="relative z-10 max-w-2xl">
        <div class="font-terminal text-xs mb-4" :style="{ color: 'var(--accent-color)' }">
          $ nexusoj --init
        </div>
        <h1 class="text-4xl md:text-5xl font-extrabold mb-6 leading-tight"
          :style="{ color: 'var(--hero-title-color)' }">
          Master Algorithms <br />
          <span class="text-transparent bg-clip-text bg-linear-to-r from-sky-400 to-sky-500 neon-text">Build the Future</span>
        </h1>
        <p class="text-lg mb-8 leading-relaxed" :style="{ color: 'var(--hero-desc-color)' }">
          加入新一代的编程竞技。在这里，你可以解决各种算法问题，参加大型的编程竞赛，并获取实时的人工智能指导，帮你理清思路、调试代码。
        </p>
        <div class="flex flex-wrap gap-4">
          <RouterLink to="/problems"
            class="hero-btn-primary text-white px-8 py-3 rounded-xl font-semibold transition-all hover:scale-105">
            开始解题
          </RouterLink>
          <RouterLink to="/contests"
            class="hero-btn-secondary px-8 py-3 rounded-xl font-semibold transition-all"
            :style="{ color: 'var(--accent-color)', border: '1px solid var(--accent-color)' }">
            浏览比赛
          </RouterLink>
        </div>
      </div>

      <!-- Glow orbs -->
      <div class="absolute top-0 right-0 -mt-20 -mr-20 w-96 h-96 rounded-full blur-3xl"
        :style="{ backgroundColor: 'var(--hero-glow-blue)' }"></div>
      <div class="absolute bottom-0 left-20 w-72 h-72 rounded-full blur-3xl"
        :style="{ backgroundColor: 'var(--hero-glow-purple)' }"></div>
    </div>

    <Stats />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12">
      <div class="lg:col-span-2">
        <ActivityChart />
      </div>

      <div class="rounded-xl p-6 cyber-glow-card" :style="{
        backgroundColor: 'var(--bg-secondary)',
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <span :style="{ color: 'var(--accent-color)' }">&#9889;</span>
            <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">每日挑战</h3>
          </div>
          <RouterLink to="/problems" class="text-xs font-terminal" :style="{ color: 'var(--accent-color)' }">
            view_all &rarr;
          </RouterLink>
        </div>

        <div class="space-y-4">
          <div v-for="problem in MOCK_PROBLEMS.slice(0, 3)" :key="problem.id"
            class="group flex items-center justify-between p-3 rounded-lg transition-all cursor-pointer"
            :style="{ border: '1px solid transparent' }"
            @click="$router.push({ name: 'ProblemDetail', params: { id: problem.id } })">
            <div class="flex items-center gap-3">
              <div>
                <h4 class="text-sm font-medium transition-colors"
                  :style="{ color: 'var(--text-primary)' }">
                  {{ problem.title }}
                </h4>
                <div class="flex items-center gap-2 mt-1">
                  <span v-for="tag in problem.tags.slice(0, 2)" :key="tag"
                    class="text-[10px] px-1.5 py-0.5 rounded" :style="{
                      backgroundColor: 'var(--surface-tertiary)',
                      color: 'var(--text-tertiary)',
                      border: '1px solid var(--border-light)'
                    }">{{ tag }}</span>
                </div>
              </div>
            </div>
            <div class="flex flex-col items-end gap-1">
              <span class="text-xs font-terminal" :style="{ color: 'var(--accent-color)' }">{{ problem.difficulty }}</span>
              <span class="text-xs" :style="{ color: 'var(--text-tertiary)' }">{{ formatAcceptance(problem.accept, problem.submission) }}</span>
            </div>
          </div>
        </div>

        <div class="mt-6 pt-6" :style="{
          borderTop: '1px solid var(--border-color)'
        }">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span :style="{ color: 'var(--neon-cyan)' }">&#9672;</span>
              <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">活跃比赛</h3>
            </div>
            <RouterLink to="/contests" class="text-xs font-terminal" :style="{ color: 'var(--accent-color)' }">
              view_all &rarr;
            </RouterLink>
          </div>
          <div class="space-y-3">
            <div v-for="contest in MOCK_CONTESTS.slice(0, 2)" :key="contest.id"
              class="p-3 rounded-lg flex items-center justify-between cyber-glow-card" :style="{
                backgroundColor: 'var(--surface-secondary)',
                borderColor: 'var(--border-color)',
                borderWidth: '1px',
                borderStyle: 'solid'
              }">
              <div>
                <h4 class="text-sm font-medium" :style="{ color: 'var(--text-primary)' }">
                  {{ contest.title }}
                </h4>
                <p class="text-xs mt-0.5" :style="{ color: 'var(--text-tertiary)' }">
                  <span v-if="contest.status === 'Live'" :style="{ color: 'var(--success-color)' }">&#9679; Live</span>
                  <span v-else>{{ contest.status }}</span> &bull;
                  {{ contest.duration }}
                </p>
              </div>
              <n-button class="text-xs text-white px-3 py-1.5 rounded transition-colors cursor-pointer font-terminal"
                @click="$router.push({ name: 'ContestDetail', params: { id: contest.id } })"
                :style="{ backgroundColor: 'var(--btn-primary)' }">
                JOIN
              </n-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Stats from '@/pages/home/Stats.vue'
import ActivityChart from '@/pages/home/ActivityChart.vue'
import { MOCK_PROBLEMS, MOCK_CONTESTS } from '@/constants/mock'
import { RouterLink } from 'vue-router'
import { formatAcceptance } from '@/utils/format'
</script>

<style scoped>
.hero-btn-primary {
  background-color: var(--hero-btn-primary-bg);
  box-shadow: 0 0 16px var(--hero-btn-primary-shadow);
}

.hero-btn-primary:hover {
  background-color: var(--hero-btn-primary-hover);
  box-shadow: 0 0 24px var(--hero-btn-primary-shadow);
}

.hero-btn-secondary:hover {
  background-color: var(--hero-btn-secondary-bg);
}
</style>
