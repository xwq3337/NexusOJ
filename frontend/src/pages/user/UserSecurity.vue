<template>
  <div class="max-w-2xl mx-auto">
    <h1 class="text-2xl font-bold mb-6" :style="{ color: 'var(--text-primary)' }">安全中心</h1>

    <!-- 修改密码模块 -->
    <div
      class="rounded-lg p-6 mb-6"
      :style="{
        backgroundColor: 'var(--card-bg)',
        border: '1px solid var(--border-color)'
      }"
    >
      <h2 class="text-lg font-semibold mb-4" :style="{ color: 'var(--text-primary)' }">修改密码</h2>

      <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
        <n-form-item label="旧密码" path="oldPassword">
          <n-input
            v-model:value="formData.oldPassword"
            type="password"

            placeholder="请输入旧密码"
            show-password-on="mousedown"
            :disabled="loading"
            class="w-full"
          />
        </n-form-item>

        <n-form-item label="新密码" path="newPassword">
          <n-input
            v-model:value="formData.newPassword"
            type="password"
            placeholder="请输入新密码"
            show-password-on="mousedown"
            :disabled="loading"
            class="w-full"
          />
        </n-form-item>

        <n-form-item label="确认新密码" path="confirmPassword">
          <n-input
            v-model:value="formData.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password-on="mousedown"
            :disabled="loading"
            class="w-full"
          />
        </n-form-item>

        <n-form-item>
          <n-space>
            <n-button
              type="primary"
              :loading="loading"
              @click="handleUpdatePassword"
            >
              确认修改
            </n-button>
            <n-button @click="handleReset">重置</n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </div>

    <!-- 第三方绑定模块（占位） -->
    <div
      class="rounded-lg p-6"
      :style="{
        backgroundColor: 'var(--card-bg)',
        border: '1px solid var(--border-color)'
      }"
    >
      <h2 class="text-lg font-semibold mb-4" :style="{ color: 'var(--text-primary)' }">第三方绑定</h2>
      <p class="text-sm" :style="{ color: 'var(--text-secondary)' }">暂无绑定的第三方账号</p>
    </div>

    <!-- 退出登录 -->
    <div class="mt-6">
      <n-button type="error" @click="handleLogout">退出登录</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NButton, NForm, NFormItem, NInput, NSpace, type FormInst, type FormRules } from 'naive-ui'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/useUserStore'
import { useLocalStorage } from '@vueuse/core'
import { inject, Ref, reactive, ref } from 'vue'
import { userApi } from '@/services/user'
import { useMessage } from 'naive-ui'

const router = useRouter()
const message = useMessage()
const { resetStore } = useUserStore()
const accessToken = useLocalStorage('access_token', '')
const refreshToken = useLocalStorage('refresh_token', '')
const isAuthorization = inject('isAuthorization') as Ref<boolean>

const formRef = ref<FormInst | null>(null)
const loading = ref(false)

interface PasswordFormData {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

const formData = reactive<PasswordFormData>({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validatePasswordSame = (_rule: any, value: string) => {
  return value === formData.newPassword
}

const rules: FormRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: validatePasswordSame,
      message: '两次输入的密码不一致',
      trigger: ['blur', 'input']
    }
  ]
}

const handleUpdatePassword = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    await userApi.updatePassword(formData.oldPassword, formData.newPassword)

    message.success('密码修改成功，请重新登录')
    handleReset()

    // 延迟后退出登录
    setTimeout(() => {
      handleLogout()
    }, 1500)
  } catch (error: any) {
    if (error?.errors) {
      // 表单验证错误，不显示消息
      return
    }
    message.error(error?.message || '密码修改失败，请检查旧密码是否正确')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formData.oldPassword = ''
  formData.newPassword = ''
  formData.confirmPassword = ''
  formRef.value?.restoreValidation()
}

const handleLogout = () => {
  resetStore()
  accessToken.value = ''
  refreshToken.value = ''
  isAuthorization.value = false
  router.push('/user/login')
}
</script>