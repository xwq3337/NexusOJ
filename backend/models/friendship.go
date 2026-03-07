package models

import (
	"fmt"
	"nexus/dao"
	"time"

	"gorm.io/gorm"
)

// 单向好友关系, 只记录用户A添加了用户B为好友, 不记录用户B是否也添加了用户A为好友
type FriendShips struct {
	ID        uint           `json:"id" gorm:"primarykey autoIncrement"`
	UserID    string         `json:"user_id" gorm:"index"`
	FriendID  string         `json:"friend_id" gorm:"index"`
	Remark    *string        `json:"remark"` // 好友备注
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// 好友请求表, 记录用户A向用户B发送好友请求的记录, 包括请求状态等信息
type FriendShipRequest struct {
	ID        uint           `json:"id" gorm:"primarykey autoIncrement"`
	UserID    string         `json:"user_id" gorm:"index"`   // 发起请求的用户ID
	FriendID  string         `json:"friend_id" gorm:"index"` // 接收请求的用户ID
	Status    string         `json:"status" `                // pending 待处理  accepted 已接受  rejected 已拒绝
	Message   *string        `json:"message"`                // 请求消息, 验证信息等
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (FriendShips) TableName() string {
	return "friendship"
}
func (FriendShipRequest) TableName() string {
	return "friendship_request"
}

func (FriendShipRequest) CreateRequest(userID, friendID string, message string) error {
	// 检查是否已经是好友关系
	var count int64

	err := dao.MysqlClient.Model(&FriendShips{}).
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Or("user_id = ? AND friend_id = ?", friendID, userID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("你们已经是好友了")
	}
	// 检查是否已经存在未处理的好友请求（非rejected)

	err = dao.MysqlClient.Model(&FriendShipRequest{}).
		Not("status = ?", "rejected").
		Where("user_id = ? AND friend_id = ?", userID, friendID).
		Or("user_id = ? AND friend_id = ?", friendID, userID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("好友请求已存在")
	}

	// 发送好友请求
	err = dao.MysqlClient.Create(&FriendShipRequest{
		UserID:   userID,
		FriendID: friendID,
		Message:  &message,
		Status:   "pending", // 默认为待处理
	}).Error

	// 返回处理结果
	return err
}

func (FriendShipRequest) HandleRequest(id uint, status string) (string, error) {
	var friendship_request FriendShipRequest
	// 处理好友请求
	err := dao.MysqlClient.Model(&FriendShipRequest{}).
		Where("id = ?", id).
		Update("status", status).Find(&friendship_request).Error

	if err != nil {
		return "", err
	}
	// 如果接受好友请求, 则在好友关系表中添加两条记录
	if status == "accepted" {
		err = dao.MysqlClient.Create(&FriendShips{
			UserID:   friendship_request.UserID,
			FriendID: friendship_request.FriendID,
		}).Error
		err = dao.MysqlClient.Create(&FriendShips{
			UserID:   friendship_request.FriendID,
			FriendID: friendship_request.UserID,
		}).Error
		if err != nil {
			return "", err
		}
	}
	// 返回处理结果
	return "好友请求已接受", nil
}

// GetFriendList 获取用户的好友列表, 包括好友的基本信息和备注等
func (FriendShips) GetFriendList(userID string) ([]FriendshipDTO, error) {
	var friends []FriendshipDTO
	// 查询用户的好友关系, 并通过JOIN查询好友的基本信息
	err := dao.MysqlClient.Table("friendship as fs").
		Select("fs.id",
			"fs.user_id",
			"fs.friend_id",
			"fs.remark",
			"fs.created_at",
			"u.username AS friend_username",
			"u.nickname AS friend_nickname",
			"u.avatar AS friend_avatar").
		Joins("INNER JOIN user as u ON u.id = fs.friend_id").
		Where("fs.user_id = ?", userID).
		Scan(&friends).Error
	return friends, err
}

// GetFriendRequestList 获取用户的好友请求列表, 包括请求的基本信息和状态等
func (FriendShipRequest) GetFriendRequestList(userID string) ([]FriendshipRequestDTO, error) {
	var requests []FriendshipRequestDTO
	// 查询用户的好友请求, 并通过JOIN查询请求发起者的基本信息
	err := dao.MysqlClient.Table("friendship_request as fr").
		Select("fr.id",
			"fr.user_id",
			"fr.friend_id",
			"fr.message",
			"fr.status",
			"fr.created_at",
			"u.username AS friend_username",
			"u.nickname AS friend_nickname",
			"u.avatar AS friend_avatar").
		Joins("INNER JOIN user as u ON u.id = fr.user_id").
		Where("fr.friend_id = ?", userID).
		Scan(&requests).Error
	return requests, err
}
