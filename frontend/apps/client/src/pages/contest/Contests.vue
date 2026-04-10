<template>
  <div class="animate-fade-in max-w-7xl mx-auto">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- 左侧：比赛列表 -->
      <div class="lg:col-span-3">
        <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">比赛</h1>
        <div class="grid gap-4">
          <div
            v-for="c in contests"
            :key="c.id"
            class="p-4 rounded-lg flex items-center justify-between cursor-pointer transition-all hover:opacity-90"
            :style="{
              backgroundColor: 'var(--surface-primary)',
              borderColor: 'var(--border-color)',
              borderWidth: '1px',
              borderStyle: 'solid'
            }"
            @click="$router.push({ name: 'ContestProblems', params: { id: c.id } })"
          >
            <div>
              <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                {{ c.title }}
              </h3>
              <p class="text-xs text-gray-400">
                {{ formatTime(c.begin_at) }} • {{ c.duration }}分钟 • {{ c.participants || 0 }} 人参加
              </p>
            </div>
            <div class="flex items-center gap-2">
              <n-tag :type="statusTagType(c.status)" size="small" round>
                {{ statusLabel(c.status) }}
              </n-tag>
              <n-tag v-if="c.is_private" size="small" :bordered="false" type="warning">
                私密
              </n-tag>
              <n-button
                v-if="c.is_registered === false"
                size="small"
                type="primary"
                @click.stop="handleRegister(c)"
              >
                报名
              </n-button>
              <n-tag v-else size="small" type="success" round>已报名</n-tag>
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
              v-for="(user, index) in userRanking"
              :key="user.username"
              class="flex items-center gap-3 p-2 rounded cursor-pointer transition-all hover:opacity-80"
              :style="{
                backgroundColor: index < 3 ? 'var(--surface-secondary)' : 'transparent'
              }"
              @click="$router.push({ name: 'UserHomePage', params: { id: String(user.id) } })"
            >
              <div
                class="w-6 h-6 flex items-center justify-center text-xs font-bold rounded"
                :style="{
                  backgroundColor: getRankColor(index).bg,
                  color: getRankColor(index).text
                }"
              >
                {{ index + 1 }}
              </div>
              <img :src="user.avatar" :alt="user.username" class="w-8 h-8 rounded-full" />
              <div class="flex-1 min-w-0">
                <div class="font-medium text-sm truncate" :style="{ color: 'var(--text-primary)' }">
                  {{ user.nickname || user.username }}
                </div>
                <div class="text-xs text-gray-400 truncate">
                  {{ user.school || user.username }}
                </div>
              </div>
              <div class="text-sm font-bold" :style="{ color: getRatingColor(user.rating) }">
                {{ user.rating }}
              </div>
            </div>
          </div>
          <div
            class="mt-4 pt-4 text-center cursor-pointer text-sm hover:opacity-80"
            :style="{ color: 'var(--text-secondary)', borderTop: '1px solid var(--border-color)' }"
          >
            Top 10 用户
          </div>
        </div>
      </div>
    </div>

    <!-- 密码弹窗 -->
    <n-modal v-model:show="showPasswordModal" preset="dialog" title="私密比赛" positive-text="确认报名"
      negative-text="取消" :loading="registering" @positive-click="submitRegister"
      @negative-click="cancelRegister">
      <n-space vertical>
        <n-text depth="3">该比赛需要密码才能报名</n-text>
        <n-input v-model:value="passwordInput" type="password" show-password-on="click"
          placeholder="请输入比赛密码" @keyup.enter="submitRegister" />
      </n-space>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NButton, NTag, NModal, NInput, NSpace, NText, useMessage } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'
import { contestApi, userApi } from '@nexusoj/server'
import type { Contest, User } from '@nexusoj/type'

const message = useMessage()
const contests = ref<(Contest & { is_registered: boolean })[]>([])
const userRanking = ref<User[]>([])
const showPasswordModal = ref(false)
const passwordInput = ref('')
const registering = ref(false)
const pendingContest = ref<Contest | null>(null)

const formatTime = (time: string) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const statusTagType = (status: string): 'success' | 'info' | 'warning' => {
  if (status === 'Live') return 'success'
  if (status === 'Upcoming') return 'info'
  return 'warning'
}

const statusLabel = (status: string) => {
  if (status === 'Live') return '进行中'
  if (status === 'Upcoming') return '即将开始'
  return '已结束'
}

const handleRegister = (contest: Contest) => {
  
  pendingContest.value = contest
  if (contest.is_private) {
    passwordInput.value = ''
    showPasswordModal.value = true
  } else {
    submitRegister()
  }
}

const cancelRegister = () => {
  showPasswordModal.value = false
  pendingContest.value = null
}

const submitRegister = async (): Promise<boolean | undefined> => {
  if (!pendingContest.value) return false
  registering.value = true
  try {
    const res = await contestApi.registerContest({
      contest_id: pendingContest.value.id,
      password: passwordInput.value || undefined,
    })
    if (res.code === 200) {
      message.success('报名成功')
      pendingContest.value = null
      registering.value = false
      // 更改比赛的is_registered
      const contest = contests.value.find(c => c.id === pendingContest.value?.id)
      if (contest) {
        contest.is_registered = true
        contest.participants = contest.participants  + 1
      }
      return true // 允许弹窗关闭
    }
    message.error(String(res.msg || '报名失败'))
  } catch (e: any) {
    message.error(e?.response?.data?.msg || '报名失败，请稍后重试')
  }
  registering.value = false
  return false // 阻止弹窗关闭，让用户可以重试
}

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
  if (rating >= 1200) return '#38bdf8'
  return '#808080'
}

onMounted(async () => {
  try {
    const res = await contestApi.getContestList(1, 20)
    const { code, info } = res
    if (code === 200 && info) {
      contests.value = info.list || []
    }
  } catch (e) {
    console.error(e)
  }
  try {
    const res = await userApi.getTopRating()
    if (res.code === 200 && res.info) {
      userRanking.value = res.info
    }
  } catch (e) {
    console.error(e)
  }
})
</script>
