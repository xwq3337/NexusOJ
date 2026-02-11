<template>
  <div :style="{ backgroundColor: 'var(--bg-primary)' }" class="min-h-screen p-6">
    <div class="max-w-7xl mx-auto">
      <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">评测记录</h1>

      <!-- 筛选和搜索区域 -->
      <div class="mb-6 flex gap-4 p-4 rounded-lg" :style="{ backgroundColor: 'var(--surface-primary)' }">
        <n-input v-model:value="searchKeyword" placeholder="搜索题目或用户..." style="max-width: 300px"
          :style="{ backgroundColor: 'var(--surface-tertiary)' }" clearable>
          <template #prefix>
            <Search :size="18" />
          </template>
        </n-input>

        <n-select v-model:value="statusFilter" :options="STATUS_OPTIONS" placeholder="状态" style="min-width: 120px"
          clearable />

        <n-select v-model:value="languageFilter"
          :options="LANGUAGE_OPTIONS.map(lang => ({ label: LANGUAGE_CONFIG[lang].label, value: lang }))"
          placeholder="语言" style="min-width: 120px" clearable />

        <n-button @click="resetFilters" type="default">重置</n-button>
      </div>

      <!-- 记录表格 -->
      <n-card :style="{ backgroundColor: 'var(--surface-primary)' }" content-style="padding: 0;">
        <n-data-table :columns="columns" :data="records" :loading="loading" size="small" :row-key="(row) => row.id" />
        <div class="flex justify-end p-4">
          <n-pagination v-model:page="pagination.page" v-model:page-size="pagination.pageSize"
            :page-count="Math.ceil(Number(totalRecords / pagination.pageSize))" :page-sizes="pagination.pageSizes"
            size="medium" show-quick-jumper show-size-picker @update-page="pagination.onUpdatePage"
            @update-page-size="pagination.onUpdatePageSize" />
        </div>
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue'
import { useMessage, NDataTable, NCard, NInput, NSelect, NButton, NTag, NAvatar, NPagination } from 'naive-ui'
import { Search, CheckCircle, XCircle, Clock, Code, User, Calendar } from 'lucide-vue-next'
import { h } from 'vue'
import Request from '@/services/api'
import { useRouter } from 'vue-router'
import { useClipboard } from '@vueuse/core'
const { copy } = useClipboard()
const router = useRouter()
const message = useMessage()
// 评测记录数据结构
interface RecordItem {
  id: number
  created_at: string
  language: string
  max_memory: number
  max_time: number
  problem_id: string
  problem_title: string
  user_id: string
  username: string
  verdict: string
}

// 筛选条件
const searchKeyword = ref('')
const statusFilter = ref<string | null>(null)
const languageFilter = ref<string | null>(null)

// 加载状态
const loading = ref(false)

// 记录数据
const records = ref<RecordItem[]>([])

// 总记录数
const totalRecords = ref(0)

// 分页状态
const currentPage = ref(1)
const pageSize = ref(10)

// 获取记录列表
const fetchRecords = async () => {
  loading.value = true
  try {
    const params: { [key: string]: any } = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    // 添加筛选参数
    if (searchKeyword.value) {
      params.search = searchKeyword.value
    }
    if (statusFilter.value) {
      params.verdict = statusFilter.value
    }
    if (languageFilter.value) {
      params.language = languageFilter.value
    }

    await Request.get('/record/list', { params }).then((res) => {
      records.value = res.info.data || []
      totalRecords.value = res.info.total || res.info?.length || 0
    })
  } catch (error) {
    console.error('获取记录失败:', error)
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  // 从路由获取初始页码
  const pageParam = router.currentRoute.value.query.page as string
  if (pageParam) {
    currentPage.value = parseInt(pageParam, 10)
  }
  fetchRecords()
})

// 重置筛选条件
const resetFilters = () => {
  searchKeyword.value = ''
  statusFilter.value = null
  languageFilter.value = null
  currentPage.value = 1
  updateRouteQuery()
  fetchRecords()
}


