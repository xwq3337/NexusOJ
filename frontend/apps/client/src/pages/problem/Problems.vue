<template>
  <div class="animate-fade-in max-w-6xl mx-auto">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Left: Problem List -->
      <div class="lg:col-span-8 mb-8 whitespace-nowrap">
        <div class="mb-4 flex items-center gap-3">
          <n-input v-model:value="searchKeyword" placeholder="搜索题目标题或 ID..." clearable @keyup.enter="handleSearch">
            <template #prefix>
              <Search class="w-4 h-4" style="color: var(--text-secondary)" />
            </template>
          </n-input>
          <n-select v-model:value="selectedTags" :options="tagOptions" multiple filterable placeholder="按标签筛选..."
            :max-tag-count="3" @update:value="handleTagChange" />
          <n-button type="primary" @click="handleSearch">搜索</n-button>
        </div>

        <div class="divide-y" :style="{
          borderBottomColor: 'var(--border-color)',
          borderBottomWidth: '1px',
          borderStyle: 'solid'
        }">
          <n-data-table :columns="columns" :data="Problems" :pagination="paginationReactive" :loading="loading"
            :bordered="false" remote @update:page="handlePageChange" @update:page-size="handlePageSizeChange" />
        </div>
      </div>

      <!-- Right: Recommendations -->
      <div class="lg:col-span-4">
        <div class="lg:sticky lg:top-24">
          <UserRecommendations />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h, renderList, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { CheckCircle2, Circle, AlertCircle, Search } from 'lucide-vue-next'
import { NTag, NDataTable, NInput, NButton, NSelect, useMessage } from 'naive-ui'
import { difficultyMap } from '@/constants'
import { formatAcceptance } from '@/utils/format'
import { problemApi } from '@nexusoj/server'
import type { ProblemListDTO } from '@nexusoj/type'
import UserRecommendations from '@/pages/user/components/UserRecommendations.vue'

const message = useMessage()
const loading = ref(false)
const Problems = ref<ProblemListDTO[]>([])
const searchKeyword = ref('')
const selectedTags = ref<string[]>([])
const allTags = ref<string[]>([])

const tagOptions = computed(() =>
  allTags.value.map(t => ({ label: t, value: t }))
)

const paginationReactive = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
})

const columns = [
  {
    title: '状态',
    key: 'status',
    width: 60,
    render(row: any) {
      let icon
      let iconClass
      let statusText

      switch (row.status) {
        case 'solved':
          icon = CheckCircle2
          iconClass = 'text-green-500'
          statusText = '已解决'
          break
        case 'attempted':
          icon = AlertCircle
          iconClass = 'text-yellow-500'
          statusText = '尝试过'
          break
        case 'unattempted':
        default:
          icon = Circle
          iconClass = 'text-gray-600'
          statusText = '未尝试'
          break
      }

      return h(icon, {
        size: 18,
        class: iconClass,
        onMouseenter: () => {
          // TODO 展示已解决/尝试过/未尝试 tooltip
          console.log(statusText)
        }
      })
    }
  },
  {
    title: '题目',
    key: 'title',
    render(row: any) {
      return h('div', { class: 'items-center space-x-2' }, [
        h(
          RouterLink,
          {
            to: `/problem/${row.id}`,
            class: 'font-medium hover:text-sky-400 transition-colors text-sm block',
            style: { color: 'var(--text-primary)' }
          },
          { default: () => `${row.id}. ${row.title}` }
        ),
        h(
          'div',
          { class: 'mt-1 flex gap-1' },
          renderList(row.tags, (tag) =>
            h(
              NTag,
              {
                type: 'success',
                size: 'small',
                class: 'text-xs'
              },
              { default: () => tag }
            )
          )
        )
      ])
    }
  },
  {
    title: '通过率',
    key: 'acceptance',
    render(row: any) {
      return h(
        'span',
        { class: 'text-sm' },
        `${formatAcceptance(row.accept, row.submission)}`
      )
    }
  },
  {
    title: '难度',
    key: 'difficulty',
    render(row: any) {
      return h(
        NTag,
        {
          type: difficultyMap[Number(row.difficulty) - 1]?.type as
            | 'success'
            | 'warning'
            | 'info'
            | 'error'
            | 'default',
          size: 'small'
        },
        { default: () => difficultyMap[Number(row.difficulty) - 1]?.text }
      )
    }
  }
]

const fetchProblems = async () => {
  loading.value = true
  try {
    const res = await problemApi.getProblemList(
      paginationReactive.page,
      paginationReactive.pageSize,
      searchKeyword.value || undefined,
      selectedTags.value.length > 0 ? selectedTags.value : undefined,
    )
    if (res.code === 200 && res.info) {
      Problems.value = res.info.problems || []
      paginationReactive.itemCount = res.info.total || 0
    }
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchTags = async () => {
  try {
    const res = await problemApi.getAllTags()
    if (res.code === 200 && res.info) {
      allTags.value = res.info
    }
  } catch (e) {
    console.error(e)
  }
}

const handlePageChange = (page: number) => {
  paginationReactive.page = page
  fetchProblems()
}

const handlePageSizeChange = (pageSize: number) => {
  paginationReactive.pageSize = pageSize
  paginationReactive.page = 1
  fetchProblems()
}

const handleTagChange = () => {
  paginationReactive.page = 1
  fetchProblems()
}

const handleSearch = () => {
  paginationReactive.page = 1
  fetchProblems()
}

onMounted(() => {
  fetchProblems()
  fetchTags()
})
</script>
