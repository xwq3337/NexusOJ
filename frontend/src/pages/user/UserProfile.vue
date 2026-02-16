<!-- 用户资料页面向他人展示的页面 -->
<script setup lang="ts">
import { NForm, NFormItem, NUpload, NButton, NSelect,NInput } from 'naive-ui'
import { reactive, ref } from 'vue'
import { Upload, User, Save } from 'lucide-vue-next'
const fileList = ref([])

import { useUserStore } from '@/stores/useUserStore'
const userStore = useUserStore()

const profileForm = reactive({
  nickname: userStore.nickname,
  gender: userStore.gender
})

const genderOptions = [
  { label: '男', value: 'male' },
  { label: '女', value: 'female' }
]

const handleAvatarChange = (options: any) => {
  // Handle avatar upload
  console.log('Avatar changed', options)
}

const updateProfile = () => {
  // Update profile logic
  console.log('Update profile', profileForm)
}
</script>

<template>
  <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
    <User class="inline mr-2" />
    个人信息
  </h3>
  <n-form :model="profileForm" label-placement="left" label-width="auto">
    <n-form-item label="昵称">
      <n-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
    </n-form-item>
    <n-form-item label="性别">
      <n-select v-model:value="profileForm.gender" :options="genderOptions" />
    </n-form-item>
    <n-form-item label="头像">
      <n-upload
        v-model:file-list="fileList"
        :show-file-list="false"
        :on-change="handleAvatarChange"
        accept="image/*"
      >
        <n-button>
          <Upload class="mr-2" />
          上传头像
        </n-button>
      </n-upload>
    </n-form-item>
    <n-form-item>
      <n-button type="primary" @click="updateProfile">
        <Save class="mr-2" />
        保存更改
      </n-button>
    </n-form-item>
  </n-form>
</template>
