#!/bin/bash

# 启动空气质量监测系统 - 自动选择Docker Compose版本

echo "🚀 启动空气质量监测系统..."
echo

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker"
    exit 1
fi

echo "📋 检查Docker Compose可用性..."

# 检查Docker Compose V2
if docker compose version > /dev/null 2>&1; then
    echo "✅ 检测到Docker Compose V2"
    echo "🚀 使用Docker Compose V2启动服务..."
    docker compose up --build -d
    
    if [ $? -eq 0 ]; then
        echo
        echo "✅ 服务启动成功！"
        echo
        echo "🌐 访问地址:"
        echo "  Web界面: http://localhost:8082"
        echo "  MySQL: localhost:3308"
        echo "  Redis: localhost:6381"
        echo "  MQTT: localhost:1883"
        echo
        echo "📋 服务状态:"
        docker compose ps
    else
        echo "❌ Docker Compose V2启动失败"
        exit 1
    fi

# 检查Docker Compose V1
elif docker-compose --version > /dev/null 2>&1; then
    echo "✅ 检测到Docker Compose V1"
    echo "🚀 使用Docker Compose V1启动服务..."
    docker-compose up --build -d
    
    if [ $? -eq 0 ]; then
        echo
        echo "✅ 服务启动成功！"
        echo
        echo "🌐 访问地址:"
        echo "  Web界面: http://localhost:8082"
        echo "  MySQL: localhost:3308"
        echo "  Redis: localhost:6381"
        echo "  MQTT: localhost:1883"
        echo
        echo "📋 服务状态:"
        docker-compose ps
    else
        echo "❌ Docker Compose V1启动失败"
        exit 1
    fi

# 都没有找到，使用Docker命令
else
    echo "⚠️  未找到Docker Compose，使用Docker命令启动..."
    echo "🚀 使用Docker命令启动服务..."
    
    # 检查Docker启动脚本是否存在
    if [ -f "scripts/docker/docker-run.sh" ]; then
        chmod +x scripts/docker/docker-run.sh
        ./scripts/docker/docker-run.sh
    else
        echo "❌ Docker启动脚本不存在"
        echo "💡 请手动安装Docker Compose或使用Docker命令"
        exit 1
    fi
fi
