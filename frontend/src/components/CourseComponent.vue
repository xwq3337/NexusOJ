<template>
  <div class="course-container p-6">
    <div class="mb-6">
      <h2 class="text-2xl font-bold" :style="{ color: 'var(--text-primary)' }">编程课程</h2>
      <p class="text-gray-500 mt-2">探索各种编程课程，提升你的技能</p>
    </div>

    <!-- 课程搜索和筛选 -->
    <div class="flex flex-col md:flex-row gap-4 mb-6">
      <n-input placeholder="搜索课程..." :style="{ backgroundColor: 'var(--surface-tertiary)' }">
        <template #prefix>
          <Search :size="18" />
        </template>
      </n-input>
      <n-select
        v-model:value="selectedCategory"
        :options="categoryOptions"
        placeholder="选择分类"
        class="w-full md:w-48"
      />
      <n-select
        v-model:value="selectedLevel"
        :options="levelOptions"
        placeholder="选择难度"
        class="w-full md:w-40"
      />
    </div>

    <!-- 课程列表 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <n-card
        v-for="course in filteredCourses"
        :key="course.id"
        class="course-card hover:shadow-lg transition-shadow cursor-pointer"
        :style="{ backgroundColor: 'var(--surface-primary)' }"
        @click="viewCourse(course.id)"
      >
        <template #cover>
          <div
            class="h-40 bg-gradient-to-r from-blue-500 to-purple-600 rounded-t-xl flex items-center justify-center"
          >
            <div class="text-white text-center p-4">
              <BookOpen :size="48" class="mx-auto mb-2" />
              <h3 class="font-bold text-lg">{{ course.title }}</h3>
            </div>
          </div>
        </template>

        <div class="mt-4">
          <div class="flex justify-between items-center mb-2">
            <span
              :class="`px-2 py-1 rounded-full text-xs font-medium ${getLevelClass(course.level)}`"
            >
              {{ course.level }}
            </span>
            <div class="flex items-center text-yellow-400">
              <Star :size="16" fill="currentColor" />
              <span class="ml-1 text-sm">{{ course.rating }}</span>
            </div>
          </div>

          <p class="text-gray-500 text-sm mb-3" :style="{ color: 'var(--text-secondary)' }">
            {{ course.description }}
          </p>

          <div class="flex justify-between items-center text-sm text-gray-500">
            <div class="flex items-center">
              <Clock :size="14" class="mr-1" />
              <span>{{ course.duration }}</span>
            </div>
            <div class="flex items-center">
              <User :size="14" class="mr-1" />
              <span>{{ course.lessons }} 课时</span>
            </div>
          </div>
        </div>
      </n-card>
    </div>

    <!-- 分页 -->
    <div class="mt-8 flex justify-center">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-size="pageSize"
        :item-count="courses.length"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  Search,
  BookOpen,
  Star,
  Clock,
  User,
  BookText,
  Code,
  Cpu,
  Palette,
  Database
} from 'lucide-vue-next'
import { NCard, NInput, NSelect, NPagination, NButton } from 'naive-ui'

// 课程数据
const courses = ref([
  {
    id: 1,
    title: 'JavaScript 基础入门',
    description: '从零开始学习 JavaScript，掌握变量、函数、对象等基础知识',
    level: '初级',
    rating: 4.8,
    duration: '12小时',
    lessons: 24,
    category: '前端开发'
  },
  {
    id: 2,
    title: 'Vue.js 深度解析',
    description: '深入学习 Vue.js 框架，掌握组件化开发和状态管理',
    level: '中级',
    rating: 4.9,
    duration: '20小时',
    lessons: 36,
    category: '前端开发'
  },
  {
    id: 3,
    title: '算法与数据结构',
    description: '掌握常用算法和数据结构，提升编程能力',
    level: '中级',
    rating: 4.7,
    duration: '25小时',
    lessons: 40,
    category: '算法'
  },
  {
    id: 4,
    title: 'Node.js 后端开发',
    description: '使用 Node.js 构建高性能后端服务',
    level: '中级',
    rating: 4.6,
    duration: '18小时',
    lessons: 30,
    category: '后端开发'
  },
  {
    id: 5,
    title: 'React 全栈开发',
    description: '使用 React 和相关技术栈构建现代 Web 应用',
    level: '高级',
    rating: 4.8,
    duration: '30小时',
    lessons: 50,
    category: '前端开发'
  },
  {
    id: 6,
    title: 'Python 数据分析',
    description: '使用 Python 进行数据分析和可视化',
    level: '中级',
    rating: 4.5,
    duration: '15小时',
    lessons: 28,
    category: '数据科学'
  }
])

// 筛选条件
const selectedCategory = ref(null)
const selectedLevel = ref(null)
const currentPage = ref(1)
const pageSize = 6
const totalPages = computed(() => Math.ceil(courses.value.length / pageSize))

// 分类选项
const categoryOptions = [
  { label: '前端开发', value: '前端开发' },
  { label: '后端开发', value: '后端开发' },
  { label: '算法', value: '算法' },
  { label: '数据科学', value: '数据科学' },
  { label: '移动开发', value: '移动开发' }
]

// 难度选项
const levelOptions = [
  { label: '初级', value: '初级' },
  { label: '中级', value: '中级' },
  { label: '高级', value: '高级' }
]

// 根据难度返回对应样式类
const getLevelClass = (level: string) => {
  switch (level) {
    case '初级':
      return 'bg-green-500/10 text-green-400'
    case '中级':
      return 'bg-yellow-500/10 text-yellow-400'
    case '高级':
      return 'bg-red-500/10 text-red-400'
    default:
      return 'bg-gray-500/10 text-gray-400'
  }
}

// 过滤后的课程
const filteredCourses = computed(() => {
  let result = courses.value

  if (selectedCategory.value) {
    result = result.filter((course) => course.category === selectedCategory.value)
  }

  if (selectedLevel.value) {
    result = result.filter((course) => course.level === selectedLevel.value)
  }

  // 分页
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return result.slice(start, end)
})

// 查看课程详情
const viewCourse = (id: number) => {
  console.log(`查看课程 ${id} 详情`)
  // 这里可以导航到课程详情页面
}
</script>

<style scoped>
.course-card {
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
}

.course-card:hover {
  transform: translateY(-4px);
}
</style>
