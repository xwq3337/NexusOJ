#!/bin/bash
git clone https://gitee.com/x_wq3337/judge.git
cp -a judge/* .
rm -rf ./judge
# Ubuntu
if ! command -v docker &> /dev/null
then
    echo "Docker could not be found, installing..."
    sudo apt-get update
    sudo apt-get install -y docker.io
fi

sudo systemctl start docker
sudo systemctl enable docker

IMAGE_NAME="sandbox_universal"
echo "Building Docker image with name: $IMAGE_NAME"
sudo docker build -t $IMAGE_NAME -f ./Judger/Dockerfile.universal .

if [ -f "requirements.txt" ]; then
   pip install -r requirements.txt --break-system-packages
fi

echo "Setup completed."