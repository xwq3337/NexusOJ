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
	Excerpt    string                      `json:"excerpt"`                     // 摘要
	Tags       datatypes.JSONSlice[string] `json:"tags"`                        // 标签
	Collection int32                       `json:"collection" gorm:"default:0"` // 收藏
	Like       int32                       `json:"like" gorm:"default:0"`       // 喜欢
	IsPrivate  bool                        `json:"is_private" gorm:"type:tinyint(1);default:0"`
	View       int32                       `json:"view" gorm:"default:0"` // 浏览量
	Status     string                      `json:"status"`
	CreatedAt  time.Time                   `json:"created_at"`
	UpdatedAt  time.Time                   `json:"updated_at"`
	DeletedAt  gorm.DeletedAt              `json:"deleted_at" gorm:"index"`
}

func (Blog) TableName() string {
	return "blog"
}
func CreateBlog(blog *Blog) error {
	err := dao.MysqlClient.Create(blog).Error
	return err
}
func UpdateBlog(blog *Blog) error {
	err := dao.MysqlClient.Model(&Blog{}).Where("id = ?", blog.ID).Select("context", "is_private", "tags", "title", "status").Updates(blog).Error
	return err
}
func QueryBlog(id string) (BlogDetail, error) {
	var blog BlogDetail
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.id", "blog.user_id", "blog.title", "blog.is_private", "blog.status",
			"blog.like", "blog.context", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username", "user.avatar").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.id = ?", id).Scan(&blog).Error
	return blog, err
}
func DeleteBlog(id string) error {
	if err := dao.MysqlClient.Delete(&Blog{ID: id}).Error; err != nil {
		return err
	}
	return nil
}

// 管理员获取博客列表
func (Blog) GetAllBlog() ([]BlogDetail, error) {
	var blogs []BlogDetail
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.id", "blog.excerpt", "blog.user_id", "blog.title", "blog.is_private", "blog.status",
			"blog.like", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.deleted_at IS NULL").
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}

// 获取所有博客的数量
func (Blog) GetBlogNumber() (int64, error) {
	var count int64
	err := dao.MysqlClient.Raw("SELECT count(*) FROM blog;").Scan(&count).Error
	return count, err
}

type BlogDetail struct {
	Blog
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// 某个用户能够查看的博客列表
func (Blog) GetAvailableBlog(user_id string, keywords string, page int, page_size int) ([]BlogDetail, int64, error) {
	var blogs []BlogDetail
	offset := (page - 1) * page_size

	// 构建基础查询
	baseQuery := dao.MysqlClient.Model(&Blog{}).
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", "Normal").
		Where("(blog.title LIKE CONCAT('%',?,'%') OR blog.tags LIKE CONCAT('%',?,'%')) AND blog.deleted_at IS NULL AND (blog.user_id = ? OR blog.is_private = 0)", keywords, keywords, user_id)

	// 执行分页查询
	err := baseQuery.
		Select("blog.id", "blog.view", "blog.excerpt", "blog.user_id", "blog.title", "blog.like", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username").
		Order("blog.created_at DESC").
		Offset(offset).
		Limit(page_size).
		Scan(&blogs).Error
	if err != nil {
		return nil, 0, err
	}

	// 执行计数查询
	var total int64
	countQuery := dao.MysqlClient.Model(&Blog{}).
		Where("status = ?", "Normal").
		Where("(title LIKE CONCAT('%',?,'%') OR tags LIKE CONCAT('%',?,'%')) AND deleted_at IS NULL AND (user_id = ? OR is_private = 0)", keywords, keywords, user_id)
	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return blogs, total, nil
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

// 已经标记删除的博客列表
func (Blog) GetRecycleBlog() ([]BlogDetail, error) {
	var blogs []BlogDetail
	err := dao.MysqlClient.Model(&Blog{}).Unscoped().
		Select("blog.*", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", "Deleted").
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}

// db.Unscoped().Delete(&user) // 永久删除
// db.Delete(&user) // 软删除
func (Blog) GetVerifyList() ([]BlogDetail, error) {
	var blogs []BlogDetail
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.*", "user.username").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", "Pending").
		Order("blog.created_at DESC").
		Scan(&blogs).Error
	return blogs, err
}
