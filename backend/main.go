package main

import (
	"nexus/config"
	"nexus/migrations"
	"nexus/router"
	"nexus/services"
	"nexus/utils/logger"
	"os"
	"path/filepath"

	"github.com/yitter/idgenerator-go/idgen"
)

var IdOptions = idgen.NewIdGeneratorOptions(1)

var LogOptions = logger.Config{
	Level:       logger.DebugLevel,
	Filename:    filepath.Join(config.LogDir, "app.log"),
	MaxSize:     100,  // MB
	MaxBackups:  7,    // 保留7天
	MaxAge:      30,   // 保留30天
	Compress:    true, // 是否压缩
	Console:     true, // 是否输出到控制台
	DailyRotate: true, // 启用按天分割日志
}

// 主函数，程序的入口点
func main() {
	// 初始化全局日志记录器，如果失败则panic
	if err := logger.InitGlobalLogger(LogOptions); err != nil {
		panic(err)
	}
	// 确保在程序退出前同步日志
	defer logger.Sync()
	// 初始化ID生成器
	idgen.SetIdGenerator(IdOptions)
	// 执行数据库迁移
	migrations.Migrate()

	// 初始化判题队列
	// 参数说明: workerNum=5(并发判题worker数量), queueSize=100(队列容量)
	services.InitJudgeQueue(5, 100)

	// 初始化画像服务
	// 参数说明: workerNum=3(画像更新worker数量), queueSize=200(队列容量)
	services.InitProfileService(3, 200)

	// 初始化比赛状态自动检查协程
	services.InitContestStatusWorker()

	// 获取路由器并启动HTTP服务器
	r := router.Router()
	// 延迟执行函数，用于捕获运行时异常
	defer func() {
		if err := recover(); err != nil {
			logger.Error("捕获异常")
		}
	}()
	// 启动服务器，监听配置的端口
	if err := r.Run(string(":" + config.Port)); err != nil {
		logger.Error("发生异常")
		os.Exit(-1)
		return
	}
}
