<template>
  <div class="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
    <!-- 背景装饰 -->
    <div class="fixed inset-0 pointer-events-none">
      <div
        class="absolute rounded-full blur-3xl opacity-60 animate-float"
        :style="{
          width: 'min(500px, 60vw)',
          height: 'min(500px, 60vw)',
          top: '-10%',
          left: '-10%',
          background: 'linear-gradient(135deg, #3b82f6 0%, #6366f1 100%)',
          animationDelay: '0s'
        }"
      />
      <div
        class="absolute rounded-full blur-3xl opacity-60 animate-float"
        :style="{
          width: 'min(400px, 50vw)',
          height: 'min(400px, 50vw)',
          bottom: '-10%',
          right: '-10%',
          background: 'linear-gradient(135deg, #8b5cf6 0%, #a855f7 100%)',
          animationDelay: '-5s'
        }"
      />
      <div
        class="absolute rounded-full blur-3xl opacity-60 animate-float"
        :style="{
          width: 'min(300px, 40vw)',
          height: 'min(300px, 40vw)',
          top: '50%',
          left: '50%',
          transform: 'translate(-50%, -50%)',
          background: 'linear-gradient(135deg, #06b6d4 0%, #3b82f6 100%)',
          animationDelay: '-10s'
        }"
      />
    </div>

    <!-- 登录卡片 -->
    <div
      class="relative w-full max-w-md rounded-3xl p-8 border backdrop-blur-xl"
      :style="{
        background: 'var(--bg-card)',
        borderColor: 'var(--border-color)',
        boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.25)'
      }"
    >
      <!-- 头部 -->
      <div class="text-center mb-8">
        <div class="flex items-center justify-center gap-3 mb-2">
          <div
            class="w-12 h-12 rounded-xl flex items-center justify-center"
            :style="{
              background: 'linear-gradient(135deg, #3b82f6 0%, #6366f1 100%)',
              boxShadow: '0 10px 20px -5px rgba(59, 130, 246, 0.4)',
              color: '#ffffff'
            }"
          >
            <Code2 :size="24" />
          </div>
          <h1
            class="text-3xl font-bold"
            :style="{
              background: 'linear-gradient(135deg, #3b82f6 0%, #6366f1 50%, #a855f7 100%)',
              WebkitBackgroundClip: 'text',
              backgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              letterSpacing: '-0.5px'
            }"
          >
            NexusOJ
          </h1>
        </div>
        <p
          class="text-sm font-medium"
          :style="{ color: 'var(--text-secondary)' }"
        >
          欢迎回来，继续你的编程之旅
        </p>
      </div>

      <!-- 表单 -->
      <n-form ref="formRef" :model="modelRef" :rules="rules" class="space-y-5">
        <!-- 用户名 -->
        <n-form-item path="username" :show-label="false">
          <n-input
            v-model:value="modelRef.username"
            placeholder="请输入用户名"
            size="large"
            @keydown.enter.prevent
          >
            <template #prefix>
              <User :size="18" :style="{ color: 'var(--text-tertiary)' }" />
            </template>
          </n-input>
        </n-form-item>

        <!-- 密码 -->
        <n-form-item path="password" :show-label="false">
          <n-input
            v-model:value="modelRef.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password-on="click"
            @input="handlePasswordInput"
            @keydown.enter.prevent
          >
            <template #prefix>
              <Lock :size="18" :style="{ color: 'var(--text-tertiary)' }" />
            </template>
          </n-input>
        </n-form-item>

        <!-- 确认密码 -->
        <n-form-item ref="rPasswordFormItemRef" first path="reenteredPassword" :show-label="false">
          <n-input
            v-model:value="modelRef.reenteredPassword"
            :disabled="!modelRef.password"
            type="password"
            placeholder="请再次输入密码"
            size="large"
            show-password-on="click"
            @keydown.enter.prevent
          >
            <template #prefix>
              <Lock :size="18" :style="{ color: 'var(--text-tertiary)' }" />
            </template>
          </n-input>
        </n-form-item>

        <!-- 登录按钮 -->
        <div class="pt-2">
          <n-button
            :disabled="modelRef.username === null"
            type="primary"
            size="large"
            block
            strong
            @click="handleValidateButtonClick"
            :style="{
              height: '48px',
              fontSize: '16px',
              fontWeight: '600',
              borderRadius: '12px',
              boxShadow: '0 4px 14px -2px rgba(59, 130, 246, 0.4)'
            }"
          >
            <template #icon>
              <LogIn :size="18" />
            </template>
            登录
          </n-button>
        </div>
      </n-form>

      <!-- 底部链接 -->
      <div class="mt-8">
        <div
          class="flex items-center gap-4 text-xs"
          :style="{ color: 'var(--text-tertiary)' }"
        >
          <div class="flex-1 h-px" :style="{ background: 'var(--border-color)' }" />
          <span>或</span>
          <div class="flex-1 h-px" :style="{ background: 'var(--border-color)' }" />
        </div>
        <div
          class="flex items-center justify-center gap-3 mt-6 text-sm font-medium"
        >
          <RouterLink
            to="/forgot-password"
            class="transition-colors hover:text-blue-400"
            :style="{ color: 'var(--text-secondary)' }"
          >
            忘记密码？
          </RouterLink>
          <span :style="{ color: 'var(--border-color)' }">|</span>
          <RouterLink
            to="/register"
            class="transition-colors hover:text-blue-400"
            :style="{ color: 'var(--text-secondary)' }"
          >
            注册账号
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FormInst, FormItemInst, FormItemRule, FormRules } from 'naive-ui'
import { NForm, NFormItem, NInput, NButton, useMessage } from 'naive-ui'
import { User, Lock, LogIn, Code2 } from 'lucide-vue-next'
import { useRoute, RouterLink, useRouter } from 'vue-router'
import { inject, ref } from 'vue'
import { HttpStatusCode } from 'axios'
import { useLocalStorage } from '@vueuse/core'
import { useUserStore } from '@/stores/useUserStore'
import { userApi } from '@/services/user'

