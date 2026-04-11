package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	jwtgo "nexus/middleware/jwt"
	"nexus/models"
	"nexus/utils"
	"os"
	"path/filepath"
	"strconv"
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
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "无效的用户ID")
		return
	}
	user, err := models.User{}.QueryUserById(userID)
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
	user.ID = uint64(idgen.NextId())
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "密码加密失败")
		return
	}
	user.Password = hash
	if err := models.CreateUser(user); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", user)
}

func (UserController) UserLogin(c *gin.Context) {
	// 解析请求参数
	var params struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, 400, err.Error())
		return
	}
	user, err := models.User{Username: params.Username}.QueryUser()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", params.Username))
			return
		} else {
			utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
	if user.ID == 0 {
		utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", params.Username))
		return
	}
	// Argon2 验证密码
	match, _ := utils.VerifyPassword(params.Password, user.Password)

	if !match {
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
	err = models.UpdateUserInfo(user)
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

	// 获取原始文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".png" // 默认使用 png
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "打开文件失败, err: "+err.Error())
		return
	}
	defer src.Close()

	// 解码图片
	img, format, err := image.Decode(src)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "解码图片失败, err: "+err.Error())
		return
	}
	filename := fmt.Sprintf("/%d%s", userID, ext)
	// 创建输出文件路径
	filePath := filepath.Join(DirPath, filename)
	// 检查文件大小，如果超过 500KB 则压缩
	const maxSize = 500 * 1024 // 500KB

	if file.Size > maxSize {
		// 压缩图片
		var buf bytes.Buffer
		var compressErr error

		// 根据格式进行压缩编码
		switch format {
		case "jpeg":
			compressErr = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
		case "png":
			compressErr = png.Encode(&buf, img)
		default:
			// 默认使用 jpeg 格式压缩
			compressErr = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85})
			ext = ".jpg"
			filename = fmt.Sprintf("/%d%s", userID, ext)
			filePath = filepath.Join(DirPath, filename)
		}

		if compressErr != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "压缩图片失败, err: "+compressErr.Error())
			return
		}
		// 如果压缩后仍然超过 500KB，进一步降低质量
		if buf.Len() > maxSize {
			buf.Reset()
			compressErr = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 70})
			if compressErr != nil {
				utils.ReturnError(c, http.StatusInternalServerError, "压缩图片失败, err: "+compressErr.Error())
				return
			}
			ext = ".jpg"
			filename = fmt.Sprintf("/%d%s", userID, ext)
			filePath = filepath.Join(DirPath, filename)
		}

		// 保存压缩后的图片
		if err = os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "保存图片失败, err: "+err.Error())
			return
		}
	} else {
		// 文件大小符合要求，直接保存
		if err = c.SaveUploadedFile(file, filePath); err != nil {
			utils.ReturnError(c, http.StatusInternalServerError, "上传头像失败, err: "+err.Error())
			return
		}
	}

	url, err := models.UpdateAvatar(userID, filename)
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

// TODO需要密码?
func (UserController) RefreshToken(c *gin.Context) {
	_id, err := ParserToken(c)
	if err != nil {
		utils.ReturnError(c, http.StatusUnauthorized, "未授权"+err.Error())
		return
	}

	user, err := models.User{}.QueryUserById(_id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ReturnError(c, http.StatusNotFound, fmt.Sprintf("未找到名为 %s 的用户或者密码错误", "ID"))
			return
		} else {
			utils.ReturnError(c, http.StatusInternalServerError, fmt.Sprintf("查询出错 %v", err))
			return
		}
	}
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

func ParserToken(c *gin.Context) (uint64, error) {
	tokenString := c.Request.Header.Get("Authorization")
	return ParserTokenByString(tokenString)
}
func ParserTokenByString(tokenString string) (uint64, error) {
	if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
		return 0, errors.New("请求未携带token或token不完整,无权限访问")
	}
	tokenString = strings.Split(tokenString, " ")[1]
	j := jwtgo.NewJWT()
	claims, err := j.ParserToken(tokenString)
	if err != nil {
		return 0, err
	}
	if claims != nil && claims.UserID == 0 {
		return 0, errors.New("token无效")
	}
	return claims.UserID, err
}

func (UserController) GetTopUsers(c *gin.Context) {
	users, err := models.GetTopUsersByRating(10)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", users)
}

func (UserController) ValidateToken(c *gin.Context) {
	UserID, err := ParserToken(c)
	if err != nil || UserID == 0 {
		utils.ReturnError(c, http.StatusUnauthorized, "token无效")
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", "token有效")
}

func (UserController) UpdateRole(c *gin.Context) {
	var params struct {
		Id   uint64 `json:"id"`
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ReturnError(c, http.StatusBadRequest, "请求参数错误"+err.Error())
		return
	}
	if params.Role != "admin" && params.Role != "user" {
		utils.ReturnError(c, http.StatusBadRequest, "角色不合法")
		return
	}
	if err := models.UpdateUserRole(params.Id, params.Role); err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, "更新角色失败"+err.Error())
		return
	}
	utils.ReturnSuccess(c, http.StatusOK, "success", nil)
}
