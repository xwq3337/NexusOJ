package controllers

import (
	"net/http"
	"pkg/models"
	"pkg/services"
	"pkg/utils"
	"pkg/utils/queue"
	"strings"

	"github.com/gin-gonic/gin"
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
		utils.ReturnSuccess(c, http.StatusOK, "suceess", problem)
		return
	}
	utils.ReturnError(c, http.StatusNotFound, err)
}
func (ProblemController) GetList(c *gin.Context) {
	result, err := models.Problem{}.GetAllProblem()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", result)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (ProblemController) ChangeProblem(c *gin.Context) {
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
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
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
		UserID:       x.UserID,
		Code:         data.Code,
		Language:     data.Language,
		TestCases:    testCases,
		JudgeConfig:  problem.JudgeConfig,
	}

	// 6. 提交到异步判题队列
	if err := services.GlobalJudgeQueue.SubmitTask(task); err != nil {
		utils.ReturnError(c, http.StatusServiceUnavailable, "判题服务暂时不可用")
		return
	}

	// 7. 立即返回提交ID(异步处理)
	utils.ReturnSuccess(c, http.StatusAccepted, "提交成功,正在判题中...", map[string]interface{}{
		"submission_id": submissionID,
		"problem_id":    data.ProblemID,
		"status":        "pending",
	})
}

// GetSubmissionStatus 查询提交状态
func (ProblemController) GetSubmissionStatus(c *gin.Context) {
	submissionID := c.Param("id")
	if submissionID == "" {
		utils.ReturnError(c, http.StatusBadRequest, "缺少提交ID")
		return
	}

	record, err := models.QueryRecordById(submissionID)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "提交记录不存在")
		return
	}

	// 检查记录是否属于当前用户
	userID, exists := c.Get("user_id")
	if exists && userID != record["user_id"] {
		utils.ReturnError(c, http.StatusForbidden, "无权访问此提交记录")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", record)
}

// GetJudgeQueueStatus 获取判题队列状态
func (ProblemController) GetJudgeQueueStatus(c *gin.Context) {
	status := services.GetQueueStatus()
	utils.ReturnSuccess(c, http.StatusOK, "success", status)
}

func Equal(a, b interface{}) bool {
	strA, okA := a.(string)
	strB, okB := b.(string)
	if okA && okB {
		return cmp.Equal(strings.TrimRight(strA, "\n"), strings.TrimRight(strB, "\n"))
	}
	return cmp.Equal(a, b)
}
