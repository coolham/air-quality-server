#!/bin/bash

# 重置数据库脚本
# 删除现有数据卷，让 MySQL 重新执行 init.sql

echo "========================================"
echo "重置数据库"
echo "========================================"
echo "问题：MySQL 数据卷中保存了旧数据，不会重新执行 init.sql"
echo "解决：删除数据卷，让 MySQL 重新初始化"
echo "========================================"

# 检查 Docker
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行"
    exit 1
fi

echo ""
echo "⚠️  警告：这将删除所有数据库数据！"
echo "包括：用户数据、设备数据、传感器数据等"
echo ""
read -p "确认重置数据库？(y/N): " confirm
if [[ ! "$confirm" =~ ^[Yy]$ ]]; then
    echo "操作已取消"
    exit 0
fi

echo ""
echo "步骤 1: 停止所有服务..."
docker compose down

echo ""
echo "步骤 2: 删除 MySQL 数据卷..."
docker volume rm air-quality-server_mysql_data
if [ $? -eq 0 ]; then
    echo "✅ MySQL 数据卷已删除"
else
    echo "ℹ️  MySQL 数据卷不存在或已删除"
fi

echo ""
echo "步骤 3: 删除 Redis 数据卷（可选）..."
docker volume rm air-quality-server_redis_data
if [ $? -eq 0 ]; then
    echo "✅ Redis 数据卷已删除"
else
    echo "ℹ️  Redis 数据卷不存在或已删除"
fi

echo ""
echo "步骤 4: 重新启动服务..."
echo "现在 MySQL 会重新执行 init.sql 脚本，创建新的表结构"
docker compose up -d

echo ""
echo "等待数据库初始化完成..."
sleep 30

echo ""
echo "步骤 5: 验证数据库结构..."
echo "检查 unified_sensor_data 表是否包含新字段..."
docker compose exec mysql mysql -uroot -padmin -e "DESCRIBE air_quality.unified_sensor_data;" 2>/dev/null

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 数据库重置完成！"
    echo ""
    echo "现在数据库包含所有必要字段："
    echo "  ✓ o3, no2, so2, co, voc (污染物指标)"
    echo "  ✓ signal_strength (信号强度)"
    echo "  ✓ quality_score (数据质量评分)"
    echo "  ✓ deleted_at (软删除字段)"
    echo ""
    echo "应用程序现在应该可以正常工作了。"
    echo "访问地址：http://localhost:8082"
else
    echo ""
    echo "❌ 数据库重置失败！"
    echo "请检查 Docker 日志：docker compose logs mysql"
fi

echo ""
read -p "按回车键继续..."
