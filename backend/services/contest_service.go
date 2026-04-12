package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nexus/dao"
	"nexus/models"
	"nexus/utils"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	redisContestPassword       = "contest:%s:password"
	redisContestParticipants   = "contest:%s:participants"
	redisContestRanking        = "contest:%s:ranking"
	redisContestRankDetail     = "contest:%s:ranking:detail"
	redisContestInfo           = "contest:%s:info"
	redisContestParticipantBit = "contest:%s:participants:bit"  // bitmap: key=contestID
	redisContestProblemStatus  = "contest:%s:problem_status:%d" // hash: contestID:userID
)

// globalMapper 全局用户ID→紧凑索引映射，供所有比赛的 bitmap 共享
var globalMapper = dao.NewMappedUserIDstrategy("global")

// SetParticipantBit 设置用户在某比赛 bitmap 中的注册位
func SetParticipantBit(contestID string, userID uint64) {
	idx := globalMapper.GetOrAssignIndex(userID)
	key := fmt.Sprintf(redisContestParticipantBit, contestID)
	dao.RedisClient.SetBit(ctx, key, int64(idx), 1)
}

// BatchCheckRegistration 批量检查用户是否报名了多个比赛
// 返回 map[contestID]bool，使用 Redis Pipeline 一次完成
func BatchCheckRegistration(contestIDs []string, userID uint64) map[string]bool {
	result := make(map[string]bool, len(contestIDs))

	idx := globalMapper.GetIndex(userID)
	if idx == 0 {
		// 用户从未映射过，不可能报名任何比赛
		for _, cid := range contestIDs {
			result[cid] = false
		}
		return result
	}

	// Pipeline: 一次发送所有 GETBIT
	pipe := dao.RedisClient.Pipeline()
	cmds := make([]*redis.IntCmd, len(contestIDs))
	for i, cid := range contestIDs {
		key := fmt.Sprintf(redisContestParticipantBit, cid)
		cmds[i] = pipe.GetBit(ctx, key, int64(idx))
	}
	pipe.Exec(ctx)

	for i, cid := range contestIDs {
		result[cid] = cmds[i].Val() == 1
	}
	return result
}

// ctx 用于 Redis 操作
var ctx = context.Background()

// ProblemStatusEntry 题目状态缓存条目
type ProblemStatusEntry struct {
	Status string `json:"status"` // accepted, wrong, unattempted
	Score  int32  `json:"score"`
}

// problemStatusEntryJSON 用于 HSet 的 JSON 序列化
type problemStatusEntryJSON struct {
	Status string `json:"status"`
	Score  int32  `json:"score"`
}

// SetProblemStatusCache 写入单题状态缓存
func SetProblemStatusCache(contestID string, userID uint64, label string, status string, score int32, ttl time.Duration) {
	key := fmt.Sprintf(redisContestProblemStatus, contestID, userID)
	val, _ := json.Marshal(problemStatusEntryJSON{Status: status, Score: score})
	pipe := dao.RedisClient.Pipeline()
	pipe.HSet(ctx, key, label, string(val))
	pipe.Expire(ctx, key, ttl)
	pipe.Exec(ctx)
}

// GetProblemStatusCache 获取用户在某比赛所有题目的状态缓存
// 返回 map[label] -> ProblemStatusEntry，缓存不存在返回 nil
func GetProblemStatusCache(contestID string, userID uint64) map[string]ProblemStatusEntry {
	key := fmt.Sprintf(redisContestProblemStatus, contestID, userID)
	result, err := dao.RedisClient.HGetAll(ctx, key).Result()
	if err != nil || len(result) == 0 {
		return nil
	}
	entries := make(map[string]ProblemStatusEntry, len(result))
	for label, val := range result {
		var entry ProblemStatusEntry
		if err := json.Unmarshal([]byte(val), &entry); err != nil {
			continue
		}
		entries[label] = entry
	}
	return entries
}

