<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { NCard, NTable, NTag, NEllipsis, NCollapse, NCollapseItem, NCode, NEmpty, NSpin, NPagination } from 'naive-ui'
import Request from '@/services/api'
import { useUserStore } from '@/stores/useUserStore'
const { id } = useUserStore()
// 测试用例数据结构
interface TestCase {
    time: number
    stdin: string
    memory: number
    status: string
    stderr: string
    stdout: string
    case_id: string
    expected: string
}

// 提交记录数据结构
interface JudgeRecord {
    id: number
    code: string
    created_at: string
    deleted_at: string | null
    judge_result: string // JSON字符串，包含 TestCase 数组
    language: string
    max_memory: number
    max_time: number
    problem_id: string
    problem_title: string
    updated_at: string
    user_id: string
    username: string
    verdict: string
}
import { ApiResponse } from '@/types/api'
const route = useRoute()

const records = ref<JudgeRecord[]>([])
const loading = ref(true)
const currentPage = ref(1)
const pageSize = 20

// 获取提交记录
const fetchRecords = async () => {
    loading.value = true
    try {
        const response = await Request.get<ApiResponse>(`/record/user/${id}`)
        if (response.code === 200) {
            records.value = response.info
        }
    } catch (error) {
        console.error('Failed to fetch submission records:', error)
    } finally {
        loading.value = false
    }
}

// 解析测试用例结果
const parseJudgeResult = (result: string): TestCase[] => {
    try {
        return JSON.parse(result)
    } catch {
        return []
    }
}

// 获取评测结果对应的标签类型
const getVerdictType = (verdict: string) => {
    const verdictMap: Record<string, 'success' | 'error' | 'warning' | 'default'> = {
        'Accepted': 'success',
        'Wrong Answer': 'error',
        'Time Limit Exceeded': 'error',
        'Memory Limit Exceeded': 'error',
        'Runtime Error': 'error',
        'Compilation Error': 'error',
        'Presentation Error': 'warning',
        'System Error': 'error',
    }
    return verdictMap[verdict] || 'default'
}

// 获取测试用例状态对应的标签类型
const getTestCaseStatusType = (status: string) => {
    const statusMap: Record<string, 'success' | 'error' | 'warning' | 'default'> = {
        'AC': 'success',
        'WA': 'error',
        'TLE': 'error',
        'MLE': 'error',
        'RE': 'error',
        'CE': 'error',
        'PE': 'warning',
    }
    return statusMap[status] || 'default'
}

import { formatMemory, formatDate, formatTime } from '@/utils/format'
// 分页数据
const paginatedRecords = computed(() => {
    const start = (currentPage.value - 1) * pageSize
    const end = start + pageSize
    return records.value.slice(start, end)
})

const totalPages = computed(() => Math.ceil(records.value.length / pageSize))

onMounted(() => {
    fetchRecords()
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
                <NCard
                    v-for="record in paginatedRecords"
                    :key="record.id"
                    class="hover:shadow-md transition-shadow cursor-pointer"
                    :style="{ backgroundColor: 'var(--card-bg)' }"
                >
                    <div class="flex items-center justify-between gap-4">
                        <div class="flex items-center gap-3 flex-1 min-w-0">
                            <NTag :type="getVerdictType(record.verdict)" size="small" :bordered="false">
                                {{ record.verdict }}
                            </NTag>
                            <span class="font-medium text-base truncate" :style="{ color: 'var(--text-primary)' }">
                                {{ record.problem_title }}
                            </span>
                        </div>
                        <div class="flex items-center gap-4 text-xs shrink-0" :style="{ color: 'var(--text-secondary)' }">
                            <span class="whitespace-nowrap">{{ formatDate(record.created_at) }}</span>
                            <span class="language whitespace-nowrap">{{ record.language }}</span>
                            <span class="time whitespace-nowrap">{{ formatTime(record.max_time) }}</span>
                            <span class="memory whitespace-nowrap">{{ formatMemory(record.max_memory) }}</span>
                        </div>
                    </div>
                </NCard>
                <!-- 分页 -->
                <div v-if="totalPages > 1" class="flex justify-center mt-8 py-4">
                    <NPagination v-model:page="currentPage" :page-count="totalPages" :page-size="pageSize"
                        show-quick-jumper />
                </div>
            </div>
        </NSpin>
    </div>
</template>
