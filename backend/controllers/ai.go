package controllers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"nexus/config"
	"nexus/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AIController 处理所有 AI 相关的代理请求
type AIController struct{}

// aiBackendURL 返回 AI Backend 的完整地址
func aiBackendURL() string {
	return config.AIBackendAddr + ":" + config.AIBackendPort
}

// ==================== SSE 流式透传 ====================

// streamProxy 核心方法：将请求体转发到 AI Backend 并 SSE 流式透传响应
func streamProxy(c *gin.Context, path string, body map[string]interface{}) {
	// 透传前端原始 Authorization header，AI Backend 内部调用 Go Backend API 需要
	forwardAuth := c.GetHeader("Authorization")
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 序列化请求体
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		c.SSEvent("error", map[string]string{"error": "请求体序列化失败"})
		c.Abort()
		return
	}

	// 发起请求到 AI Backend
	targetURL := aiBackendURL() + path
	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		c.SSEvent("error", map[string]string{"error": "创建请求失败"})
		c.Abort()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if forwardAuth != "" {
		req.Header.Set("Authorization", forwardAuth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[AI Proxy] 请求 AI Backend 失败: %v", err)
		fmt.Fprintf(c.Writer, "data: {\"error\": \"AI 服务连接失败，请稍后重试\"}\n\n")
		c.Writer.(http.Flusher).Flush()
		return
	}
	defer resp.Body.Close()

	// 检查 AI Backend 响应状态
	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		log.Printf("[AI Proxy] AI Backend 返回错误: status=%d, body=%s", resp.StatusCode, string(errBody))
		fmt.Fprintf(c.Writer, "data: {\"error\": \"AI 服务返回错误 (%d)\"}\n\n", resp.StatusCode)
		c.Writer.(http.Flusher).Flush()
		return
	}

	// 逐行读取 SSE 并透传到前端
	scanner := bufio.NewScanner(resp.Body)
	// 增大 scanner 缓冲区，避免长行被截断
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(c.Writer, "%s\n", line)
		c.Writer.(http.Flusher).Flush()
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[AI Proxy] 读取 SSE 流中断: %v", err)
	}
}

// ==================== 端点实现 ====================

// Chat SSE 流式对话代理
func (AIController) Chat(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	// 注入 userID，让 AI Backend 知道是哪个用户
	body["user_id"] = userID

	streamProxy(c, "/chat", body)
}

// AnalyzeCode SSE 流式代码分析代理
func (AIController) AnalyzeCode(c *gin.Context) {
	_, err := ParserToken(c)
	if err != nil {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	streamProxy(c, "/analyze-code", body)
}

// GenerateTests 非流式 JSON 代理（测试用例生成）
func (AIController) GenerateTests(c *gin.Context) {
	_, err := ParserToken(c)
	if err != nil {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	bodyBytes, _ := json.Marshal(body)
	targetURL := aiBackendURL() + "/generate-tests"

	req, err := http.NewRequest("POST", targetURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if auth := c.GetHeader("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "AI 服务连接失败"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// 透传 AI Backend 的响应状态码和内容
	c.Data(resp.StatusCode, "application/json", respBody)
}

// Guidance SSE 流式个性化指导代理
func (AIController) Guidance(c *gin.Context) {
	userID, err := ParserToken(c)
	if err != nil {
		return
	}

	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求体"})
		return
	}

	// 从 Go Backend 获取用户画像数据，附加到请求中
	profile, profileErr := services.EnsureProfile(userID)
	if profileErr == nil && profile != nil {
		body["user_profile"] = map[string]interface{}{
			"ability":       profile.Ability,
			"preferences":   profile.Preferences,
		}
	}

	// 注入 userID
	body["user_id"] = userID

	streamProxy(c, "/personalized-guidance", body)
}
