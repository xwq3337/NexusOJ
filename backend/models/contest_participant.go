package models

import (
	"nexus/dao"
	"time"
)

type ParticipantStatus string

const (
	ParticipantRegistered   ParticipantStatus = "registered"
	ParticipantDisqualified ParticipantStatus = "disqualified"
)

// ContestParticipant 比赛参赛者
type ContestParticipant struct {
	ID           int64             `json:"id" gorm:"primarykey"`
	ContestID    string            `json:"contest_id" gorm:"index:idx_contest_participant;not null"`
	UserID       uint64            `json:"user_id" gorm:"index:idx_contest_participant;not null"`
	RegisteredAt time.Time         `json:"registered_at"`
	Status       ParticipantStatus `json:"status" gorm:"type:varchar(20);default:'registered'"`
}

func (ContestParticipant) TableName() string {
	return "contest_participant"
}

// Register 报名比赛
func (ContestParticipant) Register(participant *ContestParticipant) error {
	participant.RegisteredAt = time.Now()
	return dao.MysqlClient.Create(participant).Error
}

// IsRegistered 检查用户是否已报名
func (ContestParticipant) IsRegistered(contestID string, userID uint64) (bool, error) {
	var count int64
	err := dao.MysqlClient.Model(&ContestParticipant{}).
		Where("contest_id = ? AND user_id = ? AND status = ?", contestID, userID, ParticipantRegistered).
		Count(&count).Error
	return count > 0, err
}

// GetParticipantCount 获取参赛人数
func (ContestParticipant) GetParticipantCount(contestID string) (int64, error) {
	var count int64
	err := dao.MysqlClient.Model(&ContestParticipant{}).
		Where("contest_id = ? AND status = ?", contestID, ParticipantRegistered).
		Count(&count).Error
	return count, err
}

// GetParticipants 获取参赛者列表(分页)
func (ContestParticipant) GetParticipants(contestID string, page, pageSize int) ([]ContestParticipantDetail, error) {
	var details []ContestParticipantDetail
	offset := (page - 1) * pageSize
	err := dao.MysqlClient.Table("contest_participant").
		Select("contest_participant.*, user.username, user.nickname, user.avatar").
		Joins("LEFT JOIN user ON user.id = contest_participant.user_id").
		Where("contest_participant.contest_id = ?", contestID).
		Order("contest_participant.registered_at ASC").
		Limit(pageSize).Offset(offset).
		Find(&details).Error
	return details, err
}

// Disqualify 取消参赛资格
func (ContestParticipant) Disqualify(contestID, userID string) error {
	return dao.MysqlClient.Model(&ContestParticipant{}).
		Where("contest_id = ? AND user_id = ?", contestID, userID).
		Update("status", ParticipantDisqualified).Error
}

// ContestParticipantDetail 参赛者详情(含用户信息)
type ContestParticipantDetail struct {
	ContestParticipant
	Username string  `json:"username"`
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}
