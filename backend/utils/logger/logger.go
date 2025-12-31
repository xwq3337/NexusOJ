package logger

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志级别类型
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// 日志配置结构
type Config struct {
	Level      LogLevel `json:"level" yaml:"level"`           // 日志级别
	Filename   string   `json:"filename" yaml:"filename"`     // 日志文件路径
	MaxSize    int      `json:"maxSize" yaml:"maxSize"`       // 文件最大大小(MB)
	MaxBackups int      `json:"maxBackups" yaml:"maxBackups"` // 最大备份数量
	MaxAge     int      `json:"maxAge" yaml:"maxAge"`         // 最大保存天数
	Compress   bool     `json:"compress" yaml:"compress"`     // 是否压缩备份
	Console    bool     `json:"console" yaml:"console"`       // 是否输出到控制台
}

// Logger 封装结构
type Logger struct {
	*zap.SugaredLogger
	atom   zap.AtomicLevel
	config Config
	mu     sync.Mutex
}

var globalLogger *Logger

// 初始化全局日志
func InitGlobalLogger(cfg Config) error {
	logger, err := NewLogger(cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// 创建新的日志实例
func NewLogger(cfg Config) (*Logger, error) {
	// 设置日志级别
	atom := zap.NewAtomicLevel()
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	atom.SetLevel(level)

	// 创建编码器配置
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建核心列表
	cores := []zapcore.Core{}

	// 文件输出核心
	if cfg.Filename != "" {
		// 确保目录存在
		if err := os.MkdirAll(filepath.Dir(cfg.Filename), 0755); err != nil {
			return nil, err
		}

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			fileWriter,
			atom,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出核心
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.Lock(os.Stdout),
			atom,
		)
		cores = append(cores, consoleCore)
	}

	// 组合核心
	core := zapcore.NewTee(cores...)

	// 创建Logger
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugarLogger := zapLogger.Sugar()

	return &Logger{
		SugaredLogger: sugarLogger,
		atom:          atom,
		config:        cfg,
	}, nil
}

// 解析日志级别
func parseLevel(level LogLevel) (zapcore.Level, error) {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel, nil
	case InfoLevel:
		return zapcore.InfoLevel, nil
	case WarnLevel:
		return zapcore.WarnLevel, nil
	case ErrorLevel:
		return zapcore.ErrorLevel, nil
	case FatalLevel:
		return zapcore.FatalLevel, nil
	default:
		return zapcore.InfoLevel, nil
	}
}

// 动态修改日志级别
func (l *Logger) SetLevel(level LogLevel) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	zapLevel, err := parseLevel(level)
	if err != nil {
		return err
	}

	l.atom.SetLevel(zapLevel)
	l.config.Level = level
	return nil
}

// 获取当前配置
func (l *Logger) GetConfig() Config {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.config
}

// 安全关闭日志
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

// ========== 全局快捷方法 ==========
func Debug(args ...interface{}) {
	globalLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	globalLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	globalLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	globalLogger.Infof(template, args...)
}
func Infow(msg string, keysAndValues ...interface{}) {
	globalLogger.Infow(msg, keysAndValues...)
}
func Warn(args ...interface{}) {
	globalLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	globalLogger.Warnf(template, args...)
}
func Warnw(msg string, keysAndValues ...interface{}) {
	globalLogger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	globalLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	globalLogger.Errorf(template, args...)
}
func Errorw(msg string, keysAndValues ...interface{}) {
	globalLogger.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	globalLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	globalLogger.Fatalf(template, args...)
}

func With(fields ...interface{}) *Logger {
	return &Logger{
		SugaredLogger: globalLogger.With(fields...),
		atom:          globalLogger.atom,
		config:        globalLogger.config,
	}
}

func SetGlobalLevel(level LogLevel) error {
	return globalLogger.SetLevel(level)
}

func Sync() error {
	return globalLogger.Sync()
}
