package controllers

import (
	"net/http"
	"pkg/models"
	"pkg/utils"

	"github.com/gin-gonic/gin"
)

type ContestController struct{}

func (ContestController) GetList(c *gin.Context) {
	results, err := models.GetAllContest()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (ContestController) GetContestProblems(c *gin.Context) {
	id := c.Param("id")
	results, err := models.QueryContest(id)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	utils.ReturnError(c, http.StatusNotFound, err)
}
