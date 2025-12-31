package ipinterceptor

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Interceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		isWebSocket := c.Request.Header.Get("Upgrade") == "websocket" &&
			strings.Contains(c.Request.Header.Get("Connection"), "Upgrade")
		ip := c.ClientIP()
		if isWebSocket {
			ip = c.Request.Header.Get("Origin")
		}
		for _, allow := range getAllowOrigins() {
			if allow == "*" {
				c.Next()
				return
			}
			if strings.Contains(ip, allow) {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{
			"msg": "该源无法访问",
		})
		c.Abort()
	}
}

func getAllowOrigins() []string {
	var mode string
	if gin.Mode() == gin.ReleaseMode {
		mode = "release"
	} else {
		mode = "debug"
	}
	var allowOrigins []string
	if mode == "release" {
		allowOrigins = []string{"47.109.57.7"}
	} else {
		allowOrigins = []string{"*"}
	}
	return allowOrigins
}
