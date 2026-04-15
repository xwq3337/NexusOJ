package controllers

import (
	"context"
	"fmt"
	"net/http"
	"nexus/dao"
	"nexus/models"
	"nexus/services"
	"nexus/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/yitter/idgenerator-go/idgen"
)

type ProblemController struct{}

func (ProblemController) CreateProblem(c *gin.Context) {
	problem := &models.Problem{}
	if err := c.BindJSON(&problem); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, err)
		return
	}

	err := models.Problem{}.CreateProblem(problem)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", problem)
}

func (ProblemController) GetProblemInfo(c *gin.Context) {
	id := c.Param("id")
	problem, err := models.Problem{}.QueryProblemById(id)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, err)
		return
	}

	// 获取用户题目状态
	myStatus := "unattempted"
	userID, err := ParserToken(c)
	if err == nil && userID > 0 {
		ctx := context.Background()
		key_solved := fmt.Sprintf(ProblemStatusSolvedBit, userID)
		key_attempted := fmt.Sprintf(ProblemStatusAttemptedBit, userID)
		solved := dao.RedisClient.GetBit(ctx, key_solved, problem.ID-1000).Val()
		attempted := dao.RedisClient.GetBit(ctx, key_attempted, problem.ID-1000).Val()
		if solved == 1 {
			myStatus = "solved"
		} else if attempted == 1 {
			myStatus = "attempted"
		}
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"problem":   problem,
		"my_status": myStatus,
	})
}
func (ProblemController) GetList(c *gin.Context) {
	userID, _ := ParserToken(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	problems, total, err := models.Problem{}.GetAllProblemPaginated(page, pageSize, search)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}

	// 获取题目列表的 ID 数组
	problemIDs := make([]int64, len(problems))
	for i, problem := range problems {
		problemIDs[i] = problem.ID - 1000
	}

	// 查询 Redis，获取用户题目状态
	ctx := c.Request.Context()
	pipe := dao.RedisClient.Pipeline()
	cmds_solved := make([]*redis.IntCmd, len(problemIDs))
	cmds_attempted := make([]*redis.IntCmd, len(problemIDs))
	if userID > 0 && len(problemIDs) > 0 {
		SolvedHash := fmt.Sprintf(ProblemStatusSolvedBit, userID)
		AttemptedHash := fmt.Sprintf(ProblemStatusAttemptedBit, userID)
		for i, problemId := range problemIDs {
			cmds_solved[i] = pipe.GetBit(ctx, SolvedHash, problemId)
			cmds_attempted[i] = pipe.GetBit(ctx, AttemptedHash, problemId)
		}
		pipe.Exec(ctx)
	}

	// 组装结果
	result := make([]models.ProblemDTO, len(problems))
	for i, p := range problems {
		var statusStr string
		statusStr = "unattempted"
		if cmds_solved[i].Val() == 1 {
			statusStr = "solved"
		} else if cmds_attempted[i].Val() == 1 {
			statusStr = "attempted"

		}
		result[i] = models.ProblemDTO{
			Problem: p.Problem,
			Status:  statusStr,
		}
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", gin.H{
		"problems": result,
		"total":    total,
	})
}
func (ProblemController) UpdateProblem(c *gin.Context) {
	problem := &models.Problem{}

	_ = c.BindJSON(&problem)
	err := models.Problem{}.UpdateProblem(problem)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", problem)
}
func (ProblemController) SearchProblem(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.ReturnError(c, http.StatusBadRequest, "关键词不能为空")
		return
	}
	result, err := models.Problem{}.QueryProblemByKeyword(keyword)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", result)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (ProblemController) GetNumber(c *gin.Context) {
	count, err := models.Problem{}.GetProblemNumber()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", count)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

func (ProblemController) SubmitProblem(c *gin.Context) {
	// 1. 参数验证
	data := models.ProblemJudgeStruct{}
	userID, _ := ParserToken(c)
	if Err := c.BindJSON(&data); Err != nil {
		utils.ReturnError(c, http.StatusBadRequest, Err)
		return
	}
	// 验证必填字段
	if data.ProblemID == 0 || data.Code == "" || data.Language == "" {
		utils.ReturnError(c, http.StatusBadRequest, "缺少必填字段")
		return
	}

	// 2. 获取题目信息
	problem, err := models.Problem{}.GetProblemInfoWithoutUsername(data.ProblemID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, err)
		return
	}

	// 3. 构建测试用例
	testCases := make([]models.JudgeTestCase, 0, len(problem.JudgeCase))
	for index, testCase := range problem.JudgeCase {
		testCases = append(testCases, models.JudgeTestCase{
			CaseID:   index + 1,
			Stdin:    testCase.Input,
			Expected: testCase.Expected,
		})
	}

	// 4. 生成提交ID
	submissionID := idgen.NextId()

	// 5. 创建判题任务
	task := &services.JudgeTask{
		SubmissionID: submissionID,
		ProblemID:    data.ProblemID,
		UserID:       userID,
		Code:         data.Code,
		Language:     data.Language,
		TestCases:    testCases,
		JudgeConfig:  problem.JudgeConfig,
	}

	// 6. 提交到判题队列并等待结果
	result, err := services.GlobalJudgeQueue.SubmitTaskSync(task)
	if err != nil {
		utils.ReturnError(c, http.StatusServiceUnavailable, fmt.Sprintf("判题失败: %v", err))
		return
	}

	// 7. 返回判题结果
	utils.ReturnSuccess(c, http.StatusOK, "判题完成", map[string]interface{}{
		"submission_id": submissionID,
		"problem_id":    data.ProblemID,
		"verdict":       result.Verdict,
		"max_time":      result.MaxTime,
		"max_memory":    result.MaxMemory,
		"result":        result.Result,
		"status":        "completed",
	})
	// 异步更新状态到 Redis和 Mysql

	go func() {
		// mysql 更新
		models.UpdateSubmission(userID, result.Verdict == "Accepted")
		// redis 更新
		ctx := context.Background()
		key_solved := fmt.Sprintf(ProblemStatusSolvedBit, userID)
		key_attempted := fmt.Sprintf(ProblemStatusAttemptedBit, userID)
		if result.Verdict == "Accepted" {
			dao.RedisClient.SetBit(ctx, key_solved, problem.ID-1000, 1)
		}
		dao.RedisClient.SetBit(ctx, key_attempted, problem.ID-1000, 1)

		// 画像增量更新
		services.SubmitProfileUpdate(&services.ProfileUpdateEvent{
			UserID:     userID,
			ProblemID:  data.ProblemID,
			Verdict:    result.Verdict,
			Language:   data.Language,
			Difficulty: problem.Difficulty,
			Tags:       problem.Tags,
			Timestamp:  time.Now(),
		})
	}()
}

func Equal(a, b interface{}) bool {
	strA, okA := a.(string)
	strB, okB := b.(string)
	if okA && okB {
		return cmp.Equal(strings.TrimRight(strA, "\n"), strings.TrimRight(strB, "\n"))
	}
	return cmp.Equal(a, b)
}
