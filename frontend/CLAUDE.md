# CLAUDE.md

本文档为 Claude Code（claude.ai/code）在此代码仓库中工作时提供指导。

## 项目概述

NexusOJ 前端 - 一个基于 Vue 3 构建的现代化在线判题平台，具备题目求解、比赛、排行榜功能，并集成了代码编辑器和 AI 辅助的知识库。

## 项目架构

### 技术栈

• 框架：Vue 3.4+，使用 Composition API 和 TypeScript
• 构建工具：Vite 6+
• 路由：Vue Router，使用基于文件的路由组件
• 状态管理：Pinia，具备 localStorage 持久化
• UI 库：Naive UI 组件库
• 样式：Tailwind CSS 4+ 与 PostCSS
• 图标：Lucide Vue Next
• 代码编辑器：Ace Editor，支持多语言
• 富文本：Markdown-it，集成 KaTeX 以支持数学公式
• 数据可视化：ECharts 及其 Vue 封装
• HTTP 客户端：Axios，配置了拦截器
• 其他：Vueuse[useClipboard, useDebounce, useLocalStorage, useToggle]

### 目录结构

src/
├── assets/          # 主题 CSS 文件
├── components/      # 可复用的 Vue 组件
├── composables/     # Vue 组合式函数 (useTheme.ts, useTitle.ts)
├── constants/       # 常量定义 (每次编写代码前，请先阅读此文件)
├── layouts/         # 布局包装器 (带侧边栏的 Layout.vue， UserLayout.vue)
├── pages/           # 路由组件，与 URL 结构匹配
├── router/          # Vue Router 配置及路由定义
├── services/        # API 服务层 (Request.ts 处理 HTTP 调用)
├── stores/          # Pinia store，用于状态管理
├── types.ts         # TypeScript 类型定义
├── utils/           # 工具函数
├── App.vue          # 根组件，包含 ThemeProvider
└── main.ts          # 应用入口点

### 关键架构模式

1. 基于页面的路由：每个主要路由在 src/pages/ 目录下都有对应的页面组件。
2. 组合式 API：所有组件均使用 `<script setup lang="ts">` 模式。
3. 类型安全：启用严格模式的完整 TypeScript，并定义了自定义类型。
4. 主题系统：通过 useTheme() 组合式函数切换明暗模式的 CSS 变量。
5. 状态持久化：用户数据通过 localStorage 和 Pinia 插件持久化。

### 主题架构

• CSS 变量：所有颜色定义在 assets/dark.css 和 assets/light.css 中。
• 命名规范：变量遵循 --{分类}-{变体} 的命名模式。
• 使用方式：组件通过内联样式访问，如  :style="{ color: 'var(--text-primary)' }"。
• 英雄区域：使用独立的 --hero-* 变量以实现独立的样式控制。

### 组件约定

• 命名：组件使用 PascalCase（例如 ActivityChart.vue）。
• Props：使用 TypeScript 接口明确定义类型。
• 事件：使用描述性名称触发事件。
• 插槽：使用具名插槽以提高清晰度。
• 样式：优先使用 Tailwind 工具类，必要时在 `<style scoped>` 中编写自定义 CSS。

### 服务层

所有 API 调用都通过 src/services/api.ts 进行：
• Request.get(url, config) - GET 请求
• Request.post(url, data, config) - POST 请求
• 基础 URL 在服务中配置，并在 Vite 配置中代理。

### 状态管理 (Pinia)

src/stores/ 目录下的 Store：
• user.ts - 用户认证和资料数据
• 通过 localStorage 和 Pinia 插件实现持久化。

## 常见工作流程

### 添加新页面

1. 在 src/pages/[页面名称].vue 创建组件。
2. 在 src/router/index.ts 中添加路由：
    {
      path: '/页面路径',
      name: '页面名称',
      component: () => import('@/pages/页面名称.vue')
    }
3. 如需，在 src/components/Navbar.vue 中添加导航链接。
4. 如果需要新的数据结构，请更新 src/types.ts。

### 创建可复用组件

1. 在 src/components/[组件名称].vue 中创建。
2. 使用 TypeScript 接口定义 props：
    interface Props {
      data: Blog[]
      loading?: boolean
    }
3. 在父组件中导入并使用。
4. 如果在多个文件中使用，请导出。

### 添加主题变量

对于新的独立主题区块（如英雄区域）：

1. 在 assets/dark.css 和 assets/light.css 中都添加变量。
2. 使用命名规范：--{区块}-{属性}。
3. 通过内联样式访问： :style="{ color: 'var(--hero-title-color)' }"。

### 代码编辑器集成

平台使用 Ace Editor，具备以下特性：
• 支持语言：JavaScript, Python, Java, C, C++, Rust, Go
• 代码执行：通过 HTTP 与后端通信
• 配置：位于父组件页面中（例如 Problem.vue）。

## 重要约束

每次修改代码时，请先阅读src/constants/index.ts下的常量，以便统一管理

1. 使用 Tailwind 类：优先使用工具类，而非自定义 CSS。
    ◦ 错误示例：class="my-custom-class"
    ◦ 正确示例：class="flex items-center gap-4"

2. 渐变类名：在 Tailwind v4 中，使用 `bg-linear-to-*` 而非 `bg-gradient-to-*`。
    ◦ 错误示例：class="bg-gradient-to-r"
    ◦ 正确示例：class="bg-linear-to-r"

3. 任意值：尽可能使用标准的间距尺寸。
    ◦ 错误示例：class="min-w-[280px]"
    ◦ 正确示例：class="min-w-70"

4. 主题变量：所有可主题化的值都必须使用 CSS 变量。
    ◦ 错误示例：class="text-white dark:text-gray-800"
    ◦ 正确示例： :style="{ color: 'var(--text-primary)' }"

5. 响应式设计：采用移动端优先策略。
    ◦ 默认：移动端样式
    ◦ 大屏幕修饰符：md:、lg:、xl:

## 测试

当前未配置测试框架。添加测试时：

1. 安装 Vitest：`npm install -D vitest @vue/test-utils`
2. 在 package.json 中添加测试脚本
3. 在组件旁创建 `*.spec.ts` 或 `*.test.ts` 文件
4. 在 src/composables/__tests__/ 中测试组合式函数

API 集成
• 基础 URL：在 Vite 代理中配置，用于开发环境。
• 端点：在 src/services/api.ts 中定义。
• 请求/响应：通过 TypeScript 接口定义类型。
• 错误处理：使用 try-catch 块并进行控制台日志记录。

代码检查与格式化
• Prettier：保存时自动格式化，使用 2 空格缩进。
• TypeScript：在 tsconfig.json 中启用了严格模式。
• 提交前请运行 npm run format。

开发提示

1. 热模块替换：适用于大多数更改。
2. 路由守卫：在 src/router/index.ts 中配置。
3. Meta 标签：通过 src/composables/useTitle.ts 管理。
4. Naive UI 组件：遵循其官方文档的最佳实践。
5. Pinia 开发者工具：可在浏览器中使用以检查状态。

性能考量

1. 懒加载：路由使用动态导入以实现代码分割。
2. 图片优化：尽可能使用 WebP 或响应式图片。
3. 包体积：构建后监控 dist/ 目录大小。
4. API 防抖：为搜索/筛选操作实现防抖。

文件命名

• 组件：PascalCase（例如 ActivityChart.vue）。
• 组合式函数：camelCase，带 use 前缀（例如 useTheme.ts）。
• 类型：camelCase（例如 types/ 目录下的 blog.ts）。
• 工具函数：camelCase（例如 formatDate.ts）。
