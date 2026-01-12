<template>
  <div :style="{ backgroundColor: 'var(--bg-primary)' }" class="min-h-screen p-6">
    <div class="max-w-6xl mx-auto">
      <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">消息中心</h1>

      <n-tabs type="line" v-model:value="activeTab" animated>
        <n-tab-pane name="chat" tab="聊天消息">
          <div class="chat-container mt-4">
            <div class="flex gap-4">
              <!-- 聊天列表侧边栏 -->
              <div
                class="w-1/3 border rounded-lg"
                :style="{
                  backgroundColor: 'var(--surface-primary)',
                  borderColor: 'var(--border-color)'
                }"
              >
                <!-- 搜索区域 -->
                <div class="p-4 border-b" :style="{ borderColor: 'var(--border-color)' }">
                  <div class="flex gap-2 mb-3">
                    <n-input
                      v-model:value="searchKeyword"
                      placeholder="搜索用户..."
                      :style="{ backgroundColor: 'var(--surface-tertiary)' }"
                      @input="searchUsers"
                    >
                      <template #prefix>
                        <Search :size="18" />
                      </template>
                    </n-input>
                    <n-button @click="searchUsers" type="primary">
                      <Search :size="18" />
                    </n-button>
                  </div>

                  <!-- 搜索结果 -->
                  <div v-if="searchKeyword && searchResults.length > 0" class="mb-3">
                    <div
                      v-for="user in searchResults"
                      :key="user.id"
                      class="p-2 rounded flex items-center gap-2 hover:bg-gray-700/30"
                    >
                      <n-avatar round :size="32" :src="user.avatar || DEFAULT_AVATAR" />
                      <div class="flex-1">
                        <div class="font-medium" :style="{ color: 'var(--text-primary)' }">
                          {{ user.nickname || user.account }}
                        </div>
                        <div class="text-xs text-gray-500">@{{ user.account }}</div>
                      </div>
                      <n-button @click="addFriend(user.id)" size="small" type="primary">
                        <UserPlus :size="16" />
                        添加
                      </n-button>
                    </div>
                  </div>

                  <div class="flex justify-between">
                    <h2 class="text-lg font-semibold" :style="{ color: 'var(--text-primary)' }">
                      聊天列表
                    </h2>
                    <div class="flex gap-2">
                      <n-button @click="fetchFriendRequests" text type="primary" size="small"
                        >好友请求({{
                          friendRequests.filter((req) => req.status === 'pending').length
                        }})</n-button
                      >
                      <n-button @click="fetchFriendList" text size="small">刷新</n-button>
                    </div>
                  </div>
                </div>

                <div class="overflow-y-auto max-h-[60vh]">
                  <div
                    v-for="chat in chatList"
                    :key="chat.id"
                    class="p-4 border-b cursor-pointer hover:bg-gray-700/20 transition-colors flex items-center gap-3"
                    :style="{
                      borderColor: 'var(--border-color)',
                      backgroundColor:
                        activeChatId === chat.id ? 'var(--surface-secondary)' : 'transparent'
                    }"
                    @click="selectChat(chat.id)"
                  >
                    <n-avatar round :size="40" :src="chat.avatar" />
                    <div class="flex-1 min-w-0">
                      <div class="flex justify-between">
                        <h3 class="font-medium truncate" :style="{ color: 'var(--text-primary)' }">
                          {{ chat.name }}
                        </h3>
                        <span class="text-xs text-gray-500">{{ chat.time }}</span>
                      </div>
                      <p class="text-sm truncate text-gray-500">
                        {{ chat.lastMessage }}
                      </p>
                    </div>
                    <n-badge :value="chat.unread" v-if="chat.unread > 0" />
                  </div>
                </div>
              </div>

              <!-- 聊天内容区域 -->
              <div
                class="flex-1 flex flex-col border rounded-lg"
                :style="{
                  backgroundColor: 'var(--surface-primary)',
                  borderColor: 'var(--border-color)'
                }"
              >
                <div
                  v-if="activeChatId"
                  class="p-4 border-b"
                  :style="{ borderColor: 'var(--border-color)' }"
                >
                  <div class="flex items-center gap-3">
                    <n-avatar round :size="40" :src="activeChat?.avatar" />
                    <div>
                      <h3 class="font-semibold" :style="{ color: 'var(--text-primary)' }">
                        {{ activeChat?.name }}
                      </h3>
                      <p class="text-xs text-gray-500">
                        {{ activeChat?.status }}
                      </p>
                    </div>
                  </div>
                </div>

                <div v-if="activeChatId" class="flex-1 p-6 overflow-y-auto max-h-[45vh] bg-gradient-to-b from-slate-50 to-white dark:from-slate-900 dark:to-slate-800">
                  <div
                    v-for="message in activeChat?.messages"
                    :key="message.id"
                    :class="['mb-6 flex gap-3', message.isOwn ? 'flex-row-reverse' : 'flex-row']"
                  >
                    <!-- 头像 -->
                    <n-avatar
                      round
                      :size="40"
                      :src="message.isOwn ? userStore.avatar : activeChat?.avatar"
                      class="flex-shrink-0 shadow-md"
                    />

                    <!-- 消息气泡 -->
                    <div
                      :class="['flex flex-col max-w-[70%]', message.isOwn ? 'items-end' : 'items-start']"
                    >
                      <!-- 发送者名称 -->
                      <div
                        v-if="!message.isOwn"
                        class="text-xs font-medium text-slate-500 dark:text-slate-400 mb-1 px-2"
                      >
                        {{ activeChat?.name }}
                      </div>

                      <!-- 消息内容 -->
                      <div
                        :class="[
                          'relative px-4 py-3 shadow-md transition-all duration-200',
                          message.isOwn
                            ? 'bg-gradient-to-br from-blue-500 to-indigo-600 text-white rounded-2xl rounded-br-sm'
                            : 'bg-white dark:bg-slate-700 text-slate-800 dark:text-slate-100 rounded-2xl rounded-bl-sm border-2 border-slate-200 dark:border-slate-600'
                        ]"
                      >
                        <p class="text-sm leading-relaxed whitespace-pre-wrap break-words">
                          {{ message.text }}
                        </p>

                        <!-- 时间戳 -->
                        <div
                          :class="[
                            'text-xs mt-2 flex items-center gap-1',
                            message.isOwn ? 'text-blue-100' : 'text-slate-400 dark:text-slate-500'
                          ]"
                        >
                          <svg
                            class="w-3 h-3"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                          >
                            <circle cx="12" cy="12" r="10"></circle>
                            <polyline points="12 6 12 12 16 14"></polyline>
                          </svg>
                          {{ message.time }}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                <div
                  v-if="activeChatId"
                  class="p-4 border-t"
                  :style="{ borderColor: 'var(--border-color)' }"
                >
                  <div class="flex gap-2">
                    <n-input
                      v-model:value="messageText"
                      placeholder="输入消息..."
                      :style="{ backgroundColor: 'var(--surface-tertiary)' }"
                      @keydown.enter="sendMessage"
                    />
                    <n-button type="primary" @click="sendMessage">发送</n-button>
                  </div>
                </div>

                <div v-else class="flex-1 flex items-center justify-center text-gray-500">
                  选择一个聊天开始对话
                </div>
              </div>
            </div>
          </div>
        </n-tab-pane>

        <n-tab-pane name="notifications" tab="系统通知">
          <div class="notifications-container mt-4">
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-xl font-semibold" :style="{ color: 'var(--text-primary)' }">
                系统通知
              </h2>
              <n-button v-if="unreadNotifications.length > 0" @click="markAllAsRead" size="small"
                >全部标记为已读</n-button
              >
            </div>

            <div class="space-y-3">
              <!-- 好友请求 -->
              <div v-if="friendRequests.length > 0">
                <h3 class="text-lg font-medium mb-3" :style="{ color: 'var(--text-primary)' }">
                  好友请求
                </h3>
                <div
                  v-for="request in friendRequests.filter((req) => req.status === 'pending')"
                  :key="request.id"
                  class="p-4 rounded-lg border flex items-start gap-3"
                  :style="{
                    backgroundColor: 'var(--surface-secondary)',
                    borderColor: 'var(--border-color)'
                  }"
                >
                  <div
                    class="p-2 rounded-full text-blue-400"
                    :style="{ backgroundColor: 'var(--surface-tertiary)/50' }"
                  >
                    <UserPlus :size="20" />
                  </div>
                  <div class="flex-1">
                    <div class="flex justify-between">
                      <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                        {{ request.senderNickname || request.senderAccount }}
                      </h3>
                      <span class="text-xs text-gray-500">{{ request.createdAt }}</span>
                    </div>
                    <p class="text-gray-500 mt-1">{{ request.message }}</p>
                    <div class="flex gap-2 mt-2">
                      <n-button
                        text
                        size="small"
                        type="primary"
                        @click="handleFriendRequest(request.id, 'accepted')"
                      >
                        <Check :size="16" />
                        同意
                      </n-button>
                      <n-button
                        text
                        size="small"
                        @click="handleFriendRequest(request.id, 'rejected')"
                      >
                        <X :size="16" />
                        拒绝
                      </n-button>
                    </div>
                  </div>
                  <n-avatar :src="request.senderAvatar || DEFAULT_AVATAR" :size="40" round />
                </div>
              </div>

              <!-- 系统通知 -->
              <div
                v-for="notification in notifications"
                :key="notification.id"
                class="p-4 rounded-lg border flex items-start gap-3"
                :style="{
                  backgroundColor: notification.isRead
                    ? 'var(--surface-primary)'
                    : 'var(--surface-secondary)',
                  borderColor: 'var(--border-color)'
                }"
              >
                <div
                  :class="[
                    'p-2 rounded-full',
                    notification.isRead ? 'text-gray-400' : 'text-blue-400'
                  ]"
                  :style="{
                    backgroundColor: notification.isRead
                      ? 'var(--surface-tertiary)'
                      : 'var(--surface-tertiary)/50'
                  }"
                >
                  <component :is="getNotificationIcon(notification.type)" :size="20" />
                </div>
                <div class="flex-1">
                  <div class="flex justify-between">
                    <h3 class="font-medium" :style="{ color: 'var(--text-primary)' }">
                      {{ notification.title }}
                    </h3>
                    <span class="text-xs text-gray-500">{{ notification.time }}</span>
                  </div>
                  <p class="text-gray-500 mt-1">{{ notification.content }}</p>
                  <div class="flex gap-2 mt-2">
                    <n-button v-if="notification.action" text size="small" type="primary">{{
                      notification.action.text
                    }}</n-button>
                    <n-button
                      v-if="!notification.isRead"
                      text
                      size="small"
                      @click="markAsRead(notification.id)"
                      >标记为已读</n-button
                    >
                  </div>
                </div>
                <n-badge dot v-if="!notification.isRead" />
              </div>

              <div
                v-if="
                  notifications.length === 0 &&
                  friendRequests.filter((req) => req.status === 'pending').length === 0
                "
                class="text-center py-8 text-gray-500"
              >
                暂无通知
              </div>
            </div>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import Request from '@/services/api'
