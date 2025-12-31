package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JsonSuccStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Info interface{} `json:"info"`
}
type JsonErrStruct struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

func ReturnSuccess(c *gin.Context, code int, msg interface{}, info interface{}) {
	json := &JsonSuccStruct{Code: code, Msg: msg, Info: info}
	c.JSON(http.StatusOK, json)
}

func ReturnError(c *gin.Context, code int, msg interface{}) {
	json := &JsonErrStruct{Code: code, Msg: msg}
	c.JSON(http.StatusOK, json)
}
