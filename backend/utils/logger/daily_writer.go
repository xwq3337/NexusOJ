package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// DailyWriter 按天分割的日志写入器
type DailyWriter struct {
	baseDir     string     // 日志文件目录
	baseName    string     // 基础文件名（不含扩展名）
	ext         string     // 文件扩展名
	maxBackups  int        // 保留天数
	currentDate string     // 当前日期（格式: 2006-01-02）
	file        *os.File   // 当前文件句柄
	mu          sync.Mutex // 保护并发写入
}

// NewDailyWriter 创建新的按天分割的日志写入器
// baseDir: 日志文件目录
// baseName: 基础文件名，如 "app"
// ext: 文件扩展名，如 ".log"
// maxBackups: 保留天数，0表示不删除旧文件
func NewDailyWriter(baseDir, baseName, ext string, maxBackups int) (*DailyWriter, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	dw := &DailyWriter{
		baseDir:    baseDir,
		baseName:   baseName,
		ext:        ext,
		maxBackups: maxBackups,
	}

	// 初始化当前文件
	if err := dw.rotateIfNeeded(); err != nil {
		return nil, err
	}

	// 启动清理过期文件的 goroutine
	go dw.cleanupOldLogs()

	return dw, nil
}

// Write 实现 io.Writer 接口
func (dw *DailyWriter) Write(p []byte) (n int, err error) {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	// 检查是否需要切换文件
	if err := dw.rotateIfNeeded(); err != nil {
		return 0, err
	}

	if dw.file == nil {
		return 0, fmt.Errorf("log file is not open")
	}

	return dw.file.Write(p)
}

// Sync 实现 io.Syncer 接口
func (dw *DailyWriter) Sync() error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	if dw.file == nil {
		return nil
	}

	return dw.file.Sync()
}

// Close 关闭当前日志文件
func (dw *DailyWriter) Close() error {
	dw.mu.Lock()
	defer dw.mu.Unlock()

	if dw.file != nil {
		if err := dw.file.Close(); err != nil {
			return err
		}
		dw.file = nil
	}

	return nil
}

// rotateIfNeeded 检查并执行日志文件切换
func (dw *DailyWriter) rotateIfNeeded() error {
	currentDate := time.Now().Format("2006-01-02")

	// 如果日期未变化，不需要切换
	if dw.currentDate == currentDate && dw.file != nil {
		return nil
	}

	// 关闭旧文件
	if dw.file != nil {
		if err := dw.file.Close(); err != nil {
			return fmt.Errorf("failed to close old log file: %w", err)
		}
	}

	// 创建新文件
	filename := dw.getFilename(currentDate)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file %s: %w", filename, err)
	}

	dw.currentDate = currentDate
	dw.file = file

	return nil
}

// getFilename 生成日志文件名
// 格式: app_2006_01_02.log
func (dw *DailyWriter) getFilename(date string) string {
	// 将 2006-01-02 格式转换为 2006_1_1 格式
	t, _ := time.Parse("2006-01-02", date)
	return filepath.Join(dw.baseDir, fmt.Sprintf("%s_%d_%d_%d%s", dw.baseName, t.Year(), t.Month(), t.Day(), dw.ext))
}

// cleanupOldLogs 清理过期的日志文件
func (dw *DailyWriter) cleanupOldLogs() {
	if dw.maxBackups <= 0 {
		return
	}

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		dw.mu.Lock()
		cutoffDate := time.Now().AddDate(0, 0, -dw.maxBackups)
		dw.mu.Unlock()

		// 扫描日志目录
		entries, err := os.ReadDir(dw.baseDir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// 解析文件名中的日期
			fileDate := dw.extractDateFromFilename(entry.Name())
			if fileDate.IsZero() {
				continue
			}

			// 删除过期文件
			if fileDate.Before(cutoffDate) {
				filePath := filepath.Join(dw.baseDir, entry.Name())
				os.Remove(filePath)
			}
		}
	}
}

// extractDateFromFilename 从文件名中提取日期
// 文件名格式: app_2006_1_1.log
func (dw *DailyWriter) extractDateFromFilename(filename string) time.Time {
	// 使用 sscanf 风格解析
	var year, month, day int
	_, err := fmt.Sscanf(filename, fmt.Sprintf("%s_%%d_%%d_%%d%s", dw.baseName, dw.ext), &year, &month, &day)
	if err != nil {
		return time.Time{}
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
