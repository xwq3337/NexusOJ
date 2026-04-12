<script setup lang="ts">
import { computed, inject } from 'vue'
import { useRouter } from 'vue-router'
import {
  NCard,
  NGrid,
  NGi,
  NStatistic,
  NProgress,
  NList,
  NListItem,
  NTag,
  NSpace,
  NThing,
  NIcon,
  NResult
} from 'naive-ui'
import {
  Users,
  FileCode,
  CheckCircle2,
  XCircle,
  Circle,
  Target,
  Megaphone
} from 'lucide-vue-next'

const router = useRouter()
const { contest , contestId, problems } = inject<any>('contestData')

const formatNumber = (num: number) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

const getProblemStatusBg = (status?: string) => {
  if (status === 'accepted') return 'var(--success-bg)'
  if (status === 'wrong') return 'var(--error-bg)'
  return 'var(--surface-tertiary)'
}

const getProblemStatusColor = (status?: string) => {
  if (status === 'accepted') return 'var(--success-color)'
  if (status === 'wrong') return 'var(--error-color)'
  return 'var(--text-tertiary)'
}

const difficultyDistribution = computed(() => {
  const total = problems.value.length
  if (total === 0) return []
  const easy = problems.value.filter((p: any) => p.difficulty <= 1).length
  const medium = problems.value.filter((p: any) => p.difficulty === 2 || p.difficulty === 3).length
  const hard = problems.value.filter((p: any) => p.difficulty >= 4).length
  return [
    { label: '简单', count: easy, percent: Math.round((easy / total) * 100), color: 'var(--success-color)' as const },
    { label: '中等', count: medium, percent: Math.round((medium / total) * 100), color: 'var(--warning-color)' as const },
    { label: '困难', count: hard, percent: Math.round((hard / total) * 100), color: 'var(--error-color)' as const }
  ]
})

const contestStats = computed(() => {
  const totalSubmissions = problems.value.reduce((acc: number, p: any) => acc + (p.submission || 0), 0)
  const totalAccepted = problems.value.reduce((acc: number, p: any) => acc + (p.accept || 0), 0)
  const acceptRate = totalSubmissions > 0 ? ((totalAccepted / totalSubmissions) * 100).toFixed(1) : '0.0'
  return { totalSubmissions, totalAccepted, acceptRate }
})
</script>

