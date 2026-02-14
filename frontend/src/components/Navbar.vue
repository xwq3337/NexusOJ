<template>
  <nav class="sticky top-0 z-50 w-full border-b backdrop-blur" :style="{
    backgroundColor: 'var(--surface-primary)',
    borderColor: 'var(--border-color)',
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

        <div class="hidden md:flex items-center gap-1">
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
          <RouterLink to="/leaderboard"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/leaderboard')">
            <BarChart2 :size="16" /> 排名
          </RouterLink>
          <RouterLink to="/KnowledgeBase"
            class="px-4 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2"
            :class="isActive('/KnowledgeBase')">
            <LibraryBig :size="16" /> 知识库
          </RouterLink>
        </div>
      </div>

      <div class="flex items-center gap-4">
        <!-- 主题切换按钮 -->
        <button @click="toggleTheme" :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }">
          <Sun v-if="currentTheme === 'dark'" :size="20" />
          <Moon v-else :size="20" />
        </button>
        <n-badge :value="101" :max="99">
          <button :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
            class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
            :style="{ backgroundColor: 'var(--surface-primary)' }" @click="$router.push({ name: 'Message' })">
            <MessageCircleMore :size="20" />
          </button>
        </n-badge>

        <button v-if="account" :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          class="hidden md:flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }" @click="$router.push({ name: 'PersonalCenter' })">
          <User :size="20" />
        </button>
        <button v-else
          class="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2 rounded-md text-sm font-medium transition-colors"
          @click="$router.push({ name: `Login` })">
          登录
        </button>

        <!-- 移动端菜单按钮 -->
        <button @click="toggleMobileMenu"
          class="md:hidden flex items-center justify-center w-10 h-10 rounded-full text-gray-400 cursor-pointer transition-colors"
          :class="`hover:${currentTheme === 'dark' ? 'text-white' : 'text-black'}`"
          :style="{ backgroundColor: 'var(--surface-primary)' }">
          <Menu v-if="!isMobileMenuOpen" :size="24" />
          <X v-else :size="24" />
        </button>
      </div>
    </div>

    <!-- 移动端菜单 -->
    <Transition name="slide-down">
      <div v-if="isMobileMenuOpen" class="md:hidden border-t" :style="{
        borderColor: 'var(--border-color)',
        borderWidth: '1px',
        borderStyle: 'solid'
      }">
        <div class="px-4 py-3 space-y-2">
          <!-- 移动端主题切换 -->
          <button @click="toggleTheme"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-400 transition-colors"
            :class="`hover:${currentTheme === 'dark' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-black'}`">
            <Sun v-if="currentTheme === 'dark'" :size="20" />
            <Moon v-else :size="20" />
            <span class="text-sm font-medium">主题切换</span>
          </button>

          <!-- 移动端菜单项 -->
          <RouterLink to="/problems" @click="isMobileMenuOpen = false"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors" :class="route.path === '/problems'
              ? 'text-blue-400 bg-blue-500/10'
              : `text-gray-400 ${currentTheme === 'dark' ? 'hover:bg-gray-800 hover:text-white' : 'hover:bg-gray-100 hover:text-black'}`
              ">
            <Hash :size="18" />
            <span class="text-sm font-medium">题目</span>
          </RouterLink>

          <RouterLink to="/records" @click="isMobileMenuOpen = false"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors" :class="route.path === '/records'
              ? 'text-blue-400 bg-blue-500/10'
              : `text-gray-400 ${currentTheme === 'dark' ? 'hover:bg-gray-800 hover:text-white' : 'hover:bg-gray-100 hover:text-black'}`
              ">
            <FileText :size="18" />
            <span class="text-sm font-medium">记录</span>
          </RouterLink>

          <RouterLink to="/contests" @click="isMobileMenuOpen = false"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors" :class="route.path === '/contests'
              ? 'text-blue-400 bg-blue-500/10'
              : `text-gray-400 ${currentTheme === 'dark' ? 'hover:bg-gray-800 hover:text-white' : 'hover:bg-gray-100 hover:text-black'}`
              ">
            <Trophy :size="18" />
            <span class="text-sm font-medium">竞赛</span>
          </RouterLink>

          <RouterLink to="/leaderboard" @click="isMobileMenuOpen = false"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors" :class="route.path === '/leaderboard'
              ? 'text-blue-400 bg-blue-500/10'
              : `text-gray-400 ${currentTheme === 'dark' ? 'hover:bg-gray-800 hover:text-white' : 'hover:bg-gray-100 hover:text-black'}`
              ">
            <BarChart2 :size="18" />
            <span class="text-sm font-medium">排名</span>
          </RouterLink>

          <RouterLink to="/KnowledgeBase" @click="isMobileMenuOpen = false"
            class="w-full flex items-center gap-3 px-4 py-3 rounded-lg transition-colors" :class="route.path === '/KnowledgeBase'
              ? 'text-blue-400 bg-blue-500/10'
              : `text-gray-400 ${currentTheme === 'dark' ? 'hover:bg-gray-800 hover:text-white' : 'hover:bg-gray-100 hover:text-black'}`
              ">
            <LibraryBig :size="18" />
            <span class="text-sm font-medium">知识库</span>
          </RouterLink>

          <!-- 移动端消息和用户 -->
          <div class="border-t pt-2 mt-2" :style="{ borderColor: 'var(--border-light)' }">
            <n-badge :value="101" :max="99" :show="isMobileMenuOpen">
              <button @click="() => { isMobileMenuOpen = false; $router.push({ name: 'Message' }) }"
                class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-400 transition-colors"
                :class="`hover:${currentTheme === 'dark' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-black'}`">
                <MessageCircleMore :size="18" />
                <span class="text-sm font-medium">消息</span>
              </button>
            </n-badge>

            <button v-if="account" @click="() => { isMobileMenuOpen = false; $router.push({ name: 'PersonalCenter' }) }"
              class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-gray-400 transition-colors mt-2"
              :class="`hover:${currentTheme === 'dark' ? 'bg-gray-800 text-white' : 'bg-gray-100 text-black'}`">
              <User :size="18" />
              <span class="text-sm font-medium">个人中心</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </nav>
</template>

<script setup lang="ts">
import { ref } from 'vue'
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
  FileText,
  Menu,
  X
} from 'lucide-vue-next'
import { useLocalStorage } from '@vueuse/core'
import { useTheme } from '@/composables/useTheme'

const account = useLocalStorage('account', '')
const route = useRoute()
const { theme: currentTheme, toggleTheme } = useTheme()

const isMobileMenuOpen = ref(false)
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const isActive = (path: string) =>
  route.path === path
    ? 'text-blue-400'
    : `text-gray-400 hover:${currentTheme.value === 'dark' ? 'text-white' : 'text-black'}`
</script>

<style scoped>
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from,
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
