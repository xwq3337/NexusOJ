<script lang="ts" setup name="AdminLayout">
import { useRouter } from 'vue-router'
import {
  House,
  Monitor,
  User,
  EditPen,
  List,
  TrophyBase,
  Reading
} from '@element-plus/icons-vue'

const isCollapse = useLocalStorage('isCollapse', false)
const router = useRouter()

const menuItems = [
  {
    index: 'home',
    title: '首页',
    icon: House
  },
  {
    index: 'performance',
    title: '性能监控',
    icon: Monitor
  },
  {
    index: 'user-info',
    title: '用户管理',
    icon: User
  },
  {
    index: 'blog-info',
    title: '博客管理',
    icon: EditPen
  },
  {
    index: 'problem-info',
    title: '题目管理',
    icon: List
  },
  {
    index: 'contest-info',
    title: '比赛管理',
    icon: TrophyBase
  },
  {
    index: 'training-info',
    title: '训练管理',
    icon: Reading
  }
]

const handleSelect = (name: string) => {
  console.log(name)
  router.push({name})
}

import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { useLocalStorage } from '@vueuse/core'
const userStore = useUserStore()
const { avatar } = storeToRefs(userStore)

</script>

<template>
  <el-container class="layout-container">
    <el-aside :width="isCollapse ? '64px' : '200px'" class="aside">
      <div class="logo-container">
        <img src="@/assets/logo.ico" alt="Logo" class="logo" />
        <span v-show="!isCollapse" class="title">OJ Admin</span>
      </div>
      <el-menu :default-active="router.currentRoute.value.path" class="el-menu-vertical" :collapse="isCollapse"
        @select="handleSelect($event)">
        <el-menu-item v-for="item in menuItems" :key="item.index" :index="item.index">
          <el-icon>
            <component :is="item.icon" />
          </el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <el-button type="primary" @click="isCollapse = !isCollapse">
          {{ isCollapse ? '展开' : '折叠' }}
        </el-button>
        <div class="header-right">
          <el-dropdown>
            <el-avatar :size="32"
            style="cursor: pointer;"
              :src="avatar ?? `https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png`" />
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="router.push('/personal-info')">个人信息</el-dropdown-item>
                <el-dropdown-item
                  @click="$router.push({ name: 'login', params: { redirect: router.currentRoute.value.path } })">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped lang="scss">
.layout-container {
  height: 100vh;
  width: 100vw;

  .aside {
    background-color: var(--el-menu-bg-color);
    transition: width 0.3s;
    overflow: hidden;

    .logo-container {
      height: 60px;
      display: flex;
      align-items: center;
      padding: 0 20px;
      overflow: hidden;

      .logo {
        height: 32px;
        width: 32px;
      }

      .title {
        margin-left: 12px;
        color: var(--el-text-color-primary);
        font-size: 18px;
        font-weight: 600;
        white-space: nowrap;
      }
    }

    .el-menu-vertical {
      border-right: none;
    }
  }

  .header {
    background-color: var(--el-bg-color);
    border-bottom: 1px solid var(--el-border-color-light);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;

    .header-right {
      display: flex;
      align-items: center;
      gap: 16px;
    }
  }

  .main {
    background-color: var(--el-bg-color-page);
    padding: 20px;
  }
}
</style>
