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

    <!-- 博客列表 - 紧凑型卡片布局 -->
    <div class="blog-list">
      <div v-if="loading" class="flex justify-center items-center py-12">
        <n-spin size="large" />
      </div>
      <template v-else>
        <div
          v-for="blog in filteredBlogs"
          :key="blog.id"
          class="blog-card"
          :style="{ backgroundColor: 'var(--surface-primary)' }"
          @click="router.push({ name: 'BlogDetail', params: { id: blog.id } })"
        >
          <!-- 用户信息区域 -->
          <div class="blog-header">
            <div class="user-info">
              <div class="avatar">{{ blog.username?.charAt(0)?.toUpperCase() || 'U' }}</div>
              <div class="user-details">
                <span class="username">{{ blog.username }}</span>
                <span class="publish-time">{{ formatTime(blog.created_at) }}</span>
              </div>
            </div>
          </div>

          <!-- 博客标题 -->
          <h3 class="blog-title">{{ blog.title }}</h3>

          <!-- 博客摘要 -->
          <p class="blog-excerpt">
            {{ blog.context.substring(0, 120) }}{{ blog.context.length > 120 ? '...' : '' }}
          </p>

          <!-- 标签 -->
          <div class="blog-tags" v-if="blog.tags && blog.tags.length">
            <span v-for="tag in blog.tags.slice(0, 3)" :key="tag" class="tag">
              {{ tag }}
            </span>
          </div>

          <!-- 底部统计信息 -->
          <div class="blog-footer">
            <div class="stats">
              <div class="stat-item">
                <Eye :size="14" />
                <span>{{ formatNumber(blog.view || 0) }}</span>
              </div>
              <div class="stat-item">
                <Heart :size="14" />
                <span>{{ formatNumber(blog.like || 0) }}</span>
              </div>
              <div class="stat-item">
                <Bookmark :size="14" />
                <span>{{ formatNumber(blog.collection || 0) }}</span>
              </div>
            </div>
          </div>
        </div>
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
import { Search, Eye, Heart, Bookmark } from 'lucide-vue-next'
import { NInput, NSelect, NPagination, NSpin } from 'naive-ui'
import Request from '@/services/api'
import type { Blog } from '@/types'

const router = useRouter()

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
  } else {
    return date.toLocaleDateString()
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

<style scoped>
.blog-list {
  display: grid;
  gap: 16px;
}

.blog-card {
  padding: 20px;
  border-radius: 12px;
  border: 1px solid var(--border-color, #e5e7eb);
  cursor: pointer;
  transition: all 0.2s ease;
  background: var(--surface-primary);
}

.blog-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-color, #6366f1);
}

/* 头部用户信息 */
.blog-header {
  margin-bottom: 12px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 16px;
  flex-shrink: 0;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.username {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.publish-time {
  font-size: 12px;
  color: var(--text-secondary);
}

/* 博客标题 */
.blog-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0 0 8px 0;
  color: var(--text-primary);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 博客摘要 */
.blog-excerpt {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0 0 12px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* 标签 */
.blog-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.tag {
  display: inline-block;
  padding: 4px 10px;
  background: var(--tag-bg, rgba(99, 102, 241, 0.1));
  color: var(--tag-color, #6366f1);
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

/* 底部统计 */
.blog-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid var(--border-color, #e5e7eb);
}

.stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--text-secondary);
}

.stat-item svg {
  stroke: var(--text-secondary);
  opacity: 0.7;
}

/* 响应式 */
@media (max-width: 768px) {
  .blog-card {
    padding: 16px;
  }

  .blog-title {
    font-size: 16px;
  }

  .blog-excerpt {
    font-size: 13px;
  }

  .stats {
    gap: 12px;
  }
}
</style>
