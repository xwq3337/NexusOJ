package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"pkg/models"
	"pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct{}

func (AdminController) AdminLogin(c *gin.Context) {
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, err := models.User{Username: data["username"], Password: data["password"], UserRole: "admin"}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
			return
		} else {
			utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
	if user.ID == "" {
		utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
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
