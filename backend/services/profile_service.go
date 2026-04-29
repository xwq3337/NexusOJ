package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nexus/dao"
	"nexus/models"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
)

// ==================== 数据结构定义 ====================

// ProfileUpdateEvent 画像更新事件（每次提交触发）
type ProfileUpdateEvent struct {
	UserID     uint64
	ProblemID  int64
	Verdict    models.JudgeVerdict
	Language   string
	Difficulty float32
	Tags       []string
	Timestamp  time.Time
}

// UserProfile 完整用户画像
type UserProfile struct {
	Ability     *AbilityDimension    `json:"ability"`
	Activity    *ActivityDimension   `json:"activity"`
	Preferences *PreferenceDimension `json:"preferences"`
	Social      *SocialDimension     `json:"social"`
}

// AbilityDimension 能力维度
type AbilityDimension struct {
	OverallScore  float32            `json:"overall_score"`
	TagScores     map[string]float32 `json:"tag_scores"`
	TagProgress   map[string]float32 `json:"tag_progress"`
	TagTotal      map[string]int     `json:"tag_total"`
	StrongestTags []string           `json:"strongest_tags"`
	WeakestTags   []string           `json:"weakest_tags"`
}

// ActivityDimension 活跃度维度
type ActivityDimension struct {
	Streak     int                       `json:"streak"`
	LastActive time.Time                 `json:"last_active"`
	Heatmaps   map[string]map[string]int `json:"heatmaps"` // year -> {"MM-DD": count}
}

// PreferenceDimension 偏好维度
type PreferenceDimension struct {
	Languages       map[string]int `json:"languages"`
	AvgDifficulty   float32        `json:"avg_difficulty"`
	TotalSubmission int            `json:"total_submission"`
}

// SocialDimension 社交维度
type SocialDimension struct {
	FriendCount  int `json:"friend_count"`
	ContestCount int `json:"contest_count"`
	BlogCount    int `json:"blog_count"`
}

// TagStatEntry Redis tag_stats Hash 中的值
type TagStatEntry struct {
	Accepted  int     `json:"accepted"`
	Attempted int     `json:"attempted"`
	AvgDiff   float32 `json:"avg_difficulty"`
}

// ==================== Redis Key 定义 ====================

const (
	profileTagStatsKey    = "profile:%d:tag_stats"
	profileActivityKey    = "profile:%d:activity"
	profileHeatmapKey     = "profile:%d:heatmap:%s" // userID, year
	profilePreferencesKey = "profile:%d:preferences"
	profileAbilityKey     = "profile:%d:ability"
	profileSocialKey      = "profile:%d:social"
	profileDirtyKey       = "profile:%d:dirty"
)

func profileTagStatsKeyF(userID uint64) string                { return fmt.Sprintf(profileTagStatsKey, userID) }
func profileActivityKeyF(userID uint64) string                { return fmt.Sprintf(profileActivityKey, userID) }
func profileHeatmapKeyF(userID uint64, year string) string    { return fmt.Sprintf(profileHeatmapKey, userID, year) }
func profilePreferencesKeyF(userID uint64) string             { return fmt.Sprintf(profilePreferencesKey, userID) }
func profileAbilityKeyF(userID uint64) string                 { return fmt.Sprintf(profileAbilityKey, userID) }
func profileSocialKeyF(userID uint64) string                  { return fmt.Sprintf(profileSocialKey, userID) }
func profileDirtyKeyF(userID uint64) string                   { return fmt.Sprintf(profileDirtyKey, userID) }

// ==================== Profile Service ====================

var GlobalProfileService *ProfileService

// ProfileService 用户画像服务
type ProfileService struct {
	updateChan chan *ProfileUpdateEvent
	workerNum  int
	wg         sync.WaitGroup
}

