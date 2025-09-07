#!/bin/bash

# 空气质量监测系统 - Docker启动脚本

set -e

echo "🚀 启动空气质量监测系统..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 检查docker-compose是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ docker-compose未安装，请先安装docker-compose"
    exit 1
fi

# 创建必要的目录
echo "📁 创建必要的目录..."
mkdir -p logs
mkdir -p config

# 检查配置文件
if [ ! -f "config/config.docker.yaml" ]; then
    echo "❌ 配置文件不存在: config/config.docker.yaml"
    exit 1
fi

# 停止现有容器
echo "🛑 停止现有容器..."
docker-compose down

# 构建并启动服务
echo "🔨 构建并启动服务..."
docker-compose up --build -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "🔍 检查服务状态..."
docker-compose ps

# 显示日志
echo "📋 显示服务日志..."
docker-compose logs --tail=50 air-quality-server

echo "✅ 空气质量监测系统启动完成！"
echo "🌐 Web界面: http://localhost:8080"
echo "📊 Dashboard: http://localhost:8080/dashboard"
echo "📡 MQTT Broker: localhost:1883"
echo ""
echo "📝 常用命令:"
echo "  查看日志: docker-compose logs -f air-quality-server"
echo "  停止服务: docker-compose down"
echo "  重启服务: docker-compose restart air-quality-server"
echo "  查看状态: docker-compose ps"