// BatchSetProblemStatusCache 批量写入题目状态缓存（从 rankItems 回填）
func BatchSetProblemStatusCache(contestID string, userID uint64, items []models.ContestRankItem, ttl time.Duration) {
	key := fmt.Sprintf(redisContestProblemStatus, contestID, userID)
	pipe := dao.RedisClient.Pipeline()
	for _, item := range items {
		status := "wrong"
		if item.IsAccepted {
			status = "accepted"
		}
		val, _ := json.Marshal(problemStatusEntryJSON{Status: status, Score: item.Score})
		pipe.HSet(ctx, key, item.ProblemLabel, string(val))
	}
	pipe.Expire(ctx, key, ttl)
	pipe.Exec(ctx)
}

// InvalidateProblemStatusCache 删除用户在某比赛的题目状态缓存
func InvalidateProblemStatusCache(contestID string, userID uint64) {
	key := fmt.Sprintf(redisContestProblemStatus, contestID, userID)
	dao.RedisClient.Del(ctx, key)
}

// CacheContestPassword 缓存比赛密码到 Redis
func CacheContestPassword(contestID, password string, endAt time.Time) error {
	key := fmt.Sprintf(redisContestPassword, contestID)
	ttl := time.Until(endAt)
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return dao.RedisClient.Set(ctx, key, password, ttl).Err()
}

// VerifyContestPassword 验证比赛密码（Argon2）
func VerifyContestPassword(contestID, password string) (bool, error) {
	key := fmt.Sprintf(redisContestPassword, contestID)
	stored, err := dao.RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return utils.VerifyPassword(password, stored)
}

// CacheParticipant 缓存参赛者到 Redis SET + Bitmap
func CacheParticipant(contestID string, userID uint64, endAt time.Time) error {
	key := fmt.Sprintf(redisContestParticipants, contestID)
	ttl := time.Until(endAt) + 24*time.Hour
	if ttl <= 0 {
		ttl = 48 * time.Hour
	}
	pipe := dao.RedisClient.Pipeline()
	pipe.SAdd(ctx, key, userID)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)

	// 同步到 bitmap
	SetParticipantBit(contestID, userID)
	return err
}

// IsParticipantCached 检查 Redis 中是否已报名
func IsParticipantCached(contestID, userID string) (bool, error) {
	key := fmt.Sprintf(redisContestParticipants, contestID)
	return dao.RedisClient.SIsMember(ctx, key, userID).Result()
}

// GetCachedParticipantCount 获取缓存的参赛人数
func GetCachedParticipantCount(contestID string) (int64, error) {
	key := fmt.Sprintf(redisContestParticipants, contestID)
	return dao.RedisClient.SCard(ctx, key).Result()
}

// UpdateACMRanking 更新 ACM 排名到 Redis Sorted Set
// ACM 排名分 = solved * 1000000 - penalty_seconds
func UpdateACMRanking(contestID string, userID uint64, solved int32, penaltySeconds int32, detail map[string]models.ContestProblemRankDetail) error {
	score := float64(solved)*1000000 - float64(penaltySeconds)

	pipe := dao.RedisClient.Pipeline()
	// 更新 Sorted Set
	rankKey := fmt.Sprintf(redisContestRanking, contestID)
	pipe.ZAdd(ctx, rankKey, &redis.Z{
		Score:  score,
		Member: userID,
	})
	pipe.Expire(ctx, rankKey, 72*time.Hour)

	// 更新详情 Hash
	detailKey := fmt.Sprintf(redisContestRankDetail, contestID)
	detailJSON, _ := json.Marshal(detail)
	pipe.HSet(ctx, detailKey, userID, string(detailJSON))
	pipe.Expire(ctx, detailKey, 72*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// UpdateOIRanking 更新 OI 排名到 Redis Sorted Set
// OI 排名分 = 总得分
func UpdateOIRanking(contestID string, userID uint64, totalScore int32, detail map[string]models.ContestProblemRankDetail) error {
	pipe := dao.RedisClient.Pipeline()
	rankKey := fmt.Sprintf(redisContestRanking, contestID)
	pipe.ZAdd(ctx, rankKey, &redis.Z{
		Score:  float64(totalScore),
		Member: userID,
	})
	pipe.Expire(ctx, rankKey, 72*time.Hour)

	detailKey := fmt.Sprintf(redisContestRankDetail, contestID)
	detailJSON, _ := json.Marshal(detail)
	pipe.HSet(ctx, detailKey, userID, string(detailJSON))
	pipe.Expire(ctx, detailKey, 72*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// GetContestRankingFromRedis 从 Redis 获取比赛排名
func GetContestRankingFromRedis(contestID string, count int64) ([]models.ContestRankSummary, error) {
	rankKey := fmt.Sprintf(redisContestRanking, contestID)
	detailKey := fmt.Sprintf(redisContestRankDetail, contestID)

	// 从 Sorted Set 获取排名(降序)
	results, err := dao.RedisClient.ZRevRangeWithScores(ctx, rankKey, 0, count-1).Result()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	// 批量获取详情
	userIDs := make([]string, len(results))
	for i, z := range results {
		userIDs[i] = z.Member.(string)
	}
	detailsJSON, err := dao.RedisClient.HMGet(ctx, detailKey, userIDs...).Result()
	if err != nil {
		return nil, err
	}

	// 批量获取用户信息
	rankings := make([]models.ContestRankSummary, 0, len(results))
	for i, z := range results {
		id, _ := strconv.ParseUint(z.Member.(string), 10, 64)
		rank := models.ContestRankSummary{
			Rank:   int64(i + 1),
			UserID: id,
		}

		// 解析详情
		if i < len(detailsJSON) && detailsJSON[i] != nil {
			detailStr, ok := detailsJSON[i].(string)
			if ok {
				json.Unmarshal([]byte(detailStr), &rank.Problems)
			}
		}

		// 计算汇总
		for _, p := range rank.Problems {
			if p.Accepted {
				rank.Solved++
			}
			rank.Score += p.Score
			rank.TotalPenalty += p.Penalty
		}

		rankings = append(rankings, rank)
	}

	// 填充用户信息
	fillUserInfo(rankings)

	return rankings, nil
}

// fillUserInfo 填充用户信息
func fillUserInfo(rankings []models.ContestRankSummary) {
	for i := range rankings {
		user, err := models.User{}.QueryUserById(rankings[i].UserID)
		if err == nil {
			rankings[i].Username = user.Username
			rankings[i].Avatar = user.Avatar
		}
	}
}

// ComputeACMRankingFromDB 从数据库重建 ACM 排名
func ComputeACMRankingFromDB(contestID string, contest models.Contest) ([]models.ContestRankSummary, error) {
	items, err := models.ContestRankItem{}.GetContestRankItems(contestID)
	if err != nil {
		return nil, err
	}

	// 按用户分组
	userMap := make(map[uint64][]models.ContestRankItem)
	for _, item := range items {
		userMap[item.UserID] = append(userMap[item.UserID], item)
	}

	rankings := make([]models.ContestRankSummary, 0, len(userMap))
	for userID, userItems := range userMap {
		summary := models.ContestRankSummary{
			UserID:   userID,
			Problems: make(map[string]models.ContestProblemRankDetail),
		}

		for _, item := range userItems {
			if item.IsAccepted {
				summary.Solved++
				summary.TotalPenalty += item.TotalPenalty
			}
			d := models.ContestProblemRankDetail{
				Attempts: item.Attempts,
				Accepted: item.IsAccepted,
				Score:    item.Score,
				Penalty:  item.TotalPenalty,
			}
			if item.AcceptedAt != nil {
				elapsed := item.AcceptedAt.Sub(contest.BeginAt)
				d.Time = formatDuration(elapsed)
			}
			summary.Problems[item.ProblemLabel] = d
		}

		// 更新 Redis
		_ = UpdateACMRanking(contestID, userID, summary.Solved, summary.TotalPenalty, summary.Problems)
		rankings = append(rankings, summary)
	}

	// ACM 排序: solved DESC, penalty ASC
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Solved != rankings[j].Solved {
			return rankings[i].Solved > rankings[j].Solved
		}
		return rankings[i].TotalPenalty < rankings[j].TotalPenalty
	})

	// 分配排名
	for i := range rankings {
		rankings[i].Rank = int64(i + 1)
	}

	fillUserInfo(rankings)
	return rankings, nil
}

// ComputeOIRankingFromDB 从数据库重建 OI 排名
func ComputeOIRankingFromDB(contestID string, contest models.Contest) ([]models.ContestRankSummary, error) {
	items, err := models.ContestRankItem{}.GetContestRankItems(contestID)
	if err != nil {
		return nil, err
	}

	userMap := make(map[uint64][]models.ContestRankItem)
	for _, item := range items {
		userMap[item.UserID] = append(userMap[item.UserID], item)
	}

	rankings := make([]models.ContestRankSummary, 0, len(userMap))
	for userID, userItems := range userMap {
		summary := models.ContestRankSummary{
			UserID:   userID,
			Problems: make(map[string]models.ContestProblemRankDetail),
		}

		for _, item := range userItems {
			summary.Score += item.Score
			if item.IsAccepted {
				summary.Solved++
			}
			d := models.ContestProblemRankDetail{
				Attempts: item.Attempts,
				Accepted: item.IsAccepted,
				Score:    item.Score,
			}
			if item.AcceptedAt != nil {
				elapsed := item.AcceptedAt.Sub(contest.BeginAt)
				d.Time = formatDuration(elapsed)
			}
			summary.Problems[item.ProblemLabel] = d
		}

		_ = UpdateOIRanking(contestID, userID, summary.Score, summary.Problems)
		rankings = append(rankings, summary)
	}

	// OI 排序: score DESC
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})

	for i := range rankings {
		rankings[i].Rank = int64(i + 1)
	}

	fillUserInfo(rankings)
	return rankings, nil
}

