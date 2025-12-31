package models

import (
	"pkg/dao"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        string `json:"id" gorm:"primarykey"`
	UserID    string `json:"userID"`
	PostID    string `json:"postID"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type SubComment struct {
	ID         string `json:"id" gorm:"primarykey"`
	UserID     string `json:"userID"`
	ParentID   string `json:"parentID"`
	AncestorID string `json:"ancestorID"`
	Content    string `json:"content"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (Comment) TableName() string {
	return "comment"
}
func (SubComment) TableName() string {
	return "subcomment"
}

func CreateComment(comment *Comment) error {
	err := dao.MysqlClient.Create(comment).Error
	return err
}
func CreateSubComment(subcomment *SubComment) error {
	err := dao.MysqlClient.Create(subcomment).Error
	return err
}

func QueryComments(id string) ([]Comment, error) {
	var comments []Comment
	err := dao.MysqlClient.Where("post_id = ?", id).First(&comments).Error
	return comments, err
}
func QuerySubComments(id string) ([]SubComment, error) {
	var subcomments []SubComment
	err := dao.MysqlClient.Where("parent_id = ?", id).First(&subcomments).Error
	return subcomments, err
}
