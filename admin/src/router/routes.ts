import type { RouteRecordRaw } from 'vue-router'
import AdminLayout from '@/views/AdminLayout.vue'
export const routes: RouteRecordRaw[] = [
  {
    path: '/user/login',
    name: 'login',
    component: () => import('@/views/User/UserLogin.vue'),
  },
  {
    path: '/',
    component: AdminLayout,
    children: [
      {
        path: '/',
        name: 'home',
        component: () => import('@/views/Home/HomeView.vue'),
      },
      {
        path: '/personal-info',
        name: 'personal-info',
        component: () => import('@/views/User/PersonalCenter.vue'),
      },
      {
        path: '/performance',
        name: 'performance',
        component: () => import('@/views/Performance/PerformanceDashboard.vue'),
        meta: {
          title: '性能监控',
          requiresAuth: true,
          roles: ['admin'],
        },
      },
      {
        path: '/blog',
        name: 'blog',
        component: () => import('@/views/Blog/BlogLayout.vue'),
        children: [
          {
            path: 'info',
            name: 'blog-info',
            component: () => import('@/views/Blog/BlogInfo.vue'),
            meta: {
              title: '博客信息管理',
            },
          },
          {
            path: 'recycle',
            name: 'blog-recycle',
            component: () => import('@/views/Blog/BlogRecycling.vue'),
            meta: {
              title: '博客回收站',
            },
          },
          {
            path: 'verify',
            name: 'blog-verify',
            component: () => import('@/views/Blog/BlogVerify.vue'),
            meta: {
              title: '博客审核',
            },
          },
        ],
      },
      {
        path: '/contest-info',
        name: 'contest-info',
        meta: {
          title: '比赛信息管理',
        },
        component: () => import('@/views/Contest/ContestInfo.vue'),
      },
      {
        path: '/contest-create',
        name: 'contest-create',
        meta: {
          title: '比赛创建',
        },
        component: () => import('@/views/Contest/ContestCreate.vue'),
      },
      {
        path: '/problem-info',
        name: 'problem-info',
        meta: {
          title: '题目信息管理',
        },
        component: () => import('@/views/Problem/ProblemInfo.vue'),
      },
      {
        path: '/problem-edit',
        name: 'problem-edit',
        meta: {
          title: '题目编辑',
        },
        component: () => import('@/views/Problem/ProblemChange.vue'),
      },
      {
        path: '/problem-create',
        name: 'problem-create',
        meta: {
          title: '题目创建',
        },
        component: () => import('@/views/Problem/ProblemCreate.vue'),
      },
      {
        path: '/training',
        name: 'training',
        component: () => import('@/views/Training/TrainingLayout.vue'),
        children: [
          {
            path: 'list',
            name: 'training-info',
            meta: {
              title: '训练信息管理',
            },
            component: () => import('@/views/Training/TrainingList.vue'),
          },
          {
            path : 'create',
            name : 'training-create',
            meta : {
              title : '训练创建',
            },
            component : () => import('@/views/Training/TrainingCreate.vue'),
          }
        ],
      },
      {
        path: '/user-info',
        name: 'user-info',
        component: () => import('@/views/User/UserList.vue'),
        meta: {
          title: '用户管理',
          requiresAuth: true,
          roles: [2, 3], // 管理员和超级管理员可访问
        },
      },
      {
        path: '/user-permission',
        name: 'user-permission',
        component: () => import('@/views/User/UserPermission.vue'),
        meta: {
          title: '权限管理',
          requiresAuth: true,
          roles: [3], // 仅超级管理员可访问
        },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFoundView.vue'),
  },
]

export default routes
