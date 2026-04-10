<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import type { ContestRankItem } from '@nexusoj/type'
import { useNexusEventSource } from '@/composables/useEventSource'
import { safeJsonParse } from '@nexusoj/utils'

const route = useRoute()
const contestId = route.params.id as string

const rankings = ref<ContestRankItem[]>([])
const contestType = ref<string>('ACM')
const loading = ref(true)

const { close, addListener } = useNexusEventSource(
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

onUnmounted(() => {
  close()
})
</script>

<template>
  <div>
    <div class="rounded-xl overflow-hidden" :style="{ background: 'var(--bg-secondary)' }">
      <div class="ranking-header" :style="{ background: 'var(--bg-tertiary)', color: 'var(--text-secondary)' }">
        <span class="w-20 text-left p-3">排名</span>
        <span class="flex-1 text-left p-3">用户</span>
        <span class="w-24 text-center p-3">通过</span>
        <span v-if="contestType === 'ACM'" class="w-28 text-center p-3">罚时</span>
        <span v-else class="w-24 text-center p-3">得分</span>
        <span class="flex-1 text-left p-3">各题详情</span>
      </div>

      <TransitionGroup name="flip-list" tag="div" class="ranking-body">
        <div
          v-for="rank in rankings"
          :key="rank.user_id"
          class="ranking-row border-t"
          :style="{ borderColor: 'var(--border-color)' }"
        >
          <div class="w-20 p-3">
            <span class="font-bold text-lg" :style="{ color: medalColor(rank.rank) }">
              {{ rank.rank }}
            </span>
          </div>
          <div class="flex-1 p-3">
            <div class="flex items-center gap-2">
              <img v-if="rank.avatar" :src="rank.avatar" class="w-7 h-7 rounded-full" />
              <span :style="{ color: 'var(--text-primary)' }">{{ rank.username }}</span>
            </div>
          </div>
          <div class="w-24 p-3 text-center">
            <span class="font-bold" :style="{ color: 'var(--text-primary)' }">{{ rank.solved }}</span>
          </div>
          <div v-if="contestType === 'ACM'" class="w-28 p-3 text-center font-mono" :style="{ color: 'var(--text-secondary)' }">
            {{ formatPenalty(rank.total_penalty) }}
          </div>
          <div v-else class="w-24 p-3 text-center font-bold" :style="{ color: 'var(--text-primary)' }">
            {{ rank.score }}
          </div>
          <div class="flex-1 p-3">
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
          </div>
        </div>
      </TransitionGroup>

      <div v-if="rankings.length === 0 && !loading" class="p-8 text-center" :style="{ color: 'var(--text-secondary)' }">
        暂无排名数据
      </div>
    </div>
  </div>
</template>

<style scoped>
.ranking-header {
  display: flex;
  font-size: 0.875rem;
  font-weight: 500;
}

.ranking-body {
  position: relative;
}

.ranking-row {
  display: flex;
  font-size: 0.875rem;
  align-items: center;
  transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

.flip-list-move {
  transition: transform 0.6s cubic-bezier(0.55, 0, 0.1, 1);
}

.flip-list-enter-active,
.flip-list-leave-active {
  transition: all 0.5s cubic-bezier(0.55, 0, 0.1, 1);
}

.flip-list-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.flip-list-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.flip-list-leave-active {
  position: absolute;
  width: 100%;
}
</style>
