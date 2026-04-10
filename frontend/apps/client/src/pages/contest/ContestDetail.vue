<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h, provide } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { useMessage } from 'naive-ui'
import {
  NButton,
  NTag,
  NCard,
  NModal,
  NInput,
  NSpace,
  NIcon,
  NMenu,
  NDivider
} from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import {
  Calendar,
  Clock,
  Users,
  Share2,
  FileCode,
  BarChart,
  Trophy,
  BookOpen
} from 'lucide-vue-next'
import { contestApi } from '@nexusoj/server'

const message = useMessage()
const route = useRoute()
const router = useRouter()
const contestId = route.params.id as string
const isRegistered = ref(false)
const contest = ref<any>()
const problems = ref<any[]>([])
const passwordModal = ref(false)
const passwordInput = ref('')

// Provide shared data to child routes
provide('contestData', { contest, isRegistered, contestId, problems })

const fetchContestDetail = async () => {
  try {
    const { code, info } = await contestApi.getContestProblems(contestId)
    if (code === 200 && info) {
      contest.value = info.contest
      problems.value = info.problems
      isRegistered.value = true
    }
  } catch (e) {
    console.error(e)
  }
}

const handleRegister = async () => {
  if (contest.value?.is_private && !passwordInput.value) {
    passwordModal.value = true
    return
  }
  try {
    const res = await contestApi.registerContest({
      contest_id: contestId,
      password: passwordInput.value || undefined
    })
    if (res.code === 200) {
      isRegistered.value = true
      passwordModal.value = false
      message.success('报名成功')
    }
  } catch (e) {
    console.error(e)
  }
}

const handleEnterContest = () => {
  router.push({ name: 'ContestProblem', params: { id: contestId, label: 'A' } })
}

const handleShare = () => {
  const url = window.location.href
  navigator.clipboard.writeText(url).then(() => {
    message.success('链接已复制到剪贴板')
  })
}

let countdownInterval: ReturnType<typeof setInterval> | null = null
const now = ref(Date.now())

const formatContestTime = (t: string) => (t ? new Date(t).toLocaleString('zh-CN') : '-')

const countdown = computed(() => {
  if (!contest.value?.begin_at) return '--:--:--'
  const endAt = new Date(contest.value.end_at).getTime()
  const beginAt = new Date(contest.value.begin_at).getTime()
  const current = now.value

  let diff: number
  if (current >= beginAt && current < endAt) {
    diff = endAt - current
  } else if (current < beginAt) {
    diff = beginAt - current
  } else {
    return '00:00:00'
  }

  if (diff <= 0) return '00:00:00'
  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)
  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
})

const activeTab = computed(() => route.name as string)

const menuOptions: MenuOption[] = [
  {
    label: () =>
      h(
        RouterLink,
        { to: { name: 'ContestProblems', params: { id: contestId } } },
        { default: () => '题目' }
      ),
    key: 'ContestProblems',
    icon: () => h(FileCode)
  },
  {
    label: () =>
      h(
        RouterLink,
        { to: { name: 'ContestSubmissions', params: { id: contestId } } },
        { default: () => '提交' }
      ),
    key: 'ContestSubmissions',
    icon: () => h(BarChart)
  },
  {
    label: () =>
      h(
        RouterLink,
        { to: { name: 'ContestRanking', params: { id: contestId } } },
        { default: () => '排名' }
      ),
    key: 'ContestRanking',
    icon: () => h(Trophy)
  },
  {
    label: () =>
      h(
        RouterLink,
        { to: { name: 'ContestEditorial', params: { id: contestId } } },
        { default: () => '赛后题解' }
      ),
    key: 'ContestEditorial',
    icon: () => h(BookOpen)
  }
]

onMounted(async () => {
  await fetchContestDetail()

  countdownInterval = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  if (countdownInterval) {
    clearInterval(countdownInterval)
  }
})
</script>

<template>
  <div class="animate-fade-in max-w-7xl mx-auto">
    <!-- Hero Section -->
    <NCard
      :style="{
        background: 'linear-gradient(135deg, var(--contest-hero-from) 0%, var(--contest-hero-to) 100%)',
        border: '1px solid var(--contest-hero-border)'
      }"
      :bordered="true"
      class="mb-6"
      content-style="padding: 32px;"
    >
      <div class="flex items-start justify-between flex-wrap gap-4">
        <div class="flex-1 min-w-0">
          <NSpace align="center" class="mb-3">
            <NTag :type="contest?.contest_type === 'ACM' ? 'info' : 'warning'" size="small" round>
              {{ contest?.contest_type || '' }} 赛制
            </NTag>
            <NTag v-if="contest?.status === 'Live'" type="success" size="small" round>
              <template #icon>
                <NIcon><span class="inline-block w-1.5 h-1.5 rounded-full animate-pulse" style="background: currentColor" /></NIcon>
              </template>
              进行中
            </NTag>
            <NTag v-else-if="contest?.status === 'Upcoming'" size="small" round>未开始</NTag>
            <NTag v-else-if="contest?.status === 'Ended'" type="default" size="small" round>已结束</NTag>
          </NSpace>

          <h1 class="text-3xl font-bold mb-3" :style="{ color: 'var(--contest-title)' }">
            {{ contest?.title || '' }}
          </h1>

          <NSpace :size="20" class="text-sm" :style="{ color: 'var(--contest-subtitle)' }">
            <NSpace :size="4" align="center">
              <NIcon :size="16"><Calendar /></NIcon>
              <span>{{ formatContestTime(contest?.begin_at) }}</span>
            </NSpace>
            <NSpace :size="4" align="center">
              <NIcon :size="16"><Clock /></NIcon>
              <span>{{ contest?.duration || 0 }} 分钟</span>
            </NSpace>
            <NSpace :size="4" align="center">
              <NIcon :size="16"><Users /></NIcon>
              <span>{{ contest?.participants }} 人参加</span>
            </NSpace>
          </NSpace>
        </div>

        <!-- Countdown Timer -->
        <div
          class="text-center px-6 py-4 rounded-xl"
          :style="{ background: 'var(--contest-timer-bg)', border: '1px solid var(--contest-timer-border)' }"
        >
          <div class="text-xs mb-1" :style="{ color: 'var(--contest-timer-label)' }">
            {{ contest?.status === 'Live' ? '剩余时间' : contest?.status === 'Upcoming' ? '距离开始' : '已结束' }}
          </div>
          <div class="text-3xl font-bold" :style="{ color: 'var(--contest-timer-value)' }">
            {{ countdown }}
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <NDivider style="margin: 20px 0 16px" />
      <NSpace>
        <NButton v-if="!isRegistered" type="primary" @click="handleRegister">立即报名</NButton>
        <NButton v-else type="primary" @click="handleEnterContest">进入比赛</NButton>
        <NButton @click="handleShare">
          <template #icon><NIcon><Share2 /></NIcon></template>
          分享
        </NButton>
      </NSpace>
    </NCard>

    <!-- Tab Navigation -->
    <NMenu mode="horizontal" :value="activeTab" :options="menuOptions" class="mb-6" />

    <!-- Sub-page content -->
    <RouterView />
  </div>

  <!-- Password Modal -->
  <NModal v-model:show="passwordModal" preset="dialog" title="输入比赛密码" positive-text="确认" negative-text="取消"
    @positive-click="handleRegister"
  >
    <NInput v-model:value="passwordInput" type="password" show-password-on="click" placeholder="请输入比赛密码" />
  </NModal>
</template>