// InitProfileService 初始化画像服务
func InitProfileService(workerNum, queueSize int) {
	GlobalProfileService = &ProfileService{
		updateChan: make(chan *ProfileUpdateEvent, queueSize),
		workerNum:  workerNum,
	}
	for i := 0; i < workerNum; i++ {
		GlobalProfileService.wg.Add(1)
		go GlobalProfileService.worker(i)
	}
	log.Printf("Profile service initialized with %d workers and queue size %d", workerNum, queueSize)
}

// SubmitProfileUpdate 非阻塞提交画像更新事件
func SubmitProfileUpdate(event *ProfileUpdateEvent) {
	if GlobalProfileService == nil {
		return
	}
	select {
	case GlobalProfileService.updateChan <- event:
	default:
		log.Printf("Profile update queue full, dropping event for user %d", event.UserID)
	}
}

func (s *ProfileService) worker(id int) {
	defer s.wg.Done()
	log.Printf("Profile worker %d started", id)
	for event := range s.updateChan {
		s.processEvent(event, id)
	}
	log.Printf("Profile worker %d stopped", id)
}

func (s *ProfileService) processEvent(event *ProfileUpdateEvent, workerID int) {
	ctx := context.Background()
	pipe := dao.RedisClient.Pipeline()

	// 1. 更新 tag_stats（Hash）
	isAccepted := event.Verdict == models.Accepted
	tagStatsKey := profileTagStatsKeyF(event.UserID)
	for _, tag := range event.Tags {
		// 读取当前值
		val, _ := dao.RedisClient.HGet(ctx, tagStatsKey, tag).Result()
		var entry TagStatEntry
		if val != "" {
			_ = json.Unmarshal([]byte(val), &entry)
		}
		entry.Attempted++
		if isAccepted {
			entry.Accepted++
			// 更新平均难度（移动平均）
			entry.AvgDiff = (entry.AvgDiff*float32(entry.Accepted-1) + event.Difficulty) / float32(entry.Accepted)
		}
		data, _ := json.Marshal(entry)
		pipe.HSet(ctx, tagStatsKey, tag, data)
	}

	// 2. 更新 heatmap（HINCRBY）
	year := event.Timestamp.Format("2006")
	dateKey := event.Timestamp.Format("01-02")
	heatmapKey := profileHeatmapKeyF(event.UserID, year)
	pipe.HIncrBy(ctx, heatmapKey, dateKey, 1)

	// 3. 更新 activity（streak、last_active）
	activityKey := profileActivityKeyF(event.UserID)
	activityData, _ := dao.RedisClient.Get(ctx, activityKey).Bytes()
	var activity struct {
		Streak     int       `json:"streak"`
		LastActive time.Time `json:"last_active"`
	}
	if activityData != nil {
		_ = json.Unmarshal(activityData, &activity)
	}
	today := time.Now().Truncate(24 * time.Hour)
	lastActive := activity.LastActive.Truncate(24 * time.Hour)
	if today.Sub(lastActive) == 24*time.Hour {
		activity.Streak++
	} else if today.Sub(lastActive) > 24*time.Hour {
		activity.Streak = 1
	}
	activity.LastActive = event.Timestamp
	activityBytes, _ := json.Marshal(activity)
	pipe.Set(ctx, activityKey, activityBytes, 0)

	// 4. 更新 preferences（语言计数、平均难度）
	prefKey := profilePreferencesKeyF(event.UserID)
	prefData, _ := dao.RedisClient.Get(ctx, prefKey).Bytes()
	var pref PreferenceDimension
	if prefData != nil {
		_ = json.Unmarshal(prefData, &pref)
	}
	if pref.Languages == nil {
		pref.Languages = make(map[string]int)
	}
	pref.Languages[event.Language]++
	pref.TotalSubmission++
	// 移动平均难度
	pref.AvgDifficulty = (pref.AvgDifficulty*float32(pref.TotalSubmission-1) + event.Difficulty) / float32(pref.TotalSubmission)
	prefBytes, _ := json.Marshal(pref)
	pipe.Set(ctx, prefKey, prefBytes, 0)

	// 5. 标记推荐缓存为脏
	pipe.Set(ctx, profileDirtyKeyF(event.UserID), "1", 0)

	// 执行 pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		log.Printf("Profile worker %d: pipeline error for user %d: %v", workerID, event.UserID, err)
	}
}

