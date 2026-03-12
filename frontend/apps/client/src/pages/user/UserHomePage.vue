<template>
  <div class="min-h-screen">
    <!-- Hero Header Section -->
    <div class="relative overflow-hidden mb-6" :style="{
      background: 'linear-gradient(135deg, var(--hero-bg-start) 0%, var(--hero-bg-end) 100%)',
      padding: '2.5rem'
    }">
      <!-- Decorative Elements -->
      <div class="absolute top-0 right-0 w-64 h-64 rounded-full opacity-10" :style="{
        background: 'radial-gradient(circle, var(--hero-accent) 0%, transparent 70%)',
        transform: 'translate(30%, -30%)'
      }"></div>
      <div class="absolute bottom-0 left-0 w-48 h-48 rounded-full opacity-10" :style="{
        background: 'radial-gradient(circle, var(--hero-accent) 0%, transparent 70%)',
        transform: 'translate(-30%, 30%)'
      }"></div>

      <!-- User Info -->
      <div class="relative z-10 flex flex-col md:flex-row items-center md:items-start gap-6">
        <!-- Avatar -->
        <div class="relative">
          <div class="w-28 h-28 rounded-full p-1" :style="{
            background: 'linear-gradient(135deg, var(--hero-accent) 0%, var(--hero-accent-secondary) 100%)'
          }">
            <img :src="user.avatar || defaultAvatar" :alt="user.nickname"
              class="w-full h-full rounded-full object-cover" :style="{ backgroundColor: 'var(--surface-primary)' }" />
          </div>
        </div>

        <!-- User Details -->
        <div class="flex-1 text-center md:text-left">
          <div class="flex flex-col md:flex-row md:items-center gap-3 mb-3">
            <h1 class="text-3xl font-bold" :style="{ color: 'var(--hero-title-color)' }">
              {{ user.nickname }}
            </h1>
            <div v-if="user.user_role === 'admin'"
              class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium" :style="{
                backgroundColor: 'rgba(239, 68, 68, 0.2)',
                color: '#ef4444'
              }">
              <Shield :size="12" />
              Admin
            </div>
          </div>

          <p class="text-base mb-3" :style="{ color: 'var(--hero-subtitle-color)' }">
            @{{ user.username }}
          </p>

          <!-- Rating Badge -->
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-xl mb-3" :style="{
            backgroundColor: formatRating(user.rating).bgColor,
            backdropFilter: 'blur(10px)'
          }">
            <Trophy :size="18" :style="{ color: formatRating(user.rating).color }" />
            <span class="font-semibold" :style="{ color: formatRating(user.rating).color }">
              {{ formatRating(user.rating).title }}
            </span>
            <span class="text-sm" :style="{ color: formatRating(user.rating).color, opacity: 0.8 }">
              · {{ user.rating }}
            </span>
          </div>

          <!-- Bio -->
          <p v-if="user.introduction" class="text-sm max-w-2xl mb-4 line-clamp-2" :style="{
            color: 'var(--hero-text-color)'
          }">
            {{ user.introduction }}
          </p>

          <!-- Quick Stats -->
          <div class="flex flex-wrap items-center justify-center md:justify-start gap-6 text-sm">
            <div class="flex items-center gap-2" :style="{ color: 'var(--hero-text-color)' }">
              <Award :size="16" />
              <span>排名 #{{ rank }}</span>
            </div>
            <div class="flex items-center gap-2" :style="{ color: 'var(--hero-text-color)' }">
              <CheckCircle :size="16" />
              <span>{{ user.accept }} 道题通过</span>
            </div>
            <div class="flex items-center gap-2" :style="{ color: 'var(--hero-text-color)' }">
              <Calendar :size="16" />
              <span>加入于 {{ formatDate(user.created_at) }}</span>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col gap-2">
          <button v-if="isOwnProfile"
            class="px-6 py-2.5 rounded-xl font-medium transition-all duration-200 hover:scale-105" :style="{
              backgroundColor: 'rgba(255, 255, 255, 0.2)',
              color: 'var(--hero-title-color)',
              backdropFilter: 'blur(10px)',
              border: '1px solid rgba(255, 255, 255, 0.3)'
            }" @click="handleEditProfile">
            <span class="flex items-center justify-center gap-2">
              <Edit2 :size="16" />
              编辑资料
            </span>
          </button>
          <button v-else class="px-6 py-2.5 rounded-xl font-medium transition-all duration-200 hover:scale-105" :style="{
            backgroundColor: 'var(--hero-accent)',
            color: '#ffffff'
          }" @click="handleAddFriend">
            <span class="flex items-center justify-center gap-2">
              <UserPlus :size="16" />
              添加好友
            </span>
          </button>
          <button class="px-6 py-2.5 rounded-xl font-medium transition-all duration-200 hover:scale-105" :style="{
            backgroundColor: 'rgba(255, 255, 255, 0.1)',
            color: 'var(--hero-title-color)',
            backdropFilter: 'blur(10px)'
          }" @click="handleShare">
            <span class="flex items-center justify-center gap-2">
              <Share2 :size="16" />
              分享主页
            </span>
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left Sidebar - User Info Cards -->
      <div class="lg:col-span-3 space-y-4">
        <!-- Personal Info Card -->
        <div class="rounded-xl p-5" :style="{
          backgroundColor: 'var(--surface-primary)',
          border: '1px solid var(--border-color)'
        }">
          <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
            <UserRound :size="16" />
            个人信息
          </h3>
          <div class="space-y-3 text-sm">
            <div v-if="user.school" class="flex items-center gap-2 p-2 rounded-lg" :style="{
              backgroundColor: 'var(--surface-secondary)'
            }">
              <GraduationCap :size="16" :style="{ color: 'var(--text-tertiary)' }" />
              <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.school }}</span>
            </div>
            <div v-if="user.birthday" class="flex items-center gap-2 p-2 rounded-lg" :style="{
              backgroundColor: 'var(--surface-secondary)'
            }">
              <Cake :size="16" :style="{ color: 'var(--text-tertiary)' }" />
              <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.birthday }}</span>
            </div>
          </div>
        </div>
        <!-- Contact Info Card -->
        <div class="rounded-xl p-5" :style="{
          backgroundColor: 'var(--surface-primary)',
          border: '1px solid var(--border-color)'
        }">
          <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
            <Info :size="16" />
            联系方式
          </h3>
          <div class="space-y-3 text-sm">
            <div v-if="user.email" class="flex items-center gap-2 p-2 rounded-lg" :style="{
              backgroundColor: 'var(--surface-secondary)'
            }">
              <Mail :size="16" :style="{ color: 'var(--text-tertiary)' }" />
              <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.email }}</span>
            </div>
            <div v-if="user.codeforces" class="flex items-center gap-2 p-2 rounded-lg" :style="{
              backgroundColor: 'var(--surface-secondary)'
            }">
              <ExternalLink :size="16" :style="{ color: 'var(--text-tertiary)' }" />
              <a :href="`https://codeforces.com/profile/${user.codeforces}`" target="_blank"
                class="truncate hover:text-blue-400 transition-colors" :style="{ color: 'var(--text-secondary)' }">
                {{ user.codeforces }}
              </a>
            </div>
          </div>
        </div>

        <!-- Badges Card -->
        <div class="rounded-xl p-5" :style="{
          backgroundColor: 'var(--surface-primary)',
          border: '1px solid var(--border-color)'
        }">
          <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
            <Medal :size="16" />
            成就徽章
          </h3>
          <div class="flex flex-wrap gap-2">
            <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
              backgroundColor: 'rgba(16, 185, 129, 0.1)',
              border: '1px solid rgba(16, 185, 129, 0.3)'
            }" title="首次解题">
              <span class="text-xl">🎯</span>
            </div>
            <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
              backgroundColor: 'rgba(59, 130, 246, 0.1)',
              border: '1px solid rgba(59, 130, 246, 0.3)'
            }" title="连续7天">
              <span class="text-xl">🔥</span>
            </div>
            <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
              backgroundColor: 'rgba(249, 115, 22, 0.1)',
              border: '1px solid rgba(249, 115, 22, 0.3)'
            }" title="百题达成">
              <span class="text-xl">💯</span>
            </div>
            <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
              backgroundColor: 'rgba(139, 92, 246, 0.1)',
              border: '1px solid rgba(139, 92, 246, 0.3)'
            }" title="活跃用户">
              <span class="text-xl">⭐</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="lg:col-span-9 space-y-6">
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Solved Problems -->
          <div class="rounded-xl p-5 relative overflow-hidden group" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <div
              class="absolute top-0 right-0 w-20 h-20 rounded-full opacity-5 transition-transform group-hover:scale-150"
              :style="{
                backgroundColor: '#10b981',
                transform: 'translate(30%, -30%)'
              }"></div>
            <div class="relative z-10">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">通过题目</span>
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                  backgroundColor: 'rgba(16, 185, 129, 0.1)'
                }">
                  <CheckCircle :size="18" :style="{ color: '#10b981' }" />
                </div>
              </div>
              <div class="flex items-end gap-2">
                <span class="text-3xl font-bold" :style="{ color: 'var(--text-primary)' }">
                  {{ user.accept }}
                </span>
                <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">题</span>
              </div>
            </div>
          </div>

          <!-- Submissions -->
          <div class="rounded-xl p-5 relative overflow-hidden group" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <div
              class="absolute top-0 right-0 w-20 h-20 rounded-full opacity-5 transition-transform group-hover:scale-150"
              :style="{
                backgroundColor: '#3b82f6',
                transform: 'translate(30%, -30%)'
              }"></div>
            <div class="relative z-10">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">总提交</span>
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                  backgroundColor: 'rgba(59, 130, 246, 0.1)'
                }">
                  <Send :size="18" :style="{ color: '#3b82f6' }" />
                </div>
              </div>
              <div class="flex items-end gap-2">
                <span class="text-3xl font-bold" :style="{ color: 'var(--text-primary)' }">
                  {{ user.submission }}
                </span>
                <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">次</span>
              </div>
            </div>
          </div>

          <!-- Acceptance Rate -->
          <div class="rounded-xl p-5 relative overflow-hidden group" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <div
              class="absolute top-0 right-0 w-20 h-20 rounded-full opacity-5 transition-transform group-hover:scale-150"
              :style="{
                backgroundColor: '#f59e0b',
                transform: 'translate(30%, -30%)'
              }"></div>
            <div class="relative z-10">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">通过率</span>
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                  backgroundColor: 'rgba(245, 158, 11, 0.1)'
                }">
                  <TrendingUp :size="18" :style="{ color: '#f59e0b' }" />
                </div>
              </div>
              <div class="flex items-end gap-2">
                <span class="text-3xl font-bold" :style="{ color: 'var(--text-primary)' }">
                  {{ formatAcceptance(user.accept, user.submission) }}
                </span>
                <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">%</span>
              </div>
            </div>
          </div>

          <!-- Rating -->
          <div class="rounded-xl p-5 relative overflow-hidden group" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <div
              class="absolute top-0 right-0 w-20 h-20 rounded-full opacity-5 transition-transform group-hover:scale-150"
              :style="{
                backgroundColor: formatRating(user.rating).bgColor,
                transform: 'translate(30%, -30%)'
              }"></div>
            <div class="relative z-10">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">Rating</span>
                <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                  backgroundColor: formatRating(user.rating).bgColor
                }">
                  <Trophy :size="18" :style="{ color: formatRating(user.rating).color }" />
                </div>
              </div>
              <div class="flex items-end gap-2">
                <span class="text-3xl font-bold" :style="{ color: formatRating(user.rating).color }">
                  {{ user.rating }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Activity Heatmap -->
        <div class="rounded-xl p-5" :style="{
          backgroundColor: 'var(--surface-primary)',
          border: '1px solid var(--border-color)'
        }">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-sm font-semibold flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <Flame :size="16" />
              活动热图
            </h3>
            <n-radio-group v-model:value="heatmapPeriod" size="small">
              <n-radio-button value="recent">最近一年</n-radio-button>
              <n-radio-button value="2025">2025</n-radio-button>
              <n-radio-button value="2024">2024</n-radio-button>
            </n-radio-group>
          </div>
          <UserHeatmap :user-id="userId" :period="heatmapPeriod" />
        </div>

        <!-- Tabs Section -->
        <div class="rounded-xl overflow-hidden" :style="{
          backgroundColor: 'var(--surface-primary)',
          border: '1px solid var(--border-color)'
        }">
          <n-tabs v-model:value="activeTab" type="line" animated>
            <n-tab-pane name="submissions" tab="提交记录">
              <div class="p-6">
                <UserSubmissions :user-id="userId" />
              </div>
            </n-tab-pane>

            <n-tab-pane name="solutions" tab="题解">
              <div class="p-6">
                <UserSolutions :user-id="userId" />
              </div>
            </n-tab-pane>

            <n-tab-pane name="blogs" tab="博客">
              <div class="p-6">
                <UserBlogs :user-id="userId" />
              </div>
            </n-tab-pane>

            <n-tab-pane name="discussions" tab="讨论">
              <div class="p-6">
                <UserDiscussions :user-id="userId" />
              </div>
            </n-tab-pane>
          </n-tabs>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NTabs, NTabPane, NRadioGroup, NRadioButton, useMessage } from 'naive-ui'
