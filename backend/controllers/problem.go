package controllers

import (
	"net/http"
	"pkg/models"
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
	_ = c.BindJSON(&problem)
	problem.ID = idgen.NextId()
	err := models.Problem{}.CreateProblem(problem)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", problem)

}

func (ProblemController) GetProblemInfo(c *gin.Context) {
	id := c.Param("id")
	problem, err := models.Problem{}.QueryProblemById(id)
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "suceess", problem)
		return
	}
	ReturnError(c, http.StatusNotFound, err)
}
func (ProblemController) GetList(c *gin.Context) {
	result, err := models.Problem{}.GetAllProblem()
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", result)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (ProblemController) ChangeProblem(c *gin.Context) {
	problem := &models.Problem{}
	_ = c.BindJSON(&problem)
	err := models.Problem{}.UpdateProblem(problem)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", problem)
}
func (ProblemController) SearchProblem(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		ReturnError(c, http.StatusBadRequest, "关键词不能为空")
		return
	}
	result, err := models.Problem{}.QueryProblemByKeyword(keyword)
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", result)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (ProblemController) GetNumber(c *gin.Context) {
	count, err := models.Problem{}.GetProblemNumber()
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", count)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}

func (ProblemController) SubmitProblem(c *gin.Context) {
	data := models.ProblemJudgeStruct{}
	if Err := c.BindJSON(&data); Err != nil {
		ReturnError(c, http.StatusInternalServerError, Err)
		return
	}
	problem, err := models.Problem{}.GetProblemInfoWithoutUsername(data.ProblemID)
	if err != nil {
		ReturnError(c, http.StatusNotFound, err)
		return
	}
	test_cases := problem.JudgeCase
	var finally_test_cases []models.JudgeTestCase
	for index, test_case := range test_cases {
		finally_test_cases = append(finally_test_cases, models.JudgeTestCase{
			CaseID:   index + 1,
			Stdin:    test_case.Input,
			Expected: test_case.Expected,
		})
	}
	config := models.JudgeInputStruct{
		SubmissionID: idgen.NextId(),
		Language:     data.Language,
		Code:         data.Code,
		TestCases:    finally_test_cases,
		ResourcesLimits: models.JudgeResourcesLimits{
			CpuTime:     1000,              // cpu 配额占比
			MemoryBytes: 128 * 1024 * 1024, // 最大内存
			StackBytes:  128 * 1024 * 1024, // 最大栈空间
			OutputBytes: 1024 * 1024,       // 输出限制
		},
		Message:        "",
		SeccompProfile: "cpp",
	}
	result, err := JudgeController{}.EvaluateCode(&config)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", result)
	record := &models.Record{
		ID:          config.SubmissionID,
		UserId:      data.UserID,
		ProblemId:   data.ProblemID,
		Code:        config.Code,
		Language:    config.Language,
		MaxTime:     result.MaxTime,
		MaxMemory:   result.MaxMemory,
		Verdict:     result.Verdict,
		JudgeResult: result.Result,
	}
	_ = models.CreateRecord(record)
}

func Equal(a, b interface{}) bool {
	strA, okA := a.(string)
	strB, okB := b.(string)
	if okA && okB {
		return cmp.Equal(strings.TrimRight(strA, "\n"), strings.TrimRight(strB, "\n"))
	}
	return cmp.Equal(a, b)
}