// ==================== 标签进度计算 ====================

// computeTagProgress 计算所有标签的真实进度（已解决/总题数）
// 扫描 Redis tag_problems:* 获取全部标签，未涉猎的标签进度为 0
func computeTagProgress(userID uint64, ability *AbilityDimension) {
	ctx := context.Background()
	solvedKey := fmt.Sprintf(dao.UserSolvedBitKey, userID)

	// 扫描所有 tag_problems:* 键，获取全部标签
	var cursor uint64
	for {
		keys, nextCursor, err := dao.RedisClient.Scan(ctx, cursor, "tag_problems:*", 100).Result()
		if err != nil {
			return
		}
		for _, key := range keys {
			tag := strings.TrimPrefix(key, "tag_problems:")
			members, err := dao.RedisClient.ZRange(ctx, key, 0, -1).Result()
			if err != nil || len(members) == 0 {
				continue
			}

			total := len(members)
			pipe := dao.RedisClient.Pipeline()
			cmds := make([]*redis.IntCmd, total)
			for i, m := range members {
				id, _ := strconv.ParseInt(m, 10, 64)
				cmds[i] = pipe.GetBit(ctx, solvedKey, id-1000)
			}
			pipe.Exec(ctx)

			var solved int
			for _, cmd := range cmds {
				if cmd.Val() == 1 {
					solved++
				}
			}
			ability.TagProgress[tag] = float32(solved) / float32(total)
			ability.TagTotal[tag] = total
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

// ==================== 画像读取方法 ====================

// GetUserProfile 从 Redis 读取完整画像
func GetUserProfile(userID uint64) (*UserProfile, error) {
	ctx := context.Background()
	profile := &UserProfile{
		Ability:     &AbilityDimension{TagScores: make(map[string]float32), TagProgress: make(map[string]float32), TagTotal: make(map[string]int)},
		Activity:    &ActivityDimension{Heatmaps: make(map[string]map[string]int)},
		Preferences: &PreferenceDimension{Languages: make(map[string]int)},
		Social:      &SocialDimension{},
	}

	// 使用 Goroutine 并发读取各维度数据
	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	// 读取 tag_stats -> 计算能力维度
	wg.Add(1)
	go func() {
		defer wg.Done()
		tagStatsKey := profileTagStatsKeyF(userID)
		result, err := dao.RedisClient.HGetAll(ctx, tagStatsKey).Result()
		if err != nil {
			errChan <- err
			return
		}
		var totalAttempted int
		var weightedScoreSum float32
		tagList := make([]string, 0, len(result))
		for tag, val := range result {
			var entry TagStatEntry
			if err := json.Unmarshal([]byte(val), &entry); err != nil {
				continue
			}
			score := float32(0)
			if entry.Attempted > 0 {
				score = float32(entry.Accepted) / float32(entry.Attempted)
			}
			profile.Ability.TagScores[tag] = score
			totalAttempted += entry.Attempted
			weightedScoreSum += score * float32(entry.Attempted)
			tagList = append(tagList, tag)
		}
		if totalAttempted > 0 {
			profile.Ability.OverallScore = weightedScoreSum / float32(totalAttempted)
		}

		// 计算所有标签真实进度（含未涉猎标签，进度为 0）
		computeTagProgress(userID, profile.Ability)

		// 基于进度排序找出最强/最弱标签
		allTags := make([]string, 0, len(profile.Ability.TagProgress))
		for tag := range profile.Ability.TagProgress {
			allTags = append(allTags, tag)
		}
		progressTags := make([]struct {
			tag      string
			progress float32
		}, 0, len(allTags))
		for _, tag := range allTags {
			progressTags = append(progressTags, struct {
				tag      string
				progress float32
			}{tag, profile.Ability.TagProgress[tag]})
		}
		sortTagsByProgress(progressTags)
		n := len(progressTags)
		if n > 3 {
			profile.Ability.StrongestTags = tagNameFromProgress(progressTags[n-3:])
			profile.Ability.WeakestTags = tagNameFromProgress(progressTags[:3])
		} else if n > 0 {
			profile.Ability.StrongestTags = tagNameFromProgress(progressTags)
		}
	}()

	// 读取 activity + heatmap
	wg.Add(1)
	go func() {
		defer wg.Done()
		activityData, _ := dao.RedisClient.Get(ctx, profileActivityKeyF(userID)).Bytes()
		if activityData != nil {
			_ = json.Unmarshal(activityData, profile.Activity)
		}
		// 读取今年和去年的 heatmap
		thisYear := time.Now().Format("2006")
		lastYear := fmt.Sprintf("%d", time.Now().Year()-1)
		yearBefore := fmt.Sprintf("%d", time.Now().Year()-2)
		for _, year := range []string{thisYear, lastYear, yearBefore} {
			result, err := dao.RedisClient.HGetAll(ctx, profileHeatmapKeyF(userID, year)).Result()
			if err == nil && len(result) > 0 {
				heatmap := make(map[string]int, len(result))
				for k, v := range result {
					var count int
					fmt.Sscanf(v, "%d", &count)
					heatmap[k] = count
				}
				profile.Activity.Heatmaps[year] = heatmap
			}
		}
	}()

	// 读取 preferences
	wg.Add(1)
	go func() {
		defer wg.Done()
		prefData, _ := dao.RedisClient.Get(ctx, profilePreferencesKeyF(userID)).Bytes()
		if prefData != nil {
			_ = json.Unmarshal(prefData, profile.Preferences)
		}
	}()

	// 读取 social
	wg.Add(1)
	go func() {
		defer wg.Done()
		socialData, _ := dao.RedisClient.Get(ctx, profileSocialKeyF(userID)).Bytes()
		if socialData != nil {
			_ = json.Unmarshal(socialData, profile.Social)
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return profile, err
		}
	}
	return profile, nil
}

// sortTagsByProgress 按进度升序排列（最弱在前）
func sortTagsByProgress(tags []struct {
	tag      string
	progress float32
}) {
	for i := 0; i < len(tags); i++ {
		for j := i + 1; j < len(tags); j++ {
			if tags[i].progress > tags[j].progress {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}
	}
}

func tagNameFromProgress(tags []struct {
	tag      string
	progress float32
}) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.tag
	}
	return names
}

// sortTagsByScore 按分数升序排列（最弱在前）
func sortTagsByScore(tags []struct {
	tag   string
	score float32
}) {
	for i := 0; i < len(tags); i++ {
		for j := i + 1; j < len(tags); j++ {
			if tags[i].score > tags[j].score {
				tags[i], tags[j] = tags[j], tags[i]
			}
		}
	}
}

func tagNames(tags []struct {
	tag   string
	score float32
}) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.tag
	}
	return names
}

