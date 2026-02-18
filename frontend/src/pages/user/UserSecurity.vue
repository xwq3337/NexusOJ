<template>
  <div class="text-2xl text-black">安全中心</div>
  <div class="text-black">
    1. 改密码
    2. 第三方绑定
  </div>
  <n-button type="info" @click="handleLogout">退出登录</n-button>
</template>


<script setup lang="ts">
import { NButton } from 'naive-ui';
import { useRouter } from 'vue-router';
const router = useRouter();
import { useUserStore } from '@/stores/useUserStore';
import { useLocalStorage } from '@vueuse/core';
import { inject, Ref } from 'vue';
const { resetStore } = useUserStore();
const accessToken = useLocalStorage('access_token', '')
const refreshToken = useLocalStorage('refresh_token', '')
const isAuthorization = inject('isAuthorization') as Ref<boolean>; // 注入认证状态
  const handleLogout = async () => {
  // 这里可以添加退出登录的逻辑，例如清除用户信息、重置状态等
  resetStore();
  accessToken.value = '', refreshToken.value = '';
  isAuthorization.value = false; // 设置认证状态为未认证
  // 重定向到登录页
  router.push('/user/login');
};
</script>