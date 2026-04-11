package models

import (
	"fmt"
	"nexus/dao"
	"nexus/utils"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 题解模型
type Solution struct {
	ID         uint64                      `json:"id" gorm:"primaryKey"`
	UserID     uint64                      `json:"user_id" gorm:"index"`                           // 提交者id
	ProblemID  uint64                      `json:"problem_id" gorm:"index"`                        // 题目id
	Title      string                      `json:"title" gorm:"type:varchar(255)"`                 // 题解标题
	Excerpt    string                      `json:"excerpt"`                                        // 题解摘要
	Context    string                      `json:"context"`                                        // 题解内容
	Tags       datatypes.JSONSlice[string] `json:"tags"`                                           // 题解标签
	Like       int32                       `json:"like" gorm:"default:0"`                          // 点赞数
	Collection int32                       `json:"collection" gorm:"default:0"`                    // 收藏数
	View       int32                       `json:"view" gorm:"default:0"`                          // 浏览量
	Status     string                      `json:"status"`                                         // 题解状态（draft、public、private）
	CreatedAt  time.Time                   `json:"created_at" gorm:"autoCreateTime;type:datetime"` // 创建时间
	UpdatedAt  time.Time                   `json:"updated_at" gorm:"autoUpdateTime;type:datetime"` // 更新时间
	DeletedAt  gorm.DeletedAt              `json:"deleted_at" gorm:"index;type:datetime"`          // 删除时间
}

func (Solution) TableName() string {
	return "solution"
}

// SolutionWithAuthor 题解 + 作者信息（用于列表查询）
type SolutionWithAuthor struct {
	Solution
	Username     string  `json:"username"`
	Avatar       *string `json:"avatar"`
	ProblemTitle string  `json:"problem_title"`
	Total        int64   `json:"-" gorm:"column:total"`
}

// ExistsActiveSolutionByUserAndProblem 检查用户是否已对某题目发布过题解（未软删除）
func ExistsActiveSolutionByUserAndProblem(userID, problemID uint64) (bool, error) {
	var count int64
	err := dao.MysqlClient.Model(&Solution{}).
		Where("user_id = ? AND problem_id = ? AND deleted_at IS NULL", userID, problemID).
		Count(&count).Error
	return count > 0, err
}

// CreateSolution 创建题解
func CreateSolution(s *Solution) error {
	return dao.MysqlClient.Create(s).Error
}

// UpdateSolution 更新题解（仅允许修改标题、摘要、内容、标签、状态）
func UpdateSolution(s *Solution) error {
	return dao.MysqlClient.Model(&Solution{}).Where("id = ? AND user_id = ?", s.ID, s.UserID).
		Select("title", "excerpt", "context", "tags", "status").
		Updates(s).Error
}

// DeleteSolution 删除题解（软删除）
func DeleteSolution(id uint64, userID uint64) error {
	return dao.MysqlClient.Where("id = ? AND user_id = ?", id, userID).
		Delete(&Solution{}).Error
}

// GetSolutionByID 根据 ID 获取题解详情（含作者信息）
func GetSolutionByID(id uint64) (*SolutionWithAuthor, error) {
	var result SolutionWithAuthor
	err := dao.MysqlClient.Raw(`
		SELECT s.*, u.username, u.avatar
		FROM solution s
		LEFT JOIN user u ON u.id = s.user_id
		WHERE s.id = ? AND s.deleted_at IS NULL
	`, id).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// IncrView 浏览量 +1
func IncrSolutionView(id uint64) {
	dao.MysqlClient.Model(&Solution{}).Where("id = ?", id).
		UpdateColumn("view", gorm.Expr("view + 1"))
}

// SolutionQueryParams 题解查询参数
type SolutionQueryParams struct {
	ProblemID uint64
	Tag       string
	Keyword   string
	Status    string
	UserID    uint64
	Page      int
	PageSize  int
}

// QuerySolutions 窗口函数分页查询题解列表 + 总数
func QuerySolutions(params SolutionQueryParams) ([]SolutionWithAuthor, int64, error) {
	// 构建基础条件
	where := "WHERE s.deleted_at IS NULL"
	args := []interface{}{}

	if params.ProblemID > 0 {
		where += " AND s.problem_id = ?"
		args = append(args, params.ProblemID)
	}
	if params.Tag != "" {
		where += " AND JSON_CONTAINS(s.tags, ?)"
		args = append(args, fmt.Sprintf(`"%s"`, params.Tag))
	}
	if params.Keyword != "" {
		where += " AND MATCH(s.title) AGAINST(? IN BOOLEAN MODE)"
		args = append(args, utils.SanitizeFTSSearch(params.Keyword))
	}
	if params.Status != "" {
		where += " AND s.status = ?"
		args = append(args, params.Status)
	} else {
		// 默认只查公开题解
		where += " AND s.status = 'public'"
	}
	if params.UserID > 0 {
		where += " AND s.user_id = ?"
		args = append(args, params.UserID)
	}

	offset := (params.Page - 1) * params.PageSize

	sql := fmt.Sprintf(`
		SELECT * FROM (
			SELECT s.*, u.username, u.avatar, p.title AS problem_title,
				COUNT(*) OVER() AS total
			FROM solution s
			LEFT JOIN user u ON u.id = s.user_id
			LEFT JOIN problem p ON p.id = s.problem_id
			%s
			ORDER BY s.created_at DESC
			LIMIT ? OFFSET ?
		) sub
	`, where)

	var results []SolutionWithAuthor
	err := dao.MysqlClient.Raw(sql, append(args, params.PageSize, offset)...).Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(results) > 0 {
		// 窗口函数的 total 会在每一行都返回，取第一行的即可
		total = results[0].Total
	}

	return results, total, nil
}
