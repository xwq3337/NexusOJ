package controllers

import (
	"fmt"
	"net/http"
	"nexus/dao"
	"nexus/models"
	"nexus/services"
	"nexus/utils"
	"nexus/utils/queue"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
	"github.com/yitter/idgenerator-go/idgen"
)

type ProblemController struct{}

var AddrQueue = queue.NewQueue(6)

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
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", problem)
		return
	}
	utils.ReturnError(c, http.StatusNotFound, err)
}
func (ProblemController) GetList(c *gin.Context) {
	userID, _ := ParserToken(c)
	problems, err := models.Problem{}.GetAllProblem()
	// 获取题目列表的 ID 数组
	problemIDs := make([]string, len(problems))
	for i, problem := range problems {
		problemIDs[i] = fmt.Sprintf("problem:%d", problem.ID)
	}
	//	查询Redis，获取用户 题目状态, 并更新到result中
	ctx := c.Request.Context()
	hash := fmt.Sprintf("user:%s:problem_status", userID)
	statuses, err := dao.RedisClient.HMGet(ctx, hash, problemIDs...).Result()
	if err != nil && err != redis.Nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	fmt.Println(statuses...)
	// 将题目状态更新到 problems 数组中
	result := make([]models.ProblemDTO, len(problems))
	for i, status := range statuses {

		var statusStr string
		switch status {
		case nil:
			statusStr = "unattempted"
		case "Accepted":
			statusStr = "solved"
		default:
			statusStr = "attempted"
		}
		result[i] = models.ProblemDTO{
			Problem: problems[i],
			Status:  statusStr,
		}
	}
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", result)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
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
	if data.ProblemID == "" || data.Code == "" || data.Language == "" {
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
	// 同步状态到 Redis
	ctx := c.Request.Context()
	hash := fmt.Sprintf("user:%s:problem_status", userID)
	key := fmt.Sprintf("problem:%s", data.ProblemID)
	// 只要是Redis 里是Accepted就不改了，其他都改成当前的结果
	currentVerdict, err := dao.RedisClient.HGet(ctx, hash, key).Result()
	if err == redis.Nil {
		// Redis 里没有这个记录，直接设置当前结果
		dao.RedisClient.HSet(ctx, hash, key, string(result.Verdict))
	} else if err != nil {
		// Redis 错误
		fmt.Println("Error getting from Redis:", err)
	} else {
		// Redis 里有记录，检查是否是 Accepted
		if currentVerdict != string(models.Accepted) {
			// 不是 Accepted，更新为当前结果
			dao.RedisClient.HSet(ctx, hash, key, string(result.Verdict))
		}
		// 如果是 Accepted，不做任何操作，保持 Accepted 状态
	}
}

func Equal(a, b interface{}) bool {
	strA, okA := a.(string)
	strB, okB := b.(string)
	if okA && okB {
		return cmp.Equal(strings.TrimRight(strA, "\n"), strings.TrimRight(strB, "\n"))
	}
	return cmp.Equal(a, b)
}
