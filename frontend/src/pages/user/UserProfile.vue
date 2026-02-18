<script setup lang="ts">
import { NDatePicker, NPopconfirm, NForm, NFormItem, NFormItemGi, NGrid, NGi, NCalendar, NUpload, NButton, NSelect, NRadioGroup, NRadio, NInput, type FormRules, useMessage } from 'naive-ui'
import { onMounted, ref, computed } from 'vue'
import { Upload, User as UserIcon, Save } from 'lucide-vue-next'
import type { User } from '@/types/user'
import { useUserStore } from '@/stores/useUserStore'
import { userApi } from '@/services/user'
import { pick } from 'lodash'
import { storeToRefs } from 'pinia'
const fileList = ref([])
const { id } = useUserStore()
const message = useMessage()

type ProfileForm = Omit<Pick<User, 'nickname' | 'gender' | 'introduction' | 'email' | 'birthday'>, 'birthday'> & {
  birthday: number | null
}
const profileForm = ref<ProfileForm>({
  nickname: '',
  gender: '0',
  introduction: '一句话介绍不了我',
  email: '',
  birthday: null,
})
onMounted(async () => {
  await userApi.getInfoById(id).then(response => {
    const data = pick(response.info, ['nickname', 'gender', 'introduction', 'email', 'birthday'] as const)
    profileForm.value = {
      ...data,
      birthday: data.birthday ? new Date(data.birthday).getTime() : null
    }

  }).catch(error => {
    console.error('获取用户信息失败:', error)
  })
})

// 将时间戳转换为字符串用于 API 提交
const birthdayForSubmit = computed(() => {
  return profileForm.value.birthday
    ? new Date(profileForm.value.birthday).toISOString().split('T')[0]
    : ''
})
const genderOptions = [
  { label: '男', value: 'male' },
  { label: '女', value: 'female' }
]
const { avatar } = storeToRefs(useUserStore())
const handleAvatarChange = async (options: any) => {
  await userApi.updateAvatar(options.file.file).then(response => {
    if (response.code == 200) {
      avatar.value = response.info + `?t=${Date.now()}`
      message.success('头像更新成功')
    }
    else {
      message.error('头像更新失败: ' + response.msg)
    }
  }).catch(error => {
    message.error('头像更新失败')
  })
}
const handleValidateButtonClick = async () => {
  await userApi.updateUser({
    ...profileForm.value,
    birthday: birthdayForSubmit.value
  }).then(response => {
    message.success('用户信息更新成功')
  }).catch(error => {
    console.error('Failed to update profile', error)
  })
}

const rules: FormRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 3, max: 20, message: '昵称长度应在 3-20 字符之间', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ]
}
const handlePositiveClick = () => {
  message.info('引用功能暂未实现')
}
</script>

<template>
  <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
    <UserIcon class="inline mr-2" />
    个人信息
  </h3>
  <n-form :model="profileForm" label-placement="left" label-width="auto" :rules="rules">
    <n-grid :cols="24" :x-gap="24">
      <n-form-item-gi :span="12" label="昵称" path="nickname">
        <n-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
      </n-form-item-gi>
      <n-form-item-gi :span="12" label="性别">
        <n-radio-group v-model:value="profileForm.gender">
          <n-radio v-for="option in genderOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </n-radio>
        </n-radio-group>
      </n-form-item-gi>
    </n-grid>
    <n-form-item label="简介">
      <n-input v-model:value="profileForm.introduction" type="textarea" placeholder="请输入简介" />
    </n-form-item>
    <n-form-item label="邮箱" path="email">
      <n-input v-model:value="profileForm.email" placeholder="请输入邮箱" />
    </n-form-item>
    <n-form-item label="生日">
      <n-date-picker type="date" v-model:value="profileForm.birthday" />
    </n-form-item>
    <!-- TODO: 添加学校 -->
    <!--  -->
    <n-form-item label="头像">
      <n-upload v-model:file-list="fileList" :show-file-list="false" :on-change="handleAvatarChange" accept="image/*">
        <n-button>
          <Upload class="mr-2" />
          上传头像
        </n-button>
      </n-upload>
    </n-form-item>
    <n-grid :cols="24" :x-gap="24">
      <n-gi :span="24">
        <div class="flex justify-end">
          <n-popconfirm @positive-click="handleValidateButtonClick">
            <template #trigger>
              <n-button type="info" ghost>
                <Save class="mr-2" />
                保存更改
              </n-button>
            </template>
            你确定要更新个人信息吗？请确保所有信息正确无误。
          </n-popconfirm>
        </div>
      </n-gi>
    </n-grid>

  </n-form>
</template>
