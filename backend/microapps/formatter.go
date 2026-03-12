package microapps

import (
	"net/http"
	"nexus/microapps/formatter"
	"nexus/utils"

	"github.com/gin-gonic/gin"
)

func FormatterHandler(c *gin.Context) {
	var params struct {
		Code     string `json:"code" binding:"required"`
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "参数格式错误"+err.Error())
	}
	service := formatter.NewService()
	res, err := service.FormatCode(params.Code, params.Language)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "格式化失败"+err.Error())
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", res)
}

// 返回可以格式化的列表
func FormatterListHandler(c *gin.Context) {
	r := formatter.NewRegistry()
	utils.ReturnSuccess(c, http.StatusOK, "success", r.SupportedLanguages())
}
