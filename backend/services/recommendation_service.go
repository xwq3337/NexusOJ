package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nexus/dao"
	"nexus/models"
	"sort"
	"sync"
	"time"
)

// ==================== 推荐数据结构 ====================

// RecommendedProblem 推荐的题目
type RecommendedProblem struct {
	ProblemID  int64    `json:"problem_id"`
	Title      string   `json:"title"`
	Difficulty float32  `json:"difficulty"`
	Tags       []string `json:"tags"`
	Score      float32  `json:"score"`  // 推荐分数 0.0~1.0
	Reason     string   `json:"reason"` // 推荐原因
}

type strategyResult struct {
	problems []RecommendedProblem
	weight   float32
}

const (
	recommendCacheKey = "recommend:%d:problems"
	recommendCacheTTL = time.Hour
)

// ==================== 推荐主入口 ====================

// GetRecommendations 获取推荐题目（带缓存）
func GetRecommendations(userID uint64, page, pageSize int, refresh bool) ([]RecommendedProblem, int, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf(recommendCacheKey, userID)

	// 尝试从缓存读取
	if !refresh {
		cached, err := dao.RedisClient.Get(ctx, cacheKey).Bytes()
		if err == nil && len(cached) > 0 {
			var all []RecommendedProblem
			if err := json.Unmarshal(cached, &all); err == nil {
				total := len(all)
				start := (page - 1) * pageSize
				if start >= total {
					return []RecommendedProblem{}, total, nil
				}
				end := start + pageSize
				if end > total {
					end = total
				}
				return all[start:end], total, nil
			}
		}
	}

	// 缓存未命中或强制刷新，计算推荐
	all, err := computeRecommendations(userID)
	if err != nil {
		return nil, 0, err
	}

	// 写入缓存
	data, _ := json.Marshal(all)
	dao.RedisClient.Set(ctx, cacheKey, data, recommendCacheTTL)

	total := len(all)
	start := (page - 1) * pageSize
	if start >= total {
		return []RecommendedProblem{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// ==================== 推荐算法核心 ====================

func computeRecommendations(userID uint64) ([]RecommendedProblem, error) {
	// 确保用户画像存在
	profile, err := EnsureProfile(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户已解决的题目 ID（用于排除）
	solvedIDs, _ := models.GetUserSolvedProblemIDs(userID)
	solvedSet := make(map[int64]bool, len(solvedIDs))
	for _, id := range solvedIDs {
		solvedSet[id] = true
	}

	// 并发执行五种推荐策略
	results := make([]strategyResult, 5)
	var wg sync.WaitGroup

	// 策略1: 难度匹配 (权重 0.30)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := difficultyBasedRecommend(profile, solvedSet)
		for i := range p {
			p[i].Score *= 0.30
		}
		results[0] = strategyResult{p, 0.30}
	}()

	// 策略2: 标签练习 (权重 0.25)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := tagBasedRecommend(profile, solvedSet)
		for i := range p {
			p[i].Score *= 0.25
		}
		results[1] = strategyResult{p, 0.25}
	}()

	// 策略3: 协作过滤 (权重 0.20)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := collaborativeFilterRecommend(userID, profile, solvedSet)
		for i := range p {
			p[i].Score *= 0.20
		}
		results[2] = strategyResult{p, 0.20}
	}()

	// 策略4: 上下文感知 (权重 0.15)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := contextualRecommend(userID, profile, solvedSet)
		for i := range p {
			p[i].Score *= 0.15
		}
		results[3] = strategyResult{p, 0.15}
	}()

	// 策略5: 新题推荐 (权重 0.10)
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := freshProblemRecommend(solvedSet)
		for i := range p {
			p[i].Score *= 0.10
		}
		results[4] = strategyResult{p, 0.10}
	}()

	wg.Wait()

	// 合并去重，按分数排序
	return mergeAndRank(results), nil
}

// ==================== 策略1: 难度匹配 ====================

func difficultyBasedRecommend(profile *UserProfile, solvedSet map[int64]bool) []RecommendedProblem {
	// 目标难度 = 当前能力 * 3 + 0.3（最近发展区）
	targetDiff := profile.Ability.OverallScore*3.0 + 0.3
	minDiff := targetDiff - 0.3
	maxDiff := targetDiff + 0.5
	if minDiff < 0 {
		minDiff = 0
	}

	excludeIDs := mapKeys(solvedSet)
	problems, err := models.Problem{}.GetProblemsByDifficultyRange(minDiff, maxDiff, excludeIDs, 20)
	if err != nil {
		return nil
	}

	result := make([]RecommendedProblem, 0, len(problems))
	for _, p := range problems {
		if solvedSet[p.ID] {
			continue
		}
		// 距离目标难度越近分数越高
		distance := abs32(p.Difficulty - targetDiff)
		score := 1.0 - distance/1.0
		if score < 0.1 {
			score = 0.1
		}
		result = append(result, RecommendedProblem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			Score:      score,
			Reason:     "difficulty_match",
		})
	}
	return result
}

