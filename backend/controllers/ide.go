package controllers

import (
	"net/http"
	"nexus/models"
	"nexus/utils"

	"github.com/gin-gonic/gin"
)

type IDEController struct{}

func (IDEController) JudgeCode(c *gin.Context) {
	data := &models.JudgeInputStruct{}
	_ = c.BindJSON(&data)
	result, _ := JudgeController{}.EvaluateCode(data)
	utils.ReturnSuccess(c, http.StatusOK, "success", result)
}