import { useUserStore } from '@/stores/useUserStore'
import {
  NTabs,
  NTabPane,
  NAvatar,
  NBadge,
  NInput,
  NButton,
  NModal,
  NInput as NInputComponent,
  NSpace,
  NCard
} from 'naive-ui'
import {
  MessageCircleMore,
  User as UserIcon,
  Settings as SettingsIcon,
  CheckCircle2,
  Calendar as CalendarIcon,
  Search,
  UserPlus,
  Check,
  X
} from 'lucide-vue-next'

// 当前激活的标签页
const activeTab = ref('chat')

// 监听标签页变化
watch(activeTab, (newTab) => {
  if (newTab === 'chat') {
    fetchFriendList()
  } else if (newTab === 'notifications') {
    fetchFriendRequests()
  }
})

// 用户存储
const userStore = useUserStore()

// 默认头像
const DEFAULT_AVATAR = 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'

// 好友列表
const friendList = ref([])

// 好友请求列表
const friendRequests = ref([])

// 搜索结果
const searchResults = ref([])

// 搜索关键词
const searchKeyword = ref('')

// 添加好友对话框
const showAddFriendModal = ref(false)

// 聊天消息相关
interface Chat {
  id: number
  name: string
  avatar: string
  lastMessage: string
  time: string
  unread: number
  status: string
  messages: Message[]
}

