<template>
  <div :style="{ backgroundColor: 'var(--bg-primary)' }" class="min-h-screen flex">
    <!-- 左侧导航菜单 -->
    <div
      class="w-40 border-r"
      :style="{
        backgroundColor: 'var(--surface-primary)',
        borderColor: 'var(--border-color)'
      }"
    >
      <div class="p-4">
        <n-menu :options="menuOptions" :value="currentRoute" @update:value="handleMenuSelect" />
      </div>
    </div>

    <!-- 右侧内容区域 -->
    <div class="flex-1 p-3">
      <div class="max-w-7xl mx-auto">
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NMenu, type MenuOption } from 'naive-ui'

const route = useRoute()
const router = useRouter()

// 当前路由名称
const currentRoute = computed(() => (route.name as string) || 'Courses')

// 页面标题
const pageTitle = computed(() => {
  switch (currentRoute.value) {
    case 'Courses':
      return '课程'
    case 'Blogs':
      return '博客'
    default:
      return '知识库'
  }
})

// 菜单选项
const menuOptions: MenuOption[] = [
  {
    label: '课程',
    key: 'Courses',
    props: {
      onClick: () => router.push({ name: 'Courses' })
    }
  },
  {
    label: '博客',
    key: 'Blogs',
    props: {
      onClick: () => router.push({ name: 'Blogs' })
    }
  }
]

// 处理菜单选择
const handleMenuSelect = (key: string) => {
  router.push({ name: key })
}
</script>
