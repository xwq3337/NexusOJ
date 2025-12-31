<template>
  <div class="blog-container p-6">
    <div class="mb-6">
      <h2 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">技术博客</h2>
      <p class="text-gray-500 mt-2">分享编程知识，交流技术心得</p>
    </div>

    <!-- 博客搜索和筛选 -->
    <div class="flex flex-col md:flex-row gap-4 mb-6">
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索博客..."
        :style="{ backgroundColor: 'var(--surface-tertiary)' }"
      >
        <template #prefix>
          <Search :size="18" />
        </template>
      </n-input>
      <n-select
        v-model:value="selectedCategory"
        :options="categoryOptions"
        placeholder="选择分类"
        class="w-full md:w-48"
      />
      <n-select
        v-model:value="selectedTag"
        :options="tagOptions"
        placeholder="选择标签"
        class="w-full md:w-40"
      />
    </div>

    <!-- 博客列表 -->
    <div class="space-y-6">
      <div v-if="loading" class="flex justify-center items-center py-12">
        <n-spin size="large" />
      </div>
      <template v-else>
        <n-card
          v-for="blog in filteredBlogs"
          :key="blog.id"
          class="blog-card hover:shadow-lg transition-shadow cursor-pointer"
          :style="{ backgroundColor: 'var(--surface-primary)' }"
          @click="$router.push({ name: 'BlogDetail', params: { id: blog.id } })"
        >
          <div class="flex flex-col md:flex-row gap-6">
            <!-- 博客图片/封面 -->
            <div class="md:w-1/4">
              <div
                class="h-40 bg-linear-to-r from-indigo-500 to-blue-600 rounded-lg flex items-center justify-center"
              >
                <FileText :size="48" class="text-white" />
              </div>
            </div>

            <!-- 博客内容 -->
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <span
                  class="px-2 py-1 bg-blue-500/10 text-blue-400 rounded-full text-xs font-medium"
                >
                  作者: {{ blog.username }}
                </span>
                <span class="text-gray-500 text-sm">{{
                  new Date(blog.created_at).toLocaleDateString()
                }}</span>
              </div>

              <h3 class="text-xl font-bold mb-2" :style="{ color: 'var(--text-primary)' }">
                {{ blog.title }}
              </h3>

              <p class="text-gray-500 mb-3" :style="{ color: 'var(--text-secondary)' }">
                {{ blog.context.substring(0, 100) }}{{ blog.context.length > 100 ? '...' : '' }}
              </p>

              <div class="flex flex-wrap gap-2 mb-3">
                <span
                  v-for="tag in blog.tags"
                  :key="tag"
                  class="px-2 py-1 bg-gray-500/10 text-gray-400 rounded text-xs"
                >
                  {{ tag }}
                </span>
              </div>

              <div class="flex items-center justify-between text-sm text-gray-500">
                <div class="flex items-center">
                  <User :size="14" class="mr-1" />
                  <span>{{ blog.username }}</span>
                </div>
                <div class="flex items-center gap-4">
                  <div class="flex items-center">
                    <MessageCircle :size="14" class="mr-1" />
                    <span>{{ blog.collection }}</span>
                  </div>
                  <div class="flex items-center">
                    <Eye :size="14" class="mr-1" />
                    <span>{{ blog.like }}</span>
                  </div>
                  <div class="flex items-center">
                    <Heart :size="14" class="mr-1" />
                    <span>{{ blog.like }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </n-card>
      </template>

      <div v-if="!loading && filteredBlogs.length === 0" class="text-center py-12 text-gray-500">
        没有找到匹配的博客
      </div>
    </div>

    <!-- 分页 -->
    <div class="mt-8 flex justify-center">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-size="pageSize"
        :item-count="blogs.length"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search,
  FileText,
  User,
  MessageCircle,
  Eye,
  Heart,
  Code,
  Cpu,
  Palette,
  Database,
  BookText
} from 'lucide-vue-next'
import { NCard, NInput, NSelect, NPagination, NSpin } from 'naive-ui'
import Request from '@/services/api'
import type { Blog } from '@/types'

// 博客数据
const blogs = ref<Blog[]>([])

// 加载博客数据
onMounted(async () => {
  try {
    const response = await Request.get('/blog/list')
    // 假设API返回的数据在info字段中
    blogs.value = response.info || response
    console.log('博客数据加载成功:', blogs.value)
  } catch (error) {
    console.error('获取博客数据失败:', error)
  } finally {
    loading.value = false
  }
})

// 筛选条件
const selectedCategory = ref(null)
const selectedTag = ref(null)
const searchQuery = ref('')
const loading = ref(true)
const currentPage = ref(1)
const pageSize = 6
const totalPages = computed(() => Math.ceil(blogs.value.length / pageSize))

// 分类选项
const categoryOptions = [
  { label: '全部', value: 'all' },
  { label: '按作者', value: 'username' },
  { label: '按ID', value: 'id' }
]

// 标签选项
const tagOptions = [
  { label: 'Vue', value: 'Vue' },
  { label: 'React', value: 'React' },
  { label: 'JavaScript', value: 'JavaScript' },
  { label: 'TypeScript', value: 'TypeScript' },
  { label: 'Node.js', value: 'Node.js' },
  { label: '算法', value: '算法' },
  { label: '面试', value: '面试' }
]

// 过滤后的博客
const filteredBlogs = computed(() => {
  let result = blogs.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(
      (blog) =>
        blog.title.toLowerCase().includes(query) ||
        blog.context.toLowerCase().includes(query) ||
        blog.username.toLowerCase().includes(query) ||
        blog.tags.some((tag: string) => tag.toLowerCase().includes(query))
    )
  }

  if (selectedCategory.value && selectedCategory.value !== 'all') {
    if (selectedCategory.value === 'username') {
      // 这里可以根据需要实现用户名过滤逻辑
      // 暂时保留所有博客
    } else if (selectedCategory.value === 'id') {
      // 这里可以根据需要实现ID过滤逻辑
      // 暂时保留所有博客
    }
  }

  if (selectedTag.value) {
    result = result.filter((blog) => blog.tags.includes(selectedTag.value))
  }

  // 分页
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return result.slice(start, end)
})
</script>

<style scoped>
.blog-card {
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.blog-card:hover {
  transform: translateY(-4px);
}
</style>
