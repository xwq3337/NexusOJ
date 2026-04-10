# CLAUDE.md

本文档为 Claude Code（claude.ai/code）在此代码仓库中工作时提供指导。

## 项目概述

NexusOJ Frontend - 基于 Turborepo 管理的现代化在线判题平台前端项目，采用 Monorepo 架构。

## Monorepo 架构

本项目使用 Turborepo + pnpm workspace 管理，采用以下结构：

```
.
├── apps/               # 应用目录
│   ├── client/        # 客户端 (@nexusoj/client)
│   └── admin/         # 管理端 (@nexusoj/admin)
├── packages/          # 共享包目录
│   ├── ui/           # 共享 UI 组件 (@nexusoj/ui)
│   ├── config/       # 共享配置 (@nexusoj/config)
│   ├── eslint-config/ # ESLint 配置 (@nexusoj/eslint-config)
│   ├── server/       # 后端 API 服务(@nexusoj/server)
│   ├── typescript/   # TypeScript 配置 (@nexusoj/typescript)
│   └── utils/        # 工具函数 (@nexusoj/utils)
├── turbo.json        # Turborepo 配置
└── pnpm-workspace.yaml # pnpm workspace 配置
```

## 技术栈

- **Monorepo 管理**: Turborepo
- **包管理器**: pnpm (workspace)
- **构建工具**: Vite 6+
- **框架**: Vue 3 + TypeScript
- **状态管理**: Pinia
- **路由**: Vue Router

## Turborepo 命令

### 根命令

从根目录执行的命令会作用于所有应用：

```bash
# 开发模式（启动所有应用）
pnpm dev

# 构建（构建所有应用）
pnpm build

# 代码检查
pnpm lint

# 代码格式化
pnpm format

# 类型检查
pnpm type-check

# 清理构建产物
pnpm clean
```

### 单个应用命令

使用 `pnpm --filter` 针对特定应用执行命令：

```bash
# 启动客户端
pnpm --filter @nexusoj/client dev

# 启动管理端
pnpm --filter @nexusoj/admin dev

# 构建客户端
pnpm --filter @nexusoj/client build
```

## 工作空间包引用

在应用中引用共享包：

```json
// apps/client/package.json
{
  "dependencies": {
    "@nexusoj/ui": "workspace:*",
    "@nexusoj/utils": "workspace:*",
    "@nexusoj/config": "workspace:*"
  }
}
```

在代码中导入：

```typescript
import { SomeComponent } from '@nexusoj/ui'
import { someUtil } from '@nexusoj/utils'
import { API_BASE_URL } from '@nexusoj/config'
```

## Turbo Pipeline

`turbo.json` 定义了任务的依赖关系和缓存策略：

- `build`: 依赖其他包的 build，输出到 `dist/`
- `dev`: 不缓存，持久化任务
- `lint`: 依赖其他包的 lint
- `type-check`: 依赖其他包的 type-check

## 开发工作流

### 添加共享包

1. 在 `packages/` 下创建目录：`packages/new-package/`
2. 创建 `package.json`：
   ```json
   {
     "name": "@nexusoj/new-package",
     "version": "0.0.0",
     "private": true,
     "main": "./src/index.ts",
     "types": "./src/index.ts"
   }
   ```
3. 创建 `src/index.ts` 导出内容
4. 在需要使用的应用中：`pnpm add @nexusoj/new-package`

### 跨包调试

Turbo 会自动处理包依赖关系。当你修改共享包时：
1. 运行 `pnpm build` 重新构建受影响的应用
2. 或在开发模式下，Turbo 会检测变更并自动重新构建

## 重要约定

1. **包命名**: 所有应用和包使用 `@nexusoj/` 作用域
2. **工作空间引用**: 使用 `workspace:*` 协议引用内部包
3. **依赖提升**: 优先将公共依赖提升到根 `package.json`
4. **类型安全**: 所有包都必须配置 TypeScript
5. **代码风格**: 所有包使用统一的 ESLint 配置

## 故障排除

### 依赖问题

```bash
# 清理所有 node_modules
pnpm -r exec rm -rf node_modules
rm -rf node_modules

# 重新安装
pnpm install
```

### Turbo 缓存问题

```bash
# 清除 Turbo 缓存
rm -rf .turbo
```

### 构建问题

确保先构建依赖的包：
```bash
pnpm build --filter=@nexusoj/dependent-package...
```

## 应用特定文档

各应用的详细开发指南请参考：

- [Client 应用文档](apps/client/CLAUDE.md)
- [Admin 应用文档](apps/admin/CLAUDE.md)
