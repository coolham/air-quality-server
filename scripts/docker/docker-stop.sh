#!/bin/bash

# 空气质量监测系统 - Docker停止脚本

set -e

echo "🛑 停止空气质量监测系统..."

# 停止生产环境
echo "停止生产环境服务..."
docker-compose down

# 停止开发环境
echo "停止开发环境服务..."
docker-compose -f docker-compose.dev.yml down

# 清理未使用的镜像和容器
echo "🧹 清理未使用的Docker资源..."
docker system prune -f

echo "✅ 空气质量监测系统已停止！"
echo ""
echo "📝 如需完全清理，可运行:"
echo "  docker-compose down -v  # 删除数据卷"
echo "  docker system prune -a  # 删除所有未使用的镜像"
