package models

import (
	"encoding/json"
	"nexus/dao"
	"nexus/utils"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Result struct {
	Error   int     `json:"error"`
	Input   string  `json:"input"`
	Output  string  `json:"output"`
	Time    float64 `json:"time"`
	Memory  float64 `json:"memory"`
	Display string  `json:"display"`
}
type Record struct {
	ID          int64                                    `json:"id" gorm:"primarykey"`
	UserId      uint64                                   `json:"user_id" gorm:"index:idx_user_problem"`
	ProblemId   int64                                    `json:"problem_id" gorm:"index:idx_user_problem"`
	Code        string                                   `json:"code" gorm:"type:longtext"`
	Language    string                                   `json:"language" gorm:"index:idx_language_verdict"`
	Verdict     JudgeVerdict                             `json:"verdict" gorm:"index:idx_language_verdict"`
	MaxTime     float32                                  `json:"max_time"`
	MaxMemory   float32                                  `json:"max_memory"`
	JudgeResult datatypes.JSONSlice[JudgeTestCaseResult] `json:"judge_result"`
	CreatedAt   time.Time                                `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	UpdatedAt   time.Time                                `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
	DeletedAt   gorm.DeletedAt                           `json:"deleted_at" gorm:"index"`
}

func (Record) TableName() string {
	return "record"
}
func CreateRecord(record *Record) error {
	err := dao.MysqlClient.Create(record).Error
	return err
}

func QueryRecord(record Record) (Record, error) {
	err := dao.MysqlClient.Find(&record, record).Error
	return record, err
}

/**
 * 根据用户id获取记录
 */
// UserRecordItem 用户提交记录列表项
type UserRecordItem struct {
	ID           int64     `json:"id"`
	ProblemID    string    `json:"problem_id"`
	Language     string    `json:"language"`
	Verdict      string    `json:"verdict"`
	MaxTime      float32   `json:"max_time"`
	MaxMemory    float32   `json:"max_memory"`
	CreatedAt    time.Time `json:"created_at"`
	ProblemTitle string    `json:"problem_title"`
	Total        int64     `json:"-" gorm:"column:total"`
}

func QueryRecordByUserId(userID string, page int, pageSize int, verdict string, language string) ([]UserRecordItem, int64, error) {
	var records []UserRecordItem

	offset := (page - 1) * pageSize

	query := dao.MysqlClient.Table("record").
		Select(`record.id AS id,
			record.problem_id AS problem_id,
			record.language AS language,
			record.verdict AS verdict,
			record.max_time AS max_time,
			record.max_memory AS max_memory,
			record.created_at AS created_at,
			problem.title AS problem_title,
			COUNT(*) OVER() AS total`).
		Joins(`INNER JOIN problem ON problem.id = record.problem_id`).
		Where("record.user_id = ?", userID)

	if verdict != "" {
		query = query.Where("record.verdict = ?", verdict)
	}
	if language != "" {
		query = query.Where("record.language = ?", language)
	}

	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Scan(&records).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(records) > 0 {
		total = records[0].Total
	}

	return records, total, nil
}

/**
 * 根据记录id获取记录
 */
type RecordDetail struct {
	Record
	ProblemTitle string `json:"problem_title" gorm:"column:problem_title"`
	Username     string `json:"username" gorm:"column:username"`
}

func QueryRecordById(id string) (*RecordDetail, error) {
	var record RecordDetail
	err := dao.MysqlClient.Table("record").
		Select("record.*", "problem.title AS problem_title", "user.username AS username").
		Joins("INNER JOIN user ON user.id = record.user_id").
		Joins("INNER JOIN problem ON problem.id = record.problem_id").
		Where("record.id = ?", id).Limit(1).Scan(&record).Error
	return &record, err
}

/**
 * 获取所有记录（支持分页和查询），使用窗口函数一次查询返回数据和总数
 */
// TagStat 标签级别的统计
type TagStat struct {
	Tag        string  `json:"tag"`
	Accepted   int     `json:"accepted"`
	Attempted  int     `json:"attempted"`
	AvgDiff    float32 `json:"avg_difficulty"` // 已通过题目的平均难度
}

// GetUserTagAttemptStats 聚合用户各标签的刷题统计（JOIN problem 获取 tags）
func GetUserTagAttemptStats(userID uint64) ([]TagStat, error) {
	var results []struct {
		Tags    string `json:"tags"`
		Verdict string `json:"verdict"`
	}
	err := dao.MysqlClient.Table("record").
		Select("problem.tags AS tags, record.verdict AS verdict").
		Joins("JOIN problem ON problem.id = record.problem_id").
		Where("record.user_id = ? AND record.deleted_at IS NULL", userID).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 在内存中按 tag 聚合
	tagMap := make(map[string]*TagStat)
	for _, r := range results {
		var tags []string
		_ = json.Unmarshal([]byte(r.Tags), &tags)
		for _, tag := range tags {
			if _, ok := tagMap[tag]; !ok {
				tagMap[tag] = &TagStat{Tag: tag}
			}
			tagMap[tag].Attempted++
			if r.Verdict == "Accepted" {
				tagMap[tag].Accepted++
			}
		}
	}

	stats := make([]TagStat, 0, len(tagMap))
	for _, v := range tagMap {
		stats = append(stats, *v)
	}
	return stats, nil
}

// GetUserDifficultyDistribution 获取用户按难度区间的通过/尝试统计
func GetUserDifficultyDistribution(userID uint64) (map[string][2]int, error) {
	var results []struct {
		Difficulty float32 `json:"difficulty"`
		Verdict    string  `json:"verdict"`
	}
	err := dao.MysqlClient.Table("record").
		Select("problem.difficulty AS difficulty, record.verdict AS verdict").
		Joins("JOIN problem ON problem.id = record.problem_id").
		Where("record.user_id = ? AND record.deleted_at IS NULL", userID).
		Find(&results).Error
	if err != nil {
		return nil, err
	}

	dist := map[string][2]int{ // [accepted, attempted]
		"easy":   {},
		"medium": {},
		"hard":   {},
	}
	for _, r := range results {
		var key string
		switch {
		case r.Difficulty < 1.0:
			key = "easy"
		case r.Difficulty < 2.0:
			key = "medium"
		default:
			key = "hard"
		}
		arr := dist[key]
		arr[1]++ // attempted
		if r.Verdict == "Accepted" {
			arr[0]++ // accepted
		}
		dist[key] = arr
	}
	return dist, nil
}

// GetUserLanguageStats 获取用户语言使用统计
func GetUserLanguageStats(userID uint64) (map[string]int, error) {
	var results []struct {
		Language string `json:"language"`
		Count    int    `json:"count"`
	}
	err := dao.MysqlClient.Table("record").
		Select("language, COUNT(*) AS count").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("language").
		Order("count DESC").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	langMap := make(map[string]int, len(results))
	for _, r := range results {
		langMap[r.Language] = r.Count
	}
	return langMap, nil
}

// GetUserActivityByDate 获取用户指定日期范围内的每日提交数
func GetUserActivityByDate(userID uint64, startDate, endDate string) (map[string]int, error) {
	var results []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
	err := dao.MysqlClient.Table("record").
		Select("DATE(created_at) AS date, COUNT(*) AS count").
		Where("user_id = ? AND created_at >= ? AND created_at <= ? AND deleted_at IS NULL", userID, startDate, endDate).
		Group("DATE(created_at)").
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	dateMap := make(map[string]int, len(results))
	for _, r := range results {
		dateMap[r.Date] = r.Count
	}
	return dateMap, nil
}

// GetRecentUserRecords 获取用户最近 N 条提交记录（用于上下文感知推荐）
func GetRecentUserRecords(userID uint64, limit int) ([]Record, error) {
	var records []Record
	err := dao.MysqlClient.Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&records).Error
	return records, err
}

