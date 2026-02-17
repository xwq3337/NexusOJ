package controllers

import (
	"net/http"
	"pkg/utils"

	"github.com/gin-gonic/gin"
)

type ContestController struct{}

func (ContestController) GetList(c *gin.Context) {
	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}

// TODO
func (ContestController) GetContestInfo(c *gin.Context) {
	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}
