<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { useRoute } from 'vue-router'
import { NTag, NEmpty, NSpin, NVirtualList, useMessage } from 'naive-ui'
import { useUserStore } from '@/stores/useUserStore'
import { formatMemory, formatDate, formatTime } from '@/utils/format'
import { GetRecordListResponse } from '@nexusoj/type'
import { LANGUAGE_CONFIG, STATUS_COLORS } from '@/constants'
import router from '@/router'
const { id } = useUserStore()
import { useClipboard } from '@vueuse/core'
import { userApi } from '@nexusoj/server'
const { copy } = useClipboard()
const route = useRoute()
const message = useMessage()

const PAGE_SIZE = 20
const itemSize = 120 // 每条记录的高度
const records = ref<GetRecordListResponse[]>([])
const loading = ref(false)
const loadingMore = ref(false)
const hasMore = ref(true)
const currentPage = ref(1)

const fetchRecords = async (isLoadMore = false) => {
  if (isLoadMore) {
    loadingMore.value = true
  } else {
    loading.value = true
  }

  try {
    const response = await userApi.getUserRecordList(id, {
      page: currentPage.value,
      page_size: PAGE_SIZE
    })

    if (response.code === 200) {
      if (isLoadMore) {
        records.value = [...records.value, ...response.info]
      } else {
        records.value = response.info
      }

      // 判断是否还有更多数据
      if (response.info.length < PAGE_SIZE) {
        hasMore.value = false
      }
    }
  } catch (error) {
    console.error('Failed to fetch submission records:', error)
    message.error('加载失败，请重试')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

const loadMore = () => {
  if (loadingMore.value || !hasMore.value) return
  currentPage.value++
  fetchRecords(true)
}

const handleScroll = (e: Event) => {
  const target = e.target as HTMLElement
  const scrollBottom = target.scrollHeight - target.scrollTop - target.clientHeight

  // 当滚动到距离底部 100px 时加载更多
  if (scrollBottom < 100 && hasMore.value && !loadingMore.value) {
    loadMore()
  }
}

const listHeight = computed(() => {
  const baseHeight = records.value.length * itemSize
  return Math.min(baseHeight, 600) + 'px'
})

// 渲染单条记录
const renderItem = (record: GetRecordListResponse) => {
  return h(
    'div',
    {
      class: 'p-4 border-b border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer',
      onClick: () => {
        router.push({ name: 'RecordDetail', params: { id: record.id } })
      }
    },
    [
      // 顶部：题目和状态
      h('div', { class: 'flex items-start justify-between gap-4 mb-3' }, [
        h('div', { class: 'flex-1 min-w-0' }, [
          h(
            'div',
            {
              class: 'font-semibold text-base truncate',
              style: { color: 'var(--text-primary)' },
              onClick: (e: Event) => {
                e.stopPropagation()
                router.push({
                  name: 'ProblemDetail',
                  params: { id: record.problem_id }
                })
              }
            },
            record.problem_title
          )
        ]),
        h(
          NTag,
          {
            size: 'small',
            color: STATUS_COLORS[record.verdict],
            bordered: false
          },
          { default: () => record.verdict }
        )
      ]),

      // 底部：语言、内存、时间、提交时间
      h('div', { class: 'flex flex-wrap items-center gap-3 text-xs' }, [
        h(
          NTag,
          {
            size: 'small',
            color: LANGUAGE_CONFIG[record.language].color
          },
          { default: () => LANGUAGE_CONFIG[record.language].label }
        ),
        h(
          'span',
          { style: { color: 'var(--text-secondary)' } },
          `内存: ${formatMemory(record.max_memory)}`
        ),
        h(
          'span',
          { style: { color: 'var(--text-secondary)' } },
          `时间: ${formatTime(record.max_time)}`
        ),
        h(
          'span',
          { style: { color: 'var(--text-secondary)' } },
          formatDate(record.created_at)
        )
      ])
    ]
  )
}

onMounted(() => {
  fetchRecords()
})
</script>

<template>
  <div class="p-4 md:p-6">
    <div class="mb-6">
      <h1 class="text-2xl md:text-3xl font-semibold mb-2" :style="{ color: 'var(--text-primary)' }">
        提交记录
      </h1>
      <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">
        共 {{ records.length }} 条记录
      </p>
    </div>

    <NSpin :show="loading">
      <div v-if="records.length === 0 && !loading" class="py-15 text-center">
        <NEmpty description="暂无提交记录" />
      </div>

      <div v-else class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden">
        <n-virtual-list
          :item-size="itemSize"
          :items="records"
          :item-resizable="false"
          :style="{ height: listHeight }"
          @scroll="handleScroll"
        >
          <template #default="{ item }">
            <component :is="renderItem(item)" />
          </template>
        </n-virtual-list>

        <!-- 加载更多提示 -->
        <div v-if="loadingMore" class="p-4 text-center border-t border-gray-200 dark:border-gray-700">
          <NSpin size="small" />
          <span class="ml-2 text-sm" :style="{ color: 'var(--text-secondary)' }">加载中...</span>
        </div>

        <div v-else-if="!hasMore && records.length > 0" class="p-4 text-center text-sm border-t border-gray-200 dark:border-gray-700" :style="{ color: 'var(--text-secondary)' }">
          已加载全部记录
        </div>
      </div>
    </NSpin>
  </div>
</template>
