package models

import (
	"fmt"
	"nexus/config"
	"nexus/dao"
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
	err := dao.MysqlClient.Model(User{}).Count(&count).Error
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
	err := dao.MysqlClient.Where("id = ? OR username LIKE CONCAT('%',?,'%') OR nickname LIKE CONCAT('%',?,'%')", key, key, key).Find(&users).Error
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
