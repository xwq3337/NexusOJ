<template>
  <div class="min-h-screen bg-transparent">
    <!-- 搜索和筛选栏 -->
    <div class="flex flex-col md:flex-row gap-3 mb-4 items-center">
      <n-button type="info" @click="$router.push({ name: 'BlogCreate' })">
        <template #icon>
          <Plus />
        </template>
        我要创作
      </n-button>
      <n-input v-model:value="searchQuery" placeholder="搜索博客标题、内容或标签..." class="flex-1 min-w-60">
        <template #prefix>
          <Search :size="16" />
        </template>
      </n-input>
      <n-select v-model:value="selectedCategory" :options="categoryOptions" placeholder="分类" class="min-w-40" />
      <n-button round type="tertiary">
        <template #icon>
          <RefreshCcw />
        </template>
        刷新
      </n-button>
    </div>
    <!-- 博客列表 -->
    <div class="flex flex-col md:flex-row gap-6">
      <div class="flex-1">
        <n-list bordered>
          <n-list-item v-for="blog in blogs" :key="blog.id" class="p-2  hover:bg-[#eeeeee]">
            <n-row>
              <span @click="$router.push({ name: 'BlogDetail', params: { id: blog.id } })"
                class="font-bold text-lg  cursor-pointer"> {{ blog.title }}</span>
            </n-row>
            <n-row>
              ...............
            </n-row>
            <n-row>
              <n-col>
                <span style="color: #8a919f;" class="justify-center align-center">
                  {{ blog.username }}
                </span>
              </n-col>
              <n-divider vertical />
              <div class="flex items-center gap-2" style="color: #8a919f;">
                {{ formatRelativeTime(blog.created_at) }}
                <n-divider vertical />
                <Eye :size="20" />
                {{ formatNumber(Number(Math.random().toPrecision(1) ) * 10) }}
                <n-divider vertical />
                <ThumbsUp :size="16" />
                {{ formatNumber(blog.like ?? 0) }}
              </div>
            </n-row>
          </n-list-item>
          <template #footer>
            <n-pagination />
          </template>
        </n-list>
      </div>

      <!-- 侧边栏-榜单 -->
      <div class="w-full md:w-80 shrink-0">
        <!-- 活跃用户榜 -->
        <n-card title="活跃用户" :bordered="false" class="mb-4">
          <n-list>
            <n-list-item v-for="user in activeUsers" :key="user.id" class="py-2!">
              <div class="flex items-center gap-3">
                <div class="flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold" :style="{
                  backgroundColor: user.rank <= 3 ? 'var(--hero-bg-from)' : 'var(--divider-color)',
                  color: user.rank <= 3 ? '#fff' : 'var(--text-secondary)'
                }">
                  {{ user.rank }}
                </div>
                <n-avatar round :size="36" :src="user.avatar" />
                <div class="flex-1 min-w-0">
                  <div class="truncate font-medium" :style="{ color: 'var(--text-primary)' }">
                    {{ user.username }}
                  </div>
                  <div class="text-xs" :style="{ color: 'var(--text-secondary)' }">
                    {{ user.solved }} 道题目
                  </div>
                </div>
                <div class="text-xs" :style="{ color: 'var(--text-tertiary)' }">
                  {{ user.rating }}
                </div>
              </div>
            </n-list-item>
          </n-list>
        </n-card>

        <!-- 热门博客 -->
        <n-card title="热门博客" :bordered="false" class="mb-4">
          <n-list>
            <n-list-item v-for="(blog, index) in hotBlogs" :key="blog.id" class="py-2!">
              <div class="flex items-start gap-3">
                <div class="flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold shrink-0" :style="{
                  backgroundColor: index + 1 <= 3 ? 'var(--hero-bg-from)' : 'var(--divider-color)',
                  color: index + 1 <= 3 ? '#fff' : 'var(--text-secondary)'
                }">
                  {{ index + 1 }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="truncate text-sm cursor-pointer hover:underline" :style="{ color: 'var(--text-primary)' }"
                    @click="router.push(`/blog/${blog.id}`)">
                    {{ blog.title }}
                  </div>
                  <div class="flex items-center gap-2 text-xs mt-1" :style="{ color: 'var(--text-secondary)' }">
                    <span>{{ blog.username }}</span>
                    <span>·</span>
                    <div class="flex items-center gap-1">
                      <Eye :size="12" />
                      {{ formatNumber(Math.random() * 100) }}
                    </div>
                  </div>
                </div>
              </div>
            </n-list-item>
          </n-list>
        </n-card>

        <!-- 最新动态 -->
        <n-card title="最新动态" :bordered="false">
          <n-list>
            <n-list-item v-for="activity in recentActivities" :key="activity.id" class="py-2!">
              <div class="flex items-start gap-3">
                <n-avatar round :size="32" :src="activity.avatar" />
                <div class="flex-1 min-w-0">
                  <div class="text-sm" :style="{ color: 'var(--text-primary)' }">
                    <span class="font-medium">{{ activity.username }}</span>
                    <span class="mx-1" :style="{ color: 'var(--text-tertiary)' }">{{ activity.action }}</span>
                    <span v-if="activity.target" class="cursor-pointer hover:underline"
                      @click="router.push(activity.targetUrl || '')">
                      {{ activity.target }}
                    </span>
                  </div>
                  <div class="text-xs mt-0.5" :style="{ color: 'var(--text-tertiary)' }">
                    {{ activity.time }}
                  </div>
                </div>
              </div>
            </n-list-item>
          </n-list>
        </n-card>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, RefreshCcw, ThumbsUp, Eye, Heart, FileX, Plus } from 'lucide-vue-next'
