#!/bin/bash

# Ubuntu Docker升级脚本
# 从Docker 27.x升级到最新版本

echo "🚀 Ubuntu Docker升级脚本"
echo "========================"

# 检查当前Docker版本
echo "📋 当前Docker版本："
docker --version
docker-compose --version 2>/dev/null || echo "docker-compose: 未安装"

# 检查系统信息
echo ""
echo "🖥️  系统信息："
lsb_release -a 2>/dev/null || cat /etc/os-release

# 确认升级
echo ""
read -p "是否继续升级Docker? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 升级已取消"
    exit 1
fi

echo ""
echo "🔧 开始升级Docker..."

# 1. 停止Docker服务
echo "⏹️  停止Docker服务..."
sudo systemctl stop docker
sudo systemctl stop docker.socket
sudo systemctl stop containerd

# 2. 卸载旧版本Docker
echo "🗑️  卸载旧版本Docker..."
sudo apt-get remove -y docker docker-engine docker.io containerd runc docker-compose

# 3. 更新包索引
echo "📦 更新包索引..."
sudo apt-get update

# 4. 安装必要的包
echo "📋 安装必要的包..."
sudo apt-get install -y \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

# 5. 添加Docker官方GPG密钥
echo "🔑 添加Docker官方GPG密钥..."
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# 6. 设置Docker仓库
echo "📚 设置Docker仓库..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 7. 更新包索引
echo "🔄 更新包索引..."
sudo apt-get update

# 8. 安装最新版本Docker
echo "⬇️  安装最新版本Docker..."
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# 9. 配置Docker镜像加速器（中国大陆用户）
echo "🌐 配置Docker镜像加速器..."
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

# 10. 启动Docker服务
echo "▶️  启动Docker服务..."
sudo systemctl start docker
sudo systemctl enable docker

# 11. 将当前用户添加到docker组
echo "👤 将当前用户添加到docker组..."
sudo usermod -aG docker $USER

# 12. 验证安装
echo ""
echo "✅ 验证Docker安装..."
sleep 3

echo "📋 新版本信息："
docker --version
docker compose version

echo ""
echo "🧪 测试Docker功能..."
if docker run --rm hello-world >/dev/null 2>&1; then
    echo "✅ Docker功能测试成功"
else
    echo "❌ Docker功能测试失败"
fi

echo ""
echo "🎉 Docker升级完成！"
echo ""
echo "📝 重要提示："
echo "   1. 请重新登录或运行 'newgrp docker' 以应用用户组更改"
echo "   2. 现在可以使用 'docker compose' 命令（注意没有连字符）"
echo "   3. 已配置中国大陆镜像加速器"
echo ""
echo "🚀 现在可以运行："
echo "   docker compose build"
echo "   docker compose up -d"
