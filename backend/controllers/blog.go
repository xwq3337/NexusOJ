package controllers

import (
	"net/http"
	"pkg/models"
	"pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BlogController struct{}

// 创建博客
func (BlogController) CreateBlog(c *gin.Context) {
	data := &models.Blog{}
	if Err := c.BindJSON(&data); Err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, Err)
		return
	}
	data.ID = uuid.New().String()
	data.UserID, _ = ParserToken(c)
	if err := (models.CreateBlog(data)); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", data)
}

// 修改博客
func (BlogController) UpdateBlog(c *gin.Context) {
	data := &models.Blog{}
	if Err := c.BindJSON(&data); Err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, Err)
		return
	}
	if err := (models.UpdateBlog(data)); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", data)
}

// 删除博客
func (BlogController) DeleteBlog(c *gin.Context) {
	id := c.Query("id")
	err := models.DeleteBlog(id)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", nil)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

// 获取博客详情
func (BlogController) GetBlogInfo(c *gin.Context) {
	id := c.Param("id")
	blog, err := models.QueryBlog(id)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", blog)
		return
	}
	utils.ReturnError(c, http.StatusNotFound, err)
}

// 获取博客数量
func (BlogController) GetNumber(c *gin.Context) {
	count, err := models.Blog{}.GetBlogNumber()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", count)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

// 获取所有博客列表
func (BlogController) GetFullList(c *gin.Context) {
	results, err := models.Blog{}.GetAllBlog()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

// 获取用户的个人博客列表
func (BlogController) GetUserBlogList(c *gin.Context) {
	userID, _ := ParserToken(c)
	results, err := models.Blog{}.GetUserBlogList(userID)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

// 获取回收站的博客列表
func (BlogController) RecycleBlog(c *gin.Context) {
	blogs, err := models.Blog{}.GetRecycleBlog()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", blogs)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (BlogController) GetVerifyList(c *gin.Context) {
	blogs, err := models.Blog{}.GetVerifyList()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", blogs)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

// 获取用户的可见博客列表（考虑其他用户私密博客）
func (BlogController) GetAvailableBlogList(c *gin.Context) {
	userID, _ := ParserToken(c)
	keyword := c.Query("keywords")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	pageInt, err := strconv.Atoi(page)
	pageSizeInt, err := strconv.Atoi(pageSize)
	results, total, err := models.Blog{}.GetAvailableBlog(userID, keyword, pageInt, pageSizeInt)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
			"data":  results,
			"total": total,
		})
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
