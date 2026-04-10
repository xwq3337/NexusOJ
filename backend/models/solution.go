package models

import (
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
	Status     string                      `json:"status"`                                         // 题解状态（草稿、公开、私密）
	CreatedAt  time.Time                   `json:"created_at" gorm:"autoCreateTime;type:datetime"` // 创建时间
	UpdatedAt  time.Time                   `json:"updated_at" gorm:"autoUpdateTime;type:datetime"` // 更新时间
	DeletedAt  gorm.DeletedAt              `json:"deleted_at" gorm:"index;type:datetime"`          // 删除时间
}
