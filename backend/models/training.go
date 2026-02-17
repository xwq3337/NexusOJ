package models

import (
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

func (Training) TableName() string {
	return "training"
}