// ==================== 策略2: 标签练习 ====================

func tagBasedRecommend(profile *UserProfile, solvedSet map[int64]bool) []RecommendedProblem {
	// 找出掌握度 0.2~0.7 的标签（练过但没掌握）
	var weakTags []string
	for tag, score := range profile.Ability.TagScores {
		if score >= 0.2 && score <= 0.7 {
			weakTags = append(weakTags, tag)
		}
	}
	if len(weakTags) == 0 {
		// 如果没有弱标签，取最弱的3个
		weakest := profile.Ability.WeakestTags
		if len(weakest) > 3 {
			weakest = weakest[:3]
		}
		weakTags = weakest
	}
	if len(weakTags) == 0 {
		return nil
	}

	excludeIDs := mapKeys(solvedSet)
	problems, err := models.Problem{}.GetProblemsByTags(weakTags, excludeIDs, 20)
	if err != nil {
		return nil
	}

	result := make([]RecommendedProblem, 0, len(problems))
	for _, p := range problems {
		if solvedSet[p.ID] {
			continue
		}
		// 弱标签的题目分数基于标签掌握度的反向
		maxWeakness := float32(0)
		for _, tag := range p.Tags {
			if score, ok := profile.Ability.TagScores[tag]; ok {
				weakness := 1.0 - score
				if weakness > maxWeakness {
					maxWeakness = weakness
				}
			}
		}
		result = append(result, RecommendedProblem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			Score:      maxWeakness,
			Reason:     "tag_practice",
		})
	}
	return result
}

// ==================== 策略3: 协作过滤 ====================

func collaborativeFilterRecommend(userID uint64, profile *UserProfile, solvedSet map[int64]bool) []RecommendedProblem {
	// 构建用户标签向量
	userVector := profile.Ability.TagScores

	// 查找相似用户：读取 cf:user_vectors 中所有用户向量
	ctx := context.Background()
	allVectors, err := dao.RedisClient.HGetAll(ctx, "cf:user_vectors").Result()
	if err != nil || len(allVectors) == 0 {
		// 索引不存在，返回空（后台 Worker 会构建索引）
		return nil
	}

	type similarUser struct {
		uid      uint64
		similarity float32
	}
	var similar []similarUser

	for uidStr, vecJSON := range allVectors {
		var uid uint64
		fmt.Sscanf(uidStr, "%d", &uid)
		if uid == userID {
			continue
		}
		var otherVec map[string]float32
		if err := json.Unmarshal([]byte(vecJSON), &otherVec); err != nil {
			continue
		}
		sim := cosineSimilarity(userVector, otherVec)
		if sim > 0.3 { // 只保留相似度 > 0.3 的用户
			similar = append(similar, similarUser{uid, sim})
		}
	}

	// 按相似度排序，取前5
	sort.Slice(similar, func(i, j int) bool {
		return similar[i].similarity > similar[j].similarity
	})
	if len(similar) > 5 {
		similar = similar[:5]
	}

	if len(similar) == 0 {
		return nil
	}

	// 获取相似用户已解决的题目，合并统计
	problemScore := make(map[int64]float32)
	for _, su := range similar {
		solved, err := models.GetUserSolvedProblemIDs(su.uid)
		if err != nil {
			continue
		}
		for _, pid := range solved {
			if !solvedSet[pid] {
				problemScore[pid] += su.similarity
			}
		}
	}

	if len(problemScore) == 0 {
		return nil
	}

	// 按分数排序取 top
	type pidScore struct {
		pid   int64
		score float32
	}
	var sorted []pidScore
	for pid, score := range problemScore {
		sorted = append(sorted, pidScore{pid, score})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})
	if len(sorted) > 20 {
		sorted = sorted[:20]
	}

	// 归一化分数
	maxScore := sorted[0].score

	result := make([]RecommendedProblem, 0, len(sorted))
	for _, ps := range sorted {
		// 查询题目信息（批量可以优化，这里简单处理）
		p, err := models.Problem{}.GetProblemInfoWithoutUsername(ps.pid)
		if err != nil {
			continue
		}
		result = append(result, RecommendedProblem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			Score:      ps.score / maxScore,
			Reason:     "similar_users",
		})
	}
	return result
}

// ==================== 策略4: 上下文感知 ====================

