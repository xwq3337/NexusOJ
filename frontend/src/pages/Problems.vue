<template>
  <div class="animate-fade-in max-w-6xl mx-auto">
    <div class="mb-8">
      <!-- TODO:搜索 -->
      <div
        class="divide-y"
        :style="{
          borderBottomColor: 'var(--border-color)',
          borderBottomWidth: '1px',
          borderStyle: 'solid'
        }"
      >
        <n-data-table
          :columns="columns"
          :data="Problems"
          :pagination="pagination"
          :bordered="false"
        >
        </n-data-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, renderList, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, Filter, CheckCircle2, Circle, ArrowRight } from 'lucide-vue-next'
import { NTag, NSpace, NDataTable, NButton, useMessage, NTab, NCol } from 'naive-ui'
const message = useMessage()
const pagination = {
  page: 1,
  pageSize: 10
}
const columns = [
  {
    title: '状态',
    key: 'status',
    width: 60,
    render(row) {
      return h(row.status === 'solved' ? CheckCircle2 : Circle, {
        size: 18,
        class: row.status === 'solved' ? 'text-green-500' : 'text-gray-600',
        onMouseenter: () => {
          // TODO 展示已解决/未解决 tooltip
          console.log(row.status === 'solved' ? '已解决' : '未解决')
        }
      })
    }
  },
  {
    title: '题目',
    key: 'title',
    render(row) {
      return h('div', { class: 'items-center space-x-2' }, [
        h(
          RouterLink,
          {
            to: `/problem/${row.id}`,
            class: 'font-medium hover:text-blue-400 transition-colors text-sm block',
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
    render(row) {
      return h(
        'span',
        { class: 'text-sm' },
        `${row.submission == 0 ? 0 : (row.accept / row.submission) * 100}%`
      )
    }
  },
  {
    title: '难度',
    key: 'difficulty',
    render(row) {
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
const difficultyMap = [
  { text: '简单', color: 'text-green-400', type: 'success' },
  { text: '容易', color: 'text-yellow-400', type: 'warning' },
  { text: '中等', color: 'text-orange-400', type: 'info' },
  { text: '困难', color: 'text-red-400', type: 'error' },
  { text: '极难', color: 'text-purple-400', type: 'error' }
]
const Problems = ref([])
import Request from '@/services/api/index'
onMounted(async () => {
  console.log('Fetching problem list...')
  await Request.get('/problem/list')
    .then((res) => {
      Problems.value = res.info
    })
    .catch((e) => {})
})
</script>
