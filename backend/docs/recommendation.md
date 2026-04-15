# NexusOJ 推荐系统与用户画像实现方案

## Context

NexusOJ 是一个在线判题平台（类似 LeetCode），当前已有完整的用户、题目、提交记录、比赛、博客等数据模型，但缺少智能化推荐和用户分析能力。目标是基于现有数据构建用户画像系统并实现个性化题目推荐，帮助用户更高效地学习和刷题。

**核心约束**：单接口最多 2 次 MySQL 查询，大量使用 Redis 缓存，利用 Goroutine 并发。

---

## 整体架构

不新建 MySQL 表。画像数据存 Redis（热数据）+ MongoDB（持久化），推荐结果缓存 Redis。

```
提交代码 → 判题完成 → 触发画像增量更新（异步 Goroutine）
                              ↓
                     Redis 画像数据更新（Pipeline）
                              ↓
                     标记推荐缓存为脏 → 定时/按需重新计算推荐
```

---

## 一、用户画像（User Profile）

### 1.1 画像维度

| 维度 | 数据 | 来源 |
|------|------|------|
| **能力** | 各标签掌握度（0~1）、综合能力分、擅长/薄弱标签 | Record + Problem 的 tag/difficulty |
| **活跃度** | 热力图（今年/去年/前年/近一年每日提交数）、连续打卡天数、活跃时段分布 | Record 的 created_at |
| **偏好** | 常用语言及占比、偏好难度范围 | Record 的 language |
| **社交** | 好友数、参赛数、博客/题解数 | 各已有模型 |

### 1.2 Redis Key 设计

```
profile:{userID}:tag_stats          # Hash: tag -> {accepted, attempted, avg_difficulty}
profile:{userID}:activity           # String: JSON {streak, last_active, preferred_hours:[24]int}
profile:{userID}:heatmap:{year}      # Hash: "MM-DD" -> count（如 "04-15" -> 5），每年一个 key
# 例：profile:123:heatmap:2026  profile:123:heatmap:2025  profile:123:heatmap:2024
# 近一年（滚动365天）在 API 层从两个相邻年份 Hash 中动态拼接
profile:{userID}:preferences        # String: JSON {languages:{}, avg_difficulty}
profile:{userID}:ability            # String: JSON {overall_score, tag_scores:{}, strongest:[], weakest:[]}
profile:{userID}:social             # String: JSON {friend_count, contest_count, blog_count}
```

### 1.3 能力评分算法

```
tag_score[tag] = accepted[tag] / max(attempted[tag], 1)
加权 tag_score = Σ(difficulty × is_accepted) / Σ(difficulty × is_attempted)
overall_ability = 各 tag_score 的加权平均（权重 = attempted_count）
```

### 1.4 增量更新（每次提交触发）

