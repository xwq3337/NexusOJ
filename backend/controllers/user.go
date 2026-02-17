package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"pkg/dao"
	jwtgo "pkg/middleware/jwt"
	"pkg/models"
	"pkg/utils"
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
		utils.ReturnSuccess(c, http.StatusOK, "suceess", user)
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
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, err := models.User{Username: data["username"], Password: data["password"]}.QueryUser()
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
func (UserController) UpdateUser(c *gin.Context) {
	user := &models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	u, _ := ParserToken(c.Request.Header.Get("Authorization"))
	user.ID = u.UserID
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
	claims, _ := ParserToken(c.Request.Header.Get("Authorization"))
	file, err := c.FormFile("avatar")
	DirPath := config.AvatarDir
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "上传头像失败, err: "+err.Error())
		return
	}

	if err = c.SaveUploadedFile(file, filepath.Join(DirPath, fmt.Sprintf("/%s.png", claims.UserID))); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "上传头像失败, err: "+err.Error())
		return
	}

	url, err := models.UpdateAvatar(claims.UserID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "更新头像失败, err: "+err.Error())
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", url)
}

// UpdatePassword 更新密码
func (UserController) UpdatePassword(c *gin.Context) {
	claims, err := ParserToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	data := make(map[string]string)
	if err := c.BindJSON(&data); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}

	oldPassword := data["old_password"]
	newPassword := data["new_password"]

	if oldPassword == "" || newPassword == "" {
		utils.ReturnError(c, http.StatusBadRequest, "旧密码和新密码不能为空")
		return
	}

	err = models.User{}.UpdatePassword(claims.UserID, oldPassword, newPassword)
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
	data := make(map[string]string)
	_ = c.BindJSON(&data)
	user, _ := models.User{Username: data["username"], Password: data["password"]}.QueryUser()
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

func ParserToken(tokenString string) (*jwtgo.CustomClaims, error) {
	j := jwtgo.NewJWT()
	claims, err := j.ParserToken(tokenString)
	return claims, err
}

// GetAllFriends 获取所有好友
func (UserController) GetAllFriends(c *gin.Context) {
	claims, err := ParserToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	friends, err := models.GetFriendList(claims.UserID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "获取好友列表失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", friends)
}

// FirendRequest 发送好友请求
func (UserController) FirendRequest(c *gin.Context) {
	claims, err := ParserToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	data := make(map[string]string)
	if err := c.BindJSON(&data); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}

	friendID := data["friend_id"]
	if friendID == "" {
		utils.ReturnError(c, http.StatusBadRequest, "好友ID不能为空")
		return
	}

	if friendID == claims.UserID {
		utils.ReturnError(c, http.StatusBadRequest, "不能添加自己为好友")
		return
	}

	// 检查是否已经是好友或已有请求
	var existingFriendship models.FriendShips
	err = dao.MysqlClient.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		claims.UserID, friendID, friendID, claims.UserID).First(&existingFriendship).Error
	if err == nil {
		if existingFriendship.Status == 1 {
			utils.ReturnError(c, http.StatusBadRequest, "已经是好友关系")
			return
		} else {
			utils.ReturnError(c, http.StatusBadRequest, "已存在待处理的好友请求")
			return
		}
	}

	friendship := &models.FriendShips{
		UserID:   claims.UserID,
		FriendID: friendID,
		Status:   0,
	}

	err = models.CreateFriendRequest(friendship)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "发送好友请求失败")
		return
	}

	utils.ReturnSuccess(c, http.StatusOK, "好友请求发送成功", nil)
}

// HandleFriendRequest 处理好友请求
func (UserController) HandleFriendRequest(c *gin.Context) {
	claims, err := ParserToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	data := make(map[string]interface{})
	if err := c.BindJSON(&data); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}

	friendshipIDFloat, ok := data["friendship_id"].(float64)
	if !ok {
		utils.ReturnError(c, http.StatusBadRequest, "friendship_id 参数错误")
		return
	}
	friendshipID := uint(friendshipIDFloat)

	accept, ok := data["accept"].(bool)
	if !ok {
		utils.ReturnError(c, http.StatusBadRequest, "accept 参数错误")
		return
	}

	// 验证该好友请求是否是发送给当前用户的
	var friendship models.FriendShips
	err = dao.MysqlClient.First(&friendship, friendshipID).Error
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, "好友请求不存在")
		return
	}

	if friendship.FriendID != claims.UserID {
		utils.ReturnError(c, http.StatusForbidden, "无权处理此好友请求")
		return
	}

	err = models.HandleFriendRequest(friendshipID, accept)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "处理好友请求失败")
		return
	}

	action := "拒绝"
	if accept {
		action = "接受"
	}
	utils.ReturnSuccess(c, http.StatusOK, "已"+action+"好友请求", nil)
}

// GetFriendRequestList 获取好友请求列表
func (UserController) GetFriendRequestList(c *gin.Context) {
	claims, err := ParserToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权")
		return
	}

	requests, err := models.GetFriendRequestList(claims.UserID)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "获取好友请求列表失败")
		return
	}

	// 查询每个请求的发送者信息
	type FriendRequestWithUser struct {
		models.FriendShips
		Username string  `json:"username"`
		Nickname *string `json:"nickname"`
		Avatar   *string `json:"avatar"`
		Rating   int16   `json:"rating"`
	}

	var result []FriendRequestWithUser
	for _, req := range requests {
		user, err := models.User{}.QueryUserById(req.UserID)
		if err != nil {
			continue
		}

		result = append(result, FriendRequestWithUser{
			FriendShips: req,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Rating:      user.Rating,
		})
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", result)
}
