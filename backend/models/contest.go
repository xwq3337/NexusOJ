package models

import (
	"pkg/dao"
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
type ContestToProblems struct {
	ContestID       string `json:"contestID"`
	ProblemID       string `json:"problemID"`
	RemoteType      uint   `json:"romotetype"`
	RemoteProblemID string `json:"remoteproblemID"`
	Submission      int32  `json:"submission"`
	Accept          int32  `json:"accept"`
}

func (Contest) TableName() string {
	return "contest"
}
func (ContestToProblems) TableName() string {
	return "contest_to_problem"
}
func CreateContest(contest *Contest) error {
	err := dao.MysqlClient.Create(contest).Error
	return err
}
func QueryContest(id string) (Contest, error) {
	var contest Contest
	err := dao.MysqlClient.Where("id = ?", id).First(&contest).Error
	return contest, err
}

func GetAllContest() ([]Contest, error) {
	var contests []Contest
	err := dao.MysqlClient.Find(&contests).Error
	return contests, err
}
