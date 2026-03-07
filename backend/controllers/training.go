package controllers

import (
	"net/http"
	"nexus/utils"

	"github.com/gin-gonic/gin"
)

type TrainingController struct{}

func (TrainingController) GetList(c *gin.Context) {
	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}
func (TrainingController) GetTrainingInfo(c *gin.Context) {
	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}
