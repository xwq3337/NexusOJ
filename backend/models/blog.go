package models

import (
	"pkg/dao"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blog struct {
	ID         string                      `json:"id" gorm:"primarykey"`        // id
	UserID     string                      `json:"user_id"`                     // 用户id
	Title      string                      `json:"title"`                       // 标题
	Context    string                      `json:"context"`                     // 文本
	Tags       datatypes.JSONSlice[string] `json:"tags"`                        // 标签
	Collection int32                       `json:"collection" gorm:"default:0"` // 收藏
	Like       int32                       `json:"like" gorm:"default:0"`       // 喜欢
	IsPrivate  bool                        `json:"is_private" gorm:"type:tinyint(1);default:0"`
	CreatedAt  time.Time                   `json:"created_at"`
	Status     int32                       `json:"status" gorm:"default:0"` // 状态 (0:待审核 1: 草稿 2: 正常 3: 违规 4: 已删除 )
	UpdatedAt  time.Time                   `json:"updated_at"`
	DeletedAt  gorm.DeletedAt              `gorm:"index"`
}

func (Blog) TableName() string {
	return "blog"
}
func CreateBlog(blog *Blog) error {
	err := dao.MysqlClient.Create(blog).Error
	return err
}
func ChangeBlog(blog *Blog) error {
	err := dao.MysqlClient.Model(&Blog{}).Where("id = ?", blog.ID).Select("context", "is_private", "tags", "title").Updates(blog).Error
	return err
}
func QueryBlog(id string) (Blog, error) {
	var blog Blog
	err := dao.MysqlClient.Where("id = ?", id).First(&blog).Error
	return blog, err
}
func DeleteBlog(id string) error {
	if err := dao.MysqlClient.Delete(&Blog{ID: id}).Error; err != nil {
		return err
	}
	return nil
}
func (Blog) GetAllBlog() ([]BlogWithUsername, error) {
	var blogs []BlogWithUsername
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.id", "blog.user_id", "blog.title", "blog.is_private", "blog.status",
			"blog.like", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.deleted_at IS NULL").
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}
func (Blog) GetBlogNumber() (int64, error) {
	var count int64
	err := dao.MysqlClient.Raw("SELECT count(*) FROM blog;").Scan(&count).Error
	return count, err
}

type BlogWithUsername struct {
	Blog
	Username string `json:"username"`
}

// 某个用户能够查看的博客列表
func (Blog) GetAvailableBlog(user_id string, keywords string) ([]BlogWithUsername, error) {
	var blogs []BlogWithUsername
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.id", "blog.user_id", "blog.title", "blog.like", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", 2). // 正常状态的博客
		Where("(blog.title LIKE CONCAT('%',?,'%') OR blog.tags LIKE CONCAT('%',?,'%')) AND blog.deleted_at IS NULL AND (blog.user_id = ? OR blog.is_private = 0)", keywords, keywords, user_id).
		Order("blog.created_at DESC").Scan(&blogs).Error
	return blogs, err
}

// 某个用户的博客列表
func (Blog) GetUserBlogList(userID string) ([]Blog, error) {
	var blogs []Blog
	err := dao.MysqlClient.Model(&Blog{}).
		Select("id", "title", "status", "`like`", "created_at", "is_private", "tags").
		Where("user_id = ?", userID).Order("created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}
func (Blog) GetRecycleBlog() ([]BlogWithUsername, error) {
	var blogs []BlogWithUsername
	err := dao.MysqlClient.Model(&Blog{}).Unscoped().
		Select("blog.*", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", 4).
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}

// db.Unscoped().Delete(&user) // 永久删除
// db.Delete(&user) // 软删除
func (Blog) GetVerifyList() ([]BlogWithUsername, error) {
	var blogs []BlogWithUsername
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.*", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", 0).
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}
