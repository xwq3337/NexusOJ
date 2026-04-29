<template>
  <div class="min-h-screen" :style="{ backgroundColor: 'var(--bg-primary)' }">
    <CyberGridCanvas />
    <div class="relative" style="z-index: 1">
      <!-- Hero Header Section -->
      <div class="overflow-hidden rounded-2xl mb-8 cyber-scanline-overlay" :style="{
        background: 'linear-gradient(135deg, rgba(14, 165, 233, 0.18) 0%, rgba(56, 189, 248, 0.14) 50%, rgba(14, 165, 233, 0.08) 100%)',
        padding: '2.5rem',
        border: '1px solid rgba(14, 165, 233, 0.2)',
        boxShadow: '0 0 40px rgba(14, 165, 233, 0.08), inset 0 1px 0 rgba(14, 165, 233, 0.1)'
      }">
        <!-- Grid pattern -->
        <div class="absolute inset-0 cyber-grid-bg opacity-40"></div>
        <!-- Glow orbs -->
        <div class="absolute top-0 right-0 w-80 h-80 rounded-full blur-3xl opacity-20" :style="{
          background: 'radial-gradient(circle, var(--neon-cyan) 0%, transparent 70%)',
          transform: 'translate(20%, -30%)'
        }"></div>
        <div class="absolute bottom-0 left-0 w-64 h-64 rounded-full blur-3xl opacity-15" :style="{
          background: 'radial-gradient(circle, var(--neon-magenta) 0%, transparent 70%)',
          transform: 'translate(-20%, 30%)'
        }"></div>

        <!-- User Info -->
        <div class="relative z-10 flex flex-col md:flex-row items-center md:items-start gap-6">
          <!-- Avatar with neon ring -->
          <div class="relative">
            <div class="w-28 h-28 rounded-full p-1" :style="{
              background: 'linear-gradient(135deg, #38bdf8 0%, #0ea5e9 50%, #0284c7 100%)',
              boxShadow: '0 0 20px var(--neon-glow-cyan), 0 0 40px rgba(14, 165, 233, 0.15)'
            }">
              <img :src="user.avatar || defaultAvatar" :alt="user.nickname"
                class="w-full h-full rounded-full object-cover"
                :style="{ backgroundColor: 'var(--surface-primary)' }" />
            </div>
          </div>

          <!-- User Details -->
          <div class="flex-1 text-center md:text-left">
            <div class="flex flex-col md:flex-row md:items-center gap-3 mb-3">
              <h1 class="text-3xl font-bold" :style="{
                color: 'var(--text-primary)',
                textShadow: '0 0 20px var(--neon-glow-cyan)'
              }">
                {{ user.nickname }}
              </h1>
              <div v-if="user.user_role === 'admin'"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium" :style="{
                  backgroundColor: 'rgba(239, 68, 68, 0.15)',
                  color: '#f87171',
                  border: '1px solid rgba(239, 68, 68, 0.25)'
                }">
                <Shield :size="12" />
                Admin
              </div>
              <div v-else-if="user.user_role === 'super_admin'"
                class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium" :style="{
                  backgroundColor: 'rgba(239, 68, 68, 0.15)',
                  color: '#f87171',
                  border: '1px solid rgba(239, 68, 68, 0.25)'
                }">
                <Shield :size="12" />
                SuperAdmin
              </div>
            </div>

            <p class="text-base mb-3 font-terminal" :style="{ color: 'var(--accent-color)' }">
              @{{ user.username }}
            </p>

            <!-- Rating Badge -->
            <div class="inline-flex items-center gap-2 px-4 py-2 rounded-xl mb-3" :style="{
              backgroundColor: formatRating(user.rating).bgColor,
              border: `1px solid ${formatRating(user.rating).color}40`
            }">
              <Trophy :size="18" :style="{ color: formatRating(user.rating).color }" />
              <span class="font-semibold" :style="{ color: formatRating(user.rating).color }">
                {{ formatRating(user.rating).title }}
              </span>
              <span class="text-sm font-terminal" :style="{ color: formatRating(user.rating).color, opacity: 0.8 }">
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
              <div class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
                <Award :size="16" :style="{ color: 'var(--accent-color)' }" />
                <span>排名 <span class="font-terminal" :style="{ color: 'var(--text-primary)' }">#{{ rank }}</span></span>
              </div>
              <div class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
                <CheckCircle :size="16" :style="{ color: 'var(--success-color)' }" />
                <span><span class="font-terminal" :style="{ color: 'var(--text-primary)' }"> {{ user.solved }}</span>
                  道题通过</span>
              </div>
              <div class="flex items-center gap-2" :style="{ color: 'var(--text-secondary)' }">
                <Calendar :size="16" :style="{ color: 'var(--text-tertiary)' }" />
                <span>加入于 {{ formatDate(user.created_at) }}</span>
              </div>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex flex-col gap-2">
            <button v-if="isOwnProfile"
              class="px-5 py-2.5 rounded-lg font-medium transition-all duration-200 hover:scale-105 text-sm" :style="{
                backgroundColor: 'var(--surface-tertiary)',
                color: 'var(--text-primary)',
                border: '1px solid var(--border-color)'
              }" @click="handleEditProfile">
              <span class="flex items-center justify-center gap-2">
                <Edit2 :size="15" />
                编辑资料
              </span>
            </button>
            <button v-else
              class="px-5 py-2.5 rounded-lg font-medium transition-all duration-200 hover:scale-105 text-sm" :style="{
                backgroundColor: 'var(--btn-primary)',
                color: '#ffffff',
                boxShadow: '0 0 12px var(--neon-glow-cyan)'
              }" @click="handleAddFriend">
              <span class="flex items-center justify-center gap-2">
                <UserPlus :size="15" />
                添加好友
              </span>
            </button>
            <button class="px-5 py-2.5 rounded-lg font-medium transition-all duration-200 hover:scale-105 text-sm"
              :style="{
                backgroundColor: 'var(--surface-tertiary)',
                color: 'var(--text-primary)',
                border: '1px solid var(--border-color)'
              }" @click="handleShare">
              <span class="flex items-center justify-center gap-2">
                <Share2 :size="15" />
                分享主页
              </span>
            </button>
          </div>
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Left Sidebar -->
        <div class="lg:col-span-3 space-y-4">
          <!-- Personal Info Card -->
          <div class="rounded-xl p-5 cyber-glow-card" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <UserRound :size="16" :style="{ color: 'var(--accent-color)' }" />
              个人信息
            </h3>
            <div class="space-y-3 text-sm">
              <div v-if="user.school" class="flex items-center gap-2 p-2.5 rounded-lg" :style="{
                backgroundColor: 'var(--surface-secondary)'
              }">
                <GraduationCap :size="15" :style="{ color: 'var(--accent-color)' }" />
                <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.school }}</span>
              </div>
              <div v-if="user.birthday" class="flex items-center gap-2 p-2.5 rounded-lg" :style="{
                backgroundColor: 'var(--surface-secondary)'
              }">
                <Cake :size="15" :style="{ color: 'var(--accent-color)' }" />
                <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.birthday }}</span>
              </div>
            </div>
          </div>

          <!-- Contact Info Card -->
          <div class="rounded-xl p-5 cyber-glow-card" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid (--border-color)'
          }">
            <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <Info :size="16" :style="{ color: 'var(--accent-color)' }" />
              联系方式
            </h3>
            <div class="space-y-3 text-sm">
              <div v-if="user.email" class="flex items-center gap-2 p-2.5 rounded-lg" :style="{
                backgroundColor: 'var(--surface-secondary)'
              }">
                <Mail :size="15" :style="{ color: 'var(--accent-color)' }" />
                <span class="truncate" :style="{ color: 'var(--text-secondary)' }">{{ user.email }}</span>
              </div>
              <div v-if="user.codeforces" class="flex items-center gap-2 p-2.5 rounded-lg" :style="{
                backgroundColor: 'var(--surface-secondary)'
              }">
                <ExternalLink :size="15" :style="{ color: 'var(--accent-color)' }" />
                <a :href="`https://codeforces.com/profile/${user.codeforces}`" target="_blank"
                  class="truncate transition-colors" :style="{ color: 'var(--accent-color)' }">
                  {{ user.codeforces }}
                </a>
              </div>
            </div>
          </div>

          <!-- Badges Card -->
          <div class="rounded-xl p-5 cyber-glow-card" :style="{
            backgroundColor: 'var(--surface-primary)',
            border: '1px solid var(--border-color)'
          }">
            <h3 class="text-sm font-semibold mb-4 flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
              <Medal :size="16" :style="{ color: 'var(--accent-color)' }" />
              成就徽章
            </h3>
            <div class="flex flex-wrap gap-2">
              <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
                backgroundColor: 'rgba(74, 222, 128, 0.08)',
                border: '1px solid rgba(74, 222, 128, 0.2)'
              }" title="首次解题">
                <span class="text-xl">🎯</span>
              </div>
              <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
                backgroundColor: 'rgba(14, 165, 233, 0.08)',
                border: '1px solid rgba(14, 165, 233, 0.2)'
              }" title="连续7天">
                <span class="text-xl">🔥</span>
              </div>
              <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
                backgroundColor: 'rgba(251, 191, 36, 0.08)',
                border: '1px solid rgba(251, 191, 36, 0.2)'
              }" title="百题达成">
                <span class="text-xl">💯</span>
              </div>
              <div class="w-12 h-12 rounded-lg flex items-center justify-center" :style="{
                backgroundColor: 'rgba(236, 72, 153, 0.08)',
                border: '1px solid rgba(236, 72, 153, 0.2)'
              }" title="活跃用户">
                <span class="text-xl">⭐</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Main Content Area -->
        <div class="lg:col-span-9">
          <div class="grid grid-cols-1 xl:grid-cols-12 gap-6">
            <!-- Left: Problem Data -->
            <div class="xl:col-span-8 space-y-6">
              <!-- Stats Cards -->
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <!-- Solved -->
                <div class="rounded-xl p-5 relative overflow-hidden group cyber-glow-card" :style="{
                  backgroundColor: 'var(--surface-primary)',
                  border: '1px solid var(--border-color)'
                }">
                  <div class="absolute -top-4 -right-4 w-16 h-16 rounded-full opacity-10 blur-md" :style="{
                    backgroundColor: '#4ade80'
                  }"></div>
                  <div class="relative z-10">
                    <div class="flex items-center justify-between mb-3">
                      <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">通过题目</span>
                      <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                        backgroundColor: 'rgba(74, 222, 128, 0.1)'
                      }">
                        <CheckCircle :size="18" :style="{ color: '#4ade80' }" />
                      </div>
                    </div>
                    <div class="flex items-end gap-2">
                      <span class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
                        {{ user.solved }}
                      </span>
                      <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">题</span>
                    </div>
                  </div>
                </div>

                <!-- Submissions -->
                <div class="rounded-xl p-5 relative overflow-hidden group cyber-glow-card" :style="{
                  backgroundColor: 'var(--surface-primary)',
                  border: '1px solid var(--border-color)'
                }">
                  <div class="absolute -top-4 -right-4 w-16 h-16 rounded-full opacity-10 blur-md" :style="{
                    backgroundColor: '#38bdf8'
                  }"></div>
                  <div class="relative z-10">
                    <div class="flex items-center justify-between mb-3">
                      <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">总提交</span>
                      <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                        backgroundColor: 'rgba(14, 165, 233, 0.1)'
                      }">
                        <Send :size="18" :style="{ color: '#38bdf8' }" />
                      </div>
                    </div>
                    <div class="flex items-end gap-2">
                      <span class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
                        {{ user.submission }}
                      </span>
                      <span class="text-sm mb-1" :style="{ color: 'var(--text-tertiary)' }">次</span>
                    </div>
                  </div>
                </div>

                <!-- Acceptance Rate -->
                <div class="rounded-xl p-5 relative overflow-hidden group cyber-glow-card" :style="{
                  backgroundColor: 'var(--surface-primary)',
                  border: '1px solid var(--border-color)'
                }">
                  <div class="absolute -top-4 -right-4 w-16 h-16 rounded-full opacity-10 blur-md" :style="{
                    backgroundColor: '#fbbf24'
                  }"></div>
                  <div class="relative z-10">
                    <div class="flex items-center justify-between mb-3">
                      <span class="text-sm font-medium" :style="{ color: 'var(--text-secondary)' }">通过率</span>
                      <div class="w-8 h-8 rounded-lg flex items-center justify-center" :style="{
                        backgroundColor: 'rgba(251, 191, 36, 0.1)'
                      }">
                        <TrendingUp :size="18" :style="{ color: '#fbbf24' }" />
                      </div>
                    </div>
                    <div class="flex items-end gap-2">
                      <span class="text-3xl font-bold font-terminal" :style="{ color: 'var(--text-primary)' }">
                        {{ formatAcceptance(user.accept, user.submission) }}
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Rating -->
                <div class="rounded-xl p-5 relative overflow-hidden group cyber-glow-card" :style="{
                  backgroundColor: 'var(--surface-primary)',
                  border: '1px solid var(--border-color)'
                }">
                  <div class="absolute -top-4 -right-4 w-16 h-16 rounded-full blur-md" :style="{
                    backgroundColor: formatRating(user.rating).bgColor,
                    opacity: 0.2
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
                      <span class="text-3xl font-bold font-terminal" :style="{ color: formatRating(user.rating).color }">
                        {{ user.rating }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Activity Heatmap -->
              <div class="rounded-xl p-5 cyber-glow-card" :style="{
                backgroundColor: 'var(--surface-primary)',
                border: '1px solid var(--border-color)'
              }">
                <div class="flex items-center justify-between mb-4">
                  <h3 class="text-sm font-semibold flex items-center gap-2" :style="{ color: 'var(--text-primary)' }">
                    <Flame :size="16" :style="{ color: 'var(--accent-color)' }" />
                    活动热图
                  </h3>
                  <n-radio-group v-model:value="heatmapPeriod" size="small">
                    <n-radio-button value="recent">最近一年</n-radio-button>
                    <n-radio-button value="2026">2026</n-radio-button>
                    <n-radio-button value="2025">2025</n-radio-button>
                  </n-radio-group>
                </div>
                <UserHeatmap :user-id="userId" :period="heatmapPeriod" :heatmap-data="heatmapRawData" />
              </div>

              <!-- Tabs Section -->
              <div class="rounded-xl overflow-hidden cyber-glow-card" :style="{
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

            <!-- Right: Knowledge & Recommendations -->
            <div class="xl:col-span-4">
              <div class="xl:sticky xl:top-24 space-y-4">
                <!-- Knowledge Radar Chart -->
                <KnowledgeRadar
                  :tag-scores="abilityData.tag_scores"
                  :strongest-tags="abilityData.strongest_tags"
                  :weakest-tags="abilityData.weakest_tags"
                />
                <!-- Language Overview -->
                <LanguageOverview :languages="abilityData.languages" />
              </div>
            </div>
          </div>
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
import CyberGridCanvas from '@/components/CyberGridCanvas.vue'
import UserHeatmap from './components/UserHeatmap.vue'
import UserSubmissions from './components/UserSubmissions.vue'
import UserSolutions from './components/UserSolutions.vue'
import UserBlogs from './components/UserBlogs.vue'
import UserDiscussions from './components/UserDiscussions.vue'
import KnowledgeRadar from './components/KnowledgeRadar.vue'
import LanguageOverview from './components/LanguageOverview.vue'
const { copy } = useClipboard()
const router = useRouter()
const route = useRoute()
const { id: currentUserId } = useUserStore()