// 获取状态图标
const getStatusIcon = (verdict: string) => {
  switch (verdict) {
    case 'Accepted':
      return CheckCircle
    case 'WrongAnswer':
    case 'RuntimeError':
    case 'TimeLimitExceeded':
    case 'MemoryLimitExceeded':
    case 'CompilationError':
      return XCircle
    default:
      return Clock
  }
}
import { LANGUAGE_OPTIONS, LANGUAGE_CONFIG, STATUS_OPTIONS, STATUS_COLORS } from '@/constants'
// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 格式化内存
const formatMemory = (memory: number) => {
  return `${Math.round(memory / 1024 / 1024 * 100) / 100} MB`
}

// 格式化时间
const formatTime = (time: number) => {
  return `${time} ms`
}

// 更新路由参数
const updateRouteQuery = () => {
  router.push({
    query: {
      ...router.currentRoute.value.query,
      page: currentPage.value.toString()
    }
  })
}
// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 150,
    render(row: RecordItem) {
      return h(
        'span',
        {
          style: { color: 'var(--text-link)', cursor: 'pointer' },
          onClick: () => {
            router.push({ name: 'RecordDetail', params: { id: row.id } })
          }
        },
        row.id
      )
    }
  },
  {
    title: '题目',
    key: 'problem_title',
    width: 200,
    render(row: RecordItem) {
      return h('div', [
        h(
          'div',
          {
            style: { fontWeight: 'bold', color: 'var(--text-primary)', cursor: 'pointer' },
            onClick: () => router.push({
              name: "ProblemDetail", params: {
                id: row.problem_id
              }
            })
          },
          row.problem_title,
        ),
        h(
          'div',
          {
            style: { fontSize: '12px', color: 'var(--text-secondary)', cursor: 'pointer' },
            onClick: () => { copy(row.problem_id), message.success('复制成功') },
          },
          `ID: ${row.problem_id}`
        )
      ])
    }
  },
  {
    title: '用户',
    key: 'username',
    width: 120,
    render(row: RecordItem) {
      return h('div', [
        h('div', { style: { fontWeight: 'bold', color: 'var(--text-primary)' } }, row.username),
        h(
          'div',
          { style: { fontSize: '12px', color: 'var(--text-secondary)' } },
          `ID: ${row.user_id.substring(0, 8)}...`
        )
      ])
    }
  },
  {
    title: '状态',
    key: 'verdict',
    width: 150,
    render(row: RecordItem) {
      return h(
        NTag,
        {
          size: 'small',
          color: STATUS_COLORS[row.verdict],
          style: { margin: '2px' }
        },
        {
          default: () => row.verdict
        }
      )
    }
  },
  {
    title: '语言',
    key: 'language',
    width: 100,
    render(row: RecordItem) {
      return h(
        NTag,
        {
          size: 'small',
          style: { margin: '2px' },
          color: LANGUAGE_CONFIG[row.language].color
        },
        {
          default: () => LANGUAGE_CONFIG[row.language].label
        }
      )
    }
  },
  {
    title: '内存',
    key: 'max_memory',
    width: 100,
    render(row: RecordItem) {
      return formatMemory(row.max_memory)
    }
  },
  {
    title: '时间',
    key: 'max_time',
    width: 100,
    render(row: RecordItem) {
      return formatTime(row.max_time)
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 180,
    render(row: RecordItem) {
      return formatDate(row.created_at)
    }
  }
]

// 分页配置
const pagination = reactive({
  page: currentPage.value,
  pageSize: pageSize.value,
  showSizePicker: true,
  pageSizes: [5, 10, 20, 50],
  itemCount: 0,
  onUpdatePage: (page: number) => {
    currentPage.value = page
    updateRouteQuery()
    fetchRecords()
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
    updateRouteQuery()
    fetchRecords()
  }
})

// 监听currentPage和pageSize的变化，更新pagination
watch(currentPage, () => {
  pagination.page = currentPage.value
})

watch(pageSize, () => {
  pagination.pageSize = pageSize.value
})

// 监听totalRecords变化，更新itemCount
watch(totalRecords, (newTotal) => {
  pagination.itemCount = newTotal
})

// 监听筛选条件变化，重新获取数据
watch([searchKeyword, statusFilter, languageFilter], () => {
  currentPage.value = 1
  updateRouteQuery()
  fetchRecords()
})
</script>