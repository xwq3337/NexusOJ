<template>
  <n-message-provider>
    <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme="naiveTheme" :theme-overrides="naiveThemeOverrides">
      <RouterView :key="theme" />
      <AiAssistant :key="theme" />
    </n-config-provider>
  </n-message-provider>
</template>

<script setup lang="ts">
import AiAssistant from './components/AiAssistant.vue'
import { NMessageProvider, NConfigProvider } from 'naive-ui'
import { zhCN, dateZhCN } from 'naive-ui'
import { useTheme } from './composables/useTheme'
import { onMounted, provide, ref } from 'vue'
import { useLocalStorage } from '@vueuse/core'
const isAuthorization = ref(false)

const checkAuth = async () => {
  const hasToken = useLocalStorage('access_token', null);

  if (!hasToken.value) {
    isAuthorization.value = false;
    return false;
  }
  console.log('Token exists, validating...', hasToken.value);
  isAuthorization.value = true;
  return true;
  // const { code } = await userApi.ValidateToken();
  // isAuthorization.value = code !== 401; // 如果返回 401，说明 token 无效
  // return isAuthorization.value;
};
provide('isAuthorization', isAuthorization); // 提供检查认证的函数
provide('checkAuth', checkAuth); // 提供检查认证的函数
onMounted(checkAuth); //组件挂载时检查认证状态
const { theme, naiveTheme, naiveThemeOverrides } = useTheme()

</script>