import { useClipboard } from '@vueuse/core'
import {
  Shield,
  Trophy,
  Award,
  CheckCircle,
  Calendar,
  Edit2,
  Share2,
  UserPlus,
  Info,
  GraduationCap,
  Mail,
  ExternalLink,
  Medal,
  Send,
  UserRound,
  TrendingUp,
  Flame,
  Cake
} from 'lucide-vue-next'
import { useUserStore } from '@/stores/useUserStore'
import { userApi } from '@nexusoj/server'
import { formatAcceptance, formatDate, formatRating } from '@/utils/format'
import type { User } from '@nexusoj/type'
import UserHeatmap from './components/UserHeatmap.vue'
import UserSubmissions from './components/UserSubmissions.vue'
import UserSolutions from './components/UserSolutions.vue'
import UserBlogs from './components/UserBlogs.vue'
import UserDiscussions from './components/UserDiscussions.vue'
const { copy } = useClipboard()
const router = useRouter()
const route = useRoute()
const { id: currentUserId } = useUserStore()

const userId = ref<string>(route.params.id as string)
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
  updated_at: ''
})

const rank = ref(42)
const activeTab = ref('submissions')
const heatmapPeriod = ref<'recent' | '2025' | '2024'>('recent')

const defaultAvatar = computed(() =>
  `https://api.dicebear.com/7.x/avataaars/svg?seed=${user.value.username}`
)

const isOwnProfile = computed(() => {
  return userId.value === currentUserId || !userId.value
})

const handleEditProfile = () => {
  router.push({ name: 'Settings' })
}
const message = useMessage()
const handleShare = async () => {
    copy(route.fullPath)
    message.success("网页地址已复制")
}

const handleAddFriend = () => {
  // TODO: Implement add friend functionality
  console.log('Add friend:', userId.value)
}

const fetchUserInfo = async () => {
  try {
    const response = await userApi.getInfoById(userId.value)
    user.value = response.info
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>
