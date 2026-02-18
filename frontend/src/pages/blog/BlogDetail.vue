<template>
  <article class="max-w-3xl px-8 py-4 mx-auto">
    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center min-h-[50vh]">
      <n-spin size="large" />
    </div>

    <!-- 博客内容 -->
    <main v-else-if="blog" class="animate-[fadeIn_0.4s_ease-out]">
      <!-- 头部信息 -->
      <header class="mb-10">
        <h1 class="tracking-tight mx-0 mt-0 mb-5 text-3xl font-bold text-(--text-primary)">{{ blog.title }}</h1>
        <div class="rounded-xl bg-(--surface-secondary) flex items-center gap-3 mb-6 p-4">
          <div
            class="font-semibold text-xl shrink-0 bg-(--color-primary-bg) text-(--color-primary) flex items-center justify-center w-11 h-11">
            <n-avatar :size="44" :src="blog?.avatar || undefined" />
          </div>
          <div class="flex flex-col gap-0.5">
            <span class="text-sm font-semibold text-(--text-primary)">{{ blog.username }}</span>
            <time class="text-sm text-(--text-tertiary)">{{ formatDateTime(blog.created_at) }}</time>
          </div>
        </div>

        <div v-if="blog.tags?.length" class="flex flex-wrap gap-1">
          <span v-for="tag in blog.tags" :key="tag"
            class="hover:bg-(--surface-tertiary) hover:text-(--text-primary) rounded-4xl text-(--text-secondary) bg-(--surface-secondary) border-(--border-secondary) transition-all duration-200 font-medium inline-flex items-center px-1 py-3 text-sm">{{
            tag }}</span>
        </div>
      </header>

      <!-- 正文内容 -->
      <section class="blog-detail__content">
        <v-md-preview :text="blog.context" />
      </section>

      <!-- 底部操作栏 -->
      <footer class="px-8 py-0">
        <div class="flex items-center">
          <button class="blog-detail__stat"
            :class="{ 'text-(--color-primary) bg-(--color-primary-bg) border-(--color-primary-border)': isLiked }"
            @click="likeBlog">
            <Heart :size="18" :fill="isLiked ? 'currentColor' : 'none'" />
            <span>{{ formatNumber(blog.like) }}</span>
          </button>

          <button class="blog-detail__stat"
            :class="{ 'text-(--color-primary) bg-(--color-primary-bg) border-(--color-primary-border)': isCollected }"
            @click="collectBlog">
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
    <div v-else class="flex flex-col items-center justify-center gap-4 min-h-[50vh] text-center
">
      <FileX :size="48" class="text-(--text-tertiary)" />
      <span class="m-0 text-2xl font-semibold text-(--text-primary)">博客未找到</span>
      <span class="m-0 text-(--text-secondary)">该博客可能已被删除或不存在</span>
      <n-button type="primary" @click="$router.push({ name: 'Blogs' })">返回博客列表</n-button>
    </div>
  </article>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { NButton, NSpin, useMessage, NAvatar } from 'naive-ui'
import { Heart, Bookmark, Share2, FileX } from 'lucide-vue-next'
import { blogApi } from '@/services/blog'

const route = useRoute()
const router = useRouter()
const message = useMessage()

const blog = ref(null)
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
      .catch(() => { })
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

.blog-detail__content {
  padding: 2rem 0;
  border-top: 1px solid var(--border-primary);
  border-bottom: 1px solid var(--border-primary);
  color: var(--text-primary);
  line-height: 1.75;
  font-size: 1.0625rem;
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
