/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2026-01-24 02:56:21
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-24 02:56:25
 * @FilePath: /backend/utils/return.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SuccessStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Info interface{} `json:"info"`
}
type ErrorStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, info interface{}) {
	json := &SuccessStruct{Code: code, Msg: msg, Info: info}
	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &ErrorStruct{Code: code, Msg: msg}
	c.JSON(http.StatusOK, json)
}
