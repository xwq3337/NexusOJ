/*
 * @Author: x_wq3337 854541540@qq.com
 * @Date: 2025-07-10 15:34:29
 * @LastEditors: x_wq3337 854541540@qq.com
 * @LastEditTime: 2026-01-24 03:09:03
 * @FilePath: /backend/middleware/ipinterceptor/index.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
