import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/Layout.vue'),
    children: [
      {
        path: '',
        component: () => import('@/pages/Home.vue'),
        name: 'Home'
      },
      {
        path: 'problems',
        component: () => import('@/pages/Problems.vue'),
        name: 'Problems'
      },
      {
        path: 'contests',
        component: () => import('@/pages/Contests.vue'),
        name: 'Contests'
      },
      {
        path: 'leaderboard',
        component: () => import('@/pages/Leaderboard.vue'),
        name: 'Leaderboard'
      },
      {
        path: 'records',
        component: () => import('@/pages/Records.vue'),
        name: 'Records'
      },
      {
        path: 'message',
        component: () => import('@/pages/Message.vue'),
        name: 'Message'
      },
      {
        path: 'knowledgebase',
        component: () => import('@/pages/KnowledgeBase.vue'),
        name: 'KnowledgeBase',
        children: [
          {
            path: 'courses',
            component: () => import('@/pages/CoursePage.vue'),
            name: 'Courses'
          },
          {
            path: 'blogs',
            component: () => import('@/pages/Blogs.vue'),
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
            component: () => import('@/pages/UserAuth.vue'),
            name: 'Login'
          },
          {
            path: 'profile',
            component: () => import('@/pages/UserProfile.vue'),
            name: 'Profile'
          },
          {
            path: 'personal-center',
            component: () => import('@/pages/UserPersonalCenter.vue'),
            name: 'PersonalCenter'
          }
        ]
      }
    ]
  },
  {
    path: '/problem/:id',
    component: () => import('@/pages/ProblemDetail.vue'),
    name: 'ProblemDetail'
  },
  {
    path: '/record/:id',
    component: () => import('@/pages/RecordDetail.vue'),
    name: 'RecordDetail'
  },
  {
    path: '/blog/:id',
    component: () => import('@/pages/BlogDetail.vue'),
    name: 'BlogDetail'
  },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
