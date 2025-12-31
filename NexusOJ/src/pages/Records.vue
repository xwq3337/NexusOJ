<template>
  <div :style="{ backgroundColor: 'var(--bg-primary)' }" class="min-h-screen p-6">
    <div class="max-w-7xl mx-auto">
      <h1 class="text-3xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">评测记录</h1>

      <!-- 筛选和搜索区域 -->
      <div
        class="mb-6 flex gap-4 p-4 rounded-lg"
        :style="{ backgroundColor: 'var(--surface-primary)' }"
      >
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索题目或用户..."
          style="max-width: 300px"
          :style="{ backgroundColor: 'var(--surface-tertiary)' }"
          clearable
        >
          <template #prefix>
            <Search :size="18" />
          </template>
        </n-input>

        <n-select
          v-model:value="statusFilter"
          :options="statusOptions"
          placeholder="状态"
          style="min-width: 120px"
          clearable
        />

        <n-select
          v-model:value="languageFilter"
          :options="languageOptions"
          placeholder="语言"
          style="min-width: 120px"
          clearable
        />

        <n-button @click="resetFilters" type="default">重置</n-button>
      </div>

      <!-- 记录表格 -->
      <n-card :style="{ backgroundColor: 'var(--surface-primary)' }" content-style="padding: 0;">
        <n-data-table
          :columns="columns"
          :data="filteredRecords"
          :loading="loading"
          :pagination="pagination"
          remote
          :row-key="(row) => row.id"
        />
      </n-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted } from 'vue'
import { NDataTable, NCard, NInput, NSelect, NButton, NTag, NAvatar } from 'naive-ui'
import { Search, CheckCircle, XCircle, Clock, Code, User, Calendar } from 'lucide-vue-next'
import { h } from 'vue'
import Request from '@/services/api'

// 评测记录数据结构
interface Record {
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
onMounted(async () => {
  await Request.get('/record/list').then((res) => {
    records.value = res.info
  })
})

// 状态选项
const statusOptions = [
  { label: '通过', value: 'Accepted' },
  { label: '答案错误', value: 'WrongAnswer' },
  { label: '超时', value: 'TimeLimitExceeded' },
  { label: '内存超限', value: 'MemoryLimitExceeded' },
  { label: '运行错误', value: 'RuntimeError' },
  { label: '编译错误', value: 'CompilationError' }
]

// 语言选项
const languageOptions = [
  { label: 'C++', value: 'cpp' },
  { label: 'Python', value: 'python' },
  { label: 'Java', value: 'java' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'C#', value: 'csharp' }
]

// 筛选条件
const searchKeyword = ref('')
const statusFilter = ref<string | null>(null)
const languageFilter = ref<string | null>(null)

// 加载状态
const loading = ref(false)

// 记录数据
const records = ref<Record[]>([])

// 计算筛选后的记录
const filteredRecords = computed(() => {
  return records.value.filter((record) => {
    // 搜索关键词筛选
    const matchesSearch =
      !searchKeyword.value ||
      record.problem_title.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      record.username.toLowerCase().includes(searchKeyword.value.toLowerCase())

    // 状态筛选
    const matchesStatus = !statusFilter.value || record.verdict === statusFilter.value

    // 语言筛选
    const matchesLanguage = !languageFilter.value || record.language === languageFilter.value

    return matchesSearch && matchesStatus && matchesLanguage
  })
})

// 重置筛选条件
const resetFilters = () => {
  searchKeyword.value = ''
  statusFilter.value = null
  languageFilter.value = null
}

// 获取状态标签类型
const getStatusType = (verdict: string) => {
  switch (verdict) {
    case 'Accepted':
      return 'success'
    case 'WrongAnswer':
    case 'RuntimeError':
    case 'TimeLimitExceeded':
    case 'MemoryLimitExceeded':
      return 'error'
    case 'CompilationError':
      return 'warning'
    default:
      return 'default'
  }
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

// 获取语言标签
const getLanguageLabel = (lang: string) => {
  switch (lang) {
    case 'cpp':
      return 'C++'
    case 'python':
      return 'Python'
    case 'java':
      return 'Java'
    case 'javascript':
      return 'JavaScript'
    case 'csharp':
      return 'C#'
    default:
      return lang
  }
}

// 获取语言类型
const getLanguageType = (lang: string) => {
  switch (lang) {
    case 'cpp':
    case 'csharp':
      return 'primary'
    case 'python':
      return 'success'
    case 'java':
    case 'rust':
      return 'warning'
    case 'go':
    case 'javascript':
      return 'info'
    default:
      return 'default'
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 格式化内存
const formatMemory = (memory: number) => {
  return `${memory} MB`
}

// 格式化时间
const formatTime = (time: number) => {
  return `${time} ms`
}
import { useRouter } from 'vue-router'
const router = useRouter()
// 表格列定义
const columns = [
  {
    title: 'ID',
    key: 'id',
    width: 150,
    render(row: Record) {
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
    render(row: Record) {
      return h('div', [
        h(
          'div',
          { style: { fontWeight: 'bold', color: 'var(--text-primary)' } },
          row.problem_title
        ),
        h(
          'div',
          { style: { fontSize: '12px', color: 'var(--text-secondary)' } },
          `ID: ${row.problem_id}`
        )
      ])
    }
  },
  {
    title: '用户',
    key: 'username',
    width: 120,
    render(row: Record) {
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
    render(row: Record) {
      return h(
        NTag,
        {
          type: getStatusType(row.verdict),
          size: 'small',
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
    render(row: Record) {
      return h(
        NTag,
        {
          type: getLanguageType(row.language),
          size: 'small',
          style: { margin: '2px' }
        },
        {
          default: () => getLanguageLabel(row.language)
        }
      )
    }
  },
  {
    title: '内存',
    key: 'max_memory',
    width: 100,
    render(row: Record) {
      return formatMemory(row.max_memory)
    }
  },
  {
    title: '时间',
    key: 'max_time',
    width: 100,
    render(row: Record) {
      return formatTime(row.max_time)
    }
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 180,
    render(row: Record) {
      return formatDate(row.created_at)
    }
  }
]

const currentPage = ref(1)
const pageSize = ref(10)

// 分页配置
const pagination = reactive({
  page: currentPage.value,
  pageSize: pageSize.value,
  showSizePicker: true,
  pageSizes: [5, 10, 20, 50],
  itemCount: 0,
  onUpdatePage: (page: number) => {
    currentPage.value = page
  },
  onUpdatePageSize: (size: number) => {
    pageSize.value = size
    currentPage.value = 1
  }
})

// 监听currentPage和pageSize的变化，更新pagination
watch(currentPage, () => {
  pagination.page = currentPage.value
})

watch(pageSize, () => {
  pagination.pageSize = pageSize.value
})

// 监听filteredRecords长度变化，更新itemCount
watch(
  () => filteredRecords.value.length,
  (newLength) => {
    pagination.itemCount = newLength
  }
)
</script>