interface Friend {
  id: number
  account: string
  nickname: string
  avatar: string
  status: string
}

interface FriendRequest {
  id: number
  senderId: number
  senderAccount: string
  senderNickname: string
  senderAvatar: string
  message: string
  status: 'pending' | 'accepted' | 'rejected'
  createdAt: string
}

interface SearchResult {
  id: number
  account: string
  nickname: string
  avatar: string
}

interface Message {
  id: number
  text: string
  time: string
  isOwn: boolean
}

const messageText = ref('')
const activeChatId = ref<number | null>(null)

const chatList = ref<Chat[]>([])

const activeChat = computed(() => {
  if (!activeChatId.value) return null
  return chatList.value.find((chat) => chat.id === activeChatId.value) || null
})

const selectChat = async (id: number) => {
  activeChatId.value = id

  try {
    // 获取与该好友的聊天记录
    const response = await Request.get(`/chat/record?id=${id}`)
    if (response && (response.code === 200 || response.code === 0)) {
      // 更新当前聊天的消息列表
      const chat = chatList.value.find((c) => c.id === id)
      if (chat) {
        console.log(response.info)
        // 转换消息格式
        chat.messages =
          response.info?.map((msg: any) => ({
            isRead : msg.Status,
            text: msg.Content,
            time: msg.CreatedAt,
            isOwn: msg.SenderID === userStore.id // 如果发送者是当前用户，则为自己的消息
          })) || []
          console.log(chat.messages)
        // 更新最后消息显示
        if (response.info && response.info.length > 0) {
          const lastMessage = response.info[response.info.length - 1]
          chat.lastMessage = lastMessage.content
          chat.time = lastMessage.time || lastMessage.createdAt
        }
      }
    } else {
      console.error('获取聊天记录失败:', response?.message || '未知错误')
    }
  } catch (error) {
    console.error('获取聊天记录失败:', error)
  }
}

