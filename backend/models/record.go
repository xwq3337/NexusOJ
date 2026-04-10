package models

import (
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
	ProblemId   string                                   `json:"problem_id" gorm:"index:idx_user_problem"`
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
	ID            int64     `json:"id"`
	ProblemID     string    `json:"problem_id"`
	Language      string    `json:"language"`
	Verdict       string    `json:"verdict"`
	MaxTime       float32   `json:"max_time"`
	MaxMemory     float32   `json:"max_memory"`
	CreatedAt     time.Time `json:"created_at"`
	ProblemTitle  string    `json:"problem_title"`
	Total         int64     `json:"-" gorm:"column:total"`
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

	for i := range results {
		delete(results[i], "total")
	}

	return results, total, nil
}
