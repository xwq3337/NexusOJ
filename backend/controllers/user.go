package controllers

import (
	"errors"
	"fmt"
	"net/http"
	jwtgo "nexus/middleware/jwt"
	"nexus/models"
	"nexus/utils"
	"path/filepath"
	"strings"
	"time"

	"nexus/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yitter/idgenerator-go/idgen"
	"gorm.io/gorm"
)

type UserController struct{}

func (UserController) GetUserInfo(c *gin.Context) {
	id := c.Param("id")
	user, err := models.User{}.QueryUserById(id)
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", user)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}

func (UserController) GetNumber(c *gin.Context) {
	count, err := models.User{}.GetUserNumber()
	if err == nil {
		utils.ReturnSuccess(c, http.StatusOK, "success", count)
		return
	}
	utils.ReturnError(c, http.StatusInternalServerError, err)
}
func (UserController) CreateUser(c *gin.Context) {
	user := &models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	user.ID = fmt.Sprintf("%d", idgen.NextId())
	if err := models.CreateUser(user); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", user)
}
func (UserController) UserLogin(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, 400, err.Error())
		return
	}
	user, err := models.User{Username: params.Username, Password: params.Password}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", params.Username))
			return
		} else {
			utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
	if user.ID == "" {
		utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", params.Username))
		return
	}
	access_token, _ := generateToken(user, 6*60*60)
	refresh_token, _ := generateToken(user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, access_token)
	tokens = append(tokens, refresh_token)
	utils.ReturnSuccess(c, http.StatusOK, tokens, user)
}
func (UserController) UpdateUser(c *gin.Context) {
	user := &models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	user.ID, _ = ParserToken(c)
	err = models.UpdateUser(user)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", user)
}

func (UserController) FuzzyQuery(c *gin.Context) {
	keyWord := c.Query("keyword")
	var users []models.User
	users, err := models.User{}.FuzzyQuery(keyWord)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", users)
}

func (UserController) UpdateAvatar(c *gin.Context) {
	userID, _ := ParserToken(c)
	file, err := c.FormFile("avatar")
	DirPath := config.AvatarDir
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "上传头像失败, err: "+err.Error())
		return
	}

	if err = c.SaveUploadedFile(file, filepath.Join(DirPath, fmt.Sprintf("/%s.png", userID))); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "上传头像失败, err: "+err.Error())
		return
	}

	url, err := models.UpdateAvatar(userID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "更新头像失败, err: "+err.Error())
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", url)
}

// UpdatePassword 更新密码
func (UserController) UpdatePassword(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}
	var params struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	if params.OldPassword == "" || params.NewPassword == "" {
		utils.ReturnError(c, http.StatusBadRequest, "旧密码和新密码不能为空")
		return
	}

	err = models.User{}.UpdatePassword(userID, params.OldPassword, params.NewPassword)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusBadRequest, "旧密码错误")
			return
		}
		utils.ReturnError(c, http.StatusInternalServerError, "更新密码失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "密码更新成功", nil)
}

func (UserController) GetAccessToken(c *gin.Context) {
	var params struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	user, _ := models.User{Username: params.Username, Password: params.Password}.QueryUser()
	token_access, _ := generateToken(user, 6*60*60)
	token_refresh, _ := generateToken(user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, token_access)
	tokens = append(tokens, token_refresh)
	utils.ReturnSuccess(c, http.StatusOK, "success", tokens)
}

func generateToken(user models.User, Time int64) (string, error) {
	j := jwtgo.NewJWT()
	claims := jwtgo.CustomClaims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000, // 签名生效时间
			ExpiresAt: time.Now().Unix() + Time, // 签名过期时间
			Issuer:    "xwq",                    // 签名颁发者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParserToken(c *gin.Context) (string, error) {
	tokenString := c.Request.Header.Get("Authorization")
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return "", errors.New("请求未携带token或token不完整,无权限访问")
	}
	tokenString = strings.Split(tokenString, " ")[1]
	j := jwtgo.NewJWT()
	claims, err := j.ParserToken(tokenString)
	if claims.UserID == "" {
		return "", errors.New("token无效")
	}
	return claims.UserID, err
}

func (UserController) ValidateToken(c *gin.Context) {
	_, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "token无效")
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", "token有效")
}
