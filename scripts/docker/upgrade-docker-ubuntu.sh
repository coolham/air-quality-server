#!/bin/bash

# Ubuntu 24.04 Docker升级脚本

echo "⬆️ 升级Docker到最新版本..."
echo

# 检查当前版本
echo "📋 当前Docker版本:"
docker --version

echo
echo "📋 当前Docker Compose版本:"
docker-compose --version 2>/dev/null || echo "Docker Compose未安装"

echo
echo "1️⃣ 停止Docker服务..."
sudo systemctl stop docker

echo "2️⃣ 备份当前Docker配置..."
sudo cp -r /etc/docker /etc/docker.backup 2>/dev/null || true

echo "3️⃣ 卸载旧版本Docker..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc docker-compose

echo "4️⃣ 更新包索引..."
sudo apt-get update

echo "5️⃣ 安装必要的包..."
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

echo "6️⃣ 添加Docker官方GPG密钥..."
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

echo "7️⃣ 设置Docker仓库..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "8️⃣ 更新包索引..."
sudo apt-get update

echo "9️⃣ 安装最新版本Docker..."
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "🔟 启动Docker服务..."
sudo systemctl start docker
sudo systemctl enable docker

echo "1️⃣1️⃣ 添加用户到docker组..."
sudo usermod -aG docker $USER

echo "1️⃣2️⃣ 验证安装..."
docker --version
docker compose version

if [ $? -eq 0 ]; then
    echo
    echo "✅ Docker升级成功！"
    echo
    echo "📋 新版本信息:"
    docker --version
    docker compose version
    echo
    echo "🚀 现在可以使用以下命令启动服务:"
    echo "  docker compose up --build -d"
    echo
    echo "⚠️  注意: 请重新登录或运行 'newgrp docker' 以使组权限生效"
    echo
    echo "🔄 或者运行以下命令立即生效:"
    echo "  newgrp docker"
else
    echo
    echo "❌ Docker升级失败"
    echo "💡 请检查错误信息或手动安装"
fi
