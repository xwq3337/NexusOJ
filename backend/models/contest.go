package models

import (
	"encoding/json"
	"fmt"
	"nexus/dao"
	"sort"
	"time"

	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

type ContestType string

const (
	ContestTypeACM ContestType = "ACM"
	ContestTypeOI  ContestType = "OI"
)

type ContestStatus string

const (
	ContestStatusUpcoming ContestStatus = "Upcoming"
	ContestStatusLive     ContestStatus = "Live"
	ContestStatusEnded    ContestStatus = "Ended"
)

type Contest struct {
	ID           string         `json:"id" gorm:"primarykey"`
	UserID       uint64         `json:"userID"`
	Title        string         `json:"title"  gorm:"type:varchar(255)"`
	Password     string         `json:"password" gorm:"-"`
	Introduction *string        `json:"introduction"`
	Like         int32          `json:"like" gorm:"default:0"`
	Collection   int32          `json:"collection" gorm:"default:0"`
	IsPrivate    bool           `json:"is_private" gorm:"default:false"`
	Participants int32          `json:"participants" gorm:"default:0"` // 参赛人数
	Duration     int            `json:"duration"`                      // 比赛时长(分钟)
	Submission   int32          `json:"submission" gorm:"default:0"`
	Accept       int32          `json:"accept" gorm:"default:0"`
	ContestType  ContestType    `json:"contest_type" gorm:"type:varchar(10);default:'ACM'"`
	Status       ContestStatus  `json:"status" gorm:"type:varchar(10);default:'Upcoming'"`
	SealRank     bool           `json:"seal_rank" gorm:"default:false"` // 是否封榜
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
	BeginAt      time.Time      `json:"begin_at" gorm:"type:datetime"`
	EndAt        time.Time      `json:"end_at" gorm:"type:datetime"`
}

func (Contest) TableName() string {
	return "contest"
}

// CreateContest 创建比赛
func (Contest) CreateContest(contest *Contest) error {
	contest.ID = fmt.Sprintf("%d", idgen.NextId())
	return dao.MysqlClient.Create(contest).Error
}

// UpdateContest 更新比赛
func (Contest) UpdateContest(contest *Contest) error {
	return dao.MysqlClient.Model(&Contest{}).Where("id = ?", contest.ID).
		Select("title", "introduction", "begin_at", "end_at", "duration",
			"is_private", "contest_type", "seal_rank").
		Updates(contest).Error
}

// QueryContestById 根据ID查询比赛
func (Contest) QueryContestById(id string) (Contest, error) {
	var contest Contest
	err := dao.MysqlClient.Where("id = ?", id).First(&contest).Error
	return contest, err
}

// GetAllContests 获取所有比赛(分页)
func (Contest) GetAllContests(page, pageSize int, search string) ([]Contest, int64, error) {
	var contests []Contest
	var total int64

	query := dao.MysqlClient.Model(&Contest{})
	if search != "" {
		query = query.Where("MATCH(title) AGAINST(? IN BOOLEAN MODE)", search)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&contests).Error
	return contests, total, err
}

// DeleteContest 删除比赛
func (Contest) DeleteContest(id string) error {
	return dao.MysqlClient.Where("id = ?", id).Delete(&Contest{}).Error
}

// UpdateContestStatus 更新比赛状态
func (Contest) UpdateContestStatus(id string, status ContestStatus) error {
	return dao.MysqlClient.Model(&Contest{}).Where("id = ?", id).Update("status", status).Error
}

// IncrParticipants 参赛人数 +1
func (Contest) IncrParticipants(id string) error {
	return dao.MysqlClient.Model(&Contest{}).Where("id = ?", id).
		UpdateColumn("participants", gorm.Expr("participants + 1")).Error
}

// UpdateContestStats 更新比赛提交/通过统计
func (Contest) UpdateContestStats(id string, isAccepted bool) {
	dao.MysqlClient.Model(&Contest{}).Where("id = ?", id).
		UpdateColumn("submission", gorm.Expr("submission + 1"))
	if isAccepted {
		dao.MysqlClient.Model(&Contest{}).Where("id = ?", id).
			UpdateColumn("accept", gorm.Expr("accept + 1"))
	}
}

// GetContestsByStatus 按状态获取比赛列表
func (Contest) GetContestsByStatus(status ContestStatus) ([]Contest, error) {
	var contests []Contest
	err := dao.MysqlClient.Where("status = ?", status).Find(&contests).Error
	return contests, err
}

// ContestWithProblemsJSON 一次 JSON 聚合查询的比赛+题目列表结果
type ContestWithProblemsJSON struct {
	ID           string           `json:"id"`
	UserID       uint64           `json:"userID"`
	Title        string           `json:"title"`
	Introduction *string          `json:"introduction"`
	Like         int32            `json:"like"`
	Collection   int32            `json:"collection"`
	IsPrivate    bool             `json:"is_private"`
	Participants int32            `json:"participants"`
	Duration     int              `json:"duration"`
	Submission   int32            `json:"submission"`
	Accept       int32            `json:"accept"`
	ContestType  ContestType      `json:"contest_type"`
	Status       ContestStatus    `json:"status"`
	SealRank     bool             `json:"seal_rank"`
	BeginAt      time.Time        `json:"begin_at"`
	EndAt        time.Time        `json:"end_at"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	ProblemsJSON string           `json:"-" gorm:"column:problems_json"`
	Problems     []ContestProblem `json:"problems" gorm:"-"`
}

// GetContestWithProblems 一次 JSON 聚合查询获取比赛信息和比赛题目列表
func (Contest) GetContestWithProblems(contestID string) (*ContestWithProblemsJSON, error) {
	var result ContestWithProblemsJSON
	err := dao.MysqlClient.Raw(`
		SELECT
			c.id, c.user_id, c.title, c.introduction, c.like, c.collection,
			c.is_private, c.participants, c.duration, c.submission, c.accept,
			c.contest_type, c.status, c.seal_rank,
			c.begin_at, c.end_at, c.created_at, c.updated_at,
			JSON_ARRAYAGG(
				JSON_OBJECT(
					'id', cp.id,
					'contest_id', cp.contest_id,
					'label', cp.label,
					'score', cp.score,
					'title', cp.title,
					'difficulty', cp.difficulty,
					'tags', cp.tags,
					'submission', cp.submission,
					'accept', cp.accept,
					'source_problem_id', cp.source_problem_id,
					'created_at', DATE_FORMAT(cp.created_at, '%Y-%m-%dT%H:%i:%sZ'),
					'updated_at', DATE_FORMAT(cp.updated_at, '%Y-%m-%dT%H:%i:%sZ')
				)
			) AS problems_json
		FROM contest c
		LEFT JOIN contest_problem cp ON c.id = cp.contest_id
		WHERE c.id = ? AND c.deleted_at IS NULL
		GROUP BY c.id
	`, contestID).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// 解析 JSON 并按 Label 排序
	if result.ProblemsJSON != "" && result.ProblemsJSON != "[null]" {
		json.Unmarshal([]byte(result.ProblemsJSON), &result.Problems)
		sort.Slice(result.Problems, func(i, j int) bool {
			return result.Problems[i].Label < result.Problems[j].Label
		})
	}

	return &result, nil
}

// ContestProblemDetailJSON 一次 JOIN 查询返回的比赛+单题结果
type ContestProblemDetailJSON struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Introduction *string       `json:"introduction"`
	Participants int32         `json:"participants"`
	Duration     int           `json:"duration"`
	ContestType  ContestType   `json:"contest_type"`
	Status       ContestStatus `json:"status"`
	SealRank     bool          `json:"seal_rank"`
	BeginAt      time.Time     `json:"begin_at"`
	EndAt        time.Time     `json:"end_at"`
	ProblemJSON  string        `json:"-" gorm:"column:problem_json"`
}

// GetContestWithProblemByLabel 一次 JOIN 查询获取比赛信息和指定标签的题目
// 使用 JSON_OBJECT 聚合，只命中 1 次 MySQL
func (Contest) GetContestWithProblemByLabel(contestID, label string) (*Contest, *ContestProblem, error) {
	var result ContestProblemDetailJSON
	err := dao.MysqlClient.Raw(`
		SELECT
			c.id, c.title, c.introduction, c.participants,
			c.duration, c.contest_type, c.status, c.seal_rank,
			c.begin_at, c.end_at,
			JSON_OBJECT(
				'id', cp.id,
				'contest_id', cp.contest_id,
				'label', cp.label,
				'score', cp.score,
				'title', cp.title,
				'context', cp.context,
				'input_description', cp.input_description,
				'output_description', cp.output_description,
				'tips', cp.tips,
				'difficulty', cp.difficulty,
				'judge_case', CAST(cp.judge_case AS JSON),
				'judge_config', CAST(cp.judge_config AS JSON),
				'judge_sample', CAST(cp.judge_sample AS JSON),
				'tags', CAST(cp.tags AS JSON),
				'submission', cp.submission,
				'accept', cp.accept,
				'source_problem_id', cp.source_problem_id,
				'created_at', DATE_FORMAT(cp.created_at, '%Y-%m-%dT%H:%i:%sZ'),
				'updated_at', DATE_FORMAT(cp.updated_at, '%Y-%m-%dT%H:%i:%sZ')
			) AS problem_json
		FROM contest c
		LEFT JOIN contest_problem cp ON c.id = cp.contest_id AND cp.label = ?
		WHERE c.id = ? AND c.deleted_at IS NULL
	`, label, contestID).Scan(&result).Error

	if err != nil {
		return nil, nil, err
	}

	contest := &Contest{
		ID:           result.ID,
		Title:        result.Title,
		Introduction: result.Introduction,
		Participants: result.Participants,
		Duration:     result.Duration,
		ContestType:  result.ContestType,
		Status:       result.Status,
		SealRank:     result.SealRank,
		BeginAt:      result.BeginAt,
		EndAt:        result.EndAt,
	}

	var problem ContestProblem
	if result.ProblemJSON != "" && result.ProblemJSON != "null" {
		if err := json.Unmarshal([]byte(result.ProblemJSON), &problem); err != nil {
			return contest, nil, err
		}
		return contest, &problem, nil
	}

	return contest, nil, nil
}
