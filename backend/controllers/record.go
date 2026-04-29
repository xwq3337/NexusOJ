package controllers

import (
	"net/http"
	"nexus/models"
	"nexus/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecodeController struct{}

func (RecodeController) GetRecodeInfo(c *gin.Context) {
	id := c.Param("id")

	record, err := models.QueryRecordById(id)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", record)
		return
	}
	utils.ReturnError(c, http.StatusNotFound, err)
}

func (RecodeController) GetList(c *gin.Context) {
	// 获取分页参数，默认第1页，每页10条
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// 转换分页参数
	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 10
	}

	// 获取查询参数
	search := c.Query("search")
	verdict := c.Query("verdict")
	language := c.Query("language")
	problemID := c.Query("problem_id")

	// 调用模型层获取数据
	results, total, err := models.Record{}.GetAllRecord(pageInt, pageSizeInt, search, verdict, language, problemID)

	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
			"data":  results,
			"total": total,
		})
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (RecodeController) CreateRecord(c *gin.Context) {
	data := &models.Record{}
	_ = c.BindJSON(&data)
}

func (RecodeController) GetRecodeListByUser(c *gin.Context) {
	userId := c.Param("id")

	// 获取分页参数，默认第1页，每页20条
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	// 转换分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 获取筛选参数
	verdict := c.Query("verdict")
	language := c.Query("language")

	records, total, err := models.QueryRecordByUserId(userId, page, pageSize, verdict, language)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"records": records,
		"total":   total,
	})
}

func (RecodeController) GetDailyActivity(c *gin.Context) {
	activity, err := models.GetGlobalDailyActivityCached(7)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", activity)
}