// ==================== 热力图读取 ====================

// GetHeatmapByYear 获取指定年份的热力图
func GetHeatmapByYear(userID uint64, year string) (map[string]int, error) {
	ctx := context.Background()
	key := profileHeatmapKeyF(userID, year)
	result, err := dao.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	heatmap := make(map[string]int, len(result))
	for k, v := range result {
		var count int
		fmt.Sscanf(v, "%d", &count)
		heatmap[k] = count
	}
	return heatmap, nil
}

// GetPastYearHeatmap 获取近一年（滚动365天）的热力图
func GetPastYearHeatmap(userID uint64) (map[string]int, error) {
	now := time.Now()
	thisYear := now.Format("2006")
	lastYear := fmt.Sprintf("%d", now.Year()-1)
	todayMMDD := now.Format("01-02")

	ctx := context.Background()
	pipe := dao.RedisClient.Pipeline()
	cmd1 := pipe.HGetAll(ctx, profileHeatmapKeyF(userID, thisYear))
	cmd2 := pipe.HGetAll(ctx, profileHeatmapKeyF(userID, lastYear))
	pipe.Exec(ctx)

	result := make(map[string]int)

	// 今年的数据全部保留
	for k, v := range cmd1.Val() {
		var count int
		fmt.Sscanf(v, "%d", &count)
		result[thisYear+"-"+k] = count
	}
	// 去年的数据只保留今天之后的部分
	for k, v := range cmd2.Val() {
		if k > todayMMDD {
			var count int
			fmt.Sscanf(v, "%d", &count)
			result[lastYear+"-"+k] = count
		}
	}
	return result, nil
}

