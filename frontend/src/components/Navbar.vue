<template>
  <nav class="sticky top-0 z-50 w-full border-b backdrop-blur" :style="{
    backgroundColor: 'var(--surface-primary)',
    borderColor: 'transparent',
    borderWidth: '1px',
    borderStyle: 'solid'
  }">
    <div class="container mx-auto px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <RouterLink to="/" class="flex items-center gap-2 font-bold text-xl text-white">
          <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-mono">
            &gt;_
          </div>
          <span :style="{ color: 'var(--text-primary)' }">NexusOJ</span>
        </RouterLink>

        <div class="hidden md:flex items-center gap-1 whitespace-nowrap">
          <RouterLink to="/problems"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/problems')">
            <Hash :size="16" /> 题目
          </RouterLink>
          <RouterLink to="/records"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/records')">
            <FileText :size="16" /> 记录
          </RouterLink>
          <RouterLink to="/contests"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/contests')">
            <Trophy :size="16" /> 竞赛
          </RouterLink>
          <RouterLink to="/courses"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/courses')">
            <BookOpen :size="16" /> 课程
          </RouterLink>
          <RouterLink to="/blogs"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/blogs')">
            <LibraryBig :size="16" /> 博客
          </RouterLink>
        </div>
      </div>


      <div class="hidden md:flex items-center gap-4">
        <!-- 主题切换按钮 -->
        <button @click="toggleTheme" :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }">
          <Sun v-if="currentTheme === 'dark'" :size="20" />
          <Moon v-else :size="20" />
        </button>
        <!-- 消息按钮 -->
        <n-badge class="hidden md:flex" :value="unRead" :max="99">
          <button :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
            class="items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
            :style="{ backgroundColor: 'var(--surface-primary)' }" @click="$router.push({ name: 'Messages' })">
            <MessageCircleMore :size="20" />
          </button>
        </n-badge>
        <!-- 登录/头像 -->
        <button v-if="isAuthorization" :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }" @click="$router.push({ name: 'Profile' })">
          <n-avatar round :src="avatar" size="small" />
        </button>
        <button v-else
          class="whitespace-nowrap bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-md text-sm font-medium transition-colors"
          @click="$router.push({ name: `Login` })">
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
import { NButton, NBadge, NAvatar, NDropdown, useMessage } from 'naive-ui'
import {
  LibraryBig,
  MessageCircleMore,
  Trophy,
  BarChart2,
  Hash,
  Sun,
  Moon,
  FileText,
  Menu,
  BookOpen
} from 'lucide-vue-next'
import { useTheme } from '@/composables/useTheme'
import { useUserStore } from '@/stores/useUserStore'
import { h, inject, onMounted, ref, Ref } from 'vue'
import { useNexusEventSource } from '@/composables/useEventSource'
import router from '@/router'
const { id, avatar } = useUserStore()
const route = useRoute()
const { theme: currentTheme, toggleTheme } = useTheme()
const isAuthorization: Ref<boolean> = inject('isAuthorization')!
const isActive = (path: string) =>
  route.path === path
    ? 'text-blue-400'
    : `text-gray-400 hover:${currentTheme.value === 'dark' ? 'text-white' : 'text-black'}`
const unRead = ref(0)
type Response = {
  data: string
  event: "message" | "heartbeat",
  id: string
  retry: number | undefined
}

onMounted(() => {
  if (isAuthorization.value) {
    const { close } = useNexusEventSource(`/service/chat/unread?id=${id}`, {
      onMessage: (msg: Response) => {
        console.log("message: ", msg)
        if (msg.event === "heartbeat") return
        unRead.value = Number(msg.data)
      },
      onError(err) {
        console.log("err", err)
      },
      onOpen() {
        console.log('open')
      },
      onClose() {
        console.log('close')
      }
    })
  }
})

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

function handleSelect(key: string | number) {
  router.push({ name: String(key) })
}

</script>