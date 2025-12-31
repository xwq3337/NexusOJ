<template>
  <nav
    class="sticky top-0 z-50 w-full border-b backdrop-blur"
    :style="{
      backgroundColor: 'var(--surface-primary)',
      borderColor: 'var(--border-color)',
      borderWidth: '1px',
      borderStyle: 'solid'
    }"
  >
    <div class="container mx-auto px-4 h-16 flex items-center justify-between">
      <div class="flex items-center gap-8">
        <RouterLink to="/" class="flex items-center gap-2 font-bold text-xl text-white">
          <div
            class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-mono"
          >
            &gt;_
          </div>
          <span :style="{ color: 'var(--text-primary)' }">NexusOJ</span>
        </RouterLink>

        <div class="hidden md:flex items-center gap-1">
          <RouterLink
            to="/problems"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/problems')"
          >
            <Hash :size="16" /> 题目
          </RouterLink>
          <RouterLink
            to="/records"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/records')"
          >
            <FileText :size="16" /> 记录
          </RouterLink>
          <RouterLink
            to="/contests"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/contests')"
          >
            <Trophy :size="16" /> 竞赛
          </RouterLink>
          <RouterLink
            to="/leaderboard"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/leaderboard')"
          >
            <BarChart2 :size="16" /> 排名
          </RouterLink>
          <RouterLink
            to="/KnowledgeBase"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/KnowledgeBase')"
          >
            <LibraryBig :size="16" /> 知识库
          </RouterLink>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <!-- 主题切换按钮 -->
        <button
          @click="toggleTheme"
          :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }"
        >
          <Sun v-if="currentTheme === 'dark'" :size="20" />
          <Moon v-else :size="20" />
        </button>
        <n-badge :value="101" :max="99">
          <button
            :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
            class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
            :style="{ backgroundColor: 'var(--surface-primary)' }"
            @click="$router.push({ name: 'Message' })"
          >
            <MessageCircleMore :size="20" />
          </button>
        </n-badge>

        <button
          v-if="account"
          :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }"
          @click="$router.push({ name: 'PersonalCenter' })"
        >
          <User :size="20" />
        </button>
        <button
          v-else
          class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-md text-sm font-medium transition-colors"
          @click="$router.push({ name: `Login` })"
        >
          登录
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { NBadge } from 'naive-ui'
import {
  LibraryBig,
  MessageCircleMore,
  User,
  Trophy,
  BarChart2,
  Hash,
  Sun,
  Moon,
  FileText
} from 'lucide-vue-next'
import { useLocalStorage } from '@vueuse/core'
const account = useLocalStorage('account', '')
const route = useRoute()
const isActive = (path: string) =>
  route.path === path
    ? 'text-blue-400'
    : `text-gray-400 hover:${currentTheme.value === 'dark' ? 'text-white' : 'text-black'}`
const currentTheme = ref<'light' | 'dark'>('light')
const toggleTheme = () => {
  const newTheme = currentTheme.value === 'light' ? 'dark' : 'light'
  document.body.classList.remove(currentTheme.value)
  document.body.classList.add(newTheme)
  document.documentElement.classList.remove(currentTheme.value === 'dark' ? 'dark' : 'light')
  document.documentElement.classList.add(newTheme === 'dark' ? 'dark' : 'light')
  localStorage.setItem('theme', newTheme)
  currentTheme.value = newTheme
}

onMounted(() => {
  const savedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  const theme = savedTheme || (prefersDark ? 'dark' : 'light')
  currentTheme.value = theme
  document.body.classList.remove('light', 'dark')
  document.body.classList.add(theme)
  document.documentElement.classList.remove('light', 'dark')
  document.documentElement.classList.add(theme === 'dark' ? 'dark' : 'light')
})
</script>
