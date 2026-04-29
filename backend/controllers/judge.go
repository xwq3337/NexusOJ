package controllers

import (
	"context"
	"fmt"
	"net/http"
	"nexus/config"
	"nexus/models"
	"nexus/utils"
	"nexus/utils/chttp"
	"time"

	"github.com/gin-gonic/gin"
)

type JudgeController struct{}

func (JudgeController) EvaluateCode(judgeStruct *models.JudgeInputStruct) (*models.JudgeOutputResult, error) {
	addr := fmt.Sprintf("%s:%s", config.JudgeServer, config.JudgeServerPort)
	client := chttp.New(
		chttp.WithBaseURL(addr),
		chttp.WithTimeout(10*time.Second),
		chttp.WithRetry(3, 3*time.Second, 5*time.Second),
		chttp.WithHeader("Accept", "application/json"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result models.JudgeOutputResult
	err := client.Post(ctx, "/submit", judgeStruct, &result)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return &result, nil
}

// Stats 代理 JudgeServer /stats 接口
func (JudgeController) Stats(c *gin.Context) {
	addr := fmt.Sprintf("%s:%s", config.JudgeServer, config.JudgeServerPort)
	client := chttp.New(
		chttp.WithBaseURL(addr),
		chttp.WithTimeout(5*time.Second),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var result interface{}
	if err := client.Get(ctx, "/stats", &result); err != nil {
		utils.ReturnError(c, http.StatusBadGateway, fmt.Errorf("JudgeServer 连接失败: %v", err))
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", result)
}

type RemoteJudge struct{}

func SendCode() {

}
