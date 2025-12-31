<template>
  <n-message-provider>
    <n-config-provider :theme="theme == 'dark' ? darkTheme : lightTheme">
      <RouterView :key="theme" />
      <AiAssistant :key="theme" />
    </n-config-provider>
  </n-message-provider>
</template>

<script setup lang="ts">
import AiAssistant from './components/AiAssistant.vue'
import { NMessageProvider, NConfigProvider, darkTheme, lightTheme } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'

const theme = useLocalStorage('theme', 'light')
</script>
