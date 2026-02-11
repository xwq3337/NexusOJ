package controllers

import (
	"net/http"
	"pkg/models"
	"pkg/utils"
	"pkg/utils/jsonx"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentController struct{}

func (CommentController) GetComment(c *gin.Context) {
	postID := c.Query("id")
	var comments []models.Comment
	comments, err := models.QueryComments(postID)
	results := []interface{}{}
	for _, v := range comments {
		result, _ := jsonx.ToMap(v)
		result["subcomment"], _ = models.QuerySubComments(v.ID)
		results = append(results, result)
	}
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", results)
}

func (CommentController) CreateComment(c *gin.Context) {
	commmet := &models.Comment{}
	_ = c.Bind(&commmet)
	commmet.ID = uuid.New().String()
	err := models.CreateComment(commmet)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", commmet)

}
func (CommentController) CreateSubComment(c *gin.Context) {
	subcommmet := &models.SubComment{}
	_ = c.Bind(&subcommmet)
	subcommmet.ID = uuid.New().String()
	err := models.CreateSubComment(subcommmet)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", subcommmet)
}