const sendMessage = async () => {
  if (!messageText.value.trim() || !activeChatId.value) return

  try {
    // 发送消息到后端
    const response = await Request.post('/chat/send', {
      recipientId: activeChatId.value,
      content: messageText.value
    })

    if (response && (response.code === 200 || response.code === 0)) {
      const chat = chatList.value.find((c) => c.id === activeChatId.value)
      if (chat) {
        const newMessage: Message = {
          id: Date.now(),
          text: messageText.value,
          time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
          isOwn: true
        }

        chat.messages.push(newMessage)
        chat.lastMessage = messageText.value
        chat.time = '刚刚'

        messageText.value = ''
      }
    }
  } catch (error) {
    console.error('发送消息失败:', error)
  }
}

// 通知相关
interface Notification {
  id: number
  title: string
  content: string
  time: string
  type: 'info' | 'warning' | 'success' | 'error'
  isRead: boolean
  action?: {
    text: string
    handler: () => void
  }
}

const notifications = ref<Notification[]>([
  {
    id: 1,
    title: '竞赛提醒',
    content: '您报名的算法竞赛将在30分钟后开始，请做好准备。',
    time: '10分钟前',
    type: 'info',
    isRead: false,
    action: {
      text: '查看详情',
      handler: () => console.log('查看竞赛详情')
    }
  },
  {
    id: 2,
    title: '系统更新',
    content: 'NexusOJ平台已完成系统更新，新增了AI助教功能。',
    time: '1小时前',
    type: 'success',
    isRead: false
  },
  {
    id: 3,
    title: '代码提交成功',
    content: '您的代码已成功提交，正在等待评测结果。',
    time: '2小时前',
    type: 'info',
    isRead: true
  },
  {
    id: 4,
    title: '新消息',
    content: '您收到了一条来自张三的消息，请及时查看。',
    time: '昨天',
    type: 'warning',
    isRead: true
  },
  {
    id: 5,
    title: '账户安全',
    content: '您的账户密码已更新，请确认是本人操作。',
    time: '2天前',
    type: 'error',
    isRead: true
  }
])

const unreadNotifications = computed(() => notifications.value.filter((n) => !n.isRead))

const markAsRead = (id: number) => {
  const notification = notifications.value.find((n) => n.id === id)
  if (notification) {
    notification.isRead = true
  }
}

const markAllAsRead = () => {
  notifications.value.forEach((notification) => {
    notification.isRead = true
  })
}

