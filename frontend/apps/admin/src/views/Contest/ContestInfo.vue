<script setup lang="ts">
import { contestApi } from '@nexusoj/server'
import type { Contest } from '@nexusoj/type'
import { onMounted, ref } from 'vue'
const contestList = ref<Contest[]>([])
onMounted(async () => {
  await contestApi.getContestList().then((res) => {
    const { info, code } = res
    if (code == 200 && info) {
      contestList.value = info
    }
  })
})
</script>

<template>
  <!-- <a-skeleton active></a-skeleton> -->
  {{ Object.keys(contestList[0] || {}) }}
</template>