func contextualRecommend(userID uint64, profile *UserProfile, solvedSet map[int64]bool) []RecommendedProblem {
	// 获取最近 10 次提交
	recentRecords, err := models.GetRecentUserRecords(userID, 10)
	if err != nil || len(recentRecords) == 0 {
		return nil
	}

	// 统计最近提交中各标签的失败率
	tagFailRate := make(map[string][2]int) // [fail_count, total_count]
	// 获取最近提交的题目标签
	for _, r := range recentRecords {
		problem, err := models.Problem{}.GetProblemInfoWithoutUsername(r.ProblemId)
		if err != nil {
			continue
		}
		isFail := r.Verdict != models.Accepted
		for _, tag := range problem.Tags {
			counts := tagFailRate[tag]
			counts[1]++ // total
			if isFail {
				counts[0]++ // fail
			}
			tagFailRate[tag] = counts
		}
	}

	// 找出失败率 > 0.7 的标签
	var struggleTags []string
	for tag, counts := range tagFailRate {
		if counts[1] >= 2 {
			failRate := float32(counts[0]) / float32(counts[1])
			if failRate > 0.7 {
				struggleTags = append(struggleTags, tag)
			}
		}
	}

	if len(struggleTags) == 0 {
		return nil
	}

	// 推荐这些标签下的简单题（难度降级）
	excludeIDs := mapKeys(solvedSet)
	// 使用较低难度范围
	problems, err := models.Problem{}.GetProblemsByDifficultyRange(0, profile.Preferences.AvgDifficulty, excludeIDs, 15)
	if err != nil {
		return nil
	}

	result := make([]RecommendedProblem, 0)
	for _, p := range problems {
		if solvedSet[p.ID] {
			continue
		}
		// 只推荐包含挣扎标签的题目
		hasStruggleTag := false
		for _, tag := range p.Tags {
			for _, st := range struggleTags {
				if tag == st {
					hasStruggleTag = true
					break
				}
			}
			if hasStruggleTag {
				break
			}
		}
		if !hasStruggleTag {
			continue
		}
		result = append(result, RecommendedProblem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			Score:      0.8,
			Reason:     "contextual",
		})
	}
	return result
}

// ==================== 策略5: 新题推荐 ====================

func freshProblemRecommend(solvedSet map[int64]bool) []RecommendedProblem {
	excludeIDs := mapKeys(solvedSet)
	problems, err := models.Problem{}.GetFreshProblems(excludeIDs, 10)
	if err != nil {
		return nil
	}

	result := make([]RecommendedProblem, 0, len(problems))
	for i, p := range problems {
		if solvedSet[p.ID] {
			continue
		}
		// 越新分数越高
		score := float32(1.0) - float32(i)*0.1
		if score < 0.1 {
			score = 0.1
		}
		result = append(result, RecommendedProblem{
			ProblemID:  p.ID,
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			Score:      score,
			Reason:     "fresh",
		})
	}
	return result
}

// ==================== 合并与排序 ====================

func mergeAndRank(results []strategyResult) []RecommendedProblem {
	merged := make(map[int64]RecommendedProblem)

	for _, sr := range results {
		for _, p := range sr.problems {
			if existing, ok := merged[p.ProblemID]; ok {
				// 已存在，累加分数
				existing.Score += p.Score
				merged[p.ProblemID] = existing
			} else {
				merged[p.ProblemID] = p
			}
		}
	}

	// 转为切片并排序
	all := make([]RecommendedProblem, 0, len(merged))
	for _, p := range merged {
		all = append(all, p)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Score > all[j].Score
	})

	return all
}

// ==================== 协作过滤索引构建 ====================

// BuildCollaborativeFilterIndex 构建用户标签向量索引（定时任务调用）
func BuildCollaborativeFilterIndex() {
	log.Println("Building collaborative filter index...")
	ctx := context.Background()

	// 获取所有用户
	users, err := models.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get users for CF index: %v", err)
		return
	}

	pipe := dao.RedisClient.Pipeline()
	for _, user := range users {
		// 确保画像存在
		profile, err := EnsureProfile(user.ID)
		if err != nil {
			continue
		}
		if len(profile.Ability.TagScores) == 0 {
			continue
		}
		vecJSON, _ := json.Marshal(profile.Ability.TagScores)
		pipe.HSet(ctx, "cf:user_vectors", user.ID, vecJSON)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Failed to save CF index: %v", err)
	} else {
		log.Printf("CF index built for %d users", len(users))
	}
}

// ==================== 辅助函数 ====================

func cosineSimilarity(vecA, vecB map[string]float32) float32 {
	var dotProduct, normA, normB float32
	for k, va := range vecA {
		if vb, ok := vecB[k]; ok {
			dotProduct += va * vb
		}
		normA += va * va
	}
	for _, vb := range vecB {
		normB += vb * vb
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (sqrt32(normA) * sqrt32(normB))
}

func sqrt32(x float32) float32 {
	return float32(sqrt(float64(x)))
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 10; i++ {
		z = z - (z*z-x)/(2*z)
	}
	return z
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func mapKeys(m map[int64]bool) []int64 {
	keys := make([]int64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
