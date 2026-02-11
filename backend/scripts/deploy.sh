#!/bin/bash
# 炫酷部署脚本配置（请根据实际情况修改）
# 颜色常量定义
COLOR_RESET='\033[0m'
COLOR_RED='\033[31m'
COLOR_GREEN='\033[32m'
COLOR_YELLOW='\033[33m'
COLOR_BLUE='\033[34m'
COLOR_PURPLE='\033[35m'
COLOR_CYAN='\033[36m'
# ASCII艺术标题
echo -e "${COLOR_PURPLE}"
echo "  _____ _           _       _   _       _       "
echo " |  __ (_)         | |     | | | |     | |      "
echo " | |__) |  ___  ___| |_ ___| |_| | ___ | |_     "
echo " |  ___/ |/ _ \/ __| __/ __|  _  |/ _ \| __|    "
echo " | |   | |  __/ (__| || (__| | | | (_) | |_     "
echo " |_|   |_|\___|\___|\__\___\_| |_|\___/ \__|    "
echo -e "${COLOR_RESET}"
echo -e "${COLOR_CYAN}===============================================${COLOR_RESET}"
echo -e "${COLOR_BLUE}          OJ 后端部署脚本 v1.0 (直接上传模式)${COLOR_RESET}"
echo -e "${COLOR_CYAN}===============================================${COLOR_RESET}"


SERVER_IP="47.109.57.7"              # 服务器IP地址
SERVER_USER="root"                   # 服务器登录用户
REMOTE_DIR="/opt/oj-backend"         # 服务器部署目录
SERVICE_NAME="oj-backend.service"    # systemd服务名
BUILD_DIR="./target/release"         # 本地构建输出目录

# 步骤0：清理并创建构建目录
echo -e "${COLOR_BLUE}[步骤0/5] 准备构建目录...${COLOR_RESET}"
rm -rf ${BUILD_DIR} && mkdir -p ${BUILD_DIR} || {
    echo -e "${COLOR_RED}[错误] 创建构建目录失败！${COLOR_RESET}"
    exit 1
}
echo -e "${COLOR_GREEN}√ 构建目录准备完成${COLOR_RESET}"

# 步骤1：构建Go二进制文件
echo -e "${COLOR_BLUE}[步骤1/5] 正在构建Go二进制文件...${COLOR_RESET}"
export GIN_MODE=release
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${BUILD_DIR}/oj-backend || {
    echo -e "${COLOR_RED}[错误] Go二进制文件构建失败!${COLOR_RESET}"
    exit 1
}
echo -e "${COLOR_GREEN}√ 二进制文件构建完成${COLOR_RESET}"

# 步骤2：复制配置文件到构建目录
echo -e "${COLOR_BLUE}[步骤2/5] 复制部署所需文件...${COLOR_RESET}"
cp config.ini ${BUILD_DIR}/ || {
    echo -e "${COLOR_RED}[警告] 配置文件复制失败，继续使用远程现有配置${COLOR_RESET}"
}
echo -e "${COLOR_GREEN}√ 文件准备完成${COLOR_RESET}"

# 步骤3：上传原始文件到服务器
echo -e "${COLOR_BLUE}[步骤3/5] 直接上传原始文件到服务器(${SERVER_IP})...${COLOR_RESET}"
# 确保远程目录存在
ssh ${SERVER_USER}@${SERVER_IP} "sudo mkdir -p ${REMOTE_DIR} && sudo chown ${SERVER_USER}:${SERVER_USER} ${REMOTE_DIR}"

# 使用rsync直接同步文件（保留时间戳，显示进度）
rsync -avz --progress ${BUILD_DIR}/oj-backend ${SERVER_USER}@${SERVER_IP}:${REMOTE_DIR}/ || {
    echo -e "${COLOR_RED}[错误] 二进制文件上传失败！${COLOR_RESET}"
    exit 1
}

# 单独上传配置文件（如果存在）
if [ -f "${BUILD_DIR}/config.ini" ]; then
    rsync -avz --progress ${BUILD_DIR}/config.ini ${SERVER_USER}@${SERVER_IP}:${REMOTE_DIR}/ || {
        echo -e "${COLOR_YELLOW}[警告] 配置文件上传失败，使用远程现有配置${COLOR_RESET}"
    }
fi

echo -e "${COLOR_GREEN}√ 文件上传完成${COLOR_RESET}"

# 步骤4：远程设置权限并重启服务
echo -e "${COLOR_BLUE}[步骤4/5] 正在配置并重启服务...${COLOR_RESET}"
ssh ${SERVER_USER}@${SERVER_IP} "
  # 设置可执行权限
  chmod +x ${REMOTE_DIR}/oj-backend

  # 配置systemd服务
  echo \"[Unit]
Description=OJ Backend Service
After=network.target

[Service]
Type=simple
User=${SERVER_USER}
WorkingDirectory=${REMOTE_DIR}
ExecStart=${REMOTE_DIR}/oj-backend
Restart=always
RestartSec=5
Environment=\"GIN_MODE=release\"

# 资源限制
MemoryMax=512M
CPUQuota=200%

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full

[Install]
WantedBy=multi-user.target\" | sudo tee /etc/systemd/system/${SERVICE_NAME} > /dev/null

  # 重新加载并重启服务
  sudo systemctl daemon-reload
  sudo systemctl restart ${SERVICE_NAME}
  sudo systemctl enable ${SERVICE_NAME}
  
  # 检查服务状态
  echo -e \"\033[34m服务状态检查:\033[0m\"
  sudo systemctl status ${SERVICE_NAME} --no-pager
"

# 步骤5：本地清理
echo -e "${COLOR_BLUE}[步骤5/5] 进行本地清理...${COLOR_RESET}"
rm -rf ${BUILD_DIR}
echo -e "${COLOR_GREEN}√ 清理完成${COLOR_RESET}"

echo -e "${COLOR_CYAN}===============================================${COLOR_RESET}"
echo -e "${COLOR_GREEN}√ 部署完成！服务已成功启动${COLOR_RESET}"
echo -e "${COLOR_BLUE}访问地址: http://${SERVER_IP}:8080${COLOR_RESET}"
echo -e "${COLOR_YELLOW}日志查看: journalctl -u ${SERVICE_NAME} -f${COLOR_RESET}"
echo -e "${COLOR_CYAN}===============================================${COLOR_RESET}"