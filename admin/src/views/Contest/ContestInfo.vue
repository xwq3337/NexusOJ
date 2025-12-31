<script setup lang="ts">
import Request from '@/api/axios'
import { onMounted, ref } from 'vue'
const contestList = ref([])
onMounted(async () => {
  await Request.get({ url: '/contest/list' }).then((res) => {
    contestList.value = res.info
  })
})
</script>

<template>
  <!-- <a-skeleton active></a-skeleton> -->
  {{ Object.keys(contestList[0] || {}) }}
</template>
