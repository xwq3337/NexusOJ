package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"nexus/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	SignKey          = "xwq200505123337" // 签名信息应该设置成动态从库中获取
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取token
		if c.Request.Header.Get("Authorization") == "" {
			utils.ReturnError(c, http.StatusUnauthorized, "请求未携带token,无权限访问")
			c.Abort()
			return
		}
		tokens := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			utils.ReturnError(c, http.StatusUnauthorized, "请求token格式有误")
			c.Abort()
			return
		}
		claims, err := NewJWT().ParserToken(tokens[1])
		if err != nil {
			utils.ReturnError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		if claims != nil && claims.UserID == 0 {
			utils.ReturnError(c, http.StatusUnauthorized, "用户不存在")
			c.Abort()
			return
		}
	}
}

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	UserID uint64 `json:"userID"`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

func (j *JWT) ParserToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}

		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, TokenInvalid

}

func (j *JWT) UpdateToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(2 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", fmt.Errorf("token获取失败:%v", err)
}

// GetUserIDFromToken 从token中解析出userID，无论是否过期
// 只要token格式正确且签名有效，就会返回userID
func GetUserIDFromToken(tokenString string) uint64 {
	if tokenString == "" {
		return 0
	}

	j := NewJWT()
	token, _ := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	// 即使解析失败（包括过期），也尝试获取claims
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims.UserID
	}

	return 0
}
