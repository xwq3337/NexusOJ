<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Lock, Message } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import Request from '@/api/axios'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
const userStore = useUserStore()
const { id, username, nickname, gender, avatar, rating } = storeToRefs(userStore)
const userForm = reactive({
  email: '',
  phone: '',
  bio: '',
  school: '',
  codeforces: '',
  birthday: '',
})
onMounted(async () => {
  await Request.get({ url: `/user/${id.value}` }).then(res => {
    username.value = res.info.username
    nickname.value = res.info.nickname
    userForm.email = res.info.email
    avatar.value = res.info.avatar
    rating.value = res.info.rating
    userForm.bio = res.info.introduction
    gender.value = res.info.gender
    userForm.school = res.info.school
    userForm.codeforces = res.info.codeforces
    userForm.birthday = res.info.birthday
  })
})

// 表单校验规则
const rules = reactive<FormRules>({
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
})

// 安全设置
const security = reactive({
  twoFactor: false
})

// 操作日志
const logs = ref([
  {
    time: '2024-03-20 10:00:00',
    type: '登录',
    ip: '192.168.1.1',
    device: 'Chrome / MacOS',
    status: 'success'
  },
  {
    time: '2024-03-19 15:30:00',
    type: '修改密码',
    ip: '192.168.1.1',
    device: 'Chrome / MacOS',
    status: 'success'
  }
])

// 修改密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 修改密码表单校验规则
const passwordRules = reactive<FormRules>({
  oldPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
})

// 状态变量
const loading = ref(false)
const formRef = ref<FormInstance>()
const passwordFormRef = ref<FormInstance>()
const passwordDialog = reactive({
  visible: false,
  loading: false
})

// 头像上传相关方法
interface UploadResponse {
  url: string
}

const handleAvatarSuccess = (response: UploadResponse) => {
  // avatar.value = response.url
  console.log(response)
  ElMessage.success('头像上传成功')
}

const beforeAvatarUpload = (file: File) => {
  const isJPG = file.type === 'image/jpeg'
  const isPNG = file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG && !isPNG) {
    ElMessage.error('头像只能是 JPG 或 PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB!')
    return false
  }
  return true
}

// 保存个人信息
const handleSave = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true
    // TODO: 调用保存API
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('保存成功')
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}

// 验证邮箱
const handleVerifyEmail = () => {
  ElMessage.info('已发送验证邮件，请查收')
}

// 验证手机
const handleVerifyPhone = () => {
  ElMessage.info('已发送验证码，请查收')
}

// 修改密码
const handleChangePassword = () => {
  passwordDialog.visible = true
}

// 提交密码修改
const submitPasswordChange = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    passwordDialog.loading = true
    // TODO: 调用修改密码API
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success('密码修改成功')
    passwordDialog.visible = false
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    passwordDialog.loading = false
  }
}

// 双因素认证开关
const handle2FAChange = (value: boolean) => {
  ElMessage.success(`${value ? '开启' : '关闭'}双因素认证`)
}
</script>

<template>
  <div class="personal-info-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
          <el-button type="primary" @click="handleSave" :loading="loading">保存更改</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :span="6">
          <div class="avatar-container">
            <el-upload class="avatar-uploader" action="/api/user/avatar" :show-file-list="false"
              :on-success="handleAvatarSuccess" :before-upload="beforeAvatarUpload">
              <img v-if="avatar" :src="avatar" class="avatar" />
              <el-icon v-else class="avatar-uploader-icon">
                <Plus />
              </el-icon>
            </el-upload>
            <div class="avatar-tips">点击更换头像</div>
          </div>
        </el-col>

        <el-col :span="18">
          <el-form ref="formRef" :model="userForm" :rules="rules" label-width="100px" class="user-form">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="username" disabled />
            </el-form-item>

            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="nickname" />
            </el-form-item>

            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userForm.email">
                <template #append>
                  <el-button @click="handleVerifyEmail">验证邮箱</el-button>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="手机号" prop="phone">
              <el-input v-model="userForm.phone">
                <template #append>
                  <el-button @click="handleVerifyPhone">验证手机</el-button>
                </template>
              </el-input>
            </el-form-item>

            <el-form-item label="个人简介" prop="bio">
              <el-input v-model="userForm.bio" type="textarea" :rows="3" placeholder="请输入个人简介" />
            </el-form-item>
          </el-form>
        </el-col>
      </el-row>
    </el-card>

    <!-- 安全设置 -->
    <el-card class="box-card security-card">
      <template #header>
        <div class="card-header">
          <span>安全设置</span>
        </div>
      </template>

      <el-timeline>
        <el-timeline-item>
          <template #dot>
            <el-icon>
              <Lock />
            </el-icon>
          </template>
          <div class="security-item">
            <div class="security-info">
              <h4>账户密码</h4>
              <p>建议您定期更换密码，设置一个包含字母、数字和特殊字符的强密码</p>
            </div>
            <el-button @click="handleChangePassword">修改密码</el-button>
          </div>
        </el-timeline-item>

        <el-timeline-item>
          <template #dot>
            <el-icon>
              <Message />
            </el-icon>
          </template>
          <div class="security-item">
            <div class="security-info">
              <h4>双因素认证</h4>
              <p>绑定手机或邮箱后，登录时将需要进行二次验证</p>
            </div>
            <el-switch v-model="security.twoFactor" @change="handle2FAChange" />
          </div>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 最近操作日志 -->
    <el-card class="box-card log-card">
      <template #header>
        <div class="card-header">
          <span>操作日志</span>
          <el-button text>查看全部</el-button>
        </div>
      </template>

      <el-table :data="logs" style="width: 100%">
        <el-table-column prop="time" label="时间" width="180" />
        <el-table-column prop="type" label="操作类型" width="120" />
        <el-table-column prop="ip" label="IP地址" width="160" />
        <el-table-column prop="device" label="设备信息" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 修改密码对话框 -->
    <el-dialog v-model="passwordDialog.visible" title="修改密码" width="500px">
      <el-form ref="passwordFormRef" :model="passwordForm" :rules="passwordRules" label-width="100px">
        <el-form-item label="当前密码" prop="oldPassword">
          <el-input v-model="passwordForm.oldPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认新密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="passwordDialog.visible = false">取消</el-button>
          <el-button type="primary" @click="submitPasswordChange" :loading="passwordDialog.loading">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>


<style scoped>
.personal-info-container {
  padding: 20px;
}

.box-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.avatar-container {
  text-align: center;
}

.avatar-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader:hover {
  border-color: var(--el-color-primary);
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}

.avatar-tips {
  font-size: 12px;
  color: #999;
  margin-top: 8px;
}

.security-card {
  margin-top: 20px;
}

.security-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
}

.security-info h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
}

.security-info p {
  margin: 0;
  color: #666;
  font-size: 13px;
}

.log-card {
  margin-top: 20px;
}

:deep(.el-timeline-item__node) {
  background-color: transparent;
  border: none;
}

:deep(.el-timeline-item__dot) {
  background-color: var(--el-color-primary-light-9);
}

:deep(.el-timeline-item__content) {
  padding-bottom: 20px;
}
</style>
