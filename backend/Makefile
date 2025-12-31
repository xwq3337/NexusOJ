SOURCE_NAME := main.go
APP_NAME := main
BUILD_DIR := target
.PHONY: build dev deploy clear

build:
	@go build -o $(BUILD_DIR)/$(APP_NAME)
	@echo "构建成功"

dev:build
	@echo "运行air在.air.toml配置下"
	@echo "If No such file, please replace path to your air path(run <alias air="$$HOME/go/bin/air">)"
	@echo "source ~/.zshrc or ~/.bashrc"
	@air

deploy: 
	@chmod +x scripts/deploy.sh
	@./scripts/deploy.sh

clear:
	@rm -rf ./tmp/*
	@rm -f main

test:
	@ab -n 1000 -c 100 http://localhost:8080/api/ip
	