package models

import (
	"fmt"
	"pkg/config"
	"pkg/dao"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password     string         `json:"password"`
	Email        *string        `json:"email"`
	Nickname     *string        `json:"nickname"`
	Introduction *string        `json:"introduction"`
	Rating       int16          `json:"rating" gorm:"default:1000"`
	School       *string        `json:"school"`
	Avatar       *string        `json:"avatar"`
	UserRole     string         `json:"user_role"`
	Gender       string         `json:"gender"`
	Submission   int32          `json:"submission" gorm:"default:0"`
	Accept       int32          `json:"accept" gorm:"default:0"`
	Codeforces   *string        `json:"codeforces"`
	Birthday     *string        `json:"birthday"`
	Status       int8           `json:"status" gorm:"default:0"` // 0 正常 1 封禁
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	BannedTo     *time.Time     `json:"banned_to"`                // 封禁到期时间
	Balance      float64        `json:"balance" gorm:"default:0"` // x币余额
}

func (User) TableName() string {
	return "user"
}

func (User) GetUserNumber() (int64, error) {
	var count int64
	err := dao.MysqlClient.Raw("SELECT count(*) FROM user;").Scan(&count).Error
	return count, err
}
func CreateUser(user *User) error {
	err := dao.MysqlClient.Create(user).Error
	return err
}

func (User) QueryUserById(id string) (User, error) {
	var user User
	err := dao.MysqlClient.Where("id = ?", id).First(&user).Error
	return user, err
}
func (User) FuzzyQuery(key string) ([]User, error) {
	var users []User
	query := `
		SELECT * 
		FROM user 
		WHERE id = ? OR username LIKE CONCAT('%',?,'%') OR nickname LIKE CONCAT('%',?,'%')	
	`
	err := dao.MysqlClient.Raw(query, key, key, key).Scan(&users).Error
	return users, err
}
func (user User) QueryUser() (User, error) {
	err := dao.MysqlClient.Find(&user, user).Error
	return user, err
}

func UpdateUser(user *User) error {
	err := dao.MysqlClient.Model(&User{}).Where("id = ?", user.ID).Omit("id", "created_at", "updated_at", "deleted_at", "banned_to", "balance", "status", "submission", "accept", "user_role", "avatar", "password").Updates(user).Error
	return err
}
func (User) UpdatePassword(userID, oldPassword, newPassword string) error {
	var user User
	err := dao.MysqlClient.Where("id = ? AND password = ?", userID, oldPassword).First(&user).Error
	if err != nil {
		return err
	}
	err = dao.MysqlClient.Model(&User{}).Where("id = ?", userID).Update("password", newPassword).Error
	return err
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := dao.MysqlClient.Find(&users).Error
	return users, err
}
func UpdateAvatar(id string) (string, error) {
	url := fmt.Sprintf("%s:%s/assets/avatar/%s.png", config.Address, config.Port, id)
	err := dao.MysqlClient.Model(&User{}).Where("id = ?", id).Update("avatar", url).Error
	return url, err
}

type FriendShips struct {
	ID        uint           `json:"id" gorm:"primarykey autoIncrement"`
	UserID    string         `json:"user_id" gorm:"index"`
	FriendID  string         `json:"friend_id" gorm:"index"`
	Status    int8           `json:"status" gorm:"default:0"` // 0 待处理 1 已接受 2 已拒绝
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (FriendShips) TableName() string {
	return "friendship"
}

// CreateFriendRequest 创建好友请求
func CreateFriendRequest(friendship *FriendShips) error {
	err := dao.MysqlClient.Create(friendship).Error
	return err
}

// GetFriendList 获取用户的所有好友（已接受的）
func GetFriendList(userID string) ([]User, error) {
	var friends []User
	// 查找所有状态为1（已接受）的好友关系
	err := dao.MysqlClient.
		Table("user").
		Joins("JOIN friendship ON (user.id = friendship.friend_id OR user.id = friendship.user_id)").
		Where("friendship.user_id = ? AND friendship.status = 1 AND user.id != ? AND user.deleted_at IS NULL", userID, userID).
		Find(&friends).Error
	return friends, err
}

// GetFriendRequestList 获取收到的好友请求列表（待处理的）
func GetFriendRequestList(userID string) ([]FriendShips, error) {
	var requests []FriendShips
	err := dao.MysqlClient.
		Where("friend_id = ? AND status = 0", userID).
		Find(&requests).Error
	return requests, err
}

// HandleFriendRequest 处理好友请求
func HandleFriendRequest(friendshipID uint, accept bool) error {
	var friendship FriendShips
	err := dao.MysqlClient.First(&friendship, friendshipID).Error
	if err != nil {
		return err
	}

	if accept {
		// 接受好友请求，创建双向好友关系
		friendship.Status = 1
		err = dao.MysqlClient.Save(&friendship).Error
		if err != nil {
			return err
		}

		// 检查是否已有反向好友关系
		var reverseFriendship FriendShips
		result := dao.MysqlClient.Where("user_id = ? AND friend_id = ?", friendship.FriendID, friendship.UserID).First(&reverseFriendship)
		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			// 创建反向好友关系
			reverseFriendship = FriendShips{
				UserID:   friendship.FriendID,
				FriendID: friendship.UserID,
				Status:   1,
			}
			err = dao.MysqlClient.Create(&reverseFriendship).Error
			if err != nil {
				return err
			}
		} else if result.Error == nil {
			// 更新反向关系状态
			reverseFriendship.Status = 1
			err = dao.MysqlClient.Save(&reverseFriendship).Error
			if err != nil {
				return err
			}
		}
	} else {
		// 拒绝好友请求
		friendship.Status = 2
		err = dao.MysqlClient.Save(&friendship).Error
		if err != nil {
			return err
		}
	}

	return nil
}
