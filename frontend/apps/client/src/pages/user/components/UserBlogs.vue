<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex items-center gap-3 mb-4 overflow-x-auto">
      <n-select
        v-model:value="filters.category"
        :options="categoryOptions"
        placeholder="分类"
        clearable
        size="small"
        style="min-width: 120px; width: 120px;"
      />
      <n-input
        v-model:value="filters.search"
        placeholder="搜索博客..."
        size="small"
        clearable
        class="flex-1"
        style="min-width: 200px;"
      >
        <template #prefix>
          <Search :size="14" />
        </template>
      </n-input>
    </div>

    <!-- Blogs List -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <n-spin size="large" />
    </div>

    <div v-else-if="blogs.length === 0" class="text-center py-12">
      <FileText :size="48" :style="{ color: 'var(--text-tertiary)' }" class="mx-auto mb-3" />
      <p :style="{ color: 'var(--text-tertiary)' }">暂无博客</p>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="blog in blogs"
        :key="blog.id"
        class="p-5 rounded-xl transition-all duration-200 hover:scale-[1.01] cursor-pointer"
        :style="{
          backgroundColor: 'var(--surface-secondary)',
          border: '1px solid var(--border-color)'
        }"
        @click="handleViewBlog(blog)"
      >
        <!-- Title -->
        <h4 class="text-lg font-semibold mb-2 line-clamp-2" :style="{ color: 'var(--text-primary)' }">
          {{ blog.title }}
        </h4>

        <!-- Summary -->
        <p class="text-sm mb-3 line-clamp-2" :style="{ color: 'var(--text-secondary)' }">
          {{ blog.summary }}
        </p>

        <!-- Tags -->
        <div class="flex flex-wrap gap-2 mb-3">
          <span
            v-for="tag in blog.tags.slice(0, 4)"
            :key="tag"
            class="px-2 py-0.5 rounded text-xs"
            :style="{
              backgroundColor: 'var(--surface-primary)',
              color: 'var(--text-tertiary)'
            }"
          >
            #{{ tag }}
          </span>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-between text-xs" :style="{ color: 'var(--text-tertiary)' }">
          <div class="flex items-center gap-4">
            <span class="flex items-center gap-1">
              <Eye :size="14" />
              {{ blog.views }}
            </span>
            <span class="flex items-center gap-1">
              <ThumbsUp :size="14" />
              {{ blog.likes }}
            </span>
            <span class="flex items-center gap-1">
              <MessageSquare :size="14" />
              {{ blog.comments }}
            </span>
          </div>
          <span>{{ formatRelativeTime(blog.created_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-6">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-size="pageSize"
        show-quick-jumper
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  NSelect,
  NInput,
  NSpin,
  NPagination,
  type SelectOption
} from 'naive-ui'
import {
  Search,
  FileText,
  Eye,
  ThumbsUp,
  MessageSquare
} from 'lucide-vue-next'
import { formatRelativeTime } from '@/utils/format'

const props = defineProps<{
  userId?: Number
}>()

const router = useRouter()

const blogs = ref<any[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const totalPages = ref(0)

const filters = ref({
  category: null as string | null,
  search: ''
})

const categoryOptions: SelectOption[] = [
  { label: '算法', value: 'algorithm' },
  { label: '数据结构', value: 'datastructure' },
  { label: '题解', value: 'solution' },
  { label: '比赛', value: 'contest' },
  { label: '经验分享', value: 'experience' }
]

const fetchBlogs = async () => {
  loading.value = true
  try {
    // TODO: 使用真实的 博客API
    // const params = {
    //   page: currentPage.value,
    //   page_size: pageSize.value,
    //   category: filters.value.category || undefined,
    //   search: filters.value.search || undefined
    // }
    // const response = await blogApi.getUserBlogs(props.userId || '', params)
    // blogs.value = response.data
    // totalPages.value = Math.ceil(response.total / pageSize.value)

    // Mock data
    blogs.value = []
    totalPages.value = 0
  } catch (error) {
    console.error('Failed to fetch blogs:', error)
  } finally {
    loading.value = false
  }
}

const handleViewBlog = (blog: any) => {
  router.push({
    name: 'BlogDetail',
    params: { id: blog.id }
  })
}

watch([currentPage, filters], () => {
  fetchBlogs()
}, { deep: true })

onMounted(() => {
  fetchBlogs()
})
</script>
