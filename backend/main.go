package main

import (
	"os"
	"path/filepath"
	"pkg/config"
	"pkg/migrations"
	"pkg/router"
	"pkg/utils/logger"

	"github.com/yitter/idgenerator-go/idgen"
)

var options = idgen.NewIdGeneratorOptions(1)

func main() {
	if err := logger.InitGlobalLogger(loggerConfig()); err != nil {
		panic(err)
	}
	defer logger.Sync()
	idgen.SetIdGenerator(options)
	migrations.Migrate()
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

func loggerConfig() logger.Config {
	cfg := logger.Config{
		Level:      logger.DebugLevel,
		Filename:   filepath.Join(config.LogDir, "app.log"),
		MaxSize:    100, // MB
		MaxBackups: 7,   // 保留7天
		MaxAge:     30,  // 保留30天
		Compress:   true,
		Console:    true,
	}
	return cfg
}
