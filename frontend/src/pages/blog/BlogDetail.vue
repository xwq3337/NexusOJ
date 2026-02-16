<template>
  <article class="blog-detail">
    <!-- 加载状态 -->
    <div v-if="loading" class="blog-detail__loading">
      <n-spin size="large" />
    </div>

    <!-- 博客内容 -->
    <main v-else-if="blog" class="blog-detail__main">
      <!-- 头部信息 -->
      <header class="blog-detail__header">
        <h1 class="blog-detail__title">{{ blog.title }}</h1>

        <div class="blog-detail__author-card">
          <div class="blog-detail__avatar">
            U
            <!-- TODO {{ blog.username?.charAt(0).toUpperCase() ?? 'U' }} -->
          </div>
          <div class="blog-detail__author-info">
            <span class="blog-detail__author-name">{{ blog.user_id }}</span>
            <time class="blog-detail__publish-time">{{ formatDateTime(blog.created_at) }}</time>
          </div>
        </div>

        <div v-if="blog.tags?.length" class="blog-detail__tags">
          <span v-for="tag in blog.tags" :key="tag" class="blog-detail__tag">{{ tag }}</span>
        </div>
      </header>

      <!-- 正文内容 -->
      <section class="blog-detail__content">
        <v-md-preview :text="blog.context" />
      </section>

      <!-- 底部操作栏 -->
      <footer class="blog-detail__footer">
        <div class="blog-detail__stats">
          <button class="blog-detail__stat" :class="{ 'blog-detail__stat--active': isLiked }" @click="likeBlog">
            <Heart :size="18" :fill="isLiked ? 'currentColor' : 'none'" />
            <span>{{ formatNumber(blog.like) }}</span>
          </button>

          <button class="blog-detail__stat" :class="{ 'blog-detail__stat--active': isCollected }" @click="collectBlog">
            <Bookmark :size="18" :fill="isCollected ? 'currentColor' : 'none'" />
            <span>{{ formatNumber(blog.collection) }}</span>
          </button>

          <button class="blog-detail__stat" @click="shareBlog">
            <Share2 :size="18" />
            <span>分享</span>
          </button>
        </div>
      </footer>
    </main>

    <!-- 空状态 -->
    <div v-else class="blog-detail__empty">
      <FileX :size="48" :style="{ color: 'var(--text-tertiary)' }" />
      <h2>博客未找到</h2>
      <p>该博客可能已被删除或不存在</p>
      <n-button type="primary" @click="$router.push({name : 'Blogs'})">返回博客列表</n-button>
    </div>
  </article>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpin, useMessage } from 'naive-ui'
import { Heart, Bookmark, Share2, FileX } from 'lucide-vue-next'
import type { Blog } from '@/types/blog'
import { blogApi } from '@/services/blog'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const blog = ref<Blog | null>(null)
const loading = ref(true)
const isLiked = ref(false)
const isCollected = ref(false)

const fetchBlogDetail = async () => {
  try {
    const response = await blogApi.getBlogById(route.params.id as string)
    blog.value = response.info
  } catch (error) {
    console.error('Failed to fetch blog:', error)
    message.error('获取博客详情失败')
  } finally {
    loading.value = false
  }
}

const likeBlog = async () => {
  if (!blog.value) return

  try {
    // TODO: 调用点赞API
    // await Request.post(`/blog/like/${blog.value.id}`)
    isLiked.value = !isLiked.value
    blog.value.like += isLiked.value ? 1 : -1
  } catch (error) {
    message.error('操作失败')
  }
}

const collectBlog = async () => {
  if (!blog.value) return

  try {
    // TODO: 调用收藏API
    // await Request.post(`/blog/collect/${blog.value.id}`)
    isCollected.value = !isCollected.value
    blog.value.collection += isCollected.value ? 1 : -1
  } catch (error) {
    message.error('操作失败')
  }
}

const shareBlog = () => {
  const url = window.location.href
  if (navigator.share) {
    navigator
      .share({ title: blog.value?.title, url })
      .catch(() => {})
  } else {
    navigator.clipboard
      .writeText(url)
      .then(() => message.success('链接已复制'))
      .catch(() => message.error('复制失败'))
  }
}

const formatDateTime = (date: string) => {
  const d = new Date(date)
  return d.toLocaleString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatNumber = (num: number) => {
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num.toString()
}

onMounted(fetchBlogDetail)
</script>

<style scoped>
.blog-detail {
  max-width: 768px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.blog-detail__loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

.blog-detail__main {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.blog-detail__header {
  margin-bottom: 2.5rem;
}

.blog-detail__meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  color: var(--text-secondary);
}

.blog-detail__separator {
  color: var(--text-tertiary);
}

.blog-detail__title {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1.3;
  color: var(--text-primary);
  margin: 0 0 1.25rem;
  letter-spacing: -0.02em;
}

.blog-detail__author-card {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 1rem;
  margin-bottom: 1.5rem;
  background: var(--surface-secondary);
  border-radius: 0.75rem;
}

.blog-detail__avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.75rem;
  height: 2.75rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-primary);
  background: var(--color-primary-bg);
  border-radius: 50%;
  flex-shrink: 0;
}

.blog-detail__author-info {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.blog-detail__author-name {
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--text-primary);
}

.blog-detail__publish-time {
  font-size: 0.8125rem;
  color: var(--text-tertiary);
}

.blog-detail__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.blog-detail__tag {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-secondary);
  background: var(--surface-secondary);
  border-radius: 9999px;
  transition: all 0.2s;
}

.blog-detail__tag:hover {
  background: var(--surface-tertiary);
  color: var(--text-primary);
}

.blog-detail__content {
  padding: 2rem 0;
  border-top: 1px solid var(--border-primary);
  border-bottom: 1px solid var(--border-primary);
  color: var(--text-primary);
  line-height: 1.75;
  font-size: 1.0625rem;
}

.blog-detail__footer {
  padding: 2rem 0;
}

.blog-detail__stats {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.blog-detail__stat {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  font-size: 0.9375rem;
  font-weight: 500;
  color: var(--text-secondary);
  background: transparent;
  border: 1px solid var(--border-primary);
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
}

.blog-detail__stat:hover {
  color: var(--text-primary);
  background: var(--surface-secondary);
  border-color: var(--border-secondary);
}

.blog-detail__stat--active {
  color: var(--color-primary);
  background: var(--color-primary-bg);
  border-color: var(--color-primary-border);
}

.blog-detail__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  min-height: 50vh;
  text-align: center;
}

.blog-detail__empty h2 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
}

.blog-detail__empty p {
  margin: 0;
  color: var(--text-secondary);
}

/* Markdown 内容样式重写 */
:deep(.v-md-editor) {
  background: transparent;
}

:deep(.v-md-editor .v-md-preview) {
  padding: 0;
  background: transparent;
}

:deep(.v-md-preview .v-md-preview-theme) {
  background: transparent;
}

/* 响应式 */
@media (max-width: 640px) {
  .blog-detail {
    padding: 1.5rem 1rem;
  }

  .blog-detail__title {
    font-size: 1.625rem;
  }

  .blog-detail__content {
    padding: 1.5rem 0;
    font-size: 1rem;
  }

  .blog-detail__stat {
    padding: 0.5rem 0.875rem;
    font-size: 0.875rem;
  }

  .blog-detail__stat span {
    display: none;
  }

  .blog-detail__stat--active span,
  .blog-detail__stat:hover span {
    display: inline;
  }
}
</style>
