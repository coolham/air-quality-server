#!/bin/bash

# 空气质量监测系统 - 中国大陆Docker构建脚本
# 解决Go模块下载超时问题

set -e

echo "🚀 为中国大陆网络环境构建Docker镜像..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

# 设置Go代理环境变量
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

echo "🔧 设置Go代理环境变量:"
echo "   GOPROXY=$GOPROXY"
echo "   GOSUMDB=$GOSUMDB"

# 创建构建参数
BUILD_ARGS=""
BUILD_ARGS="$BUILD_ARGS --build-arg GOPROXY=$GOPROXY"
BUILD_ARGS="$BUILD_ARGS --build-arg GOSUMDB=$GOSUMDB"

# 检查构建类型
BUILD_TYPE=${1:-"production"}

if [ "$BUILD_TYPE" = "dev" ]; then
    echo "🔨 构建开发环境镜像..."
    docker build $BUILD_ARGS -f Dockerfile.dev -t air-quality-server:dev .
else
    echo "🔨 构建生产环境镜像..."
    docker build $BUILD_ARGS -f Dockerfile -t air-quality-server:latest .
fi

echo "✅ Docker镜像构建完成！"

# 显示镜像信息
echo "📋 构建的镜像:"
docker images | grep air-quality-server

echo ""
echo "📝 使用说明:"
if [ "$BUILD_TYPE" = "dev" ]; then
    echo "  启动开发环境: docker-compose -f docker-compose.dev.yml up -d"
else
    echo "  启动生产环境: docker-compose up -d"
fi
echo "  查看镜像: docker images | grep air-quality-server"
echo "  删除镜像: docker rmi air-quality-server:latest"
