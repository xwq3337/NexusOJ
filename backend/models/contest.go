package models

import (
	"time"

	"gorm.io/gorm"
)

type Contest struct {
	ID           string  `json:"id" gorm:"primarykey"`
	Title        string  `json:"title"`
	Password     string  `json:"password"`
	Introduction *string `json:"introduction"`
	Like         int32   `json:"like"`
	Collection   int32   `json:"collection"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	UserID       string         `json:"userID"`
	IsPrivate    bool           `json:"isPrivate"`
	UpdatedAt    time.Time
	BeginAt      string `json:"begin_at"`
	Submission   int32  `json:"submission"`
	Accept       int32  `json:"accept"`
	Duration     int    `json:"duration"`
}

func (Contest) TableName() string {
	return "contest"
}