提交判题完成后，在现有的 `go func()` 异步块中（[problem.go:226-237](controllers/problem.go#L226-L237)）增加一行，向 Profile Worker 发送事件：

```go
services.SubmitProfileUpdate(&services.ProfileUpdateEvent{
    UserID, ProblemID, Verdict, Language, Difficulty, Tags, Timestamp
})
```

Worker 收到事件后，通过 Redis Pipeline 一次性更新：
1. tag_stats（对应 tag 的 accepted/attempted +1）
2. heatmap（`HINCRBY profile:{userID}:heatmap:{currentYear} {MM-DD} 1`，更新 streak）
3. preferences（语言计数 +1，移动平均难度）

### 1.5 热力图详细设计

**数据结构**：每年一个 Redis Hash，field 为 `"MM-DD"`，value 为当天提交次数。

```
HINCRBY profile:123:heatmap:2026 "04-15" 1    # 今天提交一次
HGETALL profile:123:heatmap:2026              # 获取2026年全部数据
```

**增量更新**：每次提交时，只需 1 次 `HINCRBY` 操作即可。

**四个时间维度的读取**：

| 维度 | 实现方式 | Redis 调用 |
|------|----------|-----------|
| 今年 | `HGETALL profile:{uid}:heatmap:2026` | 1 次 |
| 去年 | `HGETALL profile:{uid}:heatmap:2025` | 1 次 |
| 前年 | `HGETALL profile:{uid}:heatmap:2024` | 1 次 |
| 近一年 | 今年 Hash + 去年 Hash 拼接（取去年同日之后的数据） | 2 次 |

API 一次返回四个维度，通过 Pipeline 并发读取最多 2 个 Hash（今年+去年），其余年份在客户端指定年份时再查。

**冷启动**：首次访问时，从 MySQL 全量计算历史热力图：
```sql
SELECT DATE(created_at) AS date, COUNT(*) AS count
FROM records
WHERE user_id = ? AND created_at >= ?
GROUP BY DATE(created_at)
```
写入对应年份的 Redis Hash，后续只需增量更新。

**MongoDB 持久化**：定时（每小时）将所有年份 Hash dump 到 MongoDB 的 `user_profiles.activity.heatmaps` 字段，防止 Redis 数据丢失。

---

## 二、推荐系统

### 2.1 五种推荐策略

| 策略 | 权重 | 逻辑 |
|------|------|------|
| **难度匹配** | 0.30 | 推荐难度略高于当前能力的题（最近发展区 +0.3） |
| **标签练习** | 0.25 | 找掌握度 0.2~0.7 的标签（练过但没掌握），推荐该标签题目 |
| **协作过滤** | 0.20 | 基于标签向量余弦相似度找相似用户，推荐他们做过但你没做的题 |
| **上下文感知** | 0.15 | 分析最近 10 次提交，某标签失败率高则降级推荐该标签简单题 |
| **新题推荐** | 0.10 | 按创建时间推荐未尝试过的新题 |

### 2.2 Redis 辅助索引（启动时从 MySQL 加载）

```
problems_by_difficulty    # Sorted Set: score=difficulty, member=problemID
tag_problems:{tag}        # Sorted Set: score=difficulty, member=problemID
problems_by_time          # Sorted Set: score=timestamp, member=problemID
cf:user_vectors           # Hash: userID -> JSON {tag:weighted_score}
```

### 2.3 推荐生成流程

1. 从 Redis 读取用户画像（0 次 MySQL）
2. 并发执行 5 种策略（Goroutine）
3. 合并去重，按加权分数排序
4. 结果缓存到 Redis，TTL 1 小时

---

## 三、MongoDB 持久化

```
Collection: user_profiles
{
    user_id, updated_at,
    ability: { overall_score, tag_scores, strongest_tags, weakest_tags },
    activity: { streak, longest_streak, heatmap, submission_frequency },
    preferences: { languages, preferred_difficulty_range },
    social: { friend_count, contest_count }
}

Collection: user_similarity
{
    user_id,
    similar_users: [{ user_id, similarity, common_solved }],
    updated_at
}
```

用途：Redis 缓存失效时的冷启动数据源，避免回查 MySQL。

---

## 四、API 设计

路由组：`/recommend`（需 JWT 认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/recommend/profile` | 获取当前用户完整画像 |
| GET | `/recommend/profile/:id` | 获取其他用户公开画像 |
| GET | `/recommend/problems` | 获取推荐题目列表（分页） |
| POST | `/recommend/refresh` | 强制刷新推荐结果 |
| GET | `/recommend/ability` | 获取能力分析（各标签掌握度详情） |
| GET | `/recommend/activity` | 获取活跃度统计（热力图、streak） |

---

## 五、文件变更清单

### 新建文件
| 文件 | 职责 | 预估行数 |
|------|------|----------|
| `services/profile_service.go` | 画像服务：Worker 队列 + 增量更新 + 读取 | ~400 |
| `services/recommendation_service.go` | 推荐服务：5 种策略 + 合并 + 缓存 | ~500 |
| `controllers/recommendation_controller.go` | 6 个 API 端点 | ~200 |

### 修改文件
| 文件 | 变更内容 |
|------|----------|
| `models/record.go` | 新增 ~5 个聚合查询方法（标签统计、活跃度、语言统计等） |
| `models/problem.go` | 新增 ~4 个查询方法（按难度范围、按标签、最新题目等） |
| `router/router.go` | 新增 `/recommend` 路由组（~10 行） |
| `controllers/problem.go` | 提交后 goroutine 中增加 1 行画像事件发送 |
| `controllers/cache.go` | 新增题目元数据索引加载函数 |

### 参考文件（实现时对照）
- [services/judge_queue.go](services/judge_queue.go) — Worker 队列模式
- [services/contest_service.go](services/contest_service.go) — Redis Pipeline + MongoDB 持久化
- [controllers/problem.go:226-237](controllers/problem.go#L226-L237) — 提交后集成点
- [controllers/cache.go](controllers/cache.go) — Redis 批量加载模式

---

## 六、实施阶段

### Phase 1：画像数据收集（基础）
1. `models/record.go` 和 `models/problem.go` 新增聚合查询方法
2. 创建 `services/profile_service.go`：InitProfileService + Worker + 增量更新
3. 在 `controllers/problem.go` 提交流程中集成画像更新事件
4. 构建题目元数据的 Redis 索引

### Phase 2：画像读取 API
1. 创建 `controllers/recommendation_controller.go`：GetMyProfile、GetAbilityAnalysis、GetActivityStats
2. `router/router.go` 添加路由
3. 实现 MongoDB 持久化（定时 dump）

### Phase 3：推荐引擎
1. 创建 `services/recommendation_service.go`：5 种推荐策略
2. 实现 GetRecommendations 和 RefreshRecommendations 端点
3. 构建协作过滤索引 + 定时更新 Worker

### Phase 4：优化
1. 冷启动处理（新用户默认推荐热门题）
2. Admin API：重建画像缓存
3. 性能调优

---

## 七、验证方式

1. **画像验证**：提交代码后检查 Redis 中对应 key 是否正确更新
2. **推荐验证**：调用 `GET /recommend/problems` 确认返回题目符合策略（难度合理、标签匹配）
3. **性能验证**：确保单个 API 响应时间 < 100ms（Redis 缓存命中），MySQL 查询 ≤ 2 次
4. **冷启动验证**：清空 Redis 后调用接口，确认能从 MongoDB 恢复或触发 MySQL 全量计算
