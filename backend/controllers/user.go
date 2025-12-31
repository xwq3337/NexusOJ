package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	jwtgo "pkg/middleware/jwt"
	"pkg/models"
	"pkg/utils/logger"
	"time"

	"pkg/config"

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
		ReturnSuccess(c, http.StatusOK, "suceess", user)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (UserController) GetList(c *gin.Context) {
	Users, err := models.GetAllUsers()
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", Users)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (UserController) GetNumber(c *gin.Context) {
	count, err := models.User{}.GetUserNumber()
	if err == nil {
		ReturnSuccess(c, http.StatusOK, "success", count)
		return
	}
	ReturnError(c, http.StatusInternalServerError, err)
}
func (UserController) CreateUser(c *gin.Context) {
	user := &models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	user.ID = fmt.Sprintf("%d", idgen.NextId())
	if err := models.CreateUser(user); err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", user)
}
func (UserController) UserLogin(c *gin.Context) {
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, err := models.User{Username: data["username"], Password: data["password"]}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
			return
		} else {
			ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
	if user.ID == "" {
		ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
		return
	}
	access_token, _ := generateToken(c, user, 6*60*60)
	refresh_token, _ := generateToken(c, user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, access_token)
	tokens = append(tokens, refresh_token)
	ReturnSuccess(c, http.StatusOK, tokens, user)
}
func (UserController) AdminLogin(c *gin.Context) {
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, err := models.User{Username: data["username"], Password: data["password"], UserRole: 2}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
			return
		} else {
			ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
	if user.ID == "" {
		ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", data["username"]))
		return
	}
	access_token, _ := generateToken(c, user, 6*60*60)
	refresh_token, _ := generateToken(c, user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, access_token)
	tokens = append(tokens, refresh_token)
	ReturnSuccess(c, http.StatusOK, tokens, user)
}
func (UserController) GetAllFriends(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	id := x.UserID
	var users []models.User
	users, err := models.GetAllFriends(id)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", users)
}
func (UserController) FuzzyQuery(c *gin.Context) {
	keyWord := c.Query("keyword")
	logger.Debug("模糊查询关键字: %s", keyWord)
	var users []models.User
	users, err := models.User{}.FuzzyQuery(keyWord)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", users)
}

func (UserController) AddFirend(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	applicant := x.UserID
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	recipient := data["recipient"]
	verification := data["verification"]
	err := models.AddFirend(applicant, recipient, verification)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", nil)
}
func (UserController) UserLogout(c *gin.Context) {
	ReturnSuccess(c, http.StatusOK, "success", nil)
}
func (UserController) GetNewFriends(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	id := x.UserID
	items, err := models.GetNewFriends(id)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err.Error())
		return
	}

	ReturnSuccess(c, http.StatusOK, "success", items)
}
func (UserController) HandleNewFriend(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	receive_id := x.UserID
	apply_id := c.Query("apply_id")
	status := c.Query("status")
	if len(apply_id) == 0 || (status != "rejected" && status != "accepted") {
		ReturnError(c, http.StatusInternalServerError, "状态错误")
		return
	}
	err := models.HandleNewFriend(apply_id, receive_id, status)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", nil)
}
func (UserController) ChangeAvatar(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	id := x.UserID
	file, err := c.FormFile("avatar")
	DirPath := config.ImagesDir
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, "上传头像失败-1")
		return
	}

	if err = c.SaveUploadedFile(file, filepath.Join(DirPath, fmt.Sprintf("/%s.png", id))); err != nil {
		ReturnError(c, http.StatusInternalServerError, "上传头像失败-2")
		return
	}
	models.ChangeAvatar(id)
	ReturnSuccess(c, http.StatusOK, "success", nil)

}

func (UserController) ChangePassword(c *gin.Context) {
	x, _ := ParserToken(c.Request.Header.Get("Authorization"))
	id := x.UserID
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	if len(data["oldPassword"]) == 0 || len(data["newPassword"]) == 0 {
		ReturnError(c, http.StatusInternalServerError, "密码不能为空")
		return
	}
	user, err := models.User{ID: id, Password: data["oldPassword"]}.QueryUser()
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, "旧密码错误")
		return
	}
	user.Password = data["newPassword"]
	err = models.UpdateUser(&user)
	if err != nil {
		ReturnError(c, http.StatusInternalServerError, "修改密码失败")
		return
	}
	ReturnSuccess(c, http.StatusOK, "success", nil)
}
func (UserController) GetAccessToken(c *gin.Context) {
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, _ := models.User{Username: data["username"], Password: data["password"]}.QueryUser()
	token_access, _ := generateToken(c, user, 6*60*60)
	token_refresh, _ := generateToken(c, user, 7*24*60*60)
	var tokens []string
	tokens = append(tokens, token_access)
	tokens = append(tokens, token_refresh)
	ReturnSuccess(c, http.StatusOK, "success", tokens)
}

func generateToken(c *gin.Context, user models.User, Time int64) (string, error) {
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

func ParserToken(tokenString string) (*jwtgo.CustomClaims, error) {
	j := jwtgo.NewJWT()
	claims, err := j.ParserToken(tokenString)
	return claims, err
}
