<template>
  <div class="animate-fade-in max-w-7xl mx-auto">
    <!-- Hero Section -->
    <div
      class="relative overflow-hidden rounded-2xl p-8 mb-6"
      :style="{
        background: 'linear-gradient(135deg, var(--contest-hero-from) 0%, var(--contest-hero-to) 100%)',
        border: '1px solid var(--contest-hero-border)'
      }"
    >
      <div class="absolute top-0 right-0 w-96 h-96 opacity-20 blur-3xl">
        <div :style="{ background: 'var(--contest-glow)' }" class="w-full h-full rounded-full"></div>
      </div>

      <div class="relative z-10">
        <div class="flex items-start justify-between mb-6">
          <div>
            <div class="flex items-center gap-3 mb-3">
              <span
                class="px-3 py-1 rounded-full text-xs font-semibold uppercase tracking-wider"
                :style="{
                  background: 'var(--contest-type-bg)',
                  color: 'var(--contest-type-color)'
                }"
              >
                {{ contest.type }}
              </span>
              <span
                v-if="contest.status === 'Live'"
                class="flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold"
                :style="{
                  background: 'var(--contest-live-bg)',
                  color: 'var(--contest-live-color)'
                }"
              >
                <span class="w-2 h-2 rounded-full animate-pulse" :style="{ background: 'var(--contest-live-dot)' }"></span>
                进行中
              </span>
            </div>
            <h1 class="text-4xl font-bold mb-3" :style="{ color: 'var(--contest-title)' }">
              {{ contest.title }}
            </h1>
            <div class="flex items-center gap-6 text-sm" :style="{ color: 'var(--contest-subtitle)' }">
              <div class="flex items-center gap-2">
                <Calendar :size="16" />
                <span>{{ contest.startTime }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Clock :size="16" />
                <span>{{ contest.duration }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Users :size="16" />
                <span>{{ contest.registered }} 人参加</span>
              </div>
            </div>
          </div>

          <!-- Countdown Timer -->
          <div
            class="text-center px-6 py-4 rounded-xl"
            :style="{ background: 'var(--contest-timer-bg)', border: '1px solid var(--contest-timer-border)' }"
          >
            <div class="text-xs mb-1" :style="{ color: 'var(--contest-timer-label)' }">剩余时间</div>
            <div class="text-2xl font-mono font-bold" :style="{ color: 'var(--contest-timer-value)' }">
              {{ countdown }}
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-3">
          <button
            v-if="!isRegistered"
            class="px-6 py-2.5 rounded-lg font-medium transition-all"
            :style="{
              background: 'var(--contest-btn-primary)',
              color: 'var(--contest-btn-primary-text)'
            }"
          >
            立即报名
          </button>
          <button
            v-else
            class="px-6 py-2.5 rounded-lg font-medium transition-all"
            :style="{
              background: 'var(--contest-btn-secondary)',
              color: 'var(--contest-btn-secondary-text)'
            }"
          >
            进入比赛
          </button>
          <button
            class="px-4 py-2.5 rounded-lg font-medium flex items-center gap-2 transition-all"
            :style="{
              background: 'var(--contest-btn-ghost)',
              color: 'var(--contest-btn-ghost-text)'
            }"
          >
            <Share2 :size="16" />
            分享
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid lg:grid-cols-3 gap-6">
      <!-- Left Column: Problems & Standings -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Problems Section -->
        <div
          class="rounded-xl p-6"
          :style="{
            background: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }"
        >
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-xl font-semibold flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <Code2 :size="20" />
              比赛题目
            </h2>
            <span class="text-sm" :style="{ color: 'var(--text-secondary)' }">共 {{ problems.length }} 题</span>
          </div>

          <div class="space-y-3">
            <div
              v-for="problem in problems"
              :key="problem.id"
              class="flex items-center justify-between p-4 rounded-lg transition-all hover:opacity-80 cursor-pointer"
              :style="{
                background: problem.status ? 'var(--contest-problem-solved-bg)' : 'var(--surface-secondary)',
                border: '1px solid var(--border-color)'
              }"
            >
              <div class="flex items-center gap-4">
                <div
                  class="w-10 h-10 rounded-lg flex items-center justify-center font-bold text-sm"
                  :style="{
                    background: getProblemStatusBg(problem.status),
                    color: getProblemStatusColor(problem.status)
                  }"
                >
                  {{ problem.id }}
                </div>
                <div>
                  <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">{{ problem.title }}</h3>
                  <div class="flex items-center gap-3 text-xs mt-1" :style="{ color: 'var(--text-secondary)' }">
                    <span>{{ problem.difficulty }}</span>
                    <span>{{ problem.acceptRate }} 通过率</span>
                    <span>{{ problem.solved }} 人通过</span>
                  </div>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <CheckCircle2
                  v-if="problem.status === 'accepted'"
                  :size="20"
                  :style="{ color: 'var(--success-color)' }"
                />
                <XCircle v-else-if="problem.status === 'wrong'" :size="20" :style="{ color: 'var(--error-color)' }" />
                <Circle v-else :size="20" :style="{ color: 'var(--text-tertiary)' }" />
                <ChevronRight :size="18" :style="{ color: 'var(--text-tertiary)' }" />
              </div>
            </div>
          </div>
        </div>

        <!-- Standings Section -->
        <div
          class="rounded-xl p-6"
          :style="{
            background: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }"
        >
          <div class="flex items-center justify-between mb-5">
            <h2 class="text-xl font-semibold flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <Trophy :size="20" :style="{ color: '#fbbf24' }" />
              实时榜单
            </h2>
            <button class="text-sm font-medium" :style="{ color: 'var(--accent-color)' }">查看全部</button>
          </div>

          <!-- Top 3 Podium -->
          <div class="flex items-end justify-center gap-4 mb-8">
            <!-- 2nd Place -->
            <div class="flex flex-col items-center">
              <img
                :src="rankings[1].avatar"
                class="w-16 h-16 rounded-full border-4 mb-2"
                :style="{ borderColor: 'var(--contest-silver)' }"
              />
              <div class="font-semibold text-sm mb-1" :style="{ color: 'var(--text-primary)' }">
                {{ rankings[1].username }}
              </div>
              <div
                class="w-20 h-20 rounded-t-lg flex flex-col items-center justify-center"
                :style="{ background: 'linear-gradient(to bottom, var(--contest-silver), #94a3b8)' }"
              >
                <span class="text-2xl font-bold text-white">2</span>
                <span class="text-xs text-white/80">{{ rankings[1].score }}</span>
              </div>
            </div>

            <!-- 1st Place -->
            <div class="flex flex-col items-center">
              <div class="relative mb-2">
                <Crown
                  :size="24"
                  class="absolute -top-6 left-1/2 -translate-x-1/2"
                  :style="{ color: '#fbbf24' }"
                />
                <img
                  :src="rankings[0].avatar"
                  class="w-20 h-20 rounded-full border-4"
                  :style="{ borderColor: 'var(--contest-gold)' }"
                />
              </div>
              <div class="font-semibold mb-1" :style="{ color: 'var(--text-primary)' }">
                {{ rankings[0].username }}
              </div>
              <div
                class="w-24 h-28 rounded-t-lg flex flex-col items-center justify-center"
                :style="{ background: 'linear-gradient(to bottom, var(--contest-gold), #f59e0b)' }"
              >
                <span class="text-3xl font-bold text-white">1</span>
                <span class="text-sm text-white/80">{{ rankings[0].score }}</span>
              </div>
            </div>

            <!-- 3rd Place -->
            <div class="flex flex-col items-center">
              <img
                :src="rankings[2].avatar"
                class="w-16 h-16 rounded-full border-4 mb-2"
                :style="{ borderColor: 'var(--contest-bronze)' }"
              />
              <div class="font-semibold text-sm mb-1" :style="{ color: 'var(--text-primary)' }">
                {{ rankings[2].username }}
              </div>
              <div
                class="w-20 h-16 rounded-t-lg flex flex-col items-center justify-center"
                :style="{ background: 'linear-gradient(to bottom, var(--contest-bronze), #b45309)' }"
              >
                <span class="text-2xl font-bold text-white">3</span>
                <span class="text-xs text-white/80">{{ rankings[2].score }}</span>
              </div>
            </div>
          </div>

          <!-- Rankings Table -->
          <div class="space-y-2">
            <div
              v-for="rank in rankings.slice(3)"
              :key="rank.rank"
              class="flex items-center justify-between p-3 rounded-lg"
              :style="{ background: 'var(--surface-secondary)' }"
            >
              <div class="flex items-center gap-3">
                <span class="w-6 text-center font-semibold" :style="{ color: 'var(--text-secondary)' }">
                  {{ rank.rank }}
                </span>
                <img :src="rank.avatar" class="w-8 h-8 rounded-full" />
                <span class="font-medium" :style="{ color: 'var(--text-primary)' }">{{ rank.username }}</span>
              </div>
              <div class="flex items-center gap-6 text-sm">
                <div class="flex items-center gap-1" :style="{ color: 'var(--success-color)' }">
                  <CheckCircle2 :size="14" />
                  <span>{{ rank.solved }}</span>
                </div>
                <span class="font-mono" :style="{ color: 'var(--text-secondary)' }">{{ rank.time }}</span>
                <span class="font-semibold" :style="{ color: 'var(--accent-color)' }">{{ rank.score }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Stats -->
      <div class="space-y-6">
        <!-- Contest Stats -->
        <div
          class="rounded-xl p-6"
          :style="{
            background: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }"
        >
          <h3 class="text-lg font-semibold mb-4" :style="{ color: 'var(--text-primary)' }">比赛统计</h3>
          <div class="grid grid-cols-2 gap-4">
            <div
              class="p-4 rounded-lg"
              :style="{ background: 'var(--surface-secondary)', border: '1px solid var(--border-color)' }"
            >
              <div class="flex items-center gap-2 mb-2" :style="{ color: 'var(--text-secondary)' }">
                <Users :size="16" />
                <span class="text-xs">参与人数</span>
              </div>
              <div class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">{{ contest.registered }}</div>
            </div>
            <div
              class="p-4 rounded-lg"
              :style="{ background: 'var(--surface-secondary)', border: '1px solid var(--border-color)' }"
            >
              <div class="flex items-center gap-2 mb-2" :style="{ color: 'var(--text-secondary)' }">
                <FileCode :size="16" />
                <span class="text-xs">提交总数</span>
              </div>
              <div class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">{{ formatNumber(48520) }}</div>
            </div>
            <div
              class="p-4 rounded-lg"
              :style="{ background: 'var(--surface-secondary)', border: '1px solid var(--border-color)' }"
            >
              <div class="flex items-center gap-2 mb-2" :style="{ color: 'var(--text-secondary)' }">
                <CheckCircle2 :size="16" />
                <span class="text-xs">通过提交</span>
              </div>
              <div class="text-2xl font-bold" :style="{ color: 'var(--success-color)' }">
                {{ formatNumber(32150) }}
              </div>
            </div>
            <div
              class="p-4 rounded-lg"
              :style="{ background: 'var(--surface-secondary)', border: '1px solid var(--border-color)' }"
            >
              <div class="flex items-center gap-2 mb-2" :style="{ color: 'var(--text-secondary)' }">
                <Target :size="16" />
                <span class="text-xs">通过率</span>
              </div>
              <div class="text-2xl font-bold" :style="{ color: 'var(--warning-color)' }">66.2%</div>
            </div>
          </div>
        </div>

        <!-- Difficulty Distribution -->
        <div
          class="rounded-xl p-6"
          :style="{
            background: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }"
        >
          <h3 class="text-lg font-semibold mb-4" :style="{ color: 'var(--text-primary)' }">难度分布</h3>
          <div class="space-y-3">
            <div>
              <div class="flex items-center justify-between text-sm mb-1">
                <span :style="{ color: 'var(--text-secondary)' }">简单</span>
                <span class="font-medium" :style="{ color: 'var(--success-color)' }">2 题</span>
              </div>
              <div class="h-2 rounded-full overflow-hidden" :style="{ background: 'var(--surface-tertiary)' }">
                <div class="h-full rounded-full" :style="{ width: '40%', background: 'var(--success-color)' }"></div>
              </div>
            </div>
            <div>
              <div class="flex items-center justify-between text-sm mb-1">
                <span :style="{ color: 'var(--text-secondary)' }">中等</span>
                <span class="font-medium" :style="{ color: 'var(--warning-color)' }">2 题</span>
              </div>
              <div class="h-2 rounded-full overflow-hidden" :style="{ background: 'var(--surface-tertiary)' }">
                <div class="h-full rounded-full" :style="{ width: '40%', background: 'var(--warning-color)' }"></div>
              </div>
            </div>
            <div>
              <div class="flex items-center justify-between text-sm mb-1">
                <span :style="{ color: 'var(--text-secondary)' }">困难</span>
                <span class="font-medium" :style="{ color: 'var(--error-color)' }">1 题</span>
              </div>
              <div class="h-2 rounded-full overflow-hidden" :style="{ background: 'var(--surface-tertiary)' }">
                <div class="h-full rounded-full" :style="{ width: '20%', background: 'var(--error-color)' }"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Announcement -->
        <div
          class="rounded-xl p-6"
          :style="{
            background: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }"
        >
          <h3 class="text-lg font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
            <Megaphone :size="18" />
            公告
          </h3>
          <div class="space-y-3">
            <div class="p-3 rounded-lg" :style="{ background: 'var(--surface-secondary)' }">
              <div class="text-xs mb-1" :style="{ color: 'var(--text-tertiary)' }">管理员 • 5分钟前</div>
              <p class="text-sm" :style="{ color: 'var(--text-primary)' }">
                比赛进行到一半，请注意提交时间。如有问题请及时反馈。
              </p>
            </div>
            <div class="p-3 rounded-lg" :style="{ background: 'var(--surface-secondary)' }">
              <div class="text-xs mb-1" :style="{ color: 'var(--text-tertiary)' }">管理员 • 30分钟前</div>
              <p class="text-sm" :style="{ color: 'var(--text-primary)' }">
                比赛已开始，祝大家取得好成绩！
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import {
  Calendar,
  Clock,
  Users,
  Share2,
  Code2,
  Trophy,
  CheckCircle2,
  XCircle,
  Circle,
  ChevronRight,
  FileCode,
  Target,
  Crown,
  Megaphone
} from 'lucide-vue-next'
import { MOCK_CONTESTS, MOCK_CONTEST_PROBLEMS, MOCK_CONTEST_RANKING } from '@/constants'
import type { ContestProblem, ContestRanking } from '@/constants'

