# NexusOJ Frontend

基于 Turborepo 管理的现代化在线判题平台前端项目。

## 项目结构

```
.
├── apps/
│   ├── client/          # 客户端应用
│   └── admin/           # 管理端应用
├── packages/
│   ├── ui/              # 共享 UI 组件
│   ├── config/          # 共享配置
│   ├── eslint-config/   # ESLint 配置
│   ├── typescript/      # TypeScript 配置
│   ├── server/          # 后端 API 服务
│   └── utils/           # 共享工具函数
├── package.json         # 根 package.json
├── turbo.json           # Turborepo 配置
└── pnpm-workspace.yaml  # PNPM 工作空间配置
```

## 技术栈

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **包管理器**: pnpm
- **Monorepo 管理**: Turborepo
- **状态管理**: Pinia
- **路由**: Vue Router

## 快速开始

### 前置要求

- Node.js >= 22.0.0
- pnpm >= 10.0.0

### 安装依赖

```bash
pnpm install
```

### 开发

启动所有应用的开发服务器：

```bash
pnpm dev
```

启动单个应用：

```bash
# 启动客户端
pnpm dev:client

# 启动管理端
pnpm dev:admin
```

### 构建

构建所有应用：

```bash
pnpm build
```

构建单个应用：

```bash
# 构建客户端
pnpm build:client

# 构建管理端
pnpm build:admin
```

### 其他命令

```bash
# 代码检查
pnpm lint

# 代码格式化
pnpm format

# 类型检查
pnpm type-check

# 清理构建产物
pnpm clean
```

## 应用

### Client (@nexusoj/client)

用户端应用，提供题目求解、比赛、排行榜等功能。

**技术栈**:
- UI: Naive UI
- 编辑器: CodeMirror, Ace Editor
- 图表: ECharts
- 样式: Tailwind CSS

### Admin (@nexusoj/admin)

管理端应用，提供题目、比赛、博客等内容管理功能。

**技术栈**:
- UI: Element Plus, Ant Design Vue
- 编辑器: VMD Editor
- 图表: ECharts
- 样式: SCSS

## 共享包

### @nexusoj/ui

共享 UI 组件库，包含在多个应用中使用的通用组件。

### @nexusoj/config

共享配置，包括 API 端点、环境变量等。

### @nexusoj/eslint-config

共享 ESLint 配置，确保代码风格一致。

### @nexusoj/typescript

共享 TypeScript 配置。

### @nexusoj/utils

共享工具函数库。

## 开发指南

### 添加新应用

1. 在 `apps/` 目录下创建新应用
2. 初始化 `package.json`，设置名称为 `@nexusoj/your-app`
3. 在根 `package.json` 中添加相应的脚本

### 添加共享包

1. 在 `packages/` 目录下创建新包
2. 初始化 `package.json`，设置名称为 `@nexusoj/your-package`
3. 在应用中通过 `pnpm add @nexusoj/your-package` 引入

## 许可证

MIT
