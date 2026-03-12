<template>
  <div class="animate-fade-in max-w-7xl mx-auto">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- 左侧：比赛列表 -->
      <div class="lg:col-span-3">
        <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">比赛</h1>
        <div class="grid gap-4">
          <div
            v-for="c in MOCK_CONTESTS"
            :key="c.id"
            class="p-4 rounded-lg flex items-center justify-between cursor-pointer transition-all hover:opacity-90"
            :style="{
              backgroundColor: 'var(--surface-primary)',
              borderColor: 'var(--border-color)',
              borderWidth: '1px',
              borderStyle: 'solid'
            }"
            @click="$router.push({ name: 'ContestDetail', params: { id: c.id } })"
          >
            <div>
              <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                {{ c.title }}
              </h3>
              <p class="text-xs text-gray-400">
                {{ c.startTime }} • {{ c.duration }} • {{ c.registered }} 人参加
              </p>
            </div>
            <div class="flex items-center gap-2">
              <span
                class="px-2 py-1 rounded text-xs font-medium"
                :style="{
                  background: c.status === 'Live' ? 'var(--success-bg)' : 'var(--surface-secondary)',
                  color: c.status === 'Live' ? 'var(--success-color)' : 'var(--text-secondary)'
                }"
              >
                {{ c.status === 'Live' ? '进行中' : c.status === 'Upcoming' ? '即将开始' : '已结束' }}
              </span>
              <button
                class="text-xs px-3 py-1.5 rounded cursor-pointer"
                :style="{
                  color: 'var(--text-primary)',
                  backgroundColor: 'var(--btn-primary)',
                  '--tw-bg-opacity': '1'
                }"
                @click.stop
              >
                报名
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：排行榜 -->
      <div class="lg:col-span-1">
        <div
          class="p-4 rounded-lg"
          :style="{
            backgroundColor: 'var(--surface-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }"
        >
          <h2 class="text-lg font-bold mb-4" :style="{ color: 'var(--text-primary)' }">Rating 排行榜</h2>
          <div class="space-y-2">
            <div
              v-for="(user, index) in MOCK_USER_RANKING"
              :key="user.username"
              class="flex items-center gap-3 p-2 rounded cursor-pointer transition-all hover:opacity-80"
              :style="{
                backgroundColor: index < 3 ? 'var(--surface-secondary)' : 'transparent'
              }"
              @click="$router.push({ name: 'UserPersonalCenter', params: { id: user.username } })"
            >
              <!-- 排名 -->
              <div
                class="w-6 h-6 flex items-center justify-center text-xs font-bold rounded"
                :style="{
                  backgroundColor: getRankColor(index).bg,
                  color: getRankColor(index).text
                }"
              >
                {{ user.rank }}
              </div>

              <!-- 头像 -->
              <img
                :src="user.avatar"
                :alt="user.username"
                class="w-8 h-8 rounded-full"
              />

              <!-- 用户信息 -->
              <div class="flex-1 min-w-0">
                <div class="font-medium text-sm truncate" :style="{ color: 'var(--text-primary)' }">
                  {{ user.nickname || user.username }}
                </div>
                <div class="text-xs text-gray-400 truncate">
                  {{ user.school || user.username }}
                </div>
              </div>

              <!-- Rating -->
              <div class="text-sm font-bold" :style="{ color: getRatingColor(user.rating) }">
                {{ user.rating }}
              </div>
            </div>
          </div>

          <!-- 查看更多 -->
          <div
            class="mt-4 pt-4 text-center cursor-pointer text-sm hover:opacity-80"
            :style="{
              color: 'var(--text-secondary)',
              borderTop: '1px solid var(--border-color)'
            }"
          >
            Top 10 用户
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { MOCK_CONTESTS, MOCK_USER_RANKING } from '@/constants/mock'
import { ref } from 'vue'

const getRankColor = (index: number) => {
  if (index === 0) return { bg: '#FFD700', text: '#fff' }
  if (index === 1) return { bg: '#C0C0C0', text: '#fff' }
  if (index === 2) return { bg: '#CD7F32', text: '#fff' }
  return { bg: 'var(--surface-secondary)', text: 'var(--text-secondary)' }
}

const getRatingColor = (rating: number) => {
  if (rating >= 2600) return '#ff0000'
  if (rating >= 2400) return '#ff0000'
  if (rating >= 2200) return '#ff8c00'
  if (rating >= 2000) return '#ffcc00'
  if (rating >= 1800) return '#ffff00'
  if (rating >= 1600) return '#00ff00'
  if (rating >= 1400) return '#00ffff'
  if (rating >= 1200) return '#0000ff'
  return '#808080'
}
</script>
