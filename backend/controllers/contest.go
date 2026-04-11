package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nexus/dao"
	"nexus/models"
	"nexus/services"
	"nexus/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

var contestCtx = context.Background()

type ContestController struct{}

// ==================== 管理员接口 ====================

// CreateContest 创建比赛
func (ContestController) CreateContest(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil || userID == 0 {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		Title        string  `json:"title" binding:"required"`
		Introduction *string `json:"introduction"`
		ContestType  string  `json:"contest_type" binding:"required"`
		BeginAt      string  `json:"begin_at" binding:"required"`
		EndAt        string  `json:"end_at" binding:"required"`
		Duration     int     `json:"duration"`
		IsPrivate    bool    `json:"is_private"`
		Password     string  `json:"password"`
		SealRank     bool    `json:"seal_rank"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	beginAt, err := time.Parse(time.RFC3339, req.BeginAt)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "开始时间格式错误")
		return
	}
	endAt, err := time.Parse(time.RFC3339, req.EndAt)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "结束时间格式错误")
		return
	}

	if endAt.Before(beginAt) {
		utils.ReturnError(c, http.StatusBadRequest, "结束时间必须晚于开始时间")
		return
	}

	if req.Duration <= 0 {
		req.Duration = int(endAt.Sub(beginAt).Minutes())
	}

	// 哈希比赛密码
	password := req.Password
	if req.IsPrivate && req.Password != "" {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		password = hash
	}

	contest := &models.Contest{
		Title:        req.Title,
		Introduction: req.Introduction,
		UserID:       userID,
		ContestType:  models.ContestType(req.ContestType),
		BeginAt:      beginAt,
		EndAt:        endAt,
		Duration:     req.Duration,
		IsPrivate:    req.IsPrivate,
		Password:     password,
		SealRank:     req.SealRank,
	}
	contest.Status = services.CheckContestStatus(contest)

	contestModel := models.Contest{}
	if err := contestModel.CreateContest(contest); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	// 缓存密码到 Redis
	if req.IsPrivate && req.Password != "" {
		_ = services.CacheContestPassword(contest.ID, password, endAt)
	}

	_ = services.CacheContestInfo(*contest)
	utils.ReturnSuccess(c, http.StatusOK, "创建成功", contest)
}

// UpdateContest 更新比赛
func (ContestController) UpdateContest(c *gin.Context) {
	var req struct {
		ID           string  `json:"id" binding:"required"`
		Title        string  `json:"title"`
		Introduction *string `json:"introduction"`
		ContestType  string  `json:"contest_type"`
		BeginAt      string  `json:"begin_at"`
		EndAt        string  `json:"end_at"`
		Duration     int     `json:"duration"`
		IsPrivate    bool    `json:"is_private"`
		Password     string  `json:"password"`
		SealRank     bool    `json:"seal_rank"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	existing, err := (models.Contest{}).QueryContestById(req.ID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Introduction != nil {
		existing.Introduction = req.Introduction
	}
	if req.ContestType != "" {
		existing.ContestType = models.ContestType(req.ContestType)
	}
	if req.BeginAt != "" {
		beginAt, _ := time.Parse(time.RFC3339, req.BeginAt)
		existing.BeginAt = beginAt
	}
	if req.EndAt != "" {
		endAt, _ := time.Parse(time.RFC3339, req.EndAt)
		existing.EndAt = endAt
	}
	if req.Duration > 0 {
		existing.Duration = req.Duration
	}
	existing.IsPrivate = req.IsPrivate
	existing.SealRank = req.SealRank
	// 哈希比赛密码
	if req.IsPrivate && req.Password != "" {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "密码加密失败")
			return
		}
		existing.Password = hash
	} else {
		existing.Password = req.Password
	}

	contestModel := models.Contest{}
	if err := contestModel.UpdateContest(&existing); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "更新失败")
		return
	}

	// 更新密码缓存
	if existing.IsPrivate && req.Password != "" {
		_ = services.CacheContestPassword(existing.ID, existing.Password, existing.EndAt)
	} else {
		dao.RedisClient.Del(contestCtx, fmt.Sprintf("contest:password:%s", existing.ID))
	}

	services.InvalidateContestInfo(existing.ID)
	utils.ReturnSuccess(c, http.StatusOK, "更新成功", nil)
}

