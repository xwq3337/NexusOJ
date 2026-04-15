# NexusOJ

全栈在线判题系统，支持多种编程语言实时评测，适用于算法竞赛、在线刷题和教学场景。

## 功能概览

| 模块 | 说明 |
| --- | --- |
| 题库系统 | 题目管理、全文搜索（中文 ngram）、多语言代码提交与实时评测 |
| 比赛系统 | ACM/OI 赛制、公开/私密比赛、实时排行榜（SSE）、封榜、比赛报告 |
| 题解模块 | 每题每用户一篇题解，草稿/公开/私密状态，Markdown + KaTeX 编辑 |
| 聊天系统 | WebSocket 实时私聊、未读消息推送、好友管理 |
| 博客系统 | 用户博客发布、审核流程、回收站 |
| 训练模块 | 题单管理，按主题刷题 |
| AI 助手 | RAG 知识库检索 + 大模型流式对话，辅助算法学习 |
| 文件管理 | 分片上传/下载、云存储目录、Markdown 图片上传 |
| 管理后台 | 用户/比赛/题目管理、日志查看、性能监控 |

## 技术栈

### 后端 — Go

| 项 | 说明 |
| --- | --- |
| 语言 | Go 1.25 |
| 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL（主存储）、Redis（缓存/SSE/排行榜）、MongoDB（聊天记录/比赛报告） |
| 认证 | JWT（Argon2id 密码加密） |
| ID 生成 | 雪花算法（yitter） |
| 日志 | Zap + Lumberjack 轮转 |
| 热重载 | Air |

### 前端 — Vue 3 Monorepo

| 项 | 说明 |
| --- | --- |
| 框架 | Vue 3 + TypeScript |
| 构建 | Vite 6 + Turborepo（pnpm workspaces） |
| 客户端 UI | Naive UI + Tailwind CSS 4 |
| 管理端 UI | Element Plus + Ant Design Vue |
| 状态管理 | Pinia |
| 代码编辑器 | CodeMirror 6 / Ace Editor |
| Markdown | VMD Editor + KaTeX |
| 图表 | ECharts |

### AI 后端 — Node.js

| 项 | 说明 |
| --- | --- |
| 框架 | Express 5 + TypeScript |
| LLM | LangChain + OpenAI 兼容接口 |
| 向量数据库 | Milvus（RAG 知识库检索） |
| Embedding | 智谱 Embedding-2 |

### 桌面端 — Tauri

| 项 | 说明 |
| --- | --- |
| 框架 | Tauri 2（Rust 后端）+ Vue 3 |
| UI | Element Plus |

## 项目结构

```text
NexusOJ/
├── backend/                  # Go 后端服务
│   ├── config/               # INI 配置加载
│   ├── controllers/          # Gin HTTP 处理器
│   ├── dao/                  # MySQL / Redis / MongoDB 客户端
│   ├── middleware/            # JWT 认证、请求日志
│   ├── migrations/           # GORM AutoMigrate + 全文索引
│   ├── models/               # GORM 模型 + 查询方法
│   ├── router/               # 路由定义
│   ├── services/             # 业务逻辑（判题队列、比赛排名、状态工作器）
│   ├── utils/                # 工具函数（加密、日志、响应格式）
│   ├── config.ini            # 配置文件（数据库、服务端口等）
│   └── Makefile              # 构建 / 开发 / 部署命令
│
├── frontend/                 # Vue 3 前端（Turborepo monorepo）
│   ├── apps/
│   │   ├── client/           # 用户端 — 刷题、比赛、聊天
│   │   └── admin/            # 管理端 — 后台管理
│   └── packages/
│       ├── types/            # 共享 TypeScript 类型
│       ├── ui/               # 共享 UI 组件
│       ├── server/           # 后端 API 封装层
│       ├── utils/            # 共享工具函数
│       └── eslint-config/    # 共享 ESLint 配置
│
├── AI-backend/               # AI 助手后端
│   ├── src/
│   │   ├── controllers/      # 聊天控制器（SSE 流式响应）
│   │   └── tools/            # RAG 检索工具（Milvus 向量检索）
│   ├── system_prompt.txt     # 系统提示词
│   └── index.ts              # Express 入口
│
├── dekstop/                  # Tauri 桌面客户端
│   ├── src/                  # Vue 3 前端
│   └── src-tauri/            # Rust 后端
│
└── mobile/                   # 移动端（规划中）
```

## 系统架构

