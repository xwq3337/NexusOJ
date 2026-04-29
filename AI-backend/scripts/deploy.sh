#!/bin/bash
set -e

ADDRESS="root@47.109.57.7"
DIRECTORY="/opt/oj-agent"

echo "上传文件..."
scp -r app requirements.txt .env "$ADDRESS:$DIRECTORY/"

# echo "安装依赖..."
# ssh "$ADDRESS" "cd $DIRECTORY && pip3 install -r requirements.txt --break-system-packages --ignore-installed PyJWT"

echo "重启服务..."
ssh "$ADDRESS" "pkill -f 'uvicorn' || true"
ssh "$ADDRESS" "cd $DIRECTORY && nohup python3 -m uvicorn app.main:app --host 0.0.0.0 --port 5557 --workers 1 --timeout-keep-alive 120 > $DIRECTORY/app.log 2>&1 &"
# nohup uvicorn app.main:app --host 0.0.0.0 --port 5557 --workers 1 --timeout-keep-alive 120 > /opt/oj-agent/app.log 2>&1 &
echo "部署完成!"