const getNotificationIcon = (type: string) => {
  switch (type) {
    case 'success':
      return CheckCircle2
    case 'warning':
      return SettingsIcon
    case 'error':
      return UserIcon
    case 'info':
    default:
      return MessageCircleMore
  }
}

// 获取好友列表
const fetchFriendList = async () => {
  try {
    const response = await Request.get('/user/friend')
    if (response && (response.code === 200 || response.code === 0)) {
      // 转换API响应为页面需要的格式
      friendList.value =
        response.info?.map((friend: Friend) => ({
          id: friend.id,
          name: friend.nickname || friend.account,
          avatar:  DEFAULT_AVATAR,
          lastMessage: '上次在线: ' + '1970年1月1日',
          time: '',
          unread: 0,
          status: friend.status,
          messages: []
        })) || []
      // 更新聊天列表
      chatList.value = friendList.value
    } else {
      console.error('获取好友列表失败:', response?.message || '未知错误')
    }
  } catch (error) {
    console.error('获取好友列表失败:', error)
  }
}

onMounted(() => {
  fetchFriendList()
  fetchFriendRequests()

  // 设置定时刷新（每30秒）
  const refreshInterval = setInterval(() => {
    if (activeTab.value === 'chat') {
      fetchFriendList()
    } else if (activeTab.value === 'notifications') {
      fetchFriendRequests()
    }
  }, 30000)

  // 组件卸载时清除定时器
  onUnmounted(() => {
    clearInterval(refreshInterval)
  })
})

// 搜索用户
const searchUsers = async () => {
  if (!searchKeyword.value.trim()) {
    searchResults.value = []
    return
  }

  try {
    const response = await Request.get(`/user/search?keyword=${searchKeyword.value}`)
    if (response && (response.code === 200 || response.code === 0)) {
      searchResults.value =
        response.info?.map((user: SearchResult) => ({
          id: user.id,
          account: user.account,
          nickname: user.nickname,
          avatar: user.avatar || DEFAULT_AVATAR
        })) || []
    } else {
      console.error('搜索用户失败:', response?.message || '未知错误')
      searchResults.value = []
    }
  } catch (error) {
    console.error('搜索用户失败:', error)
    searchResults.value = []
  }
}

// 添加好友
const addFriend = async (userId: number) => {
  try {
    const response = await Request.post('/user/add', {
      recipient: userId,
      verification: `你好，我是${userStore.nickname || userStore.account}，很高兴认识你！`
    })

    if (response && (response.code === 200 || response.code === 0)) {
      // 假设成功返回200或0
      // 可以显示成功消息
      console.log('好友请求已发送')
      // 关闭搜索结果
      searchResults.value = []
      searchKeyword.value = ''
    }
  } catch (error) {
    console.error('添加好友失败:', error)
  }
}

// 获取好友请求
const fetchFriendRequests = async () => {
  try {
    const response = await Request.get('/user/newfriend')
    if (response && (response.code === 200 || response.code === 0)) {
      friendRequests.value =
        response.info?.map((req: any) => ({
          id: req.id,
          senderId: req.senderId,
          senderAccount: req.senderAccount,
          senderNickname: req.senderNickname,
          senderAvatar: req.senderAvatar || DEFAULT_AVATAR,
          message: req.message,
          status: req.status,
          createdAt: req.createdAt
        })) || []
    } else {
      console.error('获取好友请求失败:', response?.message || '未知错误')
    }
  } catch (error) {
    console.error('获取好友请求失败:', error)
  }
}

// 处理好友请求
const handleFriendRequest = async (applyId: number, status: 'accepted' | 'rejected') => {
  try {
    const response = await Request.get(`/user/handle?apply_id=${applyId}&status=${status}`)
    if (response && (response.code === 200 || response.code === 0)) {
      // 假设成功返回200或0
      // 更新本地请求列表
      const request = friendRequests.value.find((req) => req.id === applyId)
      if (request) {
        request.status = status
      }

      console.log(status === 'accepted' ? '已同意好友请求' : '已拒绝好友请求')

      // 重新获取好友列表以更新
      await fetchFriendList()
    } else {
      console.error('处理好友请求失败:', response?.message || '未知错误')
    }
  } catch (error) {
    console.error('处理好友请求失败:', error)
  }
}
</script>
