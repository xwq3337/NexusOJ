package models

import (
	"nexus/dao"
	"nexus/utils"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Problem struct {
	ID                int64                                `json:"id" gorm:"primarykey"`
	UserID            uint64                               `json:"user_id"`
	Title             string                               `json:"title" gorm:"type:varchar(255)"`
	Context           string                               `json:"context"`
	InputDescription  string                               `json:"input_description"`
	OutputDescription string                               `json:"output_description"`
	Tips              string                               `json:"tips"`
	Difficulty        float32                              `json:"difficulty" gorm:"index:,sort:asc"`
	JudgeCase         datatypes.JSONSlice[ProblemTestCase] `json:"judge_case"`
	JudgeConfig       JudgeConfig                          `json:"judge_config" gorm:"serializer:json"`
	JudgeSample       datatypes.JSONSlice[JudgeSample]     `json:"judge_sample"`
	Tags              datatypes.JSONSlice[string]          `json:"tags"`
	Submission        int32                                `json:"submission"`
	Accept            int32                                `json:"accept"`
	Collection        int32                                `json:"collection"`
	CreatedAt         time.Time                            `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	UpdatedAt         time.Time                            `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
	DeletedAt         gorm.DeletedAt                       `json:"deleted_at" gorm:"index"`
}
type ProblemDetail struct {
	Problem
	Username string `json:"username"`
}

// ProblemListItem 题目列表项（含总数）
type ProblemListItem struct {
	Problem
	Total int64 `json:"-" gorm:"column:total"`
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
	err := dao.MysqlClient.Model(Problem{}).Count(&count).Error
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
	err := dao.MysqlClient.Omit("created_at").Where("id = ?", problem.ID).Updates(&problem).Error
	return err
}
func (Problem) QueryProblemById(id string) (ProblemDetail, error) {
	var problem ProblemDetail
	err := dao.MysqlClient.Model(Problem{}).
		Select("problem.*", "user.username").
		Where("problem.id = ?", id).
		Joins("LEFT JOIN user ON user.id = problem.user_id").
		First(&problem).Error
	return problem, err
}
func (Problem) QueryProblemByKeyword(keyword string) ([]ProblemDetail, error) {
	var problems []ProblemDetail
	err := dao.MysqlClient.Model(Problem{}).
		Select("problem.title, problem.id", "user.username").
		Where("MATCH(problem.title) AGAINST(? IN BOOLEAN MODE)", utils.SanitizeFTSSearch(keyword)).
		Joins("LEFT JOIN user ON user.id = problem.user_id").
		Find(&problems).Error
	if err != nil {
		return []ProblemDetail{}, err
	}
	return problems, nil
}
func (Problem) GetProblemInfoWithoutUsername(problem_id int64) (Problem, error) {
	var p Problem
	err := dao.MysqlClient.Model(Problem{}).Where("id = ?", problem_id).First(&p).Error
	return p, err
}

func (Problem) GetAllProblem() ([]Problem, error) {
	var problems []Problem
	err := dao.MysqlClient.Model(Problem{}).
		Select("id", "title", "difficulty", "collection", "tags", "accept", "submission", "created_at", "updated_at").
		Where("deleted_at IS NULL").Order("id ASC").Find(&problems).Error
	return problems, err
}

// GetAllProblemPaginated 分页查询题目列表，使用窗口函数返回总数
func (Problem) GetAllProblemPaginated(page, pageSize int, search string) ([]ProblemListItem, int64, error) {
	var results []ProblemListItem
	query := dao.MysqlClient.Model(&Problem{}).
		Select("problem.id", "problem.title", "problem.difficulty", "problem.collection",
			"problem.tags", "problem.accept", "problem.submission",
			"problem.created_at", "problem.updated_at",
			"COUNT(*) OVER() AS total").
		Where("problem.deleted_at IS NULL")

	if search != "" {
		query = query.Where(
			"MATCH(problem.title) AGAINST(? IN BOOLEAN MODE) OR problem.id = ?",
			utils.SanitizeFTSSearch(search), search,
		)
	}

	offset := (page - 1) * pageSize
	err := query.Order("problem.id ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(results) > 0 {
		total = results[0].Total
	}

	return results, total, nil
}

// ProblemMeta 题目元数据（用于构建 Redis 索引）
type ProblemMeta struct {
	ID         int64    `json:"id"`
	Difficulty float32  `json:"difficulty"`
	Tags       []string `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
}

// GetAllProblemMeta 获取所有题目的元数据（ID, Difficulty, Tags, CreatedAt）
func (Problem) GetAllProblemMeta() ([]ProblemMeta, error) {
	var problems []Problem
	err := dao.MysqlClient.Model(Problem{}).
		Select("id, difficulty, tags, created_at").
		Where("deleted_at IS NULL").
		Find(&problems).Error
	if err != nil {
		return nil, err
	}
	metas := make([]ProblemMeta, len(problems))
	for i, p := range problems {
		metas[i] = ProblemMeta{
			ID:         p.ID,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
			CreatedAt:  p.CreatedAt,
		}
	}
	return metas, nil
}

// GetProblemsByDifficultyRange 获取指定难度范围内的题目（排除已解决的）
func (Problem) GetProblemsByDifficultyRange(minDiff, maxDiff float32, excludeIDs []int64, limit int) ([]Problem, error) {
	var problems []Problem
	query := dao.MysqlClient.Model(Problem{}).
		Select("id, title, difficulty, tags, accept, submission").
		Where("difficulty BETWEEN ? AND ? AND deleted_at IS NULL", minDiff, maxDiff)
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}
	err := query.Order("accept DESC").Limit(limit).Find(&problems).Error
	return problems, err
}

// GetProblemsByTags 获取包含指定标签的题目（排除指定的 ID）
func (Problem) GetProblemsByTags(tags []string, excludeIDs []int64, limit int) ([]Problem, error) {
	if len(tags) == 0 {
		return nil, nil
	}
	var problems []Problem
	query := dao.MysqlClient.Model(Problem{}).
		Select("id, title, difficulty, tags, accept, submission").
		Where("deleted_at IS NULL")
	// JSON_CONTAINS 检查 tags 字段是否包含指定标签
	for _, tag := range tags {
		query = query.Where("JSON_CONTAINS(tags, ?)", `"`+tag+`"`)
	}
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}
	err := query.Order("accept DESC").Limit(limit).Find(&problems).Error
	return problems, err
}

// GetFreshProblems 获取最新创建的题目（排除已解决的）
func (Problem) GetFreshProblems(excludeIDs []int64, limit int) ([]Problem, error) {
	var problems []Problem
	query := dao.MysqlClient.Model(Problem{}).
		Select("id, title, difficulty, tags, accept, submission, created_at").
		Where("deleted_at IS NULL")
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}
	err := query.Order("created_at DESC").Limit(limit).Find(&problems).Error
	return problems, err
}