// ProcessContestSubmissionResult 处理比赛提交的判题结果
func ProcessContestSubmissionResult(contestID string, userID uint64, problemLabel string, contest models.Contest, verdict models.JudgeVerdict, maxTime float32) {
	now := time.Now()

	// 获取已有的排名项
	var existing models.ContestRankItem
	err := dao.MysqlClient.Where("contest_id = ? AND user_id = ? AND problem_label = ?",
		contestID, userID, problemLabel).First(&existing).Error

	if contest.ContestType == models.ContestTypeACM {
		processACMSubmission(contestID, userID, problemLabel, contest, verdict, now, existing, err)
	} else {
		processOISubmission(contestID, userID, problemLabel, contest, verdict, maxTime, now, existing, err)
	}
}

func processACMSubmission(contestID string, userID uint64, problemLabel string, contest models.Contest, verdict models.JudgeVerdict, now time.Time, existing models.ContestRankItem, queryErr error) {
	isAC := verdict == models.Accepted

	if queryErr != nil {
		// 新记录
		var acceptedAt *time.Time
		var penalty int32
		if isAC {
			acceptedAt = &now
			elapsed := now.Sub(contest.BeginAt)
			penalty = int32(elapsed.Seconds())
		}
		item := &models.ContestRankItem{
			ContestID:    contestID,
			UserID:       userID,
			ProblemLabel: problemLabel,
			IsAccepted:   isAC,
			Attempts:     1,
			AcceptedAt:   acceptedAt,
			TotalPenalty: penalty,
		}
		if isAC {
			item.Score = 1
		}
		models.ContestRankItem{}.UpsertRankItem(item)
	} else if !existing.IsAccepted {
		// 更新记录
		attempts := existing.Attempts + 1
		var acceptedAt *time.Time
		var penalty int32
		if isAC {
			acceptedAt = &now
			elapsed := now.Sub(contest.BeginAt)
			// ACM 罚时: 用时 + 20分钟 * 错误次数
			penalty = int32(elapsed.Seconds()) + int32(attempts-1)*20*60
		}
		item := &models.ContestRankItem{
			ContestID:    contestID,
			UserID:       userID,
			ProblemLabel: problemLabel,
			IsAccepted:   isAC,
			Attempts:     attempts,
			AcceptedAt:   acceptedAt,
			TotalPenalty: penalty,
		}
		if isAC {
			item.Score = 1
		}
		models.ContestRankItem{}.UpsertRankItem(item)
	} else {
		// 已 AC，无需更新
		return
	}

	// 更新题目状态缓存
	status := "wrong"
	var s int32
	if isAC {
		status = "accepted"
		s = 1
	}
	ttl := time.Until(contest.EndAt) + 24*time.Hour
	if ttl <= 0 {
		ttl = 48 * time.Hour
	}
	SetProblemStatusCache(contestID, userID, problemLabel, status, s, ttl)

	// 重新计算该用户排名并更新 Redis
	rebuildUserACMRank(contestID, userID, contest)
}

