#!/bin/bash

# Docker升级验证脚本

echo "🔍 Docker升级验证脚本"
echo "===================="

# 检查Docker版本
echo "📋 Docker版本信息："
if command -v docker >/dev/null 2>&1; then
    docker --version
else
    echo "❌ Docker未安装"
    exit 1
fi

# 检查Docker Compose版本
echo ""
echo "📋 Docker Compose版本信息："
if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    docker compose version
else
    echo "❌ Docker Compose未安装或版本过旧"
fi

# 检查Docker服务状态
echo ""
echo "🔧 Docker服务状态："
if systemctl is-active --quiet docker; then
    echo "✅ Docker服务正在运行"
else
    echo "❌ Docker服务未运行"
    echo "尝试启动Docker服务..."
    sudo systemctl start docker
    sleep 3
    if systemctl is-active --quiet docker; then
        echo "✅ Docker服务启动成功"
    else
        echo "❌ Docker服务启动失败"
    fi
fi

# 检查Docker配置
echo ""
echo "🌐 Docker镜像加速器配置："
if [ -f /etc/docker/daemon.json ]; then
    echo "✅ Docker配置文件存在"
    if grep -q "registry-mirrors" /etc/docker/daemon.json; then
        echo "✅ 镜像加速器已配置"
        echo "配置的镜像源："
        grep -A 10 "registry-mirrors" /etc/docker/daemon.json
    else
        echo "❌ 镜像加速器未配置"
    fi
else
    echo "❌ Docker配置文件不存在"
fi

# 检查用户组
echo ""
echo "👤 用户组检查："
if groups $USER | grep -q docker; then
    echo "✅ 用户已在docker组中"
else
    echo "❌ 用户不在docker组中"
    echo "请运行: sudo usermod -aG docker $USER"
    echo "然后重新登录或运行: newgrp docker"
fi

# 测试Docker功能
echo ""
echo "🧪 Docker功能测试："
if docker run --rm hello-world >/dev/null 2>&1; then
    echo "✅ Docker功能测试成功"
else
    echo "❌ Docker功能测试失败"
    echo "可能需要重新登录以应用用户组更改"
fi

# 测试镜像拉取
echo ""
echo "🌐 镜像拉取测试："
if docker pull alpine:latest >/dev/null 2>&1; then
    echo "✅ 镜像拉取测试成功"
    docker rmi alpine:latest >/dev/null 2>&1
else
    echo "❌ 镜像拉取测试失败"
fi

echo ""
echo "📊 升级验证完成！"
echo ""
echo "🚀 如果所有测试都通过，现在可以运行："
echo "   docker compose build"
echo "   docker compose up -d"
