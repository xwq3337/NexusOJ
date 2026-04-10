package models

import (
	"nexus/dao"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ContestRecord 比赛提交记录（独立于题库 record）
type ContestRecord struct {
	ID           int64                                    `json:"id" gorm:"primarykey"`
	ContestID    string                                   `json:"contest_id" gorm:"index:idx_contest_record;not null"`
	UserID       uint64                                   `json:"user_id" gorm:"index:idx_contest_record;not null"`
	ProblemLabel string                                   `json:"problem_label" gorm:"index:idx_contest_record;not null"`
	Code         string                                   `json:"code" gorm:"type:longtext"`
	Language     string                                   `json:"language"`
	Verdict      JudgeVerdict                             `json:"verdict"`
	MaxTime      float32                                  `json:"max_time"`
	MaxMemory    float32                                  `json:"max_memory"`
	JudgeResult  datatypes.JSONSlice[JudgeTestCaseResult] `json:"judge_result"`
	CreatedAt    time.Time                                `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	UpdatedAt    time.Time                                `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
	DeletedAt    gorm.DeletedAt                           `json:"deleted_at" gorm:"index"`
}

func (ContestRecord) TableName() string {
	return "contest_record"
}

// CreateContestRecord 创建比赛提交记录
func CreateContestRecord(record *ContestRecord) error {
	return dao.MysqlClient.Create(record).Error
}

// ContestRecordDetail 提交记录详情（含用户名和题目标题）
type ContestRecordDetail struct {
	ContestRecord
	Username     string `json:"username" gorm:"column:username"`
	ProblemTitle string `json:"problem_title" gorm:"column:problem_title"`
}

// GetContestRecordByID 根据ID获取提交记录详情
func (ContestRecord) GetContestRecordByID(id string) (*ContestRecordDetail, error) {
	var record ContestRecordDetail
	err := dao.MysqlClient.Table("contest_record").
		Select("contest_record.*", "user.username", "cp.title AS problem_title").
		Joins("INNER JOIN user ON user.id = contest_record.user_id").
		Joins("INNER JOIN contest_problem cp ON cp.contest_id = contest_record.contest_id AND cp.label = contest_record.problem_label").
		Where("contest_record.id = ?", id).
		Scan(&record).Error
	return &record, err
}

// GetContestRecords 获取比赛提交列表（分页），使用窗口函数一次查询返回数据和总数
func (ContestRecord) GetContestRecords(contestID string, page, pageSize int, verdict, language, problemLabel string) ([]map[string]interface{}, int64, error) {
	offset := (page - 1) * pageSize

	query := dao.MysqlClient.Table("contest_record").
		Select(`contest_record.id, contest_record.user_id, contest_record.problem_label,
			contest_record.language, contest_record.verdict,
			contest_record.max_time, contest_record.max_memory, contest_record.created_at,
			user.username, cp.title AS problem_title,
			COUNT(*) OVER() AS total`).
		Joins("INNER JOIN user ON user.id = contest_record.user_id").
		Joins("INNER JOIN contest_problem cp ON cp.contest_id = contest_record.contest_id AND cp.label = contest_record.problem_label").
		Where("contest_record.contest_id = ? AND contest_record.deleted_at IS NULL", contestID)

	if verdict != "" {
		query = query.Where("contest_record.verdict = ?", verdict)
	}
	if language != "" {
		query = query.Where("contest_record.language = ?", language)
	}
	if problemLabel != "" {
		query = query.Where("contest_record.problem_label = ?", problemLabel)
	}

	var results []map[string]interface{}
	err := query.Order("contest_record.created_at DESC").
		Limit(pageSize).Offset(offset).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(results) > 0 {
		if t, ok := results[0]["total"]; ok {
			total = t.(int64)
		}
	}

	for i := range results {
		delete(results[i], "total")
	}

	return results, total, nil
}
