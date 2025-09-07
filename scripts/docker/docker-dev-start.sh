#!/bin/bash

# 空气质量监测系统 - 开发环境Docker启动脚本

set -e

echo "🚀 启动空气质量监测系统开发环境..."

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

# 停止现有容器
echo "🛑 停止现有开发环境容器..."
docker-compose -f docker-compose.dev.yml down

# 启动开发环境服务
echo "🔨 启动开发环境服务..."
docker-compose -f docker-compose.dev.yml up -d

# 等待服务启动
echo "⏳ 等待服务启动..."
sleep 10

# 检查服务状态
echo "🔍 检查开发环境服务状态..."
docker-compose -f docker-compose.dev.yml ps

echo "✅ 开发环境启动完成！"
echo "🗄️  MySQL: localhost:3307 (用户名: root, 密码: admin)"
echo "🔴  Redis: localhost:6380"
echo "📡  MQTT: localhost:1884 (可选)"
echo ""
echo "📝 常用命令:"
echo "  查看MySQL日志: docker-compose -f docker-compose.dev.yml logs -f mysql-dev"
echo "  查看Redis日志: docker-compose -f docker-compose.dev.yml logs -f redis-dev"
echo "  停止开发环境: docker-compose -f docker-compose.dev.yml down"
echo "  查看状态: docker-compose -f docker-compose.dev.yml ps"
echo ""
echo "💡 提示: 开发环境只启动数据库和缓存服务，应用服务请在本地运行"
