<script setup lang="ts">
import { ref, onMounted, reactive, h, watch } from 'vue'
import { useRoute } from 'vue-router'
import { NCard,NDataTable, NTag, NEmpty, NSpin, NPagination, useMessage } from 'naive-ui'
import Request from '@/services/api'
import { useUserStore } from '@/stores/useUserStore'
import { formatMemory, formatDate, formatTime } from '@/utils/format'
import { GetRecordListResponse } from '@/types/record'
import { LANGUAGE_CONFIG, STATUS_COLORS } from '@/constants'
import router from '@/router'
const { id } = useUserStore()
import { ApiResponse } from '@/types/api'
import { useClipboard } from '@vueuse/core'
import { userApi } from '@/services/user'
const { copy } = useClipboard()
const route = useRoute()
const message = useMessage()
const fetchRecords = async () => {
    loading.value = true
    try {
        const response = await userApi.getUserRecordList(id, {
            page: currentPage.value,
            page_size: pageSize.value
        })
        if (response.code === 200) {
            records.value = response.info
        }
    } catch (error) {
        console.error('Failed to fetch submission records:', error)
    } finally {
        loading.value = false
    }
}

const columns = [
    {
        title: 'ID',
        key: 'id',
        width: 150,
        render(row: GetRecordListResponse) {
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
        render(row: GetRecordListResponse) {
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
                        onClick: () => { copy(String(row.problem_id)), message.success('复制成功') },
                    },
                    `ID: ${row.problem_id}`
                )
            ])
        }
    },
    {
        title: '状态',
        key: 'verdict',
        width: 150,
        render(row: GetRecordListResponse) {
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
        render(row: GetRecordListResponse) {
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
        render(row: GetRecordListResponse) {
            return formatMemory(row.max_memory)
        }
    },
    {
        title: '时间',
        key: 'max_time',
        width: 100,
        render(row: GetRecordListResponse) {
            return formatTime(row.max_time)
        }
    },
    {
        title: '提交时间',
        key: 'created_at',
        width: 180,
        render(row: GetRecordListResponse) {
            return formatDate(row.created_at)
        }
    }
]
const currentPage = ref(1)
const pageSize = ref(10)
const records = ref<GetRecordListResponse[]>([])
const loading = ref(false)
const totalRecords = ref(0)
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
onMounted(() => {
    // 从路由获取初始页码
    const pageParam = router.currentRoute.value.query.page as string
    if (pageParam) {
        currentPage.value = parseInt(pageParam, 10)
    }
    fetchRecords()
})

const updateRouteQuery = () => {
    router.push({
        query: {
            ...router.currentRoute.value.query,
            page: currentPage.value.toString()
        }
    })
}

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

</script>

<template>
    <div class="p-6 mx-auto md:p-6">
        <div class="mb-6">
            <h1 class="text-3xl font-semibold mb-2" :style="{ color: 'var(--text-primary)' }">提交记录</h1>
            <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">共 {{ records.length }} 条记录</p>
        </div>

        <NSpin :show="loading">
            <div v-if="records.length === 0 && !loading" class="py-15 text-center">
                <NEmpty description="暂无提交记录" />
            </div>

            <div v-else class="min-h-100 space-y-3">
                <n-card :style="{ backgroundColor: 'var(--surface-primary)' }" content-style="padding: 0;">
                    <n-data-table :columns="columns" :data="records" :loading="loading" size="small"/>
                    <div class="flex justify-end p-4">
                        <n-pagination v-model:page="pagination.page" v-model:page-size="pagination.pageSize"
                            :page-count="Math.ceil(Number(totalRecords / pagination.pageSize))"
                            :page-sizes="pagination.pageSizes" size="medium"
                            :show-size-picker="pagination.showSizePicker" @update-page="pagination.onUpdatePage"
                            @update-page-size="pagination.onUpdatePageSize" />
                    </div>
                </n-card>
            </div>
        </NSpin>
    </div>
</template>