// DeleteContest 删除比赛
func (ContestController) DeleteContest(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	contestModel := models.Contest{}
	if err := contestModel.DeleteContest(req.ID); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "删除失败")
		return
	}

	// 清理 Redis 缓存
	dao.RedisClient.Del(contestCtx, fmt.Sprintf("contest:password:%s", req.ID))
	dao.RedisClient.Del(contestCtx, fmt.Sprintf("contest:participants:%s", req.ID))
	dao.RedisClient.Del(contestCtx, fmt.Sprintf("contest:ranking:%s", req.ID))
	dao.RedisClient.Del(contestCtx, fmt.Sprintf("contest:ranking:detail:%s", req.ID))
	services.InvalidateContestInfo(req.ID)

	utils.ReturnSuccess(c, http.StatusOK, "删除成功", nil)
}

// GetAdminContestList 管理员获取比赛列表
func (ContestController) GetAdminContestList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	contests, total, err := (models.Contest{}).GetAllContests(page, pageSize, search)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"list":  contests,
		"total": total,
	})
}

// GetAdminContestDetail 管理员获取比赛详情
func (ContestController) GetAdminContestDetail(c *gin.Context) {
	contestID := c.Param("id")
	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	problems, _ := (models.ContestProblem{}).GetContestProblems(contestID)

	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"contest":  contest,
		"problems": problems,
	})
}

