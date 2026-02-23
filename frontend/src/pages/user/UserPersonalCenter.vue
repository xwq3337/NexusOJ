<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Sidebar -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <div class="text-center">
              <n-avatar :size="80" :src="avatar + `?t=${Date.now()}`" class="mx-auto mb-4">
                <UserIcon v-if="!avatar" />
              </n-avatar>
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
                {{ nickname || username }}
              </h2>
              <p class="text-gray-500 dark:text-gray-400 text-sm">
                ID: {{ id }}
              </p>
              <p class="text-gray-500 dark:text-gray-400 text-sm">
                Rating: {{ rating || 'N/A' }}
              </p>
            </div>
            <n-divider />
            <n-menu :options="menuOptions" />
          </div>
        </div>
        <!-- Main Content -->
        <div class="lg:col-span-3">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
            <RouterView />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue'
import { useUserStore } from '@/stores/useUserStore'
const {id, username, nickname,avatar,rating} = useUserStore()
import {
  User,
  Settings,
  BarChart,
  LockKeyhole
} from 'lucide-vue-next'
import {
  NAvatar,
  NDivider,
  NMenu,
  type MenuOption
} from 'naive-ui'
import { RouterLink } from 'vue-router'

const refreshAvatar = () => {
  console.log('Refreshing avatar...')
}
const menuOptions: MenuOption[] = [
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: 'Profile',
          }
        },
        { default: () => '个人信息' }
      ),
    key: 'profile',
    icon: () => h(User)
  },
    {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: 'Security',
          }
        },
        { default: () => '安全中心' }
      ),
    key: 'security',
    icon: () => h(LockKeyhole)
  },
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: 'Settings',
          }
        },
        { default: () => '设置' }
      ),
    key: 'settings',
    icon: () => h(Settings)
  },
  {
    label: () =>
      h(
        RouterLink,
        {
          to: {
            name: 'UserRecord',
          }
        },
        { default: () => '提交记录' }
      ),
    key: 'record',
    icon: () => h(BarChart)
  }
]

</script>
