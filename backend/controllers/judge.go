/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-07-10 15:34:29
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-13 02:05:23
 * @FilePath: /backend/controllers/judge.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package controllers

import (
	"context"
	"fmt"
	"pkg/config"
	"pkg/models"
	"pkg/utils/chttp"
	"time"
)

type JudgeController struct{}

type ReturnStruct struct {
	Msg  string      `json:"msg"`
	Err  string      `json:"err"`
	Info interface{} `json:"info"`
}

func (JudgeController) EvaluateCode(judgeStruct *models.JudgeInputStruct) (*models.JudgeOutputResult, error) {
	addr := fmt.Sprintf("%s:%s", config.JudgeServer, config.JudgeServerPort)
	client := chttp.New(
		chttp.WithBaseURL(addr),
		chttp.WithTimeout(10*time.Second),
		// chttp.WithRetry(3, 3*time.Second, 5*time.Second),
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
