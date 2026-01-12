package models

import (
	"pkg/dao"
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
	UserId      string                                   `json:"user_id" gorm:"index:idx_user_problem"`
	ProblemId   string                                   `json:"problem_id" gorm:"index:idx_user_problem"`
	Code        string                                   `json:"code"`
	Language    string                                   `json:"language"`
	Verdict     JudgeVerdict                             `json:"verdict"`
	MaxTime     float32                                  `json:"max_time"`
	MaxMemory   float32                                  `json:"max_memory"`
	JudgeResult datatypes.JSONSlice[JudgeTestCaseResult] `json:"judge_result" gorm:"type:json"`
	CreatedAt   time.Time                                `json:"created_at"`
	UpdatedAt   time.Time                                `json:"updated_at"`
	DeletedAt   gorm.DeletedAt                           `gorm:"index" json:"deleted_at"`
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
func QueryRecordByUserId(userId string) ([]map[string]interface{}, error) {
	var record []map[string]interface{}
	err := dao.MysqlClient.Table("record").
		Select(`
			record.id AS id,
			record.problem_id AS problem_id,
			record.language AS language,
			record.verdict AS verdict,
			record.max_time AS max_time,
			record.max_memory AS max_memory,
			record.created_at AS created_at,
			problem.title AS problem_title
	`).Joins(`INNER JOIN user ON user.id = record.user_id`).
		Joins(`INNER JOIN problem ON problem.id = record.problem_id`).
		Where("record.user_id = ?", userId).Order("created_at DESC").Scan(&record).Error
	return record, err

}

/**
 * 根据记录id获取记录
 */
func QueryRecordById(id string) (map[string]interface{}, error) {
	var record map[string]interface{}
	err := dao.MysqlClient.Table("record").
		Select("record.*", "problem.title AS problem_title", "user.username AS username").
		Joins("INNER JOIN user ON user.id = record.user_id").
		Joins("INNER JOIN problem ON problem.id = record.problem_id").
		Where("record.id = ?", id).Limit(1).Scan(&record).Error
	return record, err
}

/**
 * 获取所有记录（支持分页和查询）
 */
func GetAllRecord(page, pageSize int, search, verdict, language string) ([]map[string]interface{}, int64, error) {
	var results []map[string]interface{}
	var total int64

	// 构建基础查询
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
				user.username AS username
		`).
		Joins(`JOIN user ON user.id = record.user_id`).
		Joins(`JOIN problem ON problem.id = record.problem_id`)

	// 添加搜索条件
	if search != "" {
		query = query.Where("problem.title LIKE ? OR user.username LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 添加状态筛选
	if verdict != "" {
		query = query.Where("record.verdict = ?", verdict)
	}

	// 添加语言筛选
	if language != "" {
		query = query.Where("record.language = ?", language)
	}

	// 获取总数
	countQuery := dao.MysqlClient.Table("record").
		Joins(`JOIN user ON user.id = record.user_id`).
		Joins(`JOIN problem ON problem.id = record.problem_id`)

	if search != "" {
		countQuery = countQuery.Where("problem.title LIKE ? OR user.username LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if verdict != "" {
		countQuery = countQuery.Where("record.verdict = ?", verdict)
	}
	if language != "" {
		countQuery = countQuery.Where("record.language = ?", language)
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 执行查询并返回分页结果
	err = query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&results).Error
	return results, total, err
}
