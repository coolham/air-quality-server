#!/bin/bash

# Docker快速升级脚本（使用国内镜像源）
# 适用于网络较慢的环境

echo "🚀 Docker快速升级脚本（国内镜像源）"
echo "=================================="

# 检查当前Docker版本
echo "📋 当前Docker版本："
docker --version

# 确认升级
echo ""
read -p "是否继续升级Docker? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 升级已取消"
    exit 1
fi

echo ""
echo "🔧 开始快速升级Docker..."

# 1. 停止Docker服务
echo "⏹️  停止Docker服务..."
sudo systemctl stop docker

# 2. 卸载旧版本
echo "🗑️  卸载旧版本..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc docker-compose

# 3. 使用阿里云镜像源安装Docker
echo "📦 使用阿里云镜像源安装Docker..."

# 添加阿里云Docker仓库
curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 更新包索引
sudo apt-get update

# 安装Docker
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 4. 配置镜像加速器
echo "🌐 配置镜像加速器..."
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json > /dev/null <<EOF
{
  "registry-mirrors": [
    "https://docker.mirrors.ustc.edu.cn",
    "https://hub-mirror.c.163.com",
    "https://mirror.baidubce.com",
    "https://ccr.ccs.tencentyun.com"
  ],
  "insecure-registries": [],
  "debug": false,
  "experimental": false,
  "features": {
    "buildkit": true
  }
}
EOF

# 5. 启动Docker
echo "▶️  启动Docker服务..."
sudo systemctl start docker
sudo systemctl enable docker

# 6. 添加用户到docker组
echo "👤 添加用户到docker组..."
sudo usermod -aG docker $USER

# 7. 验证安装
echo ""
echo "✅ 验证安装..."
sleep 3
docker --version
docker compose version

echo ""
echo "🎉 Docker快速升级完成！"
echo "📝 请重新登录或运行 'newgrp docker' 以应用用户组更改"
