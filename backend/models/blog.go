package models

import (
	"nexus/dao"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blog struct {
	ID         uuid.UUID                   `json:"id" gorm:"primarykey type:uuid;default:uuid_generate_v4()"` // id
	UserID     uint64                      `json:"user_id"`                                                   // 用户id
	Title      string                      `json:"title" gorm:"type:varchar(255)"`                            // 标题
	Context    string                      `json:"context"`                                                   // 文本
	Excerpt    string                      `json:"excerpt"`                                                   // 摘要
	Tags       datatypes.JSONSlice[string] `json:"tags"`                                                      // 标签
	Collection int32                       `json:"collection" gorm:"default:0"`                               // 收藏
	Like       int32                       `json:"like" gorm:"default:0"`                                     // 喜欢
	IsPrivate  bool                        `json:"is_private" gorm:"type:tinyint(1);default:0"`
	View       int32                       `json:"view" gorm:"default:0"` // 浏览量
	Status     string                      `json:"status"`
	CreatedAt  time.Time                   `json:"created_at" gorm:"autoCreateTime;type:datetime"`
	UpdatedAt  time.Time                   `json:"updated_at" gorm:"autoUpdateTime;type:datetime"`
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
		Select("blog.id", "blog.user_id",
			"blog.title", "blog.is_private",
			"blog.status", "blog.like",
			"blog.context", "blog.tags",
			"blog.collection", "blog.created_at",
			"blog.updated_at", "u.username", "u.avatar").
		Joins("LEFT JOIN user as u ON blog.user_id = u.id").
		Where("blog.id = ?", id).Scan(&blog).Error
	return blog, err
}
func DeleteBlog(id uuid.UUID) error {
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
	err := dao.MysqlClient.Model(Blog{}).Count(&count).Error
	return count, err
}

type BlogDetail struct {
	Blog
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// 某个用户能够查看的博客列表，使用窗口函数一次查询返回数据和总数
func (Blog) GetAvailableBlog(user_id uint64, keywords string, page int, page_size int) ([]map[string]interface{}, int64, error) {
	offset := (page - 1) * page_size

	var results []map[string]interface{}
	err := dao.MysqlClient.Model(&Blog{}).
		Select("blog.id", "blog.view", "blog.excerpt", "blog.user_id", "blog.title", "blog.like", "blog.tags", "blog.collection", "blog.created_at", "blog.updated_at", "user.username", "COUNT(*) OVER() AS total").
		Joins("LEFT JOIN user ON blog.user_id = user.id").
		Where("blog.status = ?", "Normal").
		Where("(MATCH(blog.title) AGAINST(? IN BOOLEAN MODE) OR blog.tags LIKE CONCAT('%',?,'%')) AND blog.deleted_at IS NULL AND (blog.user_id = ? OR blog.is_private = 0)", keywords, keywords, user_id).
		Order("blog.created_at DESC").
		Offset(offset).
		Limit(page_size).
		Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	if len(results) > 0 {
		if t, ok := results[0]["total"]; ok {
			total = t.(int64)
		}
	}

	for i := range results {
		delete(results[i], "total")
	}

	return results, total, nil
}

// 某个用户的博客列表
func (Blog) GetUserBlogList(userID uint64) ([]Blog, error) {
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
