package models

import (
	"fmt"
	"pkg/config"
	"pkg/dao"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string  `json:"id" gorm:"primaryKey"`
	Username     string  `json:"username" gorm:"type:varchar(50);uniqueIndex;not null"`
	Password     string  `json:"password"`
	Email        *string `json:"email"`
	Nickname     *string `json:"nickname"`
	Introduction *string `json:"introduction"`
	Rating       int16   `json:"rating" gorm:"default:1000"`
	School       *string `json:"school"`
	Avatar       *string `json:"avatar"`
	UserRole     uint    `json:"user_role"`
	Gender       string  `json:"gender"`
	Submission   int32   `json:"submission" gorm:"default:0"`
	Accept       int32   `json:"accept" gorm:"default:0"`
	Codeforces   *string `json:"codeforces"`
	Birthday     *string `json:"birthday"`
	Status       int8    `json:"status" gorm:"default:0"` // 0 正常 1 封禁
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
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
		WHERE 
			(id LIKE CONCAT('%',?,'%') OR 
			username LIKE CONCAT('%',?,'%')) AND 
			nickname LIKE CONCAT('%',?,'%')) AND 
			deleted_at IS NULL 
		ORDER BY created_at DESC`

	err := dao.MysqlClient.Raw(query, key, key).Scan(&users).Error
	return users, err
}
func (user User) QueryUser() (User, error) {
	err := dao.MysqlClient.Find(&user, user).Error
	return user, err
}
func UpdateUser(user *User) error {
	err := dao.MysqlClient.Omit("created_at").Updates(user).Error
	return err
}
func GetAllUsers() ([]User, error) {
	var users []User
	err := dao.MysqlClient.Find(&users).Error
	return users, err
}
func ChangeAvatar(id string) error {
	err := dao.MysqlClient.Model(&User{}).Where("id = ?", id).Update("avatar", fmt.Sprintf("%s/assets/image/%s.png", config.Address, id)).Error
	return err
}

type FriendShips struct {
	ID           uint `json:"id" gorm:"primarykey autoIncrement"`
	UserID1      string
	UserID2      string
	Status       string // pending accepted rejected
	Verification string // 验证消息
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (FriendShips) TableName() string {
	return "friendship"
}

func GetAllFriends(id string) ([]User, error) {
	var users []User
	err := dao.MysqlClient.Raw("SELECT id, nickname, username, avatar FROM user WHERE id IN (SELECT user_id2 as `id`  FROM friendship WHERE user_id1 = ? AND status = 'accepted' UNION ALL SELECT user_id1 as `id` FROM friendship WHERE user_id2 = ? AND status = 'accepted')", id, id).Scan(&users).Error
	return users, err
}

type FriendRequestDetail struct {
	UserID   string    `json:"id" gorm:"column:id"`
	Avatar   string    `json:"avatar"`
	Status   string    `json:"status"`    // pending/rejected
	SendTime time.Time `json:"send_time"` // 新增时间字段
	IsSender bool      `json:"is_sender"` // 是否为本方主动发送
}

func GetNewFriends(id string) ([]FriendRequestDetail, error) {
	var requests []FriendRequestDetail
	err := dao.MysqlClient.
		Table("user u").
		Select(`
        u.id, u.username, u.avatar, 
        f.status, f.created_at as send_time,
        f.user_id1 = ? as is_sender
    `, id).
		Joins(`
        JOIN friendship f ON 
            (f.user_id1 = ? AND u.id = f.user_id2) OR
            (f.user_id2 = ? AND u.id = f.user_id1)
    `, id, id).
		Where("f.status IN (?,?)", "pending", "rejected").
		Scan(&requests).Error
	return requests, err
}
func AddFirend(applicant string, recipient string, Verification string) error {
	var count int8
	err := dao.MysqlClient.Raw("SELECT count(*) FROM friendship WHERE (user_id1 = ? AND user_id2 = ?) OR (user_id1 = ? AND user_id2 = ?)", applicant, recipient, recipient, applicant).Scan(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("已经是好友或者已存在好友申请")
	}
	err = dao.MysqlClient.Create(&FriendShips{
		UserID1:      applicant,
		UserID2:      recipient,
		Status:       "pending", // 初次加好友为等待同意
		Verification: Verification,
	}).Error
	return err
}
func HandleNewFriend(applicant string, recipient string, status string) error {
	if status == "accepted" {
		return AcceptToFriendRequest(recipient, applicant)
	} else {
		return RejectToFriendRequest(recipient, applicant)
	}
}
func AcceptToFriendRequest(recipient string, applicant string) error {
	err := dao.MysqlClient.Model(&FriendShips{}).Where("user_id1 = ?", applicant).Where("user_id2  = ?", recipient).Where("status =  ? ", "pending").Updates(FriendShips{
		UserID1: applicant,
		UserID2: recipient,
		Status:  "accepted",
	}).Error
	return err
}

func RejectToFriendRequest(recipient string, applicant string) error {
	err := dao.MysqlClient.Model(&FriendShips{}).Where("user_id1 = ?", applicant).Where("user_id2  = ?", recipient).Where("status = pending").Updates(FriendShips{
		UserID1: applicant,
		UserID2: recipient,
		Status:  "rejected",
	}).Error
	return err
}
