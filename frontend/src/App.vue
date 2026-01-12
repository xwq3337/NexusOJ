<!--
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-12-16 13:21:42
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-12 22:46:06
 * @FilePath: /frontend/src/App.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <n-message-provider>
    <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="theme == 'dark' ? darkTheme : lightTheme">
      <RouterView :key="theme" />
      <AiAssistant :key="theme" />
    </n-config-provider>
  </n-message-provider>
</template>

<script setup lang="ts">
import AiAssistant from './components/AiAssistant.vue'
import { NMessageProvider, NConfigProvider, darkTheme, lightTheme } from 'naive-ui'
 import { zhCN, dateZhCN } from 'naive-ui'
import { useLocalStorage } from '@vueuse/core'

const theme = useLocalStorage('theme', 'light')
</script>
