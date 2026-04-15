#!/bin/bash
set -e

ADDRESS="root@47.109.57.7"
DIRECTORY="/opt/oj-agent"

echo "上传文件..."
scp -r ai_backend chat manage.py requirements.txt .env system_prompt.txt "$ADDRESS:$DIRECTORY/"

echo "安装依赖..."
ssh "$ADDRESS" "cd $DIRECTORY && pip3 install -r requirements.txt --break-system-packages"

echo "重启服务..."
ssh "$ADDRESS" "pkill -f 'gunicorn' || true"
ssh "$ADDRESS" "cd $DIRECTORY && nohup gunicorn ai_backend.wsgi:application --bind 0.0.0.0:5557 --worker-class gevent --workers 1 --timeout 120 > $DIRECTORY/app.log 2>&1 &"

echo "部署完成!"
