package models

import (
	"nexus/dao"
	"time"
)

// ContestRankItem 比赛排名数据(每用户每题一条记录)
type ContestRankItem struct {
	ID           int64      `json:"id" gorm:"primarykey"`
	ContestID    string     `json:"contest_id" gorm:"index:idx_contest_rank;not null"`
	UserID       uint64     `json:"user_id" gorm:"index:idx_contest_rank_user;not null"`
	ProblemLabel string     `json:"problem_label"`
	IsAccepted   bool       `json:"is_accepted" gorm:"default:false"`
	Attempts     int32      `json:"attempts" gorm:"default:0"`      // 提交次数(AC前)
	AcceptedAt   *time.Time `json:"accepted_at"`                    // 首次AC时间
	Score        int32      `json:"score" gorm:"default:0"`         // OI: 最高分; ACM: 0或1
	TotalPenalty int32      `json:"total_penalty" gorm:"default:0"` // ACM 罚时(秒)
}

func (ContestRankItem) TableName() string {
	return "contest_rank_item"
}

// UpsertRankItem 创建或更新排名项
func (ContestRankItem) UpsertRankItem(item *ContestRankItem) error {
	var existing ContestRankItem
	err := dao.MysqlClient.Where("contest_id = ? AND user_id = ? AND problem_label = ?",
		item.ContestID, item.UserID, item.ProblemLabel).First(&existing).Error

	if err != nil {
		// 记录不存在，创建新记录
		return dao.MysqlClient.Create(item).Error
	}

	// 已AC则不再更新
	if existing.IsAccepted {
		return nil
	}

	// 更新已有记录
	err = dao.MysqlClient.Model(&existing).Updates(map[string]interface{}{
		"is_accepted":   item.IsAccepted,
		"attempts":      item.Attempts,
		"accepted_at":   item.AcceptedAt,
		"score":         item.Score,
		"total_penalty": item.TotalPenalty,
	}).Error
	return err
}

// GetContestRankItems 获取比赛所有排名项
func (ContestRankItem) GetContestRankItems(contestID string) ([]ContestRankItem, error) {
	var items []ContestRankItem
	err := dao.MysqlClient.Where("contest_id = ?", contestID).Find(&items).Error
	return items, err
}

// GetUserRankItems 获取用户在某比赛的所有排名项
func (ContestRankItem) GetUserRankItems(contestID string, userID uint64) ([]ContestRankItem, error) {
	var items []ContestRankItem
	err := dao.MysqlClient.Where("contest_id = ? AND user_id = ?", contestID, userID).Find(&items).Error
	return items, err
}

// ContestRankSummary 排名汇总
type ContestRankSummary struct {
	Rank         int64                               `json:"rank"`
	UserID       uint64                              `json:"user_id"`
	Username     string                              `json:"username"`
	Avatar       *string                             `json:"avatar"`
	Solved       int32                               `json:"solved"`
	TotalPenalty int32                               `json:"total_penalty"` // 秒
	Score        int32                               `json:"score"`         // OI 总分
	Problems     map[string]ContestProblemRankDetail `json:"problems"`
}

// ContestProblemRankDetail 每题排名详情
type ContestProblemRankDetail struct {
	Attempts int32  `json:"attempts"`
	Accepted bool   `json:"accepted"`
	Time     string `json:"time,omitempty"`  // AC 时间格式化
	Score    int32  `json:"score,omitempty"` // OI 分数
}
