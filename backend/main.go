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

func main() {
	if err := logger.InitGlobalLogger(LogOptions); err != nil {
		panic(err)
	}
	defer logger.Sync()
	idgen.SetIdGenerator(IdOptions)
	migrations.Migrate()

	// 初始化判题队列
	// 参数说明: workerNum=5(并发判题worker数量), queueSize=100(队列容量)
	services.InitJudgeQueue(5, 100)

	// 初始化画像服务
	// 参数说明: workerNum=3(画像更新worker数量), queueSize=200(队列容量)
	services.InitProfileService(3, 200)

	// 初始化比赛状态自动检查协程
	services.InitContestStatusWorker()

	r := router.Router()
	defer func() {
		if err := recover(); err != nil {
			logger.Error("捕获异常")
		}
	}()
	if err := r.Run(string(":" + config.Port)); err != nil {
		logger.Error("发生异常")
		os.Exit(-1)
		return
	}
}