func processOISubmission(contestID string, userID uint64, problemLabel string, contest models.Contest, verdict models.JudgeVerdict, maxTime float32, now time.Time, existing models.ContestRankItem, queryErr error) {
	// OI 模式: 取最高分，需要根据通过率或部分分计算
	// 简化实现: AC=100分, 否则保持原分
	var score int32
	if verdict == models.Accepted {
		score = 100
	} else if queryErr != nil {
		score = 0
	} else {
		score = existing.Score // 保持已有最高分
	}

	if queryErr != nil {
		item := &models.ContestRankItem{
			ContestID:    contestID,
			UserID:       userID,
			ProblemLabel: problemLabel,
			IsAccepted:   verdict == models.Accepted,
			Attempts:     1,
			Score:        score,
		}
		if verdict == models.Accepted {
			item.AcceptedAt = &now
		}
		models.ContestRankItem{}.UpsertRankItem(item)
	} else {
		// 取最高分
		if score > existing.Score {
			item := &models.ContestRankItem{
				ContestID:    contestID,
				UserID:       userID,
				ProblemLabel: problemLabel,
				IsAccepted:   existing.IsAccepted || verdict == models.Accepted,
				Attempts:     existing.Attempts + 1,
				Score:        score,
			}
			if verdict == models.Accepted {
				item.AcceptedAt = &now
			}
			models.ContestRankItem{}.UpsertRankItem(item)
		} else {
			// 只增加尝试次数
			dao.MysqlClient.Model(&models.ContestRankItem{}).
				Where("contest_id = ? AND user_id = ? AND problem_label = ?",
					contestID, userID, problemLabel).
				UpdateColumn("attempts", existing.Attempts+1)
		}
	}

	// 更新题目状态缓存
	status := "wrong"
	if verdict == models.Accepted {
		status = "accepted"
	}
	ttl := time.Until(contest.EndAt) + 24*time.Hour
	if ttl <= 0 {
		ttl = 48 * time.Hour
	}
	SetProblemStatusCache(contestID, userID, problemLabel, status, score, ttl)

	rebuildUserOIRank(contestID, userID)
}

func rebuildUserACMRank(contestID string, userID uint64, contest models.Contest) {
	items, err := models.ContestRankItem{}.GetUserRankItems(contestID, userID)
	if err != nil {
		return
	}
	var solved int32
	var penalty int32
	detail := make(map[string]models.ContestProblemRankDetail)
	for _, item := range items {
		if item.IsAccepted {
			solved++
			penalty += item.TotalPenalty
		}
		d := models.ContestProblemRankDetail{
			Attempts: item.Attempts,
			Accepted: item.IsAccepted,
			Score:    item.Score,
			Penalty:  item.TotalPenalty,
		}
		if item.AcceptedAt != nil {
			elapsed := item.AcceptedAt.Sub(contest.BeginAt)
			d.Time = formatDuration(elapsed)
		}
		detail[item.ProblemLabel] = d
	}
	_ = UpdateACMRanking(contestID, userID, solved, penalty, detail)
}

