<script setup lang="ts">
import { usePerformanceStore } from '@/stores/performance'

const performanceStore = usePerformanceStore()
performanceStore.initializeMonitoring()
</script>

<template>
  <router-view />
</template>

<style>
html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
}

#app {
  height: 100%;
  width: 100%;
}
</style>
