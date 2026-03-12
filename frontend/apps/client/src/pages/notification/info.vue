<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">系统通知</h1>
      <n-button
        v-if="hasUnread"
        type="primary"
        size="small"
        @click="markAllAsRead"
      >
        全部已读
      </n-button>
    </div>

    <n-spin :show="loading">
      <div class="space-y-4">
        <n-empty v-if="!loading && notifications.length === 0" description="暂无通知" />

        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="p-4 rounded-lg border transition-all cursor-pointer hover:shadow-md"
          :class="{ 'opacity-60': notification.read }"
          :style="{
            backgroundColor: notification.read ? 'var(--hover-bg)' : 'var(--card-bg)',
            borderColor: notification.read ? 'var(--border-color)' : 'var(--primary-color)'
          }"
          @click="markAsRead(notification.id)"
        >
          <div class="flex items-start gap-4">
            <div
              class="shrink-0 w-10 h-10 rounded-full flex items-center justify-center"
              :style="{ backgroundColor: getNotificationColor(notification.type) }"
            >
              <component :is="getNotificationIcon(notification.type)" :size="20" />
            </div>

            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span
                  class="text-xs px-2 py-0.5 rounded"
                  :style="{
                    backgroundColor: getNotificationColor(notification.type),
                    color: 'var(--text-primary)'
                  }"
                >
                  {{ getNotificationTypeLabel(notification.type) }}
                </span>
                <span
                  v-if="!notification.read"
                  class="w-2 h-2 rounded-full"
                  :style="{ backgroundColor: 'var(--primary-color)' }"
                />
              </div>

              <h3
                class="font-medium mb-1"
                :style="{ color: 'var(--text-primary)' }"
              >
                {{ notification.title }}
              </h3>

              <p
                class="text-sm mb-2 line-clamp-2"
                :style="{ color: 'var(--text-secondary)' }"
              >
                {{ notification.content }}
              </p>

              <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                {{ notification.time }}
              </p>
            </div>

            <n-button
              v-if="notification.action"
              text
              size="small"
              @click.stop="handleAction(notification)"
            >
              {{ notification.action }}
            </n-button>
          </div>
        </div>
      </div>
    </n-spin>

    <div v-if="notifications.length > 0" class="flex justify-center mt-6">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-size="pageSize"
        show-quick-jumper
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Bell,
  Trophy,
  AlertCircle,
  CheckCircle,
  Info,
  Gift
} from 'lucide-vue-next'
import { NButton, NEmpty, NSpin, NPagination } from 'naive-ui'
import { useRouter } from 'vue-router'

const router = useRouter()

// Types
interface SystemNotification {
  id: string
  type: 'system' | 'contest' | 'achievement' | 'warning' | 'reward'
  title: string
  content: string
  time: string
  read: boolean
  action?: string
  link?: string
}

// MOCK Data
const MOCK_NOTIFICATIONS: SystemNotification[] = [
  {
    id: '1',
    type: 'contest',
    title: '比赛即将开始',
    content: '周赛 #123 将在 30 分钟后开始，请做好准备！比赛时长为 2 小时，共 5 道题目。',
    time: '10 分钟前',
    read: false,
    action: '查看详情',
    link: '/contest/123'
  },
  {
    id: '2',
    type: 'achievement',
    title: '恭喜获得新成就！',
    content: '您已成功解决 100 道题目，解锁成就「百题斩」！获得 50 豆豆奖励。',
    time: '1 小时前',
    read: false,
    action: '查看成就',
    link: '/user/personal-center'
  },
  {
    id: '3',
    type: 'system',
    title: '系统维护通知',
    content: '系统将于 2026-02-18 02:00 - 04:00 进行服务器维护，届时将暂停服务。',
    time: '2 小时前',
    read: true
  },
  {
    id: '4',
    type: 'reward',
    title: '签到成功',
    content: '今日签到成功！获得 10 豆豆，连续签到 7 天可获得额外奖励。',
    time: '5 小时前',
    read: false,
    action: '立即签到',
    link: '/user/personal-center'
  },
  {
    id: '5',
    type: 'warning',
    title: '账号安全提醒',
    content: '检测到您的账号在新设备上登录，如非本人操作请及时修改密码。',
    time: '1 天前',
    read: true,
    action: '修改密码'
  },
  {
    id: '6',
    type: 'contest',
    title: '比赛结果公布',
    content: '双周赛 #122 已结束，您的排名为第 42 名，获得 20 豆豆奖励。',
    time: '2 天前',
    read: true,
    action: '查看详情',
    link: '/contest/122'
  },
  {
    id: '7',
    type: 'achievement',
    title: '新题目上线',
    content: '数据结构专题新增 10 道题目，快来挑战吧！',
    time: '3 天前',
    read: true,
    action: '开始做题',
    link: '/problems'
  }
]

// State
const loading = ref(false)
const currentPage = ref(1)
const pageSize = 10
const notifications = ref<SystemNotification[]>([])
const totalPages = computed(() => Math.ceil(MOCK_NOTIFICATIONS.length / pageSize))

const hasUnread = computed(() => notifications.value.some(n => !n.read))

// Methods
const getNotificationIcon = (type: string) => {
  const icons = {
    system: Bell,
    contest: Trophy,
    achievement: CheckCircle,
    warning: AlertCircle,
    reward: Gift
  }
  return icons[type as keyof typeof icons] || Info
}

const getNotificationColor = (type: string) => {
  const colors = {
    system: 'rgba(32, 128, 240, 0.1)',
    contest: 'rgba(240, 160, 32, 0.1)',
    achievement: 'rgba(24, 160, 88, 0.1)',
    warning: 'rgba(208, 48, 88, 0.1)',
    reward: 'rgba(250, 173, 20, 0.1)'
  }
  return colors[type as keyof typeof colors] || 'rgba(32, 128, 240, 0.1)'
}

const getNotificationTypeLabel = (type: string) => {
  const labels = {
    system: '系统',
    contest: '比赛',
    achievement: '成就',
    warning: '警告',
    reward: '奖励'
  }
  return labels[type as keyof typeof labels] || '通知'
}

const markAsRead = (id: string) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification && !notification.read) {
    notification.read = true
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
}

const handleAction = (notification: SystemNotification) => {
  markAsRead(notification.id)
  if (notification.link) {
    router.push(notification.link)
  }
}

const loadNotifications = () => {
  loading.value = true
  // Simulate API call
  setTimeout(() => {
    const start = (currentPage.value - 1) * pageSize
    const end = start + pageSize
    notifications.value = MOCK_NOTIFICATIONS.slice(start, end)
    loading.value = false
  }, 300)
}

onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
