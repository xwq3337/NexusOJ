package models

import (
	"pkg/dao"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Problem struct {
	ID                int64                                `json:"id" gorm:"primarykey"`
	UserID            string                               `json:"user_id"`
	Title             string                               `json:"title" gorm:"index"`
	Context           string                               `json:"context"`
	InputDescription  string                               `json:"input_description"`
	OutputDescription string                               `json:"output_description"`
	Tips              string                               `json:"tips"`
	Difficulty        float32                              `json:"difficulty"`
	JudgeCase         datatypes.JSONSlice[ProblemTestCase] `json:"judge_case"`
	JudgeConfig       JudgeConfig                          `json:"judge_config" gorm:"serializer:json"`
	JudgeSample       datatypes.JSONSlice[JudgeSample]     `json:"judge_sample"`
	Tags              datatypes.JSONSlice[string]          `json:"tags"`
	Submission        int32                                `json:"submission"`
	Accept            int32                                `json:"accept"`
	Collection        int32                                `json:"collection"`
	CreatedAt         time.Time                            `json:"created_at"`
	UpdatedAt         time.Time                            `json:"updated_at"`
	DeletedAt         gorm.DeletedAt                       `json:"deleted_at" gorm:"index"`
}
type ProblemDetail struct {
	Problem
	Username string `json:"username"`
}

type JudgeSample struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}
type JudgeConfig struct {
	TimeLimit   int32 `json:"time_limit"`
	MemoryLimit int32 `json:"memory_limit"`
}
type ProblemTestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

func (Problem) GetProblemNumber() (int64, error) {
	var count int64
	err := dao.MysqlClient.Raw("SELECT count(*) FROM problem;").Scan(&count).Error
	return count, err
}

func (Problem) TableName() string {
	return "problem"
}

func (Problem) CreateProblem(problem *Problem) error {
	// 获取当前最大的ID
	var maxID int64
	err := dao.MysqlClient.Model(&Problem{}).Select("COALESCE(MAX(id), 1000)").Scan(&maxID).Error
	if err != nil {
		return err
	}
	// 设置为最大ID+1
	problem.ID = maxID + 1

	// 创建题目
	err = dao.MysqlClient.Create(problem).Error
	return err
}
func (Problem) UpdateProblem(problem *Problem) error {
	err := dao.MysqlClient.Omit("created_at").Updates(&problem).Error
	return err
}
func (Problem) QueryProblemById(id string) (ProblemDetail, error) {
	var problem ProblemDetail
	err := dao.MysqlClient.Model(Problem{}).
		Select("problem.*", "user.username").
		Where("problem.id = ?", id).
		Joins("left join user on user.id = problem.user_id").
		First(&problem).Error
	return problem, err
}
func (Problem) QueryProblemByKeyword(keyword string) ([]ProblemDetail, error) {
	var problems []ProblemDetail
	err := dao.MysqlClient.Model(Problem{}).
		Select("problem.title, problem.id", "user.username").
		Where("problem.title LIKE ?", "%"+keyword+"%").
		// Or("problem.context LIKE ?", "%"+keyword+"%").
		Joins("left join user on user.id = problem.user_id").
		Find(&problems).Error
	if err != nil {
		return nil, err
	}
	return problems, nil
}
func (Problem) GetProblemInfoWithoutUsername(id string) (Problem, error) {
	var p Problem
	err := dao.MysqlClient.Model(Problem{}).Where("id = ?", id).First(&p).Error
	return p, err
}

func (Problem) GetAllProblem() ([]Problem, error) {
	var problems []Problem
	err := dao.MysqlClient.Raw("SELECT id ,title, difficulty,collection, tags, accept,submission, created_at, updated_at FROM problem WHERE deleted_at IS NULL ORDER BY id;").Scan(&problems).Error
	return problems, err
}
