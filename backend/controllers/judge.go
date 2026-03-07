package controllers

import (
	"context"
	"fmt"
	"nexus/config"
	"nexus/models"
	"nexus/utils/chttp"
	"time"
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

type RemoteJudge struct{}

func SendCode() {

}