// SetContestProblems 设置比赛题目（接收完整题目数据）
func (ContestController) SetContestProblems(c *gin.Context) {
	contestID := c.Param("id")

	var req struct {
		Problems []struct {
			Label             string                   `json:"label" binding:"required"`
			Score             int32                    `json:"score"`
			Title             string                   `json:"title" binding:"required"`
			Context           string                   `json:"context"`
			InputDescription  string                   `json:"input_description"`
			OutputDescription string                   `json:"output_description"`
			Tips              string                   `json:"tips"`
			Difficulty        float32                  `json:"difficulty"`
			JudgeCase         []models.ProblemTestCase `json:"judge_case"`
			JudgeConfig       models.JudgeConfig       `json:"judge_config"`
			JudgeSample       []models.JudgeSample     `json:"judge_sample"`
			Tags              []string                 `json:"tags"`
			SourceProblemID   *int64                   `json:"source_problem_id"` // 从题库导入时的来源ID
		} `json:"problems" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 先删除旧数据
	_ = (models.ContestProblem{}).DeleteContestProblems(contestID)

	// 创建新题目
	problems := make([]models.ContestProblem, 0, len(req.Problems))
	for _, p := range req.Problems {
		score := p.Score
		if score <= 0 {
			score = 100
		}
		problems = append(problems, models.ContestProblem{
			ContestID:         contestID,
			Label:             p.Label,
			Score:             score,
			Title:             p.Title,
			Context:           p.Context,
			InputDescription:  p.InputDescription,
			OutputDescription: p.OutputDescription,
			Tips:              p.Tips,
			Difficulty:        p.Difficulty,
			JudgeCase:         p.JudgeCase,
			JudgeConfig:       p.JudgeConfig,
			JudgeSample:       p.JudgeSample,
			Tags:              p.Tags,
			SourceProblemID:   p.SourceProblemID,
		})
	}

	if err := (models.ContestProblem{}).CreateContestProblems(problems); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "设置题目失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "设置成功", nil)
}

// ManageParticipant 管理参赛者
func (ContestController) ManageParticipant(c *gin.Context) {
	contestID := c.Param("id")

	var req struct {
		UserID string `json:"user_id" binding:"required"`
		Action string `json:"action" binding:"required"` // disqualify
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	switch req.Action {
	case "disqualify":
		if err := (models.ContestParticipant{}).Disqualify(contestID, req.UserID); err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "操作失败")
			return
		}
		dao.RedisClient.SRem(contestCtx, fmt.Sprintf("contest:participants:%s", contestID), req.UserID)
	default:
		utils.ReturnError(c, http.StatusBadRequest, "不支持的操作")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "操作成功", nil)
}

// GenerateReport 生成比赛报告
func (ContestController) GenerateReport(c *gin.Context) {
	contestID := c.Param("id")

	report, err := services.GetContestReport(contestID)
	if err != nil {
		// 尝试生成
		if err := services.GenerateContestReport(contestID); err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "生成报告失败")
			return
		}
		report, _ = services.GetContestReport(contestID)
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", report)
}

// ==================== 用户接口 ====================

// GetList 获取公开比赛列表
func (ContestController) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	contests, total, err := (models.Contest{}).GetAllContests(page, pageSize, search)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}

	// 为每个比赛补充参赛人数和当前用户报名状态
	userID, err := ParserToken(c)

	type ContestWithCount struct {
		models.Contest
		IsRegistered bool `json:"is_registered"`
	}

	result := make([]ContestWithCount, 0, len(contests))

	// 批量检查报名状态 (Redis Pipeline，替代 N 次 MySQL 查询)
	contestIDs := make([]string, len(contests))
	for i, ct := range contests {
		contestIDs[i] = ct.ID
	}
	regMap := services.BatchCheckRegistration(contestIDs, userID)

	for _, ct := range contests {
		services.ComputeContestStatus(&ct)
		result = append(result, ContestWithCount{
			Contest:      ct,
			IsRegistered: regMap[ct.ID],
		})
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"list":  result,
		"total": total,
	})
}

// GetContestInfo 获取比赛详情
/**
* 包含比赛基本信息、题目列表、参赛人数、当前用户报名状态
 */
func (ContestController) GetContestInfo(c *gin.Context) {
	contestID := c.Param("id")

	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 实时计算状态
	services.ComputeContestStatus(&contest)

	userID, _ := ParserToken(c)
	// 检查是否已报名 (Redis Bitmap)
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"contest":       contest,
		"is_registered": regMap[contestID],
	})
}

// RegisterContest 报名比赛
func (ContestController) RegisterContest(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		ContestID string `json:"contest_id" binding:"required"`
		Password  string `json:"password"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	contest, err := (models.Contest{}).QueryContestById(req.ContestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 检查是否已报名 (Redis Bitmap)
	regMap := services.BatchCheckRegistration([]string{req.ContestID}, userID)
	if regMap[req.ContestID] {
		utils.ReturnError(c, http.StatusBadRequest, "已报名")
		return
	}

	// 检查密码
	if contest.IsPrivate {
		if req.Password == "" {
			utils.ReturnError(c, http.StatusBadRequest, "需要密码")
			return
		}
		// 先从 Redis 查，找不到再从请求中的 Password 字段比对
		valid, redisErr := services.VerifyContestPassword(req.ContestID, req.Password)
		if redisErr != nil || !valid {
			utils.ReturnError(c, http.StatusBadRequest, "密码错误")
			return
		}
	}

	// 写入数据库
	participant := &models.ContestParticipant{
		ContestID: req.ContestID,
		UserID:    userID,
		Status:    models.ParticipantRegistered,
	}
	if err := (models.ContestParticipant{}).Register(participant); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "报名失败")
		return
	}

	// 缓存到 Redis
	_ = services.CacheParticipant(req.ContestID, userID, contest.EndAt)

	// 参赛人数 +1
	_ = (models.Contest{}).IncrParticipants(req.ContestID)

	utils.ReturnSuccess(c, http.StatusOK, "报名成功", nil)
}

// GetContestProblems 获取比赛题目(需已报名且比赛进行中)