func rebuildUserOIRank(contestID string, userID uint64) {
	items, err := models.ContestRankItem{}.GetUserRankItems(contestID, userID)
	if err != nil {
		return
	}
	var totalScore int32
	detail := make(map[string]models.ContestProblemRankDetail)
	for _, item := range items {
		totalScore += item.Score
		detail[item.ProblemLabel] = models.ContestProblemRankDetail{
			Attempts: item.Attempts,
			Accepted: item.IsAccepted,
			Score:    item.Score,
		}
	}
	_ = UpdateOIRanking(contestID, userID, totalScore, detail)
}

// GenerateContestReport 生成比赛报告(存入 MongoDB)
func GenerateContestReport(contestID string) error {
	contest, err := models.Contest{}.QueryContestById(contestID)
	if err != nil {
		return err
	}

	items, _ := models.ContestRankItem{}.GetContestRankItems(contestID)

	// 按题目统计
	problemStats := make(map[string]*ProblemStat)
	for _, item := range items {
		if _, ok := problemStats[item.ProblemLabel]; !ok {
			problemStats[item.ProblemLabel] = &ProblemStat{Label: item.ProblemLabel}
		}
		ps := problemStats[item.ProblemLabel]
		ps.Attempts++
		if item.IsAccepted {
			ps.Accepted++
		}
	}

	// 获取排名
	var rankings []models.ContestRankSummary
	if contest.ContestType == models.ContestTypeACM {
		rankings, _ = ComputeACMRankingFromDB(contestID, contest)
	} else {
		rankings, _ = ComputeOIRankingFromDB(contestID, contest)
	}

	// 构建 Top 参赛者
	topParticipants := make([]map[string]interface{}, 0, len(rankings))
	for _, r := range rankings {
		if r.Rank > 20 {
			break
		}
		topParticipants = append(topParticipants, map[string]interface{}{
			"user_id":         r.UserID,
			"username":        r.Username,
			"rank":            r.Rank,
			"score":           r.Score,
			"problems_solved": r.Solved,
			"total_penalty":   r.TotalPenalty,
		})
	}

	// 构建题目统计
	psList := make([]map[string]interface{}, 0)
	for _, ps := range problemStats {
		acceptRate := float64(0)
		if ps.Attempts > 0 {
			acceptRate = float64(ps.Accepted) / float64(ps.Attempts) * 100
		}
		psList = append(psList, map[string]interface{}{
			"label":       ps.Label,
			"attempts":    ps.Attempts,
			"accepted":    ps.Accepted,
			"accept_rate": acceptRate,
		})
	}

	report := map[string]interface{}{
		"contest_id":         contestID,
		"title":              contest.Title,
		"type":               contest.ContestType,
		"generated_at":       time.Now(),
		"total_participants": contest.Participants,
		"total_submissions":  contest.Submission,
		"problem_stats":      psList,
		"top_participants":   topParticipants,
	}

	return dao.InsertDocument("nexus", "contest_reports", report)
}

// GetContestReport 从 MongoDB 获取比赛报告
func GetContestReport(contestID string) (map[string]interface{}, error) {
	results, err := dao.QueryDocument("nexus", "contest_reports", map[string]interface{}{
		"contest_id": contestID,
	})
	if err != nil || len(results) == 0 {
		return nil, fmt.Errorf("report not found")
	}
	return results[0], nil
}

type ProblemStat struct {
	Label    string
	Attempts int
	Accepted int
}

func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

// CacheContestInfo 缓存比赛信息
func CacheContestInfo(contest models.Contest) error {
	key := fmt.Sprintf(redisContestInfo, contest.ID)
	data, _ := json.Marshal(contest)
	return dao.RedisClient.Set(ctx, key, string(data), 5*time.Minute).Err()
}

