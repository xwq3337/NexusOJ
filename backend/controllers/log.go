package controllers

import (
	"fmt"
	"net/http"
	"nexus/config"
	"nexus/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LogController struct{}

// LogDateResponse 日志日期响应
type LogDateResponse struct {
	Dates []string `json:"dates"` // 日期列表，格式：2006-01-02
}

// LogListResponse 日志列表响应
type LogListResponse struct {
	Content  string `json:"content"`   // 日志内容
	Total    int64  `json:"total"`     // 总字节数
	Page     int    `json:"page"`      // 当前页码
	PageSize int    `json:"pageSize"`  // 每页大小
	HasMore  bool   `json:"hasMore"`   // 是否有更多数据
}

// GetDate 获取有记录的日志的日期
func (LogController) GetDate(c *gin.Context) {
	// 扫描日志目录
	entries, err := os.ReadDir(config.LogDir)
	if err != nil {
		utils.ReturnError(c, http.StatusInternalServerError, err)
		return
	}
	// 用于存储找到的日期
	dateMap := make(map[string]bool)
	// 日期格式: app_2006_1_1.log 或 app_2026_3_12.log
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		// 检查文件名格式
		if !strings.HasPrefix(filename, "app_") || !strings.HasSuffix(filename, ".log") {
			continue
		}
		// 解析日期: app_2026_3_12.log -> 2026_3_12
		dateStr := strings.TrimPrefix(filename, "app_")
		dateStr = strings.TrimSuffix(dateStr, ".log")

		// 分割日期字符串: 2026_3_12 -> [2026, 3, 12]
		parts := strings.Split(dateStr, "_")
		if len(parts) != 3 {
			fmt.Println("invalid date format")
			continue
		}

		// 构建标准日期格式: 2026-03-12
		year, month, day := parts[0], parts[1], parts[2]

		// 补零格式化月份
		if len(month) == 1 {
			month = "0" + month
		}
		// 补零格式化日期
		if len(day) == 1 {
			day = "0" + day
		}

		standardDate := fmt.Sprintf("%s-%s-%s", year, month, day)
		dateMap[standardDate] = true
	}
	// 转换为切片并排序（降序，最新的日期在前）
	dates := make([]string, 0, len(dateMap))
	for date := range dateMap {
		dates = append(dates, date)
	}

	// 按日期降序排序
	sort.Slice(dates, func(i, j int) bool {
		return dates[i] > dates[j]
	})

	response := LogDateResponse{
		Dates: dates,
	}

	utils.ReturnSuccess(c, http.StatusOK, "success", response)
}

// GetList 获取某一天的日志内容
func (LogController) GetList(c *gin.Context) {
	// 获取日期参数，格式: 2006-01-02
	date := c.Query("date")
	if date == "" {
		utils.ReturnError(c, http.StatusBadRequest, fmt.Errorf("date parameter is required"))
		return
	}

	// 验证日期格式
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		utils.ReturnError(c, http.StatusBadRequest, fmt.Errorf("invalid date format, expected: 2006-01-02"))
		return
	}

	// 构建文件名: app_2006_1_1.log
	filename := fmt.Sprintf("app_%d_%d_%d.log", parsedDate.Year(), parsedDate.Month(), parsedDate.Day())
	filepath := filepath.Join(config.LogDir, filename)

	// 读取文件内容
	content, err := os.ReadFile(filepath)
	if err != nil {
		utils.ReturnError(c, http.StatusNotFound, fmt.Errorf("log file not found for date: %s", date))
		return
	}

	// 返回日志内容（纯文本）
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, string(content))
}