// ==================== 冷启动：全量计算画像 ====================

// FullProfileRecompute 从 MySQL 全量计算用户画像（冷启动或数据丢失时使用）
func FullProfileRecompute(userID uint64) error {
	log.Printf("Starting full profile recompute for user %d", userID)
	ctx := context.Background()

	// 1. 聚合 tag 统计
	tagStats, err := models.GetUserTagAttemptStats(userID)
	if err != nil {
		return fmt.Errorf("failed to get tag stats: %w", err)
	}
	tagStatsKey := profileTagStatsKeyF(userID)
	pipe := dao.RedisClient.Pipeline()
	for _, ts := range tagStats {
		entry := TagStatEntry{
			Accepted:  ts.Accepted,
			Attempted: ts.Attempted,
			AvgDiff:   ts.AvgDiff,
		}
		data, _ := json.Marshal(entry)
		pipe.HSet(ctx, tagStatsKey, ts.Tag, data)
	}

	// 2. 聚合热力图（近3年）
	now := time.Now()
	for i := 0; i < 3; i++ {
		year := now.Year() - i
		startDate := fmt.Sprintf("%d-01-01", year)
		endDate := fmt.Sprintf("%d-12-31", year)
		dateMap, err := models.GetUserActivityByDate(userID, startDate, endDate)
		if err != nil {
			continue
		}
		heatmapKey := profileHeatmapKeyF(userID, fmt.Sprintf("%d", year))
		for date, count := range dateMap {
			// date 格式 "2026-04-15" -> 提取 "04-15"
			parts := splitDate(date)
			if len(parts) == 3 {
				mmdd := parts[1] + "-" + parts[2]
				pipe.HSet(ctx, heatmapKey, mmdd, count)
			}
		}
	}

	// 3. 聚合语言统计
	langStats, err := models.GetUserLanguageStats(userID)
	if err != nil {
		return fmt.Errorf("failed to get language stats: %w", err)
	}
	pref := PreferenceDimension{
		Languages: langStats,
	}
	// 计算平均难度
	diffDist, err := models.GetUserDifficultyDistribution(userID)
	if err == nil {
		var totalAccepted, totalAttempted int
		for _, arr := range diffDist {
			totalAccepted += arr[0]
			totalAttempted += arr[1]
		}
		pref.TotalSubmission = totalAttempted
	}
	prefBytes, _ := json.Marshal(pref)
	pipe.Set(ctx, profilePreferencesKeyF(userID), prefBytes, 0)

	// 4. 计算初始 streak
	activity := struct {
		Streak     int       `json:"streak"`
		LastActive time.Time `json:"last_active"`
	}{}
	// 从最近提交记录推算 streak
	today := time.Now().Truncate(24 * time.Hour)
	for i := 0; i < 365; i++ {
		date := today.AddDate(0, 0, -i)
		heatmapKey := profileHeatmapKeyF(userID, date.Format("2006"))
		mmdd := date.Format("01-02")
		val, _ := dao.RedisClient.HGet(ctx, heatmapKey, mmdd).Result()
		if val != "" && val != "0" {
			if i == 0 || activity.Streak > 0 {
				activity.Streak++
				activity.LastActive = date
			}
		} else if activity.Streak > 0 {
			break
		}
	}
	activityBytes, _ := json.Marshal(activity)
	pipe.Set(ctx, profileActivityKeyF(userID), activityBytes, 0)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to save profile to Redis: %w", err)
	}

	log.Printf("Full profile recompute completed for user %d", userID)
	return nil
}

