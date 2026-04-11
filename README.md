# NexusOJ

一个全栈在线判题系统，支持多种编程语言，适用于算法竞赛、在线刷题和教学场景。

## 功能特性

- **题库系统** — 题目管理、全文搜索（支持中文 ngram）、代码提交与实时评测
- **比赛系统** — 支持公开/私密比赛、实时排行榜（SSE 推送）、封榜、比赛报告生成
- **题解模块** — 每题每用户一篇题解，支持草稿/公开/私密状态，Markdown 编辑
- **聊天系统** — 基于 WebSocket 的实时私聊、未读消息推送、好友管理
- **博客系统** — 用户博客发布、审核、回收站
- **训练模块** — 题单管理，按主题刷题
- **文件管理** — 分片上传/下载、云存储、Markdown 图片上传
- **管理后台** — 用户管理、比赛管理、题目管理、日志查看、性能监控

## 技术栈

### 后端技术栈
| --- | --- |
| 语言 | Go 1.25 |
| 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL（主存储）、Redis（缓存/SSE）、MongoDB（聊天记录/比赛报告） |
| 认证 | JWT |
| 密码加密 | Argon2id |
| ID 生成 | 雪花算法（yitter） |
| 日志 | Zap + Lumberjack 日志轮转 |
| 热重载 | Air |

### 前端技术栈
| --- | --- |
| 框架 | Vue 3 + TypeScript |
| 构建 | Vite 6 + Turborepo（monorepo） |
| 包管理 | pnpm workspaces |
| UI 库 | Naive UI（客户端）、Element Plus（管理端） |
| 状态管理 | Pinia |
| 代码编辑器 | CodeMirror 6 / Ace Editor |
| Markdown | VMD Editor + KaTeX |
| 图表 | ECharts |
| 样式 | Tailwind CSS 4（客户端）、SCSS（管理端） |

## 项目结构

```
NexusOJ/
├── backend/                # 后端服务（Go）
│   ├── config/             # INI 配置加载
│   ├── controllers/        # Gin HTTP 处理器
│   ├── dao/                # MySQL/Redis/MongoDB 客户端
│   ├── middleware/          # JWT 认证、请求日志
│   ├── migrations/         # GORM AutoMigrate
│   ├── models/             # GORM 模型 + 查询方法
│   ├── router/             # 路由定义
│   ├── services/           # 业务逻辑（判题队列、比赛状态工作器）
│   ├── utils/              # 工具函数（加密、日志、响应格式）
│   ├── config.ini          # 配置文件
│   └── Makefile            # 构建与部署命令
├── frontend/               # 前端应用（monorepo）
│   ├── apps/
│   │   ├── client/         # 用户端（刷题、比赛、聊天）
│   │   └── admin/          # 管理端（后台管理）
│   └── packages/
│       ├── types/          # 共享 TypeScript 类型
│       ├── ui/             # 共享 UI 组件
│       ├── utils/          # 共享工具函数
│       └── eslint-config/  # 共享 ESLint 配置
├── agent/                  # AI 代码助手
├── judge-docker/           # 判题沙箱（已弃用，参考 sandbox-cpp）
└── mobile/                 # 移动端（开发中）
```

## 快速开始

### 环境要求

- Go >= 1.25
- Node.js >= 22.0.0
- pnpm
- MySQL 8.0+
- Redis 6.0+
- MongoDB 5.0+

### 后端

```bash
cd backend

# 安装 Go 依赖
go mod download

# 编辑配置文件（填写数据库连接信息）
# vim config.ini

# 开发模式（热重载，需要 ~/go/bin/air）
make dev

# 或直接构建运行
make build && ./target/main
```

后端启动后监听 `http://localhost:8080`。

### 前端

```bash
cd frontend

# 安装依赖
pnpm install

# 启动全部应用
pnpm dev

# 或单独启动
pnpm dev:client    # 用户端，默认 http://localhost:3000
pnpm dev:admin     # 管理端，默认 http://localhost:8888
```

前端开发服务器会自动代理 API 请求到后端。

## 后端架构

### 分层结构

```
config/       → INI 配置加载，导出全局变量
dao/          → MySQL、Redis、MongoDB 客户端连接
models/       → GORM 模型结构体 + 数据库查询方法（存储库层）
services/     → 业务逻辑（判题队列、比赛排名、状态工作器）
controllers/  → Gin 处理器，调用 models/services
router/       → 路由定义
middleware/    → JWT 认证、请求日志
utils/        → 响应工具、日志、加密、JSON 工具
```

### 后台工作器

- **判题队列** — 5 个并发工作器，队列容量 100，处理代码提交评测
- **比赛状态工作器** — 每 60 秒轮询，自动更新比赛状态，触发报告生成和缓存初始化

### API 路由组

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

## 部署

```bash
cd backend

# 交叉编译到 Linux 并通过 rsync 部署到服务器
make deploy
```

部署脚本（`scripts/deploy.sh`）会自动完成编译、上传和 systemd 服务重启。生产环境配置以 `config.ini` 中的 `DataDir` 为文件存储根目录。

## 许可证

MIT License
