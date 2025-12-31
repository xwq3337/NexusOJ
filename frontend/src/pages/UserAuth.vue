<template>
  <div style="max-width: 480px; margin: 0 auto; padding: 24px">
    <n-form ref="formRef" :model="modelRef" :rules="rules">
      <n-form-item path="account" label="用户名">
        <n-input
          v-model:value="modelRef.account"
          @keydown.enter.prevent
          :placeholder="'请输入用户名'"
        />
      </n-form-item>
      <n-form-item path="password" label="密码">
        <n-input
          v-model:value="modelRef.password"
          type="password"
          @input="handlePasswordInput"
          @keydown.enter.prevent
          :placeholder="'请输入密码'"
        />
      </n-form-item>
      <n-form-item ref="rPasswordFormItemRef" first path="reenteredPassword" label="重复密码">
        <n-input
          v-model:value="modelRef.reenteredPassword"
          :disabled="!modelRef.password"
          type="password"
          @keydown.enter.prevent
          :placeholder="'请再次输入密码'"
        />
      </n-form-item>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <div style="display: flex; justify-content: flex-end">
            <n-button
              :disabled="modelRef.account === null"
              round
              type="primary"
              @click="handleValidateButtonClick"
            >
              验证
            </n-button>
          </div>
        </n-col>
      </n-row>
    </n-form>
  </div>
</template>

<script setup lang="ts">
import type { FormInst, FormItemInst, FormItemRule, FormRules } from 'naive-ui'
import { useMessage, NForm, NRow, NCol, NFormItem, NButton, NInput } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import { ref } from 'vue'
import Request from '@/services/api/index'
import { HttpStatusCode } from 'axios'
import { useLocalStorage } from '@vueuse/core'
const route = useRoute()
const router = useRouter()
const regex = /redirect=([^&]+)/
const redirectPath = route.fullPath.match(regex)?.slice(1)[0] ?? '/'
interface ModelType {
  account: string | null
  password: string | null
  reenteredPassword: string | null
}

const formRef = ref<FormInst | null>(null)
const rPasswordFormItemRef = ref<FormItemInst | null>(null)
const message = useMessage()
const modelRef = ref<ModelType>({
  account: null,
  password: null,
  reenteredPassword: null
})

function validatePasswordStartWith(rule: FormItemRule, value: string): boolean {
  return (
    !!modelRef.value.password &&
    modelRef.value.password.startsWith(value) &&
    modelRef.value.password.length >= value.length
  )
}

function validatePasswordSame(rule: FormItemRule, value: string): boolean {
  return value === modelRef.value.password
}

const rules: FormRules = {
  account: [
    {
      required: true,
      validator(rule: FormItemRule, value: string) {
        if (!value) {
          return new Error('需要账户名')
        } else if (value.length < 3 || value.length > 20) {
          return new Error('账户名长度应在3到20个字符之间')
        } else if (!/^[a-zA-Z0-9_]+$/.test(value)) {
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

function handlePasswordInput() {
  if (modelRef.value.reenteredPassword) {
    rPasswordFormItemRef.value?.validate({ trigger: 'password-input' })
  }
}
const AccessToken = useLocalStorage('access_token', null)
const RefreshToken = useLocalStorage('refresh_token', null)
const account = useLocalStorage('account', '')
const password = useLocalStorage('password', '')
import { useUserStore } from '@/stores/useUserStore'
const { initStore } = useUserStore()
// 上传登陆事件
function handleValidateButtonClick(e: MouseEvent) {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      Request.post('/user/login', {
        // TODO:字段重命名
        // TODO:密码加密
        username: modelRef.value.account,
        password: modelRef.value.password
      })
        .then((response: any) => {
          if (response.code == HttpStatusCode.Ok) {
            AccessToken.value = response.msg[0]
            RefreshToken.value = response.msg[1]
            account.value = modelRef.value.account
            password.value = modelRef.value.password
            initStore({
              Id: response.info.id,
              Account: modelRef.value.account,
              Nickname: response.info.nickname,
              Gender: response.info.gender,
              Avatar: response.info.avatar,
              Rating: response.info.rating
            })
            router.push({ path: redirectPath, replace: true })
          } else {
            console.log(response.msg)
          }
        })
        .catch((err: any) => {
          message.error('登录失败')
          console.log(err)
        })
    } else {
      console.log(errors)
      message.error('验证失败')
    }
  })
}
</script>

<style scoped>
:deep(.n-form-item-label) {
  color: var(--text-primary);
}

:deep(.n-input-wrapper) {
  background-color: var(--input-bg);
  border-color: var(--border-color);
  color: #fff;
}
</style>
