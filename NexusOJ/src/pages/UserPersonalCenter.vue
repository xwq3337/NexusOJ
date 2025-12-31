<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <div class="text-center">
              <n-avatar :size="80" :src="userStore.avatar || undefined" class="mx-auto mb-4">
                <UserIcon v-if="!userStore.avatar" />
              </n-avatar>
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                {{ userStore.nickname || userStore.account }}
              </h2>
              <p class="text-gray-500 dark:text-gray-400 text-sm">
                Rating: {{ userStore.rating || 'N/A' }}
              </p>
            </div>
            <n-divider />
            <n-menu
              v-model:value="activeTab"
              :options="menuOptions"
              @update:value="handleMenuClick"
            />
          </div>
        </div>

        <!-- Main Content -->
        <div class="lg:col-span-3">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <!-- Profile Tab -->
            <div v-if="activeTab === 'profile'">
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                <UserIcon class="inline mr-2" />
                个人信息
              </h3>
              <n-form :model="profileForm" label-placement="left" label-width="auto">
                <n-form-item label="昵称">
                  <n-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
                </n-form-item>
                <n-form-item label="性别">
                  <n-select v-model:value="profileForm.gender" :options="genderOptions" />
                </n-form-item>
                <n-form-item label="头像">
                  <n-upload
                    v-model:file-list="fileList"
                    :show-file-list="false"
                    :on-change="handleAvatarChange"
                    accept="image/*"
                  >
                    <n-button>
                      <UploadIcon class="mr-2" />
                      上传头像
                    </n-button>
                  </n-upload>
                </n-form-item>
                <n-form-item>
                  <n-button type="primary" @click="updateProfile">
                    <SaveIcon class="mr-2" />
                    保存更改
                  </n-button>
                </n-form-item>
              </n-form>
            </div>

            <!-- Settings Tab -->
            <div v-if="activeTab === 'settings'">
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                <SettingsIcon class="inline mr-2" />
                设置
              </h3>
              <n-form label-placement="left" label-width="auto">
                <n-form-item label="语言">
                  <n-select v-model:value="language" :options="languageOptions" />
                </n-form-item>
              </n-form>
            </div>

            <!-- Stats Tab -->
            <div v-if="activeTab === 'stats'">
              <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                <BarChartIcon class="inline mr-2" />
                统计信息
              </h3>
              <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                <n-card>
                  <n-statistic title="总提交数" :value="stats.totalSubmissions" />
                </n-card>
                <n-card>
                  <n-statistic title="通过数" :value="stats.accepted" />
                </n-card>
                <n-card>
                  <n-statistic title="通过率" :value="stats.acceptanceRate" suffix="%" />
                </n-card>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useUserStore } from '@/stores/useUserStore'
import {
  User as UserIcon,
  Settings as SettingsIcon,
  BarChart as BarChartIcon,
  Upload as UploadIcon,
  Save as SaveIcon
} from 'lucide-vue-next'
import {
  NAvatar,
  NButton,
  NCard,
  NDivider,
  NForm,
  NFormItem,
  NInput,
  NMenu,
  NSelect,
  NStatistic,
  NUpload
} from 'naive-ui'

const userStore = useUserStore()

const activeTab = ref('profile')
const theme = ref(localStorage.getItem('theme') || 'light')
const language = ref('zh')

const profileForm = reactive({
  nickname: userStore.nickname,
  gender: userStore.gender
})

const fileList = ref([])

const menuOptions = [
  {
    label: '个人信息',
    key: 'profile',
    icon: () => h(UserIcon)
  },
  {
    label: '设置',
    key: 'settings',
    icon: () => h(SettingsIcon)
  },
  {
    label: '统计信息',
    key: 'stats',
    icon: () => h(BarChartIcon)
  }
]

const genderOptions = [
  { label: '男', value: 1 },
  { label: '女', value: 2 },
  { label: '其他', value: 0 }
]

const themeOptions = [
  { label: '浅色', value: 'light' },
  { label: '深色', value: 'dark' }
]

const languageOptions = [
  { label: '中文', value: 'zh' },
  { label: 'English', value: 'en' }
]

const stats = reactive({
  totalSubmissions: 0,
  accepted: 0,
  acceptanceRate: 0
})

const handleMenuClick = (key: string) => {
  activeTab.value = key
}

const handleAvatarChange = (options: any) => {
  // Handle avatar upload
  console.log('Avatar changed', options)
}

const updateProfile = () => {
  // Update profile logic
  console.log('Update profile', profileForm)
}

const changeTheme = (value: string) => {
  theme.value = value
  localStorage.setItem('theme', value)
  document.body.classList.remove('light', 'dark')
  document.body.classList.add(value)
  document.documentElement.classList.remove('light', 'dark')
  document.documentElement.classList.add(value)
}

onMounted(() => {
  // Load stats
  // This would typically be an API call
  stats.totalSubmissions = 100
  stats.accepted = 75
  stats.acceptanceRate = Math.round((stats.accepted / stats.totalSubmissions) * 100)
})
</script>
