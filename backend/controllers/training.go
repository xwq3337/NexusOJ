package controllers

import (
	"net/http"
	"pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
)

type TrainingController struct{}

func (TrainingController) GetList(c *gin.Context) {
	results, err := models.GetAllTraining()
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (TrainingController) GetTrainingInfo(c *gin.Context) {
	id := c.Param("id")
	results, err := models.GetTrainingProblems(id)
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	ReturnError(c, http.StatusNotFound, err)
}

func (TrainingController) CreateTraining(c *gin.Context) {
	var json models.CreateTrainingDTO
	if err := c.ShouldBindJSON(&json); err != nil {
		ReturnError(c, http.StatusBadRequest, "参数错误")
		return
	}
	json.ID = idgen.NextId()
	err := models.CreateTraining(&json)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	ReturnSuccess(c, http.StatusOK, "创建成功", nil)
}
