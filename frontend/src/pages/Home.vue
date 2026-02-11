<template>
  <div class="animate-fade-in">
    <div
      class="relative overflow-hidden rounded-3xl bg-gradient-to-r from-blue-900/40 to-indigo-900/40 border border-white/10 p-8 md:p-12 mb-12"
    >
      <div class="relative z-10 max-w-2xl">
        <h1
          class="text-4xl md:text-5xl font-extrabold mb-6 leading-tight"
          :style="{ color: 'var(--text-primary)' }"
        >
          Master Algorithms <br />
          <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 to-indigo-400"
            >Build the Future</span
          >
        </h1>
        <p class="text-gray-300 text-lg mb-8 leading-relaxed">
          加入新一代的编程竞技。在这里，你可以解决各种算法问题，参加大型的编程竞赛，并获取实时的人工智能指导，帮你理清思路、调试代码。
        </p>
        <div class="flex flex-wrap gap-4">
          <RouterLink
            to="/problems"
            class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-xl font-semibold transition-all hover:scale-105 shadow-lg shadow-blue-500/25"
            >开始解题</RouterLink
          >
          <RouterLink
            to="/contests"
            class="text-white px-8 py-3 rounded-xl font-semibold transition-all"
            :style="{ backgroundColor: 'var(--btn-secondary)' }"
            :class="['hover:bg-gray-700']"
            >浏览比赛</RouterLink
          >
        </div>
      </div>

      <div
        class="absolute top-0 right-0 -mt-20 -mr-20 w-96 h-96 bg-blue-500/20 rounded-full blur-3xl"
      ></div>
      <div class="absolute bottom-0 left-20 w-72 h-72 bg-purple-500/10 rounded-full blur-3xl"></div>
    </div>

    <Stats />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-12">
      <div class="lg:col-span-2">
        <ActivityChart />
      </div>

      <div
        class="rounded-xl p-6"
        :style="{
          backgroundColor: 'var(--bg-secondary)',
          borderColor: 'var(--border-color)',
          borderWidth: '1px',
          borderStyle: 'solid'
        }"
      >
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <span class="text-yellow-500">⚡</span>
            <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">每日挑战</h3>
          </div>
          <RouterLink to="/problems" class="text-xs text-blue-400 hover:text-blue-300"
            >详情</RouterLink
          >
        </div>

        <div class="space-y-4">
          <div
            v-for="problem in MOCK_PROBLEMS.slice(0, 3)"
            :key="problem.id"
            class="group flex items-center justify-between p-3 rounded-lg transition-colors cursor-pointer border border-transparent"
            @click="$router.push({ name: 'ProblemDetail', params: { id: problem.id } })"
          >
            <div class="flex items-center gap-3">
              <component
                :is="problem.status === 'solved' ? CheckCircle2 : Circle"
                :size="18"
                :class="
                  problem.status === 'solved'
                    ? 'text-green-500'
                    : 'text-gray-600 group-hover:text-gray-400'
                "
              />
              <div>
                <h4
                  class="text-sm font-medium group-hover:text-blue-400 transition-colors"
                  :style="{ color: 'var(--text-primary)' }"
                >
                  {{ problem.title }}
                </h4>
                <div class="flex items-center gap-2 mt-1">
                  <span
                    v-for="tag in problem.tags.slice(0, 2)"
                    :key="tag"
                    class="text-[10px] px-1.5 py-0.5 rounded text-gray-400 border"
                    :style="{
                      backgroundColor: 'var(--surface-tertiary)',
                      borderColor: 'var(--border-light)',
                      borderWidth: '1px',
                      borderStyle: 'solid'
                    }"
                    >{{ tag }}</span
                  >
                </div>
              </div>
            </div>
            <div class="flex flex-col items-end gap-1">
              <span
                :class="[
                  `text-xs px-2 py-0.5 rounded-full font-medium`,
                  problem.difficulty === 'Easy'
                    ? 'bg-green-500/10 text-green-400'
                    : problem.difficulty === 'Medium'
                      ? 'bg-yellow-500/10 text-yellow-400'
                      : 'bg-red-500/10 text-red-400'
                ]"
                >{{ problem.difficulty }}</span
              >
              <span class="text-[10px] text-gray-500">{{ problem.acceptance }}</span>
            </div>
          </div>
        </div>

        <div
          class="mt-6 pt-6 border-t"
          :style="{
            borderColor: 'var(--border-color)',
            borderTopWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <span class="text-blue-400">🌐</span>
              <h3 class="text-lg font-bold" :style="{ color: 'var(--text-primary)' }">活跃比赛</h3>
            </div>
            <RouterLink to="/contests" class="text-xs text-blue-400 hover:text-blue-300"
              >详情</RouterLink
            >
          </div>
          <div class="space-y-3">
            <div
              v-for="contest in MOCK_CONTESTS.slice(0, 2)"
              :key="contest.id"
              class="p-3 rounded-lg flex items-center justify-between"
              :style="{
                backgroundColor: 'var(--surface-secondary)',
                borderColor: 'var(--border-color)',
                borderWidth: '1px',
                borderStyle: 'solid'
              }"
            >
              <div>
                <h4 class="text-sm font-medium" :style="{ color: 'var(--text-primary)' }">
                  {{ contest.title }}
                </h4>
                <p class="text-xs text-gray-500 mt-0.5">
                  <span v-if="contest.status === 'Live'" class="text-red-400 animate-pulse"
                    >● Live</span
                  >
                  <span v-else>{{ contest.status }}</span> •
                  {{ contest.duration }}
                </p>
              </div>
              <button
                class="text-xs text-white px-3 py-1.5 rounded transition-colors"
                :style="{ backgroundColor: 'var(--btn-secondary)' }"
                :class="['hover:bg-gray-700']"
              >
                Details
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import Stats from '@/components/Stats.vue'
import ActivityChart from '@/components/ActivityChart.vue'
import { MOCK_PROBLEMS, MOCK_CONTESTS } from '@/constants'
import { CheckCircle2, Circle } from 'lucide-vue-next'
import { RouterLink } from 'vue-router'
import { ref } from 'vue'
</script>
