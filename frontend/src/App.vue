<template>
  <n-message-provider>
    <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="naiveTheme">
      <RouterView :key="theme" />
      <AiAssistant :key="theme" />
    </n-config-provider>
  </n-message-provider>
</template>

<script setup lang="ts">
import AiAssistant from './components/AiAssistant.vue'
import { NMessageProvider, NConfigProvider } from 'naive-ui'
import { zhCN, dateZhCN,NInput,NButton } from 'naive-ui'
import { useTheme } from './composables/useTheme'

const { theme, naiveTheme } = useTheme()
</script>
