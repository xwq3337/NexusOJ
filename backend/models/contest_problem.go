package models

import (
	"nexus/dao"
	"time"

	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ContestProblem 比赛题目（独立于题库，包含完整题目数据）
type ContestProblem struct {
	ID                int64                                `json:"id" gorm:"primarykey"`
	ContestID         string                               `json:"contest_id" gorm:"index:idx_contest_problem;not null"`
	Label             string                               `json:"label" gorm:"index:idx_contest_problem;not null"`
	Score             int32                                `json:"score"`
	Title             string                               `json:"title" gorm:"type:varchar(255)"`
	Context           string                               `json:"context" gorm:"type:longtext"`
	InputDescription  string                               `json:"input_description" gorm:"type:longtext"`
	OutputDescription string                               `json:"output_description" gorm:"type:longtext"`
	Tips              string                               `json:"tips" gorm:"type:longtext"`
	Difficulty        float32                              `json:"difficulty"`
	JudgeCase         datatypes.JSONSlice[ProblemTestCase] `json:"judge_case"`
	JudgeConfig       JudgeConfig                          `json:"judge_config" gorm:"serializer:json"`
	JudgeSample       datatypes.JSONSlice[JudgeSample]     `json:"judge_sample"`
	Tags              datatypes.JSONSlice[string]          `json:"tags"`
	Submission        int32                                `json:"submission" gorm:"default:0"`
	Accept            int32                                `json:"accept" gorm:"default:0"`
	SourceProblemID   *int64                               `json:"source_problem_id" gorm:"index:idx_source_problem"`
	CreatedAt         time.Time                            `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	UpdatedAt         time.Time                            `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
}

func (ContestProblem) TableName() string {
	return "contest_problem"
}

// CreateContestProblems 批量创建比赛题目
func (ContestProblem) CreateContestProblems(problems []ContestProblem) error {
	for i := range problems {
		if problems[i].ID == 0 {
			problems[i].ID = idgen.NextId()
		}
	}
	return dao.MysqlClient.Create(&problems).Error
}

// GetContestProblemSummaries 获取比赛题目摘要（不含 context、tips、judge_case、judge_config、judge_sample、input/output_description）
func (ContestProblem) GetContestProblemSummaries(contestID string) ([]ContestProblem, error) {
	var problems []ContestProblem
	err := dao.MysqlClient.Select(
		"id", "contest_id", "label", "score", "title",
		"difficulty", "tags", "submission", "accept",
		"source_problem_id", "created_at", "updated_at",
	).Where("contest_id = ?", contestID).
		Order("label ASC").Find(&problems).Error
	return problems, err
}

// GetContestProblems 获取比赛的所有题目
func (ContestProblem) GetContestProblems(contestID string) ([]ContestProblem, error) {
	var problems []ContestProblem
	err := dao.MysqlClient.Where("contest_id = ?", contestID).
		Order("label ASC").Find(&problems).Error
	return problems, err
}

// GetContestProblemByLabel 根据比赛ID和标签获取单个题目
func (ContestProblem) GetContestProblemByLabel(contestID, label string) (*ContestProblem, error) {
	var problem ContestProblem
	err := dao.MysqlClient.Where("contest_id = ? AND label = ?", contestID, label).First(&problem).Error
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// GetContestProblemByID 根据ID获取单个题目
func (ContestProblem) GetContestProblemByID(id int64) (*ContestProblem, error) {
	var problem ContestProblem
	err := dao.MysqlClient.Where("id = ?", id).First(&problem).Error
	if err != nil {
		return nil, err
	}
	return &problem, nil
}

// DeleteContestProblems 删除比赛的所有题目关联
func (ContestProblem) DeleteContestProblems(contestID string) error {
	return dao.MysqlClient.Where("contest_id = ?", contestID).
		Delete(&ContestProblem{}).Error
}

// IncrSubmission 递增提交/通过计数
func (ContestProblem) IncrSubmission(id int64, isAccepted bool) {
	dao.MysqlClient.Model(&ContestProblem{}).Where("id = ?", id).
		UpdateColumn("submission", gorm.Expr("submission + 1"))
	if isAccepted {
		dao.MysqlClient.Model(&ContestProblem{}).Where("id = ?", id).
			UpdateColumn("accept", gorm.Expr("accept + 1"))
	}
}

// UpdateSourceProblemID 回写导入后的题库ID
func (ContestProblem) UpdateSourceProblemID(id int64, problemID int64) error {
	return dao.MysqlClient.Model(&ContestProblem{}).Where("id = ?", id).
		Update("source_problem_id", problemID).Error
}