const route = useRoute()
const router = useRouter()
const { initStore } = useUserStore()

const regex = /redirect=([^&]+)/
const redirectPath = route.fullPath.match(regex)?.slice(1)[0] ?? '/'

const AccessToken = useLocalStorage('access_token', null)
const RefreshToken = useLocalStorage('refresh_token', null)
const checkAuth = inject('checkAuth') as () => Promise<void>
const message = useMessage()

interface ModelType {
  username: string | null
  password: string | null
  reenteredPassword: string | null
}

const formRef = ref<FormInst | null>(null)
const rPasswordFormItemRef = ref<FormItemInst | null>(null)
const modelRef = ref<ModelType>({
  username: null,
  password: null,
  reenteredPassword: null
})

const validatePasswordStartWith = (rule: FormItemRule, value: string): boolean => {
  return (
    !!modelRef.value.password &&
    modelRef.value.password.startsWith(value) &&
    modelRef.value.password.length >= value.length
  )
}

const validatePasswordSame = (rule: FormItemRule, value: string): boolean => {
  return value === modelRef.value.password
}

const rules: FormRules = {
  username: [
    {
      required: true,
      validator(rule: FormItemRule, value: string) {
        if (!value) {
          return new Error('需要账户名')
        }
        if (value.length < 3 || value.length > 20) {
          return new Error('账户名长度应在3到20个字符之间')
        }
        if (!/^[a-zA-Z0-9_]+$/.test(value)) {
          return new Error('账户名应该为字母、数字或下划线')
        }
        return true
      },
      trigger: ['input', 'blur']
    }
  ],
  password: [
    {
      required: true,
      message: '请输入密码'
    }
  ],
  reenteredPassword: [
    {
      required: true,
      message: '请再次输入密码',
      trigger: ['input', 'blur']
    },
    {
      validator: validatePasswordStartWith,
      message: '两次密码输入不一致',
      trigger: 'input'
    },
    {
      validator: validatePasswordSame,
      message: '两次密码输入不一致',
      trigger: ['blur', 'password-input']
    }
  ]
}

const handlePasswordInput = () => {
  if (modelRef.value.reenteredPassword) {
    rPasswordFormItemRef.value?.validate({ trigger: 'password-input' })
  }
}

const handleValidateButtonClick = (e: MouseEvent) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      userApi
        .Login(modelRef.value.username, modelRef.value.password)
        .then((response) => {
          const { msg, info, code } = response
          if (code === HttpStatusCode.Ok) {
            AccessToken.value = msg[0]
            RefreshToken.value = msg[1]
            initStore({
              Id: info.id,
              Username: modelRef.value.username,
              Nickname: info.nickname,
              Gender: info.gender,
              Avatar: info.avatar,
              Rating: info.rating
            })
            router.push({ path: redirectPath, replace: true })
          } else {
            message.error('登录失败，请检查用户名和密码')
          }
        })
        .catch(() => {
          message.error('登录失败')
        })
        .finally(async () => {
          await checkAuth()
        })
    } else {
      message.error('验证失败')
    }
  })
}
</script>

<style scoped>
@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  25% {
    transform: translate(5%, 5%) scale(1.05);
  }
  50% {
    transform: translate(-5%, 10%) scale(0.95);
  }
  75% {
    transform: translate(-10%, -5%) scale(1.02);
  }
}

.animate-float {
  animation: float 20s ease-in-out infinite;
}

:deep(.n-input) {
  --n-border: 1px solid var(--input-border);
  --n-border-hover: var(--input-focus);
  --n-border-focus: var(--input-focus);
  --n-color: var(--input-bg);
  --n-color-focus: var(--input-bg);
  --n-text-color: var(--text-primary);
  --n-placeholder-color: var(--input-placeholder);
  --n-border-radius: 12px;
  --n-height: 48px;
  --n-padding: 0 16px;
}

:deep(.n-input__prefix) {
  color: var(--text-tertiary);
  margin-right: 8px;
}

:deep(.n-button--primary) {
  background: #3b82f6;
}

:deep(.n-button--primary:hover) {
  background: #2563eb;
}

:deep(.n-button--primary:active) {
  background: #1d4ed8;
}

@media (max-width: 640px) {
  .max-w-md {
    max-width: 100%;
  }

  .p-8 {
    padding: 1.5rem;
  }

  .text-3xl {
    font-size: 1.5rem;
  }
}
</style>
