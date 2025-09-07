#!/bin/bash

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================"
echo -e "   启动空气质量监测服务"
echo -e "   (包含内置MQTT服务器)"
echo -e "========================================${NC}"
echo

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到Go，请先安装Go${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Go已安装: $(go version)${NC}"
echo

# 设置环境变量
export AIR_QUALITY_CONFIG="config/config.yaml"
export SERVER_HOST="0.0.0.0"
export SERVER_PORT="8080"
export LOG_LEVEL="info"

echo -e "${YELLOW}📋 启动配置:${NC}"
echo "  - HTTP服务器: http://localhost:8080"
echo "  - MQTT服务器: tcp://localhost:1883"
echo "  - 用户名: admin"
echo "  - 密码: password"
echo

echo -e "${GREEN}🚀 正在启动服务...${NC}"
echo "按 Ctrl+C 停止服务"
echo "----------------------------------------"

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 启动Go应用程序
go run cmd/air-quality-server/main.go

echo
echo -e "${GREEN}✅ 服务已停止${NC}"
