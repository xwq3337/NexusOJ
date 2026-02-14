<template>
  <div class="min-h-screen bg-transparent">
    <!-- 标题 -->
    <div class="text-center mb-6">
      <h1 class="text-3xl font-bold mb-2" :style="{ color: 'var(--text-primary)' }">
        排行榜
      </h1>
      <p :style="{ color: 'var(--text-secondary)' }">
        查看平台最活跃的用户
      </p>
    </div>

    <!-- 排行榜列表 -->
    <n-card :bordered="false" class="max-w-4xl mx-auto">
      <n-list>
        <n-list-item v-for="user in leaderboard" :key="user.id" class="py-4! px-6!">
          <div class="flex items-center gap-4">
            <!-- 排名 -->
            <div class="flex items-center justify-center w-10 h-10 rounded-full text-sm font-bold shrink-0" :style="{
              backgroundColor: user.rank <= 3 ? 'var(--hero-bg-from)' : 'var(--divider-color)',
              color: user.rank <= 3 ? '#fff' : 'var(--text-secondary)'
            }">
              {{ user.rank }}
            </div>

            <!-- 头像 -->
            <n-avatar round :size="48" :src="user.avatar" class="shrink-0" />

            <!-- 用户信息 -->
            <div class="flex-1 min-w-0">
              <div class="font-semibold text-base truncate" :style="{ color: 'var(--text-primary)' }">
                {{ user.username }}
              </div>
              <div class="text-sm mt-0.5" :style="{ color: 'var(--text-secondary)' }">
                解决了 {{ user.solved }} 道题目
              </div>
            </div>

            <!-- 评分 -->
            <div class="text-right shrink-0">
              <div class="text-2xl font-bold" :style="{ color: 'var(--hero-bg-from)' }">
                {{ user.rating }}
              </div>
              <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                Rating
              </div>
            </div>
          </div>
        </n-list-item>
      </n-list>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NCard, NList, NListItem, NAvatar } from 'naive-ui'
import Request from '@/services/api'

interface User {
  id: number
  username: string
  avatar: string
  rank: number
  solved: number
  rating: number
}

const leaderboard = ref<User[]>([])
const loading = ref(true)
const MOCK_LEADERBOARD = [
  { id: 1, username: '张三', avatar: 'https://picsum.photos/200', rank: 1, solved: 156, rating: 2450 },
  { id: 2, username: '李四', avatar: 'https://picsum.photos/200', rank: 2, solved: 142, rating: 2380 },
  { id: 3, username: '王五', avatar: 'https://picsum.photos/200', rank: 3, solved: 138, rating: 2310 },
  { id: 4, username: '赵六', avatar: 'https://picsum.photos/200', rank: 4, solved: 125, rating: 2250 },
  { id: 5, username: '钱七', avatar: 'https://picsum.photos/200', rank: 5, solved: 118, rating: 2190 },
  { id: 6, username: '孙八', avatar: 'https://picsum.photos/200', rank: 6, solved: 112, rating: 2150 },
  { id: 7, username: '周九', avatar: 'https://picsum.photos/200', rank: 7, solved: 108, rating: 2110 },
  { id: 8, username: '吴十', avatar: 'https://picsum.photos/200', rank: 8, solved: 102, rating: 2070 },
  { id: 9, username: '郑一', avatar: 'https://picsum.photos/200', rank: 9, solved: 98, rating: 2030 },
  { id: 10, username: '王二', avatar: 'https://picsum.photos/200', rank: 10, solved: 95, rating: 1990 }
]
onMounted(async () => {
  const response = await Request.get('/user/leaderboard')
  leaderboard.value = response.info || MOCK_LEADERBOARD
  loading.value = false
})
</script>
