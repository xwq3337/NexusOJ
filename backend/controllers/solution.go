package controllers

import (
	"net/http"
	"nexus/models"
	"nexus/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SolutionController struct{}

// GetSolutions 根据题目ID和条件查询题解列表（窗口函数分页）
func (SolutionController) GetSolutions(c *gin.Context) {
	problemID, _ := strconv.ParseUint(c.Query("problem_id"), 10, 64)
	tag := c.Query("tag")
	keyword := c.Query("keyword")
	status := c.Query("status")
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 50 {
		pageSize = 20
	}

	solutions, total, err := models.QuerySolutions(models.SolutionQueryParams{
		ProblemID: problemID,
		Tag:       tag,
		Keyword:   keyword,
		Status:    status,
		UserID:    userID,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"solutions": solutions,
		"total":     total,
	})
}

// GetSolutionDetail 获取题解详情
func (SolutionController) GetSolutionDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	solution, err := models.GetSolutionByID(id)
	if err != nil || solution == nil {
		utils.ReturnError(c, http.StatusNotFound, "题解不存在")
		return
	}

	// 浏览量 +1（异步）
	go models.IncrSolutionView(id)

	utils.ReturnSuccess(c, http.StatusOK, "success", solution)
}

// CreateSolution 创建题解
func (SolutionController) CreateSolution(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		ProblemID uint64   `json:"problem_id" binding:"required"`
		Title     string   `json:"title" binding:"required"`
		Excerpt   string   `json:"excerpt"`
		Context   string   `json:"context" binding:"required"`
		Tags      []string `json:"tags"`
		Status    string   `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	if req.Status == "" {
		req.Status = "public"
	}

	solution := &models.Solution{
		UserID:    userID,
		ProblemID: req.ProblemID,
		Title:     req.Title,
		Excerpt:   req.Excerpt,
		Context:   req.Context,
		Tags:      req.Tags,
		Status:    req.Status,
	}

	exists, _ := models.ExistsActiveSolutionByUserAndProblem(userID, req.ProblemID)
	if exists {
		utils.ReturnError(c, http.StatusConflict, "你已经发布过该题目的题解")
		return
	}

	if err := models.CreateSolution(solution); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "创建失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", solution)
}

// UpdateSolution 更新题解
func (SolutionController) UpdateSolution(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	var req struct {
		Title   string   `json:"title"`
		Excerpt string   `json:"excerpt"`
		Context string   `json:"context"`
		Tags    []string `json:"tags"`
		Status  string   `json:"status"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	solution := &models.Solution{
		ID:       id,
		UserID:   userID,
		Title:    req.Title,
		Excerpt:  req.Excerpt,
		Context:  req.Context,
		Tags:     req.Tags,
		Status:   req.Status,
	}

	if err := models.UpdateSolution(solution); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "更新失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}

// DeleteSolution 删除题解
func (SolutionController) DeleteSolution(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := models.DeleteSolution(id, userID); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "删除失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}