import { NInput, NButton, NDivider, NCol, NRow, NSelect, NPagination, NList, NListItem, NCard, NAvatar } from 'naive-ui'
import type { Blog } from '@/types/blog'
import { blogApi } from '@/services/blog'
import { formatRelativeTime } from '@/utils/format'

const router = useRouter()
// 博客数据
const blogs = ref<(Blog & { user_id: string, username: string })[]>()

// 加载博客数据
onMounted(async () => {
  try {
    const response = await blogApi.getList()
    blogs.value = response.info
  } catch (error) {
    console.error('获取博客数据失败:', error)
  } finally {
    loading.value = false
  }
})

// 筛选条件
const selectedCategory = ref('all')
const searchQuery = ref('')
const loading = ref(true)

// 分类选项
const categoryOptions = [
  { label: '全部', value: 'all' },
  { label: '最新', value: 'latest' },
  { label: '最热', value: 'hot' }
]

// 侧边栏数据
const activeUsers = ref([
  { id: 1, username: '张三', avatar: '', rank: 1, solved: 156, rating: 2450 },
  { id: 2, username: '李四', avatar: '', rank: 2, solved: 142, rating: 2380 },
  { id: 3, username: '王五', avatar: '', rank: 3, solved: 138, rating: 2310 },
  { id: 4, username: '赵六', avatar: '', rank: 4, solved: 125, rating: 2250 },
  { id: 5, username: '钱七', avatar: '', rank: 5, solved: 118, rating: 2190 }
])

const hotBlogs = blogs.value
const recentActivities = ref([
  { id: 1, username: '张三', action: '发布了', target: 'Vue 3 组合式 API 详解', targetUrl: '/blog/1', avatar: '', time: '5分钟前' },
  { id: 2, username: '李四', action: '解决了', target: '动态规划入门', targetUrl: '/problem/1001', avatar: '', time: '10分钟前' },
  { id: 3, username: '王五', action: '评论了', target: 'TypeScript 高级类型', targetUrl: '/blog/2', avatar: '', time: '15分钟前' },
  { id: 4, username: '赵六', action: '通过了', target: '二叉树遍历', targetUrl: '/problem/1002', avatar: '', time: '20分钟前' },
  { id: 5, username: '钱七', action: '点赞了', target: 'React Hooks 实战', targetUrl: '/blog/3', avatar: '', time: '25分钟前' }
])

// 生成博客封面图片



// 格式化时间显示
const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) {
    const hours = Math.floor(diff / (1000 * 60 * 60))
    if (hours === 0) {
      const minutes = Math.floor(diff / (1000 * 60))
      return minutes === 0 ? '刚刚' : `${minutes}分钟前`
    }
    return `${hours}小时前`
  } else if (days === 1) {
    return '昨天'
  } else if (days < 7) {
    return `${days}天前`
  } else if (days < 30) {
    return `${Math.floor(days / 7)}周前`
  } else {
    return date.toLocaleDateString('zh-CN', { year: 'numeric', month: 'numeric', day: 'numeric' })
  }
}

// 格式化数字显示
const formatNumber = (num: number) => {
  if (num >= 10000) {
    return `${(num / 10000).toFixed(1)}w`
  } else if (num >= 1000) {
    return `${(num / 1000).toFixed(1)}k`
  }
  return num.toString()
}
</script>
