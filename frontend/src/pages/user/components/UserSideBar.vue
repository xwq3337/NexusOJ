<template>
  <!-- Sticky Sidebar Container -->
  <div class="lg:sticky lg:top-2">
    <!-- Profile Card -->
    <div class="rounded-xl p-1.5" :style="{
      backgroundColor: 'var(--surface-primary)',
      borderColor: 'var(--border-color)',
      borderWidth: '1px',
      borderStyle: 'solid'
    }">
      <!-- Avatar Section -->
      <div class="flex flex-col items-center text-center mb-6">
        <img :src="user.avatar || 'https://api.dicebear.com/7.x/avataaars/svg?seed=' + user.username"
          :alt="user.nickname" class="w-20 h-20 rounded-full border-4 mb-3"
          :style="{ borderColor: 'var(--border-color)' }" />

        <h1 class="text-xl font-bold mb-1" :style="{ color: 'var(--text-primary)' }">
          {{ user.nickname }}
          <span v-if="user.user_role === 'admin'" class="ml-1 px-2 py-0.5 rounded text-xs bg-red-500/10 text-red-400">
            Admin
          </span>
        </h1>

        <p class="text-sm mb-2" :style="{ color: 'var(--text-secondary)' }">
          @{{ user.username }}
        </p>

        <!-- Rating Badge -->
        <div class="px-3 py-1 rounded-full text-sm font-medium mb-3" :style="{
          backgroundColor: getRatingBgColor(user.rating),
          color: getRatingColor(user.rating)
        }">
          {{ getRatingTitle(user.rating) }} · {{ user.rating }}
        </div>

        <!-- Bio -->
        <p v-if="user.introduction" class="text-sm mb-4 line-clamp-2" :style="{ color: 'var(--text-secondary)' }">
          {{ user.introduction }}
        </p>

        <!-- Action Buttons -->
        <div v-if="isOwnProfile" class="flex w-full gap-2">
          <button class="flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors" :style="{
            backgroundColor: 'var(--btn-primary)',
            color: '#ffffff'
          }" @click="handleEditProfile">
            编辑
          </button>
          <button class="flex-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors" :style="{
            backgroundColor: 'var(--surface-secondary)',
            color: 'var(--text-primary)',
            borderColor: 'var(--border-color)',
            borderWidth: '1px',
            borderStyle: 'solid'
          }" @click="handleShare">
            分享
          </button>
        </div>
      </div>

      <!-- Stats -->
      <div class="grid grid-cols-2 gap-3 mb-6 pb-6" :style="{ borderBottom: '1px solid var(--border-color)' }">
        <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
          <p class="text-lg font-bold text-blue-400">#{{ rank }}</p>
          <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">排名</p>
        </div>
        <div class="text-center p-3 rounded-lg" :style="{ backgroundColor: 'var(--surface-secondary)' }">
          <p class="text-lg font-bold text-green-400">{{ user.accept }}</p>
          <p class="text-xs" :style="{ color: 'var(--text-tertiary)' }">通过</p>
        </div>
      </div>

      <!-- Info -->
      <div class="space-y-3 text-sm">
        <div v-if="user.school" class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
          <component :is="GraduationCap" :size="16" class="shrink-0" />
          <span class="truncate">{{ user.school }}</span>
        </div>

        <div class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
          <component :is="Calendar" :size="16" class="shrink-0" />
          <span>加入于 {{ formatDate(user.created_at) }}</span>
        </div>

        <div v-if="user.codeforces" class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
          <component :is="ExternalLink" :size="16" class="shrink-0" />
          <a :href="`https://codeforces.com/profile/${user.codeforces}`" target="_blank"
            class="truncate hover:text-blue-400 transition-colors">
            {{ user.codeforces }}
          </a>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GraduationCap, Calendar, ExternalLink } from 'lucide-vue-next'
import type { User } from '@/types/user'
import { userApi } from '@/services/user';
import { formatDate } from '@/utils/format';
const props = defineProps<{
  userId?: string
}>()

const router = useRouter()
const user = ref<User>({
  id: '',
  username: '',
  avatar: '',
  introduction: '',
  codeforces: '',
  school: '',
  accept: 0,
  rating: 0,
  submission: 0,
  status: 0,
  email: '',
  nickname: '',
  gender: '0',
  user_role: '',
  birthday: '',
  banned_to: '',
  balance: 0,
  created_at: '',
  updated_at: ""
})
onMounted(async () => {
  const response = await userApi.getInfoById(props.userId)
  user.value = response.info
  // TODO:添加错误处理
})


const rank = ref(42)
import { useUserStore } from '@/stores/useUserStore';
const { id } = useUserStore()
// 判断是否为自己的主页
const isOwnProfile = computed(() => {
  const currentUserId = id
  return props.userId === currentUserId || !props.userId
})

// Rating 颜色
const getRatingColor = (rating: number) => {
  if (rating >= 2400) return '#ef4444'
  if (rating >= 2100) return '#f97316'
  if (rating >= 1900) return '#eab308'
  if (rating >= 1700) return '#8b5cf6'
  if (rating >= 1500) return '#3b82f6'
  if (rating >= 1200) return '#10b981'
  return '#6b7280'
}

// Rating 背景颜色
const getRatingBgColor = (rating: number) => {
  if (rating >= 2400) return 'rgba(239, 68, 68, 0.1)'
  if (rating >= 2100) return 'rgba(249, 115, 22, 0.1)'
  if (rating >= 1900) return 'rgba(234, 179, 8, 0.1)'
  if (rating >= 1700) return 'rgba(139, 92, 246, 0.1)'
  if (rating >= 1500) return 'rgba(59, 130, 246, 0.1)'
  if (rating >= 1200) return 'rgba(16, 185, 129, 0.1)'
  return 'rgba(107, 114, 128, 0.1)'
}

// Rating 称号
const getRatingTitle = (rating: number) => {
  if (rating >= 2400) return 'Grandmaster'
  if (rating >= 2100) return 'Master'
  if (rating >= 1900) return 'Candidate Master'
  if (rating >= 1700) return 'Expert'
  if (rating >= 1500) return 'Specialist'
  if (rating >= 1200) return 'Pupil'
  return 'Newbie'
}

// 编辑资料
const handleEditProfile = () => {
  router.push({ name: 'Settings' })
}

// 分享主页
const handleShare = () => {
  const url = window.location.href
  navigator.clipboard.writeText(url)
  // TODO: 显示提示消息
}

// 从 API 获取用户信息
const fetchUserInfo = async () => {
  try {
    // TODO: 从 API 获取实际数据
    // const response = await Request.get(`/api/user/${props.userId}`)
    // user.value = response.data
    // rank.value = response.rank
    // rating.value = response.rating
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>
