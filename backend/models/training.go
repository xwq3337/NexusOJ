package models

import (
	"pkg/dao"
	"time"

	"gorm.io/gorm"
)

type Training struct {
	ID           int64          `json:"id" gorm:"primaryKey"`
	Title        string         `json:"title"`
	Introduction *string        `json:"introduction"`
	UserID       string         `json:"user_id"`    // 创建人ID
	Password     *string        `json:"password"`   // 密码 (为空则不需要密码)
	IsPrivate    bool           `json:"is_private"` // 是否私有
	Like         int32          `json:"like"`       // 点赞数
	Collection   int32          `json:"collection"` // 收藏数
	CreatedAt    time.Time      `json:"created_at"` // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"` // 更新时间
	Submission   int32          `json:"submission"` // 提交数
	Accept       int32          `json:"accept"`     // 通过数
	DeletedAt    gorm.DeletedAt `gorm:"index"`      // 软删除
}

// 训练的题目都要来自题库，所以这里的题目ID可以直接使用题目的ID
type TrainingToProblem struct {
	ID         int64 `json:"id"` // 主键
	TrainingID int64 `json:"training_id"`
	ProblemID  int64 `json:"problem_id"`
	Submission int32 `json:"submission"`
	Accept     int32 `json:"accept"`
}

// 训练的题目提交记录
type TrainingToProblemToRecord struct {
	ID                  int64 `json:"id"`          // 主键
	TrainingToProblemID int64 `json:"training_id"` // 标识训练和题目关系
	RecordID            int64 `json:"record_id"`
}

func (Training) TableName() string {
	return "training"
}
func (TrainingToProblem) TableName() string {
	return "training_to_problem" // 训练-题目 关联表
}
func (TrainingToProblemToRecord) TableName() string {
	return "training_to_problem_to_record" // 训练-题目-提交记录 关联表
}

type CreateTrainingDTO struct {
	ID           int64   `json:"id" binding:"required"`
	UserID       string  `json:"user_id" binding:"required"`
	Title        string  `json:"title" binding:"required"`
	Introduction *string `json:"introduction"`
	Password     *string `json:"password"`
	IsPrivate    bool    `json:"is_private"`
	Problems     []int64 `json:"problems" binding:"required"` // 题目ID列表
}

func CreateTraining(value *CreateTrainingDTO) error {
	training := &Training{
		ID:           value.ID,
		Title:        value.Title,
		Introduction: value.Introduction,
		UserID:       value.UserID,
		Password:     value.Password,
		IsPrivate:    value.IsPrivate,
	}
	err := dao.MysqlClient.Create(training).Error
	if err != nil {
		return err
	}
	for _, problemID := range value.Problems {
		problem := &TrainingToProblem{
			TrainingID: value.ID,
			ProblemID:  problemID,
		}
		err = dao.MysqlClient.Create(problem).Error
		if err != nil {
			return err
		}
	}
	return nil
}
func QueryTraining(id string) (Training, error) {
	var training Training
	err := dao.MysqlClient.Where("id = ?", id).First(&training).Error
	return training, err
}

func GetAllTraining() ([]Training, error) {
	var trainings []Training
	err := dao.MysqlClient.Find(&trainings).Error
	return trainings, err
}
func GetTrainingProblems(id string) ([]TrainingToProblem, error) {
	var problems []TrainingToProblem
	err := dao.MysqlClient.Where("training_id = ?", id).First(&problems).Error
	return problems, err
}