func splitDate(date string) []string {
	var parts []string
	current := ""
	for _, c := range date {
		if c == '-' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// ==================== MongoDB 持久化 ====================

const (
	profileDB         = "nexus"
	profileCollection = "user_profiles"
)

// PersistProfileToMongoDB 将画像持久化到 MongoDB
func PersistProfileToMongoDB(userID uint64) error {
	profile, err := GetUserProfile(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": userID}
	update := bson.M{
		"user_id":     userID,
		"updated_at":  time.Now(),
		"ability":     profile.Ability,
		"activity":    profile.Activity,
		"preferences": profile.Preferences,
		"social":      profile.Social,
	}

	return dao.UpdateDocument(profileDB, profileCollection, filter, update)
}

// LoadProfileFromMongoDB 从 MongoDB 加载画像（冷启动用）
func LoadProfileFromMongoDB(userID uint64) (*UserProfile, error) {
	filter := bson.M{"user_id": userID}
	results, err := dao.QueryDocument(profileDB, profileCollection, filter)
	if err != nil || len(results) == 0 {
		return nil, fmt.Errorf("profile not found in MongoDB for user %d", userID)
	}

	data, err := json.Marshal(results[0])
	if err != nil {
		return nil, err
	}
	var profile UserProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// ==================== 辅助方法 ====================

// EnsureProfile 确保用户画像存在于 Redis 中，不存在则触发冷启动
func EnsureProfile(userID uint64) (*UserProfile, error) {
	// 先检查 Redis 是否有画像数据
	ctx := context.Background()
	exists := dao.RedisClient.Exists(ctx, profileTagStatsKeyF(userID)).Val()
	if exists > 0 {
		return GetUserProfile(userID)
	}

	// 尝试从 MongoDB 加载
	profile, err := LoadProfileFromMongoDB(userID)
	if err == nil {
		// 写回 Redis
		restoreProfileToRedis(userID, profile)
		return profile, nil
	}

	// 从 MySQL 全量计算
	if err := FullProfileRecompute(userID); err != nil {
		return nil, err
	}
	return GetUserProfile(userID)
}

// restoreProfileToRedis 将 MongoDB 中的画像恢复到 Redis
func restoreProfileToRedis(userID uint64, profile *UserProfile) {
	ctx := context.Background()
	pipe := dao.RedisClient.Pipeline()

	// 恢复 tag_stats
	if profile.Ability != nil {
		tagStatsKey := profileTagStatsKeyF(userID)
		for tag, score := range profile.Ability.TagScores {
			entry := TagStatEntry{Accepted: int(score * 100), Attempted: 100}
			data, _ := json.Marshal(entry)
			pipe.HSet(ctx, tagStatsKey, tag, data)
		}
	}

	// 恢复 heatmap
	if profile.Activity != nil && profile.Activity.Heatmaps != nil {
		for year, heatmap := range profile.Activity.Heatmaps {
			key := profileHeatmapKeyF(userID, year)
			for mmdd, count := range heatmap {
				pipe.HSet(ctx, key, mmdd, count)
			}
		}
	}

	// 恢复 preferences
	if profile.Preferences != nil {
		data, _ := json.Marshal(profile.Preferences)
		pipe.Set(ctx, profilePreferencesKeyF(userID), data, 0)
	}

	// 恢复 activity
	if profile.Activity != nil {
		data, _ := json.Marshal(profile.Activity)
		pipe.Set(ctx, profileActivityKeyF(userID), data, 0)
	}

	pipe.Exec(ctx)
}