const userId = ref<Number>(Number(route.params.id as string))
const user = ref<User & { solved: number }>({
  id: 0,
  username: '',
  avatar: '',
  introduction: '',
  password: '',
  codeforces: '',
  school: '',
  accept: 0,
  rating: 0,
  submission: 0,
  status: 0,
  solved: 0,
  email: '',
  nickname: '',
  gender: '0',
  user_role: 'user',
  birthday: '',
  banned_to: '',
  balance: 0,
  created_at: '',
  updated_at: '',
})

const rank = ref(42)
const activeTab = ref('submissions')
const heatmapPeriod = ref<'recent' | '2026' | '2025'>('recent')

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

const heatmapRawData = ref<{
  heatmaps?: Record<string, Record<string, number>>
  past_year_heatmap?: Record<string, number>
}>({})

// 知识掌握度数据
const abilityData = ref<{
  overall_score: number
  tag_scores: Record<string, number>
  strongest_tags: string[]
  weakest_tags: string[]
  languages: Record<string, number>
  avg_difficulty: number
}>({
  overall_score: 0,
  tag_scores: {},
  strongest_tags: [],
  weakest_tags: [],
  languages: {},
  avg_difficulty: 0,
})
const handleAddFriend = () => {
  // TODO: Implement add friend functionality
  console.log('Add friend:', userId.value)
}

const fetchUserInfo = async () => {
  try {
    const { code, info } = await userApi.getHomePageInfoById(userId.value)
    if (code === 200 && info) {
      user.value = info
      heatmapRawData.value = {
        heatmaps: info.heatmaps,
        past_year_heatmap: info.past_year_heatmap
      }
    }
  } catch (error) {
    console.error('Failed to fetch user info:', error)
  }
}

const fetchAbilityData = async () => {
  // TODO: 仅在自己的主页获取能力分析
  // if (!isOwnProfile.value) return
  try {
    const { code, info } = await userApi.getAbilityAnalysis()
    if (code === 200 && info) {
      abilityData.value = info
    }
  } catch (error) {
    console.error('Failed to fetch ability data:', error)
  }
}

onMounted(() => {
  fetchUserInfo()
  fetchAbilityData()
})
</script>
