<template>
  <div class="blog-detail-container p-6 max-w-4xl mx-auto">
    <div v-if="loading" class="flex justify-center items-center py-12">
      <n-spin size="large" />
    </div>

    <div v-else-if="blog" class="blog-content">
      <!-- 博客标题和元信息 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold mb-4" :style="{ color: 'var(--text-primary)' }">
          {{ blog.title }}
        </h1>

        <div class="flex flex-wrap items-center gap-4 text-sm text-gray-500 mb-4">
          <div class="flex items-center">
            <User :size="16" class="mr-1" />
            <span>{{ blog.username }}</span>
          </div>
          <div class="flex items-center">
            <Calendar :size="16" class="mr-1" />
            <span>{{ new Date(blog.created_at).toLocaleString() }}</span>
          </div>
          <div class="flex items-center">
            <Heart :size="16" class="mr-1" />
            <span>{{ blog.like }} 点赞</span>
          </div>
          <div class="flex items-center">
            <Folder :size="16" class="mr-1" />
            <span>{{ blog.collection }} 收藏</span>
          </div>
        </div>

        <div class="flex flex-wrap gap-2 mt-4">
          <n-tag
            v-for="tag in blog.tags"
            :key="tag"
            type="info"
            size="small"
            :style="{ backgroundColor: 'var(--surface-tertiary)' }"
          >
            {{ tag }}
          </n-tag>
          <n-tag
            v-if="blog.status !== undefined"
            :type="blog.status === 2 ? 'success' : blog.status === 1 ? 'warning' : 'info'"
            size="small"
          >
            {{ blog.status === 2 ? '已发布' : blog.status === 1 ? '草稿' : '其他' }}
          </n-tag>
        </div>
      </div>

      <!-- 博客内容 -->
      <div class="blog-body" :style="{ backgroundColor: 'var(--surface-primary)' }">
        <MarkdownPreview :text="blog.context" />
      </div>

      <!-- 操作按钮 -->
      <div class="mt-8 flex gap-3">
        <n-button @click="likeBlog" :loading="likeLoading" type="primary">
          <template #icon>
            <Heart :size="18" />
          </template>
          点赞 ({{ blog.like }})
        </n-button>
        <n-button @click="collectBlog" :loading="collectLoading" type="tertiary">
          <template #icon>
            <Folder :size="18" />
          </template>
          {{ isCollected ? '取消收藏' : '收藏' }} ({{ blog.collection }})
        </n-button>
        <n-button @click="shareBlog" type="tertiary">
          <template #icon>
            <Share :size="18" />
          </template>
          分享
        </n-button>
      </div>
    </div>

    <div v-else class="text-center py-12 text-gray-500">
      <h2 class="text-xl mb-2">博客未找到</h2>
      <p>该博客可能已被删除或不存在</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { NButton, NTag, NSpin, useMessage } from 'naive-ui'
import { User, Calendar, Heart, Folder, Share } from 'lucide-vue-next'
import MarkdownPreview from '@/components/MarkdownPreview.vue'
import Request from '@/services/api'
import type { Blog } from '@/types'

const route = useRoute()
const message = useMessage()

// 响应式数据
const blog = ref<Blog | null>(null)
const loading = ref(true)
const likeLoading = ref(false)
const collectLoading = ref(false)
const isCollected = ref(false)

// 获取博客详情
const fetchBlogDetail = async () => {
  try {
    const response = await Request.get(`/blog/${route.params.id}`)
    blog.value = response.info || response

    // 检查是否已收藏（这里可以调用检查收藏状态的API）
    isCollected.value = false // 默认未收藏，实际项目中应从API获取
  } catch (error) {
    console.error('获取博客详情失败:', error)
    message.error('获取博客详情失败')
  } finally {
    loading.value = false
  }
}

// 点赞博客
const likeBlog = async () => {
  if (!blog.value) return

  likeLoading.value = true
  try {
    // 这里调用点赞API
    // await Request.post(`/blog/like/${blog.value.id}`);

    // 临时更新本地数据
    blog.value.like += 1
    message.success('点赞成功')
  } catch (error) {
    console.error('点赞失败:', error)
    message.error('点赞失败')
  } finally {
    likeLoading.value = false
  }
}

// 收藏/取消收藏博客
const collectBlog = async () => {
  if (!blog.value) return

  collectLoading.value = true
  try {
    if (isCollected.value) {
      // 取消收藏API
      // await Request.delete(`/blog/collect/${blog.value.id}`);
      message.success('已取消收藏')
    } else {
      // 收藏API
      // await Request.post(`/blog/collect/${blog.value.id}`);
      message.success('收藏成功')
    }

    // 临时更新本地状态
    isCollected.value = !isCollected.value
    blog.value.collection = isCollected.value
      ? blog.value.collection + 1
      : blog.value.collection - 1
  } catch (error) {
    console.error('收藏失败:', error)
    message.error(isCollected.value ? '取消收藏失败' : '收藏失败')
  } finally {
    collectLoading.value = false
  }
}

// 分享博客
const shareBlog = () => {
  if (navigator.share) {
    navigator
      .share({
        title: blog.value?.title,
        text:
          blog.value?.context.substring(0, 100) + (blog.value?.context.length > 100 ? '...' : ''),
        url: window.location.href
      })
      .catch(console.error)
  } else {
    // 复制链接到剪贴板
    navigator.clipboard
      .writeText(window.location.href)
      .then(() => message.success('链接已复制到剪贴板'))
      .catch(() => message.error('复制失败'))
  }
}

// 页面加载时获取博客详情
onMounted(() => {
  fetchBlogDetail()
})
</script>

<style scoped>
.blog-body {
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 24px;
}

:deep(.v-md-editor) {
  background: transparent;
}

:deep(.v-md-editor .v-md-preview) {
  padding: 0;
}
</style>