const route = useRoute()
const contestId = route.params.id as string

const contest = computed(
  () => MOCK_CONTESTS.find((c) => c.id === contestId) || MOCK_CONTESTS[0]
)

const problems = ref<ContestProblem[]>(MOCK_CONTEST_PROBLEMS)
const rankings = ref<ContestRanking[]>(MOCK_CONTEST_RANKING)
const isRegistered = ref(false)

let countdownInterval: ReturnType<typeof setInterval> | null = null
const remainingSeconds = ref(47 * 60 + 32) // Mock: 47:32 remaining

const countdown = computed(() => {
  const hours = Math.floor(remainingSeconds.value / 3600)
  const minutes = Math.floor((remainingSeconds.value % 3600) / 60)
  const seconds = remainingSeconds.value % 60
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds
    .toString()
    .padStart(2, '0')}`
})

const formatNumber = (num: number) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

const getProblemStatusBg = (status?: string) => {
  if (status === 'accepted') return 'var(--success-bg)'
  if (status === 'wrong') return 'var(--error-bg)'
  return 'var(--surface-tertiary)'
}

const getProblemStatusColor = (status?: string) => {
  if (status === 'accepted') return 'var(--success-color)'
  if (status === 'wrong') return 'var(--error-color)'
  return 'var(--text-tertiary)'
}

onMounted(() => {
  countdownInterval = setInterval(() => {
    if (remainingSeconds.value > 0) {
      remainingSeconds.value--
    }
  }, 1000)
})

onUnmounted(() => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
  }
})
</script>