// GetUserSolvedProblemIDs 获取用户已解决的所有题目ID
func GetUserSolvedProblemIDs(userID uint64) ([]int64, error) {
	var ids []int64
	err := dao.MysqlClient.Table("record").
		Select("DISTINCT problem_id").
		Where("user_id = ? AND verdict = 'Accepted' AND deleted_at IS NULL", userID).
		Pluck("problem_id", &ids).Error
	return ids, err
}

func (Record) GetAllRecord(page, pageSize int, search, verdict, language, problemID string) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}

	query := dao.MysqlClient.Table("record").
		Select(`record.id AS id,
				record.user_id AS user_id,
				record.problem_id AS problem_id,
				record.language AS language,
				record.verdict AS verdict,
				record.max_time AS max_time,
				record.max_memory AS max_memory,
				record.created_at AS created_at,
				problem.title AS problem_title,
				user.username AS username,
				COUNT(*) OVER() AS total
		`).
		Joins(`JOIN user ON user.id = record.user_id`).
		Joins(`JOIN problem ON problem.id = record.problem_id`)

	if search != "" {
		query = query.Where("MATCH(problem.title) AGAINST(? IN BOOLEAN MODE) OR MATCH(user.username) AGAINST(? IN BOOLEAN MODE)", utils.SanitizeFTSSearch(search), utils.SanitizeFTSSearch(search))
	}
	if verdict != "" {
		query = query.Where("record.verdict = ?", verdict)
	}
	if language != "" {
		query = query.Where("record.language = ?", language)
	}
	if problemID != "" {
		query = query.Where("record.problem_id = ?", problemID)
	}

	offset := (page - 1) * pageSize

	err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(results) > 0 {
		if t, ok := results[0]["total"]; ok {
			total = t.(int64)
		}
	}
	// TODO:无需删除total字段
	for i := range results {
		delete(results[i], "total")
	}

	return results, total, nil
}
