import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/Layout.vue'),
    children: [
      {
        path: '',
        component: () => import('@/pages/home/Home.vue'),
        name: 'Home'
      },
      {
        path: 'problems',
        component: () => import('@/pages/problem/Problems.vue'),
        name: 'Problems'
      },
      {
        path: 'contests',
        component: () => import('@/pages/contest/Contests.vue'),
        name: 'Contests'
      },
      {
        path: 'leaderboard',
        component: () => import('@/pages/Leaderboard.vue'),
        name: 'Leaderboard'
      },
      {
        path: 'records',
        component: () => import('@/pages/record/Records.vue'),
        name: 'Records'
      },
      {
        path: 'notifications',
        component: () => import('@/pages/notification/notifications.vue'),
        name: 'Notifications',
        redirect: { name: 'Info' },
        children: [
          {
            path: 'info',
            component: () => import('@/pages/notification/info.vue'),
            name: 'Info',
            meta: {
              description: '系统通知'
            }
          },
          {
            path: 'message',
            component: () => import('@/pages/notification/message.vue'),
            name: 'Message',
            meta: {
              description: '私信'
            }
          }
        ]
      },
      {
        path: 'knowledgebase',
        component: () => import('@/pages/KnowledgeBase.vue'),
        name: 'KnowledgeBase',
        children: [
          {
            path: 'courses',
            component: () => import('@/pages/Courses.vue'),
            name: 'Courses'
          },
          {
            path: 'blogs',
            component: () => import('@/pages/blog/Blogs.vue'),
            name: 'Blogs'
          },
          {
            path: '',
            redirect: { name: 'Courses' }
          }
        ]
      },
      {
        path: 'user',
        component: () => import('@/layouts/UserLayout.vue'),
        children: [
          {
            path: 'login',
            component: () => import('@/pages/user/UserAuth.vue'),
            name: 'Login',
            meta: {
              desciption: '登录与注册'
            }
          },
          {
            path: 'personal-center',
            component: () => import('@/pages/user/UserPersonalCenter.vue'),
            name: 'PersonalCenter',
            children: [
              {
                path: 'profile',
                component: () => import('@/pages/user/UserProfile.vue'),
                name: 'Profile',
                meta: {
                  description: '向他人展示的个人中心'
                }
              },
              {
                path: 'security',
                component: () => import('@/pages/user/UserSecurity.vue'),
                name: 'Security',
                meta: {
                  description: '账号安全'
                }
              },
              {
                path: 'settings',
                component: () => import('@/pages/user/UserSettings.vue'),
                name: 'Settings'
              },
              {
                path: 'user-record',
                component: () => import('@/pages/user/UserRecord.vue'),
                name: 'UserRecord'
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/problem/:id',
    component: () => import('@/pages/problem/ProblemDetail.vue'),
    name: 'ProblemDetail'
  },
  {
    path: '/record/:id',
    component: () => import('@/pages/record/RecordDetail.vue'),
    name: 'RecordDetail'
  },
  {
    path: '/blog/:id',
    component: () => import('@/pages/blog/BlogDetail.vue'),
    name: 'BlogDetail',
  },
  {
    path: '/blog/create',
    component: () => import('@/pages/blog/BlogCreate.vue'),
    name: 'BlogCreate'
  },
  {
    path: '/contest/:id',
    component: () => import('@/pages/contest/ContestDetail.vue'),
    name: 'ContestDetail'
  },
  {
    path: '/user/:id',
    component: () => import('@/pages/user/UserHomePage.vue'),
    name: 'UserHomePage'
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
