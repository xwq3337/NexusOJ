<template>
  <nav class="sticky top-0 z-50 w-full backdrop-blur-md" :style="{
    backgroundColor: 'color-mix(in srgb, var(--surface-primary) 85%, transparent)',
    boxShadow: '0 1px 0 var(--border-color), 0 4px 20px rgba(14, 165, 233, 0.06)'
  }">
    <div class="container mx-auto px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <RouterLink to="/" class="flex items-center gap-2 font-bold text-xl">
          <div
            class="w-8 h-8 rounded-md flex items-center justify-center text-xs font-terminal font-bold neon-border"
            :style="{
              backgroundColor: 'var(--neon-cyan)',
              color: '#0a0e1a',
              boxShadow: '0 0 10px var(--neon-glow-cyan)'
            }">
            &gt;_
          </div>
          <span :style="{ color: 'var(--text-primary)' }">NexusOJ</span>
        </RouterLink>

        <div class="hidden md:flex items-center gap-1 whitespace-nowrap">
          <RouterLink to="/problems"
            class="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center gap-2"
            :style="isActiveStyle('/problems')">
            <Hash :size="16" /> 题目
          </RouterLink>
          <RouterLink to="/records"
            class="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center gap-2"
            :style="isActiveStyle('/records')">
            <FileText :size="16" /> 记录
          </RouterLink>
          <RouterLink to="/contests"
            class="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center gap-2"
            :style="isActiveStyle('/contests')">
            <Trophy :size="16" /> 竞赛
          </RouterLink>
          <RouterLink to="/courses"
            class="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center gap-2"
            :style="isActiveStyle('/courses')">
            <BookOpen :size="16" /> 课程
          </RouterLink>
          <RouterLink to="/blogs"
            class="px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 flex items-center gap-2"
            :style="isActiveStyle('/blogs')">
            <LibraryBig :size="16" /> 博客
          </RouterLink>
        </div>
      </div>


      <div class="hidden md:flex items-center gap-4">
        <!-- 主题切换按钮 -->
        <button @click="toggleTheme"
          class="flex items-center justify-center w-9 h-9 rounded-md transition-all duration-200 cursor-pointer"
          :style="{
            color: 'var(--text-secondary)',
            // backgroundColor: 'var(--surface-secondary)',
            // border: '1px solid var(--border-color)'
          }">
          <Sun v-if="currentTheme === 'dark'" :size="18" />
          <Moon v-else :size="18" />
        </button>
        <!-- 消息按钮 -->
        <n-dropdown trigger="hover" :options="messageOptions" placement="bottom-start" @select="handleMessageSelect">
          <n-badge class="hidden md:flex w-7 mr-2 h-10 rounded-full" :value="unRead" :max="99">
            <button
              class="flex items-center justify-center w-9 h-9 rounded-md transition-all duration-200 cursor-pointer"
              :style="{
                color: 'var(--text-secondary)',
                // backgroundColor: 'var(--surface-secondary)',
                // border: '1px solid var(--border-color)'
              }">
              <MessageSquareText :size="18" />
            </button>
          </n-badge>
        </n-dropdown>
        <!-- 登录/头像 -->
        <button v-if="isAuthorization"
          class="flex items-center justify-center w-9 h-9 rounded-full cursor-pointer transition-all duration-200"
          :style="{
            border: '2px solid var(--accent-color)',
            boxShadow: '0 0 8px var(--neon-glow-cyan)'
          }"
          @click="$router.push({ name: 'Profile' })">
          <n-avatar round :src="avatar" size="small" />
        </button>
        <button v-else
          class="whitespace-nowrap text-white px-5 py-2 rounded-md text-sm font-medium transition-all duration-200 cursor-pointer"
          :style="{
            backgroundColor: 'var(--btn-primary)',
            boxShadow: '0 0 12px var(--neon-glow-cyan)'
          }"
          @click="$router.push({ name: `Auth` })">
          登录
        </button>
      </div>
      <div class="hidden max-md:flex items-center gap-4 whitespace-nowrap">
        <n-dropdown placement="bottom-start" animated class="color-(--text-primary)" trigger="click" :options="options"
          @select="handleSelect">
          <n-button type="info" ghost round>
            <template #icon>
              <Menu :size="16" />
            </template>
          </n-button>
        </n-dropdown>
      </div>

    </div>
  </nav>
</template>

<script setup lang="ts">
import { useRoute, RouterLink } from 'vue-router'
import { NButton, NBadge, NAvatar, NDropdown } from 'naive-ui'
import {
  LibraryBig,
  MessageCircleMore,
  Trophy,
  Hash,
  Sun,
  Moon,
  FileText,
  Menu,
  MessageSquareText,
  BookOpen,
  Bell
} from 'lucide-vue-next'
import { useTheme } from '@/composables/useTheme'
import { useUserStore } from '@/stores/useUserStore'
import { h, inject, onMounted, ref, Ref, watch } from 'vue'
import { useNexusEventSource } from '@/composables/useEventSource'
import router from '@/router'
const { id, avatar } = useUserStore()
const route = useRoute()
const { theme: currentTheme, toggleTheme } = useTheme()
const isAuthorization: Ref<boolean> = inject('isAuthorization')!
const isActiveStyle = (path: string) => {
  if (route.path === path) {
    return {
      color: 'var(--accent-color)',
      backgroundColor: 'var(--surface-secondary)',
      boxShadow: '0 0 8px var(--neon-glow-cyan)'
    }
  }
  return {
    color: 'var(--text-secondary)'
  }
}
const unRead = ref(0)
type Response = {
  data: string
  event: "message" | "heartbeat",
  id: string
  retry: number | undefined
}

watch(isAuthorization, (auth) => {
  if (auth) {
    const { close, addListener } = useNexusEventSource(`/sse/chat/unread?id=${id}`, {
      onError(err) {
        console.log("unread err", err)
      },
      onOpen() {
        console.log('unread open')
      },
      onClose() {
        console.log('unread close')
      }
    })
    addListener('message', (msg : Response) => {
      unRead.value = Number(msg.data)
     })
    addListener('heartbeat', (msg : Response) => {
      console.log('unread heartbeat')
     })
  }
}, { immediate: true })


const options = [
  {
    label: '题库',
    key: 'Problems',
  },
  {
    label: '竞赛',
    key: 'Contests',
  },
  {
    label: '记录',
    key: 'Records',
  },
  {
    label: '课程',
    key: 'Courses'
  }, {
    label: '博客',
    key: 'Blogs'
  },
  {
    label: '消息',
    key: 'Messages'
  },
  {
    label: "个人主页",
    key: "Profile"
  }
]


const messageOptions = [
  {
    label: '私信',
    key: 'message',
    icon: () => h(MessageCircleMore, { size: 16 })
  },
  {
    label: '系统通知',
    key: 'info',
    icon: () => h(Bell, { size: 16 })
  }
]

function handleMessageSelect(key: string | number) {
  if (key === 'message') {
    router.push({ name: 'Messages' })
  } else if (key === 'info') {
    router.push({ name: 'Infos' })
  }
}
function handleSelect(key: string | number) {
  router.push({ name: String(key) })
}

</script>