<template>
  <NGrid :x-gap="24" :y-gap="24" cols="1 l:3" responsive="screen">
    <!-- Left Column -->
    <NGi span="0 l:2">
      <NCard title="比赛题目" :bordered="true" :segmented="{ content: true }">
        <template #header-extra>
          <NTag size="small" round>共 {{ problems.length }} 题</NTag>
        </template>

        <NList v-if="problems.length > 0" hoverable clickable>
          <NListItem
            v-for="problem in problems"
            :key="problem.id"
            @click="router.push({ name: 'ContestProblem', params: { id: contestId, label: problem.label } })"
          >
            <template #prefix>
              <div
                class="w-10 h-10 rounded-lg flex items-center justify-center font-bold text-sm"
                :style="{
                  background: getProblemStatusBg(problem.my_status),
                  color: getProblemStatusColor(problem.my_status)
                }"
              >
                {{ problem.label }}
              </div>
            </template>

            <NThing>
              <template #header>
                <span :style="{ color: 'var(--text-primary)' }">{{ problem.title }}</span>
              </template>
              <template #description>
                <NSpace :size="12" class="text-xs">
                  <NTag :bordered="false" size="tiny" :type="problem.difficulty <= 1 ? 'success' : problem.difficulty <= 3 ? 'warning' : 'error'">
                    {{ problem.difficulty <= 1 ? '简单' : problem.difficulty <= 3 ? '中等' : '困难' }}
                  </NTag>
                  <span :style="{ color: 'var(--text-secondary)' }">
                    通过率 {{ problem.submission > 0 ? ((problem.accept / problem.submission) * 100).toFixed(1) : 0 }} %
                  </span>
                  <span :style="{ color: 'var(--text-secondary)' }">{{ problem.accept || 0 }} 人通过</span>
                </NSpace>
              </template>
            </NThing>

            <template #suffix>
              <NSpace :size="4" align="center">
                <NIcon v-if="problem.my_status === 'accepted'" :size="20" style="color: var(--success-color)"><CheckCircle2 /></NIcon>
                <NIcon v-else-if="problem.my_status === 'wrong'" :size="20" style="color: var(--error-color)"><XCircle /></NIcon>
                <NIcon v-else :size="20" style="color: var(--text-tertiary)"><Circle /></NIcon>
              </NSpace>
            </template>
          </NListItem>
        </NList>

        <NResult v-else description="暂无题目数据" />
      </NCard>
    </NGi>

    <!-- Right Column -->
    <NGi span="0 l:1">
      <NSpace vertical :size="24">
        <!-- Contest Stats -->
        <NCard title="比赛统计" :bordered="true" :segmented="{ content: true }">
          <NGrid :x-gap="12" :y-gap="12" cols="2">
            <NGi>
              <NStatistic label="参与人数">
                <template #prefix>
                  <NIcon :size="16" style="color: var(--text-secondary)"><Users /></NIcon>
                </template>
                {{ contest?.participants || 0 }}
              </NStatistic>
            </NGi>
            <NGi>
              <NStatistic label="提交总数">
                <template #prefix>
                  <NIcon :size="16" style="color: var(--text-secondary)"><FileCode /></NIcon>
                </template>
                {{ formatNumber(contestStats.totalSubmissions) }}
              </NStatistic>
            </NGi>
            <NGi>
              <NStatistic label="通过提交" style="color: var(--success-color)">
                <template #prefix>
                  <NIcon :size="16" style="color: var(--success-color)"><CheckCircle2 /></NIcon>
                </template>
                {{ formatNumber(contestStats.totalAccepted) }}
              </NStatistic>
            </NGi>
            <NGi>
              <NStatistic label="通过率" style="color: var(--warning-color)">
                <template #prefix>
                  <NIcon :size="16" style="color: var(--warning-color)"><Target /></NIcon>
                </template>
                {{ contestStats.acceptRate }}%
              </NStatistic>
            </NGi>
          </NGrid>
        </NCard>

        <!-- Difficulty Distribution -->
        <NCard title="难度分布" :bordered="true" :segmented="{ content: true }">
          <NSpace vertical :size="16">
            <div v-for="d in difficultyDistribution" :key="d.label">
              <div class="flex items-center justify-between text-sm mb-1">
                <span :style="{ color: 'var(--text-secondary)' }">{{ d.label }}</span>
                <span class="font-medium" :style="{ color: d.color }">{{ d.count }} 题</span>
              </div>
              <NProgress
                type="line"
                :percentage="d.percent"
                :show-indicator="false"
                :height="8"
                :color="d.color"
                rail-color="var(--surface-tertiary)"
              />
            </div>
            <NResult v-if="difficultyDistribution.length === 0" description="暂无数据" />
          </NSpace>
        </NCard>

        <!-- Announcement -->
        <NCard :bordered="true" :segmented="{ content: true }">
          <template #header>
            <NSpace :size="8" align="center">
              <NIcon :size="18"><Megaphone /></NIcon>
              <span>公告</span>
            </NSpace>
          </template>
          <NList>
            <NListItem>
              <NThing>
                <template #header>比赛进行中提醒</template>
                <template #description>
                  <span class="text-xs" style="color: var(--text-tertiary)">管理员 · 5分钟前</span>
                </template>
                比赛进行到一半，请注意提交时间。如有问题请及时反馈。
              </NThing>
            </NListItem>
            <NListItem>
              <NThing>
                <template #header>比赛已开始</template>
                <template #description>
                  <span class="text-xs" style="color: var(--text-tertiary)">管理员 · 30分钟前</span>
                </template>
                比赛已开始，祝大家取得好成绩！
              </NThing>
            </NListItem>
          </NList>
        </NCard>
      </NSpace>
    </NGi>
  </NGrid>
</template>
