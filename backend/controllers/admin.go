package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"nexus/models"
	"nexus/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct{}

func (AdminController) AdminLogin(c *gin.Context) {
	var param struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		utils.ReturnError(c, 400, err.Error())
		return
	}
	user, err := models.User{Username: param.Username, UserRole: "admin"}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", param.Username))
			return
		} else {
			utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}

	if user.ID == 0 {
		utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户", param.Username))
		return
	}
	match, err := utils.VerifyPassword(param.Password, user.Password)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("密码验证出错 %v", err))
		return
	}
	if !match {
		utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("%s: 密码错误", param.Username))
		return
	}
	access_token, _ := generateToken(user, 6*60*60)
	refresh_token, _ := generateToken(user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, access_token)
	tokens = append(tokens, refresh_token)
	utils.ReturnSuccess(c, http.StatusOK, tokens, user)
}

func (AdminController) GetUserList(c *gin.Context) {
	Users, err := models.GetAllUsers()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", Users)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
