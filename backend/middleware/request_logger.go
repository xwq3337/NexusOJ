package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"pkg/utils/logger"

	"github.com/gin-gonic/gin"
)

// RequestLoggerConfig 日志中间件配置
type RequestLoggerConfig struct {
	SkipPaths        []string // 跳过日志的路径
	MaxBodyLogSize   int      // 记录请求体的最大大小（字节）
	LogRequestHeader bool     // 是否记录请求头
	LogResponseBody  bool     // 是否记录响应体
	LogRequestBody   bool     // 是否记录请求体
}

// DefaultRequestLoggerConfig 默认配置
func DefaultRequestLoggerConfig() RequestLoggerConfig {
	return RequestLoggerConfig{
		SkipPaths:        []string{"/health", "/metrics", "/favicon.ico"},
		MaxBodyLogSize:   1024, // 1KB
		LogRequestHeader: true,
		LogResponseBody:  false, // 默认不记录响应体
		LogRequestBody:   true,
	}
}

// RequestLogger 请求日志中间件
func RequestLogger(config RequestLoggerConfig) gin.HandlerFunc {
	skip := make(map[string]bool, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skip[path] = true
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 跳过指定路径的日志记录
		if skip[path] {
			c.Next()
			return
		}

		// 记录请求信息
		var requestBody []byte
		if config.LogRequestBody && c.Request.Body != nil {
			// 读取请求体
			body, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 限制日志记录的大小
				if len(body) > config.MaxBodyLogSize {
					requestBody = body[:config.MaxBodyLogSize]
				} else {
					requestBody = body
				}

				// 恢复请求体以供后续处理
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			} else {
				logger.Warnf("Failed to read request body: %v", err)
			}
		}

		// 记录响应体
		responseWriter := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		// 处理请求
		c.Next()

		// 计算处理时间
		latency := time.Since(start)

		// 收集日志字段
		fields := []interface{}{
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"ip", c.ClientIP(),
			"latency", latency.String(),
			"user_agent", c.Request.UserAgent(),
			"time", start.Format(time.RFC3339),
		}

		if raw != "" {
			fields = append(fields, "query", raw)
		}

		if len(requestBody) > 0 {
			// 尝试将JSON格式化，否则直接记录为字符串
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, requestBody, "", "  "); err == nil {
				fields = append(fields, "request_body", prettyJSON.String())
			} else {
				fields = append(fields, "request_body", string(requestBody))
			}
		}

		if config.LogRequestHeader {
			headers := make(map[string]string)
			for k, v := range c.Request.Header {
				headers[k] = v[0] // 只记录第一个值
			}
			fields = append(fields, "headers", headers)
		}

		if config.LogResponseBody && responseWriter.body.Len() > 0 {
			responseBody := responseWriter.body.Bytes()
			// 限制日志记录的大小
			if len(responseBody) > config.MaxBodyLogSize {
				responseBody = responseBody[:config.MaxBodyLogSize]
			}

			// 尝试将JSON格式化，否则直接记录为字符串
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, responseBody, "", "  "); err == nil {
				fields = append(fields, "response_body", prettyJSON.String())
			} else {
				fields = append(fields, "response_body", string(responseBody))
			}
		}

		// 获取错误信息（如果有）
		if len(c.Errors) > 0 {
			errMsgs := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errMsgs[i] = err.Error()
			}
			fields = append(fields, "errors", errMsgs)
		}

		// 根据状态码选择日志级别
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Errorw("Server error", fields...)
		case status >= 400:
			logger.Warnw("Client error", fields...)
		default:
			logger.Infow("Request processed", fields...)
		}
	}
}

// responseBodyWriter 用于捕获响应体
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