// GetCachedContestInfo 获取缓存的比赛信息
func GetCachedContestInfo(contestID string) (*models.Contest, error) {
	key := fmt.Sprintf(redisContestInfo, contestID)
	data, err := dao.RedisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var contest models.Contest
	if err := json.Unmarshal(data, &contest); err != nil {
		return nil, err
	}
	return &contest, nil
}

// InvalidateContestInfo 使缓存失效
func InvalidateContestInfo(contestID string) {
	key := fmt.Sprintf(redisContestInfo, contestID)
	dao.RedisClient.Del(ctx, key)
}

// CheckContestStatus 根据时间计算比赛状态
func CheckContestStatus(contest *models.Contest) models.ContestStatus {
	now := time.Now()
	if now.Before(contest.BeginAt) {
		return models.ContestStatusUpcoming
	}
	if now.After(contest.EndAt) {
		return models.ContestStatusEnded
	}
	return models.ContestStatusLive
}

// ComputeContestStatus 实时计算状态并覆盖 struct 字段，不写数据库
func ComputeContestStatus(contest *models.Contest) {
	contest.Status = CheckContestStatus(contest)
}

// InitContestCache 初始化比赛缓存(从数据库加载到 Redis)
func InitContestCache(contest *models.Contest) {
	// 缓存密码
	if contest.IsPrivate && contest.Password != "" {
		_ = CacheContestPassword(contest.ID, contest.Password, contest.EndAt)
	}

	// 缓存比赛信息
	_ = CacheContestInfo(*contest)

	// 加载参赛者到 Redis
	participants, err := models.ContestParticipant{}.GetParticipants(contest.ID, 1, 10000)
	if err == nil {
		for _, p := range participants {
			_ = CacheParticipant(contest.ID, p.UserID, contest.EndAt)
		}
	}

	// 重建排名缓存
	if contest.ContestType == models.ContestTypeACM {
		ComputeACMRankingFromDB(contest.ID, *contest)
	} else {
		ComputeOIRankingFromDB(contest.ID, *contest)
	}
}

// SSESubscriber SSE 订阅管理
type SSESubscriber struct {
	ContestID string
	UserID    uint64
	Ch        chan []byte
}

var sseSubscribers = make(map[string][]*SSESubscriber) // contestID -> subscribers

// SubscribeRanking 订阅排名更新
func SubscribeRanking(contestID string, userID uint64) *SSESubscriber {
	sub := &SSESubscriber{
		ContestID: contestID,
		UserID:    userID,
		Ch:        make(chan []byte, 10),
	}
	sseSubscribers[contestID] = append(sseSubscribers[contestID], sub)
	return sub
}

// UnsubscribeRanking 取消订阅
func UnsubscribeRanking(sub *SSESubscriber) {
	subs := sseSubscribers[sub.ContestID]
	for i, s := range subs {
		if s == sub {
			sseSubscribers[sub.ContestID] = append(subs[:i], subs[i+1:]...)
			close(sub.Ch)
			break
		}
	}
}

// NotifyRankingUpdate 通知所有订阅者排名更新
func NotifyRankingUpdate(contestID string, data []byte) {
	subs, ok := sseSubscribers[contestID]
	if !ok {
		return
	}
	for _, sub := range subs {
		select {
		case sub.Ch <- data:
		default:
			log.Printf("SSE subscriber channel full for contest %s user %s", contestID, sub.UserID)
		}
	}
}

// PublishRankingUpdate 发布排名更新(供 judge queue 回调)
func PublishRankingUpdate(contestID string) {
	rankings, err := GetContestRankingFromRedis(contestID, 100)
	if err != nil {
		return
	}
	data, _ := json.Marshal(map[string]interface{}{
		"type":     "full",
		"ranking":  rankings,
		"updateAt": time.Now().Unix(),
	})
	NotifyRankingUpdate(contestID, data)
}
