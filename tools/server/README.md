# 服务器启动脚本

这个目录包含了用于启动空气质量监测服务器的脚本。

## 文件说明

### Windows启动脚本
- **`start_server.bat`** - Windows批处理脚本
  - 设置环境变量
  - 启动Go应用程序
  - 包含错误处理

### Linux/Mac启动脚本
- **`start_server.sh`** - Shell脚本
  - 设置环境变量
  - 启动Go应用程序
  - 支持信号处理

## 使用方法

### Windows
```cmd
# 直接运行
start_server.bat

# 或者在PowerShell中
.\start_server.bat
```

### Linux/Mac
```bash
# 添加执行权限
chmod +x start_server.sh

# 运行脚本
./start_server.sh
```

## 环境变量

脚本会设置以下环境变量：

- `AIR_QUALITY_CONFIG` - 配置文件路径
- `AIR_QUALITY_WEB_ROOT` - Web资源根目录
- `GIN_MODE` - Gin框架模式（development/release）

## 服务端口

启动后，服务器将在以下端口提供服务：

- **HTTP API**: `http://localhost:8080`
- **MQTT Broker**: `mqtt://localhost:1883`
- **Web界面**: `http://localhost:8080/dashboard`

## 注意事项

1. 确保已安装Go 1.19+
2. 确保数据库和Redis服务已启动
3. 配置文件路径正确
4. 端口1883和8080未被占用

## 故障排除

### 端口被占用
```bash
# 检查端口占用
netstat -an | findstr :8080
netstat -an | findstr :1883

# 杀死占用进程
taskkill /PID <进程ID> /F
```

### 配置文件错误
- 检查`config/config.yaml`文件是否存在
- 验证配置文件格式是否正确
- 确认数据库连接参数

### 依赖服务未启动
- 确保MySQL/PostgreSQL数据库服务运行
- 确保Redis服务运行
- 检查网络连接
