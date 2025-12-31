package controllers

import (
	"net/http"
	"pkg/models"

	"github.com/gin-gonic/gin"
)

type RecodeController struct{}

func (RecodeController) GetRecodeInfo(c *gin.Context) {
	id := c.Param("id")

	record, err := models.QueryRecordById(id)
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", record)
		return
	}
	ReturnError(c, http.StatusNotFound, err)
}
func (RecodeController) GetList(c *gin.Context) {
	results, err := models.GetAllRecord()

	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", results)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (RecodeController) CreateRecord(c *gin.Context) {
	data := &models.Record{}
	_ = c.BindJSON(&data)
}

func (RecodeController) GetRecodeListByUser(c *gin.Context) {
	results, err := models.QueryRecordByUserId(c.Param("id"))
	if err != nil {
		ReturnError(c, http.StatusNotFound, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", results)

}
