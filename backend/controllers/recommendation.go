package controllers

import (
	"net/http"
	"nexus/services"
	"nexus/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecommendationController struct{}

// GetMyProfile 获取当前用户完整画像
func (RecommendationController) GetMyProfile(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未登录")
		return
	}

	profile, err := services.EnsureProfile(userID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", profile)
}

// GetUserProfile 获取其他用户的公开画像
func (RecommendationController) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	profile, err := services.EnsureProfile(id)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", profile)
}

// GetAbilityAnalysis 获取能力分析（各标签掌握度详情）
func (RecommendationController) GetAbilityAnalysis(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未登录")
		return
	}

	profile, err := services.EnsureProfile(userID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"overall_score":  profile.Ability.OverallScore,
		"tag_scores":     profile.Ability.TagScores,
		"tag_progress":   profile.Ability.TagProgress,
		"tag_total":      profile.Ability.TagTotal,
		"strongest_tags": profile.Ability.StrongestTags,
		"weakest_tags":   profile.Ability.WeakestTags,
		"languages":      profile.Preferences.Languages,
		"avg_difficulty": profile.Preferences.AvgDifficulty,
	})
}

// GetActivityStats 获取活跃度统计（热力图、streak）
func (RecommendationController) GetActivityStats(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未登录")
		return
	}

	profile, err := services.EnsureProfile(userID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	// 如果请求指定了年份，返回该年的热力图
	year := c.Query("year")
	if year != "" {
		heatmap, err := services.GetHeatmapByYear(userID, year)
		if err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, err)
			return
		}
		utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
			"year":    year,
			"heatmap": heatmap,
		})
		return
	}

	// 默认返回近一年的热力图
	pastYearHeatmap, _ := services.GetPastYearHeatmap(userID)

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"streak":             profile.Activity.Streak,
		"last_active":        profile.Activity.LastActive,
		"past_year_heatmap":  pastYearHeatmap,
		"heatmaps_by_year":   profile.Activity.Heatmaps,
	})
}

// GetRecommendations 获取推荐题目列表
func (RecommendationController) GetRecommendations(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	refresh := c.Query("refresh") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 10
	}

	problems, total, err := services.GetRecommendations(userID, page, pageSize, refresh)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"problems": problems,
		"total":    total,
	})
}

// RefreshRecommendations 强制刷新推荐结果
func (RecommendationController) RefreshRecommendations(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未登录")
		return
	}

	problems, total, err := services.GetRecommendations(userID, 1, 10, true)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "推荐已刷新", gin.H{
		"problems": problems,
		"total":    total,
	})
}
