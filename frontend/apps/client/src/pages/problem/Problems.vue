<template>
  <div class="animate-fade-in max-w-6xl mx-auto">
    <div class="mb-8 whitespace-nowrap">
      <!-- TODO:搜索 -->
      <div class="divide-y" :style="{
        borderBottomColor: 'var(--border-color)',
        borderBottomWidth: '1px',
        borderStyle: 'solid'
      }">
        <n-data-table :columns="columns" :data="Problems" :pagination="pagination" :bordered="false">
        </n-data-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, renderList, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { CheckCircle2, Circle, AlertCircle, ArrowRight } from 'lucide-vue-next'
import { NTag, NDataTable, useMessage, NTab, NCol } from 'naive-ui'
const message = useMessage()
import { difficultyMap } from '@/constants'
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
        `${formatAcceptance(row.accept, row.submission)}`
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

const Problems = ref<ProblemListDTO[]>([])
import { formatAcceptance } from '@/utils/format'
import { problemApi } from '@nexusoj/server'
import { ProblemListDTO } from '@nexusoj/type'
onMounted(async () => {
  await problemApi.getProblemList()
    .then((res) => {
      Problems.value = res.info
    })
    .catch((e) => { })
})
</script>
