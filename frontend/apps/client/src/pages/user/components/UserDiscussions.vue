<template>
  <div class="space-y-4">
    <!-- Filters -->
    <div class="flex items-center gap-3 mb-4 overflow-x-auto">
      <n-select
        v-model:value="filters.status"
        :options="statusOptions"
        placeholder="状态"
        clearable
        size="small"
        style="min-width: 120px; width: 120px;"
      />
      <n-input
        v-model:value="filters.search"
        placeholder="搜索讨论..."
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

    <!-- Discussions List -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <n-spin size="large" />
    </div>

    <div v-else-if="discussions.length === 0" class="text-center py-12">
      <MessagesSquare :size="48" :style="{ color: 'var(--text-tertiary)' }" class="mx-auto mb-3" />
      <p :style="{ color: 'var(--text-tertiary)' }">暂无讨论</p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="discussion in discussions"
        :key="discussion.id"
        class="p-4 rounded-xl transition-all duration-200 hover:scale-[1.01] cursor-pointer"
        :style="{
          backgroundColor: 'var(--surface-secondary)',
          border: '1px solid var(--border-color)'
        }"
        @click="handleViewDiscussion(discussion)"
      >
        <div class="flex items-start gap-3">
          <!-- Avatar -->
          <img
            :src="discussion.authorAvatar || defaultAvatar"
            :alt="discussion.authorName"
            class="w-10 h-10 rounded-full shrink-0"
          />

          <!-- Content -->
          <div class="flex-1 min-w-0">
            <!-- Title -->
            <h4 class="font-medium mb-1 line-clamp-1" :style="{ color: 'var(--text-primary)' }">
              {{ discussion.title }}
            </h4>

            <!-- Preview -->
            <p class="text-sm mb-2 line-clamp-2" :style="{ color: 'var(--text-secondary)' }">
              {{ discussion.content }}
            </p>

            <!-- Footer -->
            <div class="flex items-center justify-between text-xs" :style="{ color: 'var(--text-tertiary)' }">
              <div class="flex items-center gap-3">
                <span>{{ discussion.authorName }}</span>
                <span class="flex items-center gap-1">
                  <MessageSquare :size="12" />
                  {{ discussion.replies }}
                </span>
                <span class="flex items-center gap-1">
                  <Eye :size="12" />
                  {{ discussion.views }}
                </span>
              </div>
              <div class="flex items-center gap-2">
                <span
                  v-if="discussion.solved"
                  class="px-2 py-0.5 rounded text-xs"
                  :style="{
                    backgroundColor: 'rgba(16, 185, 129, 0.1)',
                    color: '#10b981'
                  }"
                >
                  已解决
                </span>
                <span>{{ formatRelativeTime(discussion.created_at) }}</span>
              </div>
            </div>
          </div>
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
import { ref, computed, onMounted, watch } from 'vue'
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
  MessagesSquare,
  MessageSquare,
  Eye
} from 'lucide-vue-next'
import { formatRelativeTime } from '@/utils/format'

const props = defineProps<{
  userId?: Number
}>()

const router = useRouter()

const discussions = ref<any[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(15)
const totalPages = ref(0)

const filters = ref({
  status: null as string | null,
  search: ''
})

const defaultAvatar = computed(() =>
  `https://api.dicebear.com/7.x/avataaars/svg?seed=default`
)

const statusOptions: SelectOption[] = [
  { label: '已解决', value: 'solved' },
  { label: '未解决', value: 'unsolved' },
  { label: '活跃', value: 'active' }
]

const fetchDiscussions = async () => {
  loading.value = true
  try {
    // TODO: 使用真实的 讨论API 
    // const params = {
    //   page: currentPage.value,
    //   page_size: pageSize.value,
    //   status: filters.value.status || undefined,
    //   search: filters.value.search || undefined
    // }
    // const response = await discussionApi.getUserDiscussions(props.userId || '', params)
    // discussions.value = response.data
    // totalPages.value = Math.ceil(response.total / pageSize.value)

    // Mock data
    discussions.value = []
    totalPages.value = 0
  } catch (error) {
    console.error('Failed to fetch discussions:', error)
  } finally {
    loading.value = false
  }
}

const handleViewDiscussion = (discussion: any) => {
  router.push({
    name: 'DiscussionDetail',
    params: { id: discussion.id }
  })
}

watch([currentPage, filters], () => {
  fetchDiscussions()
}, { deep: true })

onMounted(() => {
  fetchDiscussions()
})
</script>
