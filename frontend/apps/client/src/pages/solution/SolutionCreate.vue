<template>
  <div class="max-w-5xl mx-auto p-4 md:p-6">
    <n-card :bordered="false" class="mb-4">
      <div class="flex items-center justify-between mb-4">
        <h1 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">创作题解</h1>
        <div class="flex items-center gap-2">
          <n-button @click="handleSaveDraft" :loading="loading">
            <template #icon><Save :size="16" /></template>
            保存草稿
          </n-button>
          <n-button type="primary" @click="handlePublish" :loading="loading">
            <template #icon><Send :size="16" /></template>
            发布
          </n-button>
        </div>
      </div>

      <n-input v-model:value="form.title" placeholder="请输入题解标题（必填）" size="large" :maxlength="100" show-count
        class="mb-4" />
      <n-input v-model:value="form.excerpt" placeholder="请输入题解摘要" size="large" :maxlength="200" show-count
        class="mb-4" />
      <div class="mb-4">
        <div class="flex items-center gap-2 mb-2">
          <TagIcon :size="16" :style="{ color: 'var(--text-secondary)' }" />
          <span class="text-sm" :style="{ color: 'var(--text-secondary)' }">标签（用回车分隔）</span>
        </div>
        <n-dynamic-tags v-model:value="form.tags" :max="10" />
      </div>
    </n-card>

    <n-card :bordered="false">
      <v-md-editor v-model="form.context" height="600px"
        left-toolbar="undo redo clear | h bold italic strikethrough quote | ul ol table hr | link image code" />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage } from 'naive-ui'
import { NCard, NInput, NButton, NDynamicTags } from 'naive-ui'
import { Save, Send, Tag as TagIcon } from 'lucide-vue-next'
import { solutionApi } from '@nexusoj/server'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const problemId = route.query.problem_id as string
const problemTitle = route.query.problem_title as string

const form = ref({
  title: problemTitle ? `${problemTitle} 题解` : '',
  excerpt: '',
  context: '',
  tags: [] as string[],
  status: 'public' as string,
})

const loading = ref(false)

const handleSaveDraft = async () => {
  if (!form.value.title.trim()) {
    message.warning('请输入题解标题')
    return
  }
  loading.value = true
  try {
    await solutionApi.createSolution({
      problem_id: Number(problemId),
      title: form.value.title.trim(),
      excerpt: form.value.excerpt,
      context: form.value.context,
      tags: form.value.tags,
      status: 'draft',
    })
    message.success('草稿保存成功')
  } catch (e: any) {
    message.error(e.message || '保存失败')
  } finally {
    loading.value = false
  }
}

const handlePublish = async () => {
  if (!form.value.title.trim()) {
    message.warning('请输入题解标题')
    return
  }
  if (!form.value.context.trim()) {
    message.warning('请输入题解内容')
    return
  }
  loading.value = true
  try {
    const { code, info } = await solutionApi.createSolution({
      problem_id: Number(problemId),
      title: form.value.title.trim(),
      excerpt: form.value.excerpt,
      context: form.value.context,
      tags: form.value.tags,
      status: 'public',
    })
    if (code === 200) {
      message.success('题解发布成功')
      router.back()
    }
  } catch (e: any) {
    message.error(e.message || '发布失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
:deep(.v-md-editor) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.v-md-editor--fullscreen) {
  z-index: 9999;
}
</style>
