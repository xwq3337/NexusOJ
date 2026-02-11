<script setup lang="ts">
import Request from '@/api/axios';
import { HttpStatusCode } from 'axios';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/stores/user';
import { storeToRefs } from 'pinia';
import { reactive, ref } from 'vue';
import { useLocalStorage } from '@vueuse/core';
import { useRouter } from 'vue-router';
import type { FormInstance, FormRules } from 'element-plus';

const userstore = useUserStore();
const { id, username, avatar, nickname, gender } = storeToRefs(userstore);
const AccessToken = useLocalStorage("access_token", "");
const RefreshToken = useLocalStorage("refresh_token", "");
const router = useRouter();
const loginFormRef = ref<FormInstance>();
const loading = ref(false);

const form = reactive({
  username: '',
  password: '',
  rememberMe: false
});

// 表单验证规则
const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
});

const login = async () => {
  if (!loginFormRef.value) return;

  await loginFormRef.value.validate((valid) => {
    if (valid) {
      performLogin();
    } else {
      ElMessage.error('请检查输入信息');
    }
  });
};

const performLogin = async () => {
  loading.value = true;
  try {
    const res = await Request.post({
      url: '/admin/login',
      data: {
        username: form.username,
        password: form.password
      }
    });

    if (res.code === HttpStatusCode.Ok) {
      ElMessage.success('登录成功，欢迎回来！');

      // 更新用户信息
      id.value = res.info.id;
      username.value = res.info.username;
      avatar.value = res.info.avatar;
      nickname.value = res.info.nickname;
      gender.value = res.info.gender;
      AccessToken.value = res.msg[0];
      RefreshToken.value = res.msg[1];

      // 记住密码功能
      if (form.rememberMe) {
        localStorage.setItem('rememberedUsername', form.username);
      } else {
        localStorage.removeItem('rememberedUsername');
      }

      router.push({ name: 'home' });
    } else {
      ElMessage.error('登录失败：用户名或密码错误');
    }
  } catch (error) {
    ElMessage.error('登录失败，请检查网络连接');
    console.error('登录错误:', error);
  } finally {
    loading.value = false;
  }
};

// 页面加载时检查是否有记住的用户名
const rememberedUsername = localStorage.getItem('rememberedUsername');
if (rememberedUsername) {
  form.username = rememberedUsername;
  form.rememberMe = true;
}
</script>
<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="login-background">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- 头部logo和标题 -->
      <div class="login-header">
        <div class="logo">
          <el-icon size="48" color="#409eff">
            <Management />
          </el-icon>
        </div>
        <h1 class="login-title">管理后台</h1>
        <p class="login-subtitle">欢迎回来，请登录您的账户</p>
      </div>

      <!-- 登录表单 -->
      <div class="login-form">
        <el-form
          ref="loginFormRef"
          :model="form"
          :rules="rules"
          size="large"
          @keyup.enter="login"
        >
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              prefix-icon="User"
              clearable
              class="login-input"
            />
          </el-form-item>

          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              clearable
              class="login-input"
            />
          </el-form-item>

          <el-form-item>
            <div class="login-options">
              <el-checkbox v-model="form.rememberMe">
                记住密码
              </el-checkbox>
              <el-link type="primary" :underline="false">
                忘记密码？
              </el-link>
            </div>
          </el-form-item>

          <el-form-item>
            <el-button
              type="primary"
              size="large"
              class="login-button"
              :loading="loading"
              @click="login"
            >
              <span v-if="!loading">立即登录</span>
              <span v-else>登录中...</span>
            </el-button>
          </el-form-item>
        </el-form>
      </div>

      <!-- 底部信息 -->
      <div class="login-footer">
        <p>© 2024 管理系统. 保留所有权利.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

/* 背景装饰元素 */
.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  z-index: 0;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  top: 10%;
  right: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(10deg);
  }
}

/* 登录卡片 */
.login-card {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  padding: 40px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

/* 头部样式 */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  margin-bottom: 20px;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: #2c3e50;
  margin: 0 0 10px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-subtitle {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
  font-weight: 400;
}

/* 表单样式 */
.login-form {
  margin-bottom: 30px;
}

.login-input {
  margin-bottom: 20px;
}

:deep(.login-input .el-input__wrapper) {
  border-radius: 12px;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
  transition: all 0.3s ease;
}

:deep(.login-input .el-input__wrapper:hover) {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

:deep(.login-input .el-input__wrapper.is-focus) {
  border-color: #409eff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.login-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 20px 0;
}

:deep(.login-options .el-checkbox) {
  color: #606266;
}

:deep(.login-options .el-link) {
  font-size: 14px;
}

.login-button {
  width: 100%;
  height: 50px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.login-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s;
}

.login-button:hover::before {
  left: 100%;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 25px rgba(102, 126, 234, 0.4);
}

.login-button:active {
  transform: translateY(0);
}

/* 底部样式 */
.login-footer {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.login-footer p {
  color: #95a5a6;
  font-size: 14px;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-card {
    max-width: 90%;
    padding: 30px 20px;
    margin: 20px;
  }

  .login-title {
    font-size: 28px;
  }

  .bg-shape {
    display: none;
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 20px;
  }

  .login-card {
    padding: 25px 15px;
  }

  .login-title {
    font-size: 24px;
  }

  .login-subtitle {
    font-size: 14px;
  }
}

/* 加载动画 */
.login-button.is-loading {
  position: relative;
}

:deep(.el-loading-spinner) {
  color: white;
}

/* 表单验证错误样式 */
:deep(.el-form-item.is-error .el-input__wrapper) {
  border-color: #f56c6c;
  box-shadow: 0 2px 8px rgba(245, 108, 108, 0.2);
}

:deep(.el-form-item__error) {
  font-size: 12px;
  margin-top: 5px;
}

/* 输入框图标颜色 */
:deep(.el-input__prefix-inner .el-icon) {
  color: #909399;
}

:deep(.el-input__wrapper.is-focus .el-input__prefix-inner .el-icon) {
  color: #409eff;
}
</style>