```text
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  Client Web │  │  Admin Web  │  │  Desktop App│
│  :3000      │  │  :8888      │  │  (Tauri)    │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │
       └────────────────┼────────────────┘
                        │  REST API / WebSocket / SSE
                 ┌──────▼──────┐
                 │  Go Backend │
                 │  :8080      │
                 └──┬───┬───┬──┘
                    │   │   │
          ┌─────────┘   │   └─────────┐
          ▼             ▼             ▼
    ┌──────────┐  ┌──────────┐  ┌──────────┐
    │  MySQL   │  │  Redis   │  │  MongoDB │
    │  主存储   │  │  缓存/SSE │  │  聊天/报告│
    └──────────┘  └──────────┘  └──────────┘

       ┌──────────────────────────┐
       │  AI Backend  :3000       │
       │  Express + LangChain     │
       └──────┬───────────────────┘
              │
       ┌──────▼──────┐
       │   Milvus    │
       │  向量知识库   │
       └─────────────┘
```

## 快速开始

### 环境要求

- Go >= 1.25
- Node.js >= 22.0.0、pnpm >= 10.0.0
- MySQL 8.0+、Redis 6.0+、MongoDB 5.0+
- （可选）Milvus 向量数据库（AI 功能需要）

### 1. 后端

```bash
cd backend

# 安装依赖
go mod download

# 编辑配置（填写数据库连接信息）
cp config.example.ini config.ini
vim config.ini

# 开发模式（热重载，需要 air）
make dev

# 或直接构建运行
make build && ./target/main
```

后端监听 `http://localhost:8080`。

### 2. 前端

```bash
cd frontend

# 安装依赖
pnpm install

# 启动全部应用
pnpm dev

# 或单独启动
pnpm dev:client    # 用户端 http://localhost:3000
pnpm dev:admin     # 管理端 http://localhost:8888
```

开发服务器会自动代理 API 请求到后端。

### 3. AI 后端

```bash
cd AI-backend

# 安装依赖
npm install

# 配置环境变量
cp .env.example .env
# 编辑 .env，填入 Milvus / LLM API 等配置

# 开发模式
npm run dev

# 或编译运行
npm run build && npm start
```

AI 后端监听 `http://localhost:3000`（端口可通过 `.env` 配置）。

### 4. 桌面端

```bash
cd dekstop

# 安装依赖
pnpm install

# 开发模式（需要 Tauri CLI）
pnpm tauri:dev

# 构建桌面应用
pnpm tauri build
```

## 后端架构

### 分层结构

```text
config/       → INI 配置加载，导出全局变量
dao/          → MySQL / Redis / MongoDB 客户端连接
models/       → GORM 模型结构体 + 数据库查询方法（存储库层）
services/     → 业务逻辑（判题队列、比赛排名、状态工作器）
controllers/  → Gin 处理器，调用 models / services
router/       → 路由定义
middleware/    → JWT 认证、请求日志
utils/        → 响应工具、日志、加密
```

### 后台工作器

- **判题队列** — 5 个并发 Worker，队列容量 100，处理代码提交评测
- **比赛状态工作器** — 每 60 秒轮询，自动更新比赛状态（未开始 → 进行中 → 已结束），触发报告生成和缓存初始化

### API 路由

| 路由前缀 | 说明 |
| --- | --- |
| `/user` | 用户注册、登录、信息管理 |
| `/problem` | 题目管理、代码提交 |
| `/record` | 提交记录查询 |
| `/contest` | 比赛列表、报名、排行榜、SSE 实时推送 |
| `/blog` | 博客发布与管理 |
| `/solution` | 题解模块 |
| `/training` | 训练题单 |
| `/file`、`/upload`、`/download` | 文件管理、分片上传下载 |
| `/ws/chat` | WebSocket 聊天 |
| `/admin` | 管理后台接口 |

## AI 助手工作流

1. 前端发送聊天消息到 AI 后端 `/chat`
2. AI 后端从用户消息中提取关键词，通过 Embedding 模型向量化
3. 在 Milvus 中检索相关知识库文档（RAG）
4. 将检索结果注入 System Prompt，调用大模型流式生成回答
5. 以 SSE 方式实时推送给前端

## 部署

```bash
cd backend

# 交叉编译到 Linux 并通过 rsync 部署
make deploy
```

部署脚本（`scripts/deploy.sh`）会自动完成编译、上传和 systemd 服务重启。生产环境以 `config.ini` 中的 `DataDir` 为文件存储根目录。

## 许可证

MIT License