func (ContestController) GetContestProblems(c *gin.Context) {
	contestID := c.Param("id")
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 一次 JSON 聚合查询获取比赛信息和比赛题目
	contestWithProblems, err := (models.Contest{}).GetContestWithProblems(contestID)
	if err != nil || contestWithProblems == nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	problems := contestWithProblems.Problems

	// 检查报名 (Redis Bitmap)
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	if !regMap[contestID] {
		utils.ReturnError(c, http.StatusForbidden, "未报名")
		return
	}

	// 检查比赛状态(未开始不允许查看题目)
	actualStatus := services.CheckContestStatus(&models.Contest{
		BeginAt: contestWithProblems.BeginAt,
		EndAt:   contestWithProblems.EndAt,
	})
	if actualStatus == models.ContestStatusUpcoming {
		utils.ReturnError(c, http.StatusForbidden, "比赛未开始")
		return
	}

	// 查询用户提交状态（优先 Redis 缓存）
	statusMap := services.GetProblemStatusCache(contestID, userID)
	if statusMap == nil {
		// 缓存未命中，回源 MySQL
		rankItems, _ := (models.ContestRankItem{}).GetUserRankItems(contestID, userID)
		statusMap = make(map[string]services.ProblemStatusEntry)
		for _, item := range rankItems {
			status := "wrong"
			if item.IsAccepted {
				status = "accepted"
			}
			statusMap[item.ProblemLabel] = services.ProblemStatusEntry{
				Status: status,
				Score:  item.Score,
			}
		}
		// 回填缓存
		if len(rankItems) > 0 {
			ttl := time.Until(contestWithProblems.EndAt) + 24*time.Hour
			if ttl <= 0 {
				ttl = 48 * time.Hour
			}
			services.BatchSetProblemStatusCache(contestID, userID, rankItems, ttl)
		}
	}

	type ProblemWithStatus struct {
		models.ContestProblem
		MyStatus string `json:"my_status"` // accepted, wrong, unattempted
		MyScore  int32  `json:"my_score"`
	}
	result := make([]ProblemWithStatus, 0, len(problems))
	for _, p := range problems {
		ps := ProblemWithStatus{ContestProblem: p}
		if entry, ok := statusMap[p.Label]; ok {
			ps.MyStatus = entry.Status
			ps.MyScore = entry.Score
		} else {
			ps.MyStatus = "unattempted"
		}
		result = append(result, ps)
	}
	contestWithProblems.Problems = nil
	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"contest":  contestWithProblems,
		"problems": result,
	})
}

