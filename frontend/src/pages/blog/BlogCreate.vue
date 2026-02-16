<template>
  <div class="max-w-5xl mx-auto p-4 md:p-6">
    <!-- 标题栏 -->
    <n-card :bordered="false" class="mb-4">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">创作博客</h1>
        <div class="flex items-center gap-2">
          <n-button @click="handleSaveDraft" :loading="loading">
            <template #icon>
              <Save :size="16" />
            </template>
            保存草稿
          </n-button>
          <n-button type="primary" @click="handlePublish" :loading="loading">
            <template #icon>
              <Send :size="16" />
            </template>
            发布
          </n-button>
        </div>
      </div>

      <!-- 标题输入 -->
      <n-input v-model:value="blogForm.title" placeholder="请输入文章标题（必填）" size="large" :maxlength="100" show-count
        class="mb-4" @blur="validateTitle" />

      <!-- 标签输入 -->
      <div class="mb-4">
        <div class="flex items-center gap-2 mb-2">
          <Tag :size="16" :style="{ color: 'var(--text-secondary)' }" />
          <span class="text-sm" :style="{ color: 'var(--text-secondary)' }">标签（用回车分隔）</span>
        </div>
        <n-dynamic-tags v-model:value="blogForm.tags" :max="10" />
      </div>

      <!-- 设置选项 -->
      <div class="flex items-center gap-6">
        <n-checkbox v-model:checked="blogForm.is_private">
          <span :style="{ color: 'var(--text-primary)' }">私密文章</span>
        </n-checkbox>
      </div>
    </n-card>

    <!-- 编辑器 -->
    <n-card :bordered="false">
      <v-md-editor v-model="blogForm.context" height="600px"
        left-toolbar="undo redo clear | h bold italic strikethrough quote | ul ol table hr | link image code" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { Save, Send, Tag } from 'lucide-vue-next'
import { NCard, NInput, NButton, NCheckbox, NDynamicTags } from 'naive-ui'
import type { Blog } from '@/types/blog'
import { createBlogParam } from '@/types/blog'
import { blogApi } from '@/services/blog'

const router = useRouter()
const message = useMessage()

// 表单数据
const blogForm = ref<createBlogParam>({
  title: '',
  context: '',
  tags: [],
  status: "Pending",
  is_private: false
})

const loading = ref(false)
const isEdit = ref(false)
const blogId = ref<string>('1')

// 验证标题
const validateTitle = () => {
  if (!blogForm.value.title.trim()) {
    message.warning('请输入文章标题')
    return false
  }
  if (blogForm.value.title.trim().length < 2) {
    message.warning('标题至少需要2个字符')
    return false
  }
  return true
}

// 验证内容
const validateContent = () => {
  if (!blogForm.value.context.trim()) {
    message.warning('请输入文章内容')
    return false
  }
  if (blogForm.value.context.trim().length < 10) {
    message.warning('内容至少需要10个字符')
    return false
  }
  return true
}

// 保存草稿
const handleSaveDraft = async () => {
  if (!validateTitle()) return

  loading.value = true
  try {
    const data: createBlogParam = {
      title: blogForm.value.title.trim(),
      context: blogForm.value.context,
      tags: blogForm.value.tags,
      is_private: blogForm.value.is_private,
      status: "Draft"
    }

    if (isEdit.value && blogId.value) {
      await blogApi.updateBlog({
        id: blogId.value,
        ...data
      })
      message.success('草稿保存成功')
    } else {
      const response = await blogApi.createBlog(data)
      message.success('草稿保存成功')
      if (response.info.id) {
        blogId.value = response.info.id
        isEdit.value = true
      }
    }
  } catch (error: any) {
    console.error('保存草稿失败:', error)
    message.error(error.message || '保存草稿失败，请重试')
  } finally {
    loading.value = false
  }
}

// 发布文章
const handlePublish = async () => {
  if (!validateTitle()) return
  if (!validateContent()) return

  loading.value = true
  try {
    const data: createBlogParam = {
      title: blogForm.value.title.trim(),
      context: blogForm.value.context,
      tags: blogForm.value.tags,
      is_private: blogForm.value.is_private,
      status: "Pending"
    }

    if (isEdit.value && blogId.value) {
      await blogApi.updateBlog({
        id: blogId.value,
        ...data
      })
      message.success('文章发布成功')
    } else {
      const response = await blogApi.createBlog(data)
      message.success('文章发布成功')
      blogId.value = response.info?.id || ''
    }

    // 跳转到文章详情页
    setTimeout(() => {
      router.push({ name: 'BlogDetail', params: { id: blogId.value } })
    }, 500)
  } catch (error: any) {
    console.error('发布失败:', error)
    message.error(error.message || '发布失败，请重试')
  } finally {
    loading.value = false
  }
}

// 加载博客数据（编辑模式）
const loadBlogData = async (id: string) => {
  try {
    const response = await blogApi.getBlogById(id)
    const blog: Blog = response.info
    const { title, context, tags, is_private, status } = blog
    blogForm.value = { title, context, tags: tags ?? [], is_private, status}
    blogId.value = blog.id
    isEdit.value = true
  } catch (error: any) {
    console.error('加载博客数据失败:', error)
    message.error('加载博客数据失败')
    router.push({ name: 'Blogs' })
  }
}

// 初始化
onMounted(() => {
  // 检查是否是编辑模式
  const blogIdParam = router.currentRoute.value.params.id as string
  if (blogIdParam) {
    loadBlogData(blogIdParam)
  }
})
</script>

<style scoped>
/* Markdown 编辑器样式覆盖 */
:deep(.v-md-editor) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.v-md-editor--fullscreen) {
  z-index: 9999;
}
</style>
