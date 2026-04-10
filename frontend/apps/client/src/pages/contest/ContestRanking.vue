<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { contestApi } from '@nexusoj/server'
import type { ContestRankItem } from '@nexusoj/type'
import { useNexusEventSource } from '@/composables/useEventSource'
import { safeJsonParse } from '@nexusoj/utils'

const route = useRoute()
const contestId = route.params.id as string

const rankings = ref<ContestRankItem[]>([])
const contestType = ref<string>('ACM')
const loading = ref(true)

const fetchRanking = async () => {
  loading.value = true
  try {
    const res = await contestApi.getContestRanking(contestId)
    if (res.code === 200 && res.info) {
      rankings.value = res.info
    }
  } catch (e) {
    console.error(e)
  }
  loading.value = false
}

const { close: closeSSE, addListener } = useNexusEventSource(
  `/service/contest/${contestId}/ranking/stream`,
  {
    onError(err) {
      console.error('rank error:', err)
    },
    onOpen() {
      console.log('rank open')
    },
    onClose() {
      console.log('rank close')
    },
  }
)
addListener('ranking-update', (msg) => {
  console.log('Received ranking-update event:', msg)
  try {
    const { data, err } = safeJsonParse(msg.data)
    if (!err && data.type === 'full' && data.ranking) {
      rankings.value = data.ranking
    }
  } catch (e) {
    console.error('Error parsing ranking-update message:', e)
  }
})

const formatPenalty = (seconds: number) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = seconds % 60
  return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
}

const medalColor = (rank: number) => {
  if (rank === 1) return '#FFD700'
  if (rank === 2) return '#C0C0C0'
  if (rank === 3) return '#CD7F32'
  return ''
}

onMounted(() => {
  fetchRanking()
})

</script>

<template>
  <div>
    <div class="rounded-xl overflow-hidden" :style="{ background: 'var(--bg-secondary)' }">
      <table class="w-full text-sm">
        <thead>
          <tr :style="{ background: 'var(--bg-tertiary)', color: 'var(--text-secondary)' }">
            <th class="p-3 text-left w-20">排名</th>
            <th class="p-3 text-left">用户</th>
            <th class="p-3 text-center w-24">通过</th>
            <th v-if="contestType === 'ACM'" class="p-3 text-center w-28">罚时</th>
            <th v-else class="p-3 text-center w-24">得分</th>
            <th class="p-3 text-left">各题详情</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(rank, index) in rankings" :key="index" class="border-t"
            :style="{ borderColor: 'var(--border-color)' }">
            <td class="p-3">
              <span class="font-bold text-lg" :style="{ color: medalColor(rank.rank) }">
                {{ rank.rank }}
              </span>
            </td>
            <td class="p-3">
              <div class="flex items-center gap-2">
                <img v-if="rank.avatar" :src="rank.avatar" class="w-7 h-7 rounded-full" />
                <span :style="{ color: 'var(--text-primary)' }">{{ rank.username }}</span>
              </div>
            </td>
            <td class="p-3 text-center">
              <span class="font-bold" :style="{ color: 'var(--text-primary)' }">{{ rank.solved }}</span>
            </td>
            <td v-if="contestType === 'ACM'" class="p-3 text-center font-mono"
              :style="{ color: 'var(--text-secondary)' }">
              {{ formatPenalty(rank.total_penalty) }}
            </td>
            <td v-else class="p-3 text-center font-bold" :style="{ color: 'var(--text-primary)' }">
              {{ rank.score }}
            </td>
            <td class="p-3">
              <div class="flex gap-1 flex-wrap">
                <span v-for="(detail, label) in rank.problems" :key="label"
                  class="px-2 py-0.5 rounded text-xs font-medium" :style="{
                    background: detail.accepted ? 'var(--success-bg, #dcfce7)' : 'var(--bg-tertiary)',
                    color: detail.accepted ? 'var(--success-text, #166534)' : 'var(--text-secondary)'
                  }">
                  {{ label }}
                  <span v-if="detail.accepted && detail.time" class="ml-1 opacity-70">{{ detail.time }}</span>
                  <span v-else-if="detail.attempts > 0" class="ml-1 opacity-70">(-{{ detail.attempts }})</span>
                </span>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="rankings.length === 0 && !loading" class="p-8 text-center" :style="{ color: 'var(--text-secondary)' }">
        暂无排名数据
      </div>
    </div>
  </div>
</template>
