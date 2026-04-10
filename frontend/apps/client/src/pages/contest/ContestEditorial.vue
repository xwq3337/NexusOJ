<script setup lang="ts">
import { ref, computed, onMounted, inject } from 'vue'
import { NResult, NCard, NSpace, NTag, NIcon, NDivider } from 'naive-ui'
import { BookOpen, Lock } from 'lucide-vue-next'

const { contest } = inject<any>('contestData')

const isContestEnded = computed(() => contest.value?.status === 'Ended')
</script>

<template>
  <div>
    <!-- Contest not ended -->
    <NResult v-if="!isContestEnded" status="info" title="题解暂未开放" description="比赛结束后将开放题解查看">
      <template #icon>
        <NIcon :size="48"><Lock /></NIcon>
      </template>
    </NResult>

    <!-- Contest ended but no editorial -->
    <NResult v-else status="info" title="暂无题解" description="该比赛暂未发布题解">
      <template #icon>
        <NIcon :size="48"><BookOpen /></NIcon>
      </template>
    </NResult>
  </div>
</template>
