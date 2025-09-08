#!/bin/bash

# Docker 快速启动脚本
# 解决 "Unknown column 'o3' in 'field list'" 错误

echo "========================================"
echo "Docker 快速启动"
echo "========================================"

# 检查 Docker
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行"
    exit 1
fi

echo "正在启动服务..."
echo "数据库将通过 init 脚本自动初始化，包含所有必要字段"

docker-compose up -d

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 服务启动成功！"
    echo ""
    echo "服务地址："
    echo "  Web 界面: http://localhost:8082"
    echo "  MQTT 端口: 1883"
    echo "  MySQL 端口: 3308"
    echo ""
    echo "数据库已包含所有必要字段，错误已解决。"
else
    echo ""
    echo "❌ 服务启动失败！"
fi

echo ""
read -p "按回车键继续..."