// SubmitContestProblem 比赛中提交代码
func (ContestController) SubmitContestProblem(c *gin.Context) {
	contestID := c.Param("id")
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req struct {
		ProblemLabel string `json:"problem_label" binding:"required"`
		Code         string `json:"code" binding:"required"`
		Language     string `json:"language" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取比赛信息
	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 检查比赛状态
	actualStatus := services.CheckContestStatus(&contest)
	if actualStatus != models.ContestStatusLive {
		utils.ReturnError(c, http.StatusForbidden, "比赛未在进行中")
		return
	}

	// 检查报名 (Redis Bitmap)
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	if !regMap[contestID] {
		utils.ReturnError(c, http.StatusForbidden, "未报名")
		return
	}

	// 查找比赛题目
	contestProblem, err := (models.ContestProblem{}).GetContestProblemByLabel(contestID, req.ProblemLabel)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "题目标签不存在")
		return
	}

	// 构建测试用例
	testCases := make([]models.JudgeTestCase, 0, len(contestProblem.JudgeCase))
	for index, tc := range contestProblem.JudgeCase {
		testCases = append(testCases, models.JudgeTestCase{
			CaseID:   index + 1,
			Stdin:    tc.Input,
			Expected: tc.Expected,
		})
	}

	// 生成提交ID
	submissionID := idgen.NextId()

	// 创建判题任务
	task := &services.JudgeTask{
		SubmissionID: submissionID,
		ProblemID:    fmt.Sprintf("%d", contestProblem.ID),
		UserID:       userID,
		Code:         req.Code,
		Language:     req.Language,
		TestCases:    testCases,
		JudgeConfig:  contestProblem.JudgeConfig,
		ContestID:    contestID,
		ProblemLabel: req.ProblemLabel,
	}

	// 提交到判题队列
	result, err := services.GlobalJudgeQueue.SubmitTaskSync(task)
	if err != nil {
		utils.ReturnError(c, http.StatusServiceUnavailable, "判题失败: "+err.Error())
		return
	}

	// 更新比赛统计
	(models.Contest{}).UpdateContestStats(contestID, result.Verdict == models.Accepted)

	// 处理排名更新
	services.ProcessContestSubmissionResult(contestID, userID, req.ProblemLabel, contest, result.Verdict, result.MaxTime)

	// 通知 SSE 订阅者
	services.PublishRankingUpdate(contestID)

	utils.ReturnSuccess(c, http.StatusOK, "判题完成", map[string]interface{}{
		"submission_id": submissionID,
		"problem_label": req.ProblemLabel,
		"verdict":       result.Verdict,
		"max_time":      result.MaxTime,
		"max_memory":    result.MaxMemory,
		"result":        result.Result,
	})
}

// GetContestRanking 获取比赛排名
func (ContestController) GetContestRanking(c *gin.Context) {
	contestID := c.Param("id")

	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 实时计算状态
	services.ComputeContestStatus(&contest)

	// 封榜检查(非管理员)
	if contest.SealRank && contest.Status == models.ContestStatusLive {
		userID, _ := ParserToken(c)
		// 管理员可查看，普通用户只能看到自己的
		_ = userID
		// TODO: 检查管理员角色，非管理员封榜时返回部分数据
	}

	// 先从 Redis 获取
	rankings, err := services.GetContestRankingFromRedis(contestID, 100)
	if err != nil || rankings == nil {
		// Redis 没有，从数据库重建
		if contest.ContestType == models.ContestTypeACM {
			rankings, _ = services.ComputeACMRankingFromDB(contestID, contest)
		} else {
			rankings, _ = services.ComputeOIRankingFromDB(contestID, contest)
		}
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", rankings)
}

// StreamContestRanking SSE 实时排名推送
func (ContestController) StreamContestRanking(c *gin.Context) {
	contestID := c.Param("id")
	userID, _ := ParserToken(c)

	_, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 设置 SSE 头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 订阅排名更新
	sub := services.SubscribeRanking(contestID, userID)
	defer services.UnsubscribeRanking(sub)

	// 发送初始数据
	contest, _ := (models.Contest{}).QueryContestById(contestID)
	var initialData interface{}
	if contest.ContestType == models.ContestTypeACM {
		initialRankings, _ := services.ComputeACMRankingFromDB(contestID, contest)
		initialData = initialRankings
	} else {
		initialRankings, _ := services.ComputeOIRankingFromDB(contestID, contest)
		initialData = initialRankings
	}
	initJSON, _ := json.Marshal(map[string]interface{}{
		"type":    "full",
		"ranking": initialData,
	})
	fmt.Fprintf(c.Writer, "event: ranking-update\ndata: %s\n\n", string(initJSON))
	c.Writer.(http.Flusher).Flush()

	// 监听更新或断开
	ctx := c.Request.Context()
	for {
		select {
		case data, ok := <-sub.Ch:
			if !ok {
				return
			}
			fmt.Fprintf(c.Writer, "event: ranking-update\ndata: %s\n\n", string(data))
			c.Writer.(http.Flusher).Flush()
		case <-ctx.Done():
			return
		case <-time.After(30 * time.Second):
			// 心跳
			fmt.Fprintf(c.Writer, "event: heartbeat\ndata: ping\n\n")
			c.Writer.(http.Flusher).Flush()
		}
	}
}

// GetMyContestStatus 获取当前用户参赛状态
func (ContestController) GetMyContestStatus(c *gin.Context) {
	contestID := c.Param("id")
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}
	services.ComputeContestStatus(&contest)

	isRegistered, _ := (models.ContestParticipant{}).IsRegistered(contestID, userID)
	rankItems, _ := (models.ContestRankItem{}).GetUserRankItems(contestID, userID)

	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"is_registered":  isRegistered,
		"contest_status": contest.Status,
		"contest_type":   contest.ContestType,
		"problems":       rankItems,
	})
}

// GetContestProblemDetail 获取比赛单题详情（需已报名且比赛已开始）

func (ContestController) GetContestProblemDetail(c *gin.Context) {
	contestID := c.Param("id")
	label := c.Param("label")
	userID, err := ParserToken(c)
	if err != nil || userID == 0 {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 一次 JOIN 查询获取比赛信息 + 指定标签的题目
	contest, problem, err := (models.Contest{}).GetContestWithProblemByLabel(contestID, label)
	if err != nil || contest == nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	// 检查报名
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	if !regMap[contestID] {
		utils.ReturnError(c, http.StatusForbidden, "未报名")
		return
	}

	// 检查比赛状态
	actualStatus := services.CheckContestStatus(contest)
	if actualStatus == models.ContestStatusUpcoming {
		utils.ReturnError(c, http.StatusForbidden, "比赛未开始")
		return
	}

	if problem == nil {
		utils.ReturnError(c, http.StatusNotFound, "题目不存在")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"contest": contest,
		"problem": problem,
	})
}

// GetContestSubmissions 获取比赛提交列表
func (ContestController) GetContestSubmissions(c *gin.Context) {
	contestID := c.Param("id")
	userID, err := ParserToken(c)
	if err != nil || userID == 0 {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 检查报名
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	if !regMap[contestID] {
		utils.ReturnError(c, http.StatusForbidden, "未报名")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	verdict := c.Query("verdict")
	language := c.Query("language")
	problemLabel := c.Query("problem_label")

	records, total, err := (models.ContestRecord{}).GetContestRecords(contestID, page, pageSize, verdict, language, problemLabel)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", map[string]interface{}{
		"list":  records,
		"total": total,
	})
}

// GetContestSubmissionDetail 获取比赛单条提交详情
func (ContestController) GetContestSubmissionDetail(c *gin.Context) {
	contestID := c.Param("id")
	recordID := c.Param("rid")
	userID, err := ParserToken(c)
	if err != nil || userID == 0 {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 检查报名
	regMap := services.BatchCheckRegistration([]string{contestID}, userID)
	if !regMap[contestID] {
		utils.ReturnError(c, http.StatusForbidden, "未报名")
		return
	}

	record, err := (models.ContestRecord{}).GetContestRecordByID(recordID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "提交记录不存在")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", record)
}

// GetImportPreview 预览可导入的题目（管理员，比赛结束后）
func (ContestController) GetImportPreview(c *gin.Context) {
	contestID := c.Param("id")

	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	services.ComputeContestStatus(&contest)
	if contest.Status != models.ContestStatusEnded {
		utils.ReturnError(c, http.StatusBadRequest, "比赛尚未结束")
		return
	}

	problems, err := (models.ContestProblem{}).GetContestProblems(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "查询失败")
		return
	}

	type ImportPreviewItem struct {
		models.ContestProblem
		Imported bool `json:"imported"` // 是否已导入
	}

	result := make([]ImportPreviewItem, 0, len(problems))
	for _, p := range problems {
		result = append(result, ImportPreviewItem{
			ContestProblem: p,
			Imported:       p.SourceProblemID != nil,
		})
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", result)
}

// ImportContestProblems 将比赛题目导入题库（管理员）
func (ContestController) ImportContestProblems(c *gin.Context) {
	contestID := c.Param("id")

	contest, err := (models.Contest{}).QueryContestById(contestID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "比赛不存在")
		return
	}

	services.ComputeContestStatus(&contest)
	if contest.Status != models.ContestStatusEnded {
		utils.ReturnError(c, http.StatusBadRequest, "比赛尚未结束")
		return
	}

	var req struct {
		ProblemIDs []int64 `json:"problem_ids" binding:"required"` // contest_problem ID 列表
	}
	if err := c.BindJSON(&req); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}

	type ImportResult struct {
		ContestProblemID int64 `json:"contest_problem_id"`
		NewProblemID     int64 `json:"new_problem_id"`
	}
	results := make([]ImportResult, 0, len(req.ProblemIDs))

	for _, cpID := range req.ProblemIDs {
		cp, err := (models.ContestProblem{}).GetContestProblemByID(cpID)
		if err != nil {
			continue
		}
		// 跳过已导入的
		if cp.SourceProblemID != nil {
			results = append(results, ImportResult{
				ContestProblemID: cpID,
				NewProblemID:     *cp.SourceProblemID,
			})
			continue
		}

		// 创建题库题目
		problem := &models.Problem{
			Title:             cp.Title,
			Context:           cp.Context,
			InputDescription:  cp.InputDescription,
			OutputDescription: cp.OutputDescription,
			Tips:              cp.Tips,
			Difficulty:        cp.Difficulty,
			JudgeCase:         cp.JudgeCase,
			JudgeConfig:       cp.JudgeConfig,
			JudgeSample:       cp.JudgeSample,
			Tags:              cp.Tags,
		}
		if err := (models.Problem{}).CreateProblem(problem); err != nil {
			log.Printf("Failed to import contest problem %d: %v", cpID, err)
			continue
		}

		// 回写 SourceProblemID
		_ = (models.ContestProblem{}).UpdateSourceProblemID(cpID, problem.ID)

		results = append(results, ImportResult{
			ContestProblemID: cpID,
			NewProblemID:     problem.ID,
		})
	}

	utils.ReturnSuccess(c, http.StatusOK, "导入完成", results)
}
