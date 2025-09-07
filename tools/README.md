# 空气质量监测系统 - 工具集

这个目录包含了用于测试、部署和管理空气质量监测系统的各种工具。

## 目录结构

```
tools/
├── mqtt/                    # MQTT测试工具
│   ├── basic_test.py       # 基础MQTT测试
│   ├── advanced_test.py    # 高级MQTT测试
│   ├── config_driven_test.py # 配置驱动测试
│   ├── demo.py             # 演示程序
│   ├── test_config.json    # 测试配置
│   └── README.md           # MQTT工具说明
├── server/                  # 服务器启动脚本
│   ├── start_server.sh     # Linux/Mac启动脚本
│   ├── start_server.bat    # Windows启动脚本
│   └── README.md           # 服务器启动说明
├── docs/                    # 文档
│   ├── QUICKSTART.md       # 快速开始指南
│   └── README.md           # 工具总览
├── requirements.txt         # Python依赖
└── README.md               # 本文件
```

## 快速开始

### 1. 安装依赖
```bash
pip install -r requirements.txt
```

### 2. 启动服务器
```bash
# Linux/Mac
./server/start_server.sh

# Windows
server\start_server.bat
```

### 3. 测试MQTT功能
```bash
# 基础测试
python mqtt/basic_test.py --count 5

# 高级测试
python mqtt/advanced_test.py --config

# 交互式演示
python mqtt/demo.py
```

## 工具分类

### MQTT测试工具 (`mqtt/`)
用于测试MQTT消息收发功能：
- **基础测试** - 发送传感器数据和设备状态
- **高级测试** - 配置消息、命令消息、设备响应
- **配置驱动测试** - 基于JSON配置的批量测试
- **演示程序** - 交互式学习和演示

### 服务器启动脚本 (`server/`)
用于启动空气质量监测服务器：
- **Windows脚本** - 批处理文件，设置环境变量
- **Linux/Mac脚本** - Shell脚本，支持信号处理

### 文档 (`docs/`)
提供详细的使用说明：
- **快速开始指南** - 环境准备和基本使用
- **工具总览** - 所有工具的详细说明

## 使用流程

1. **环境准备** - 参考`docs/QUICKSTART.md`
2. **安装依赖** - 运行`pip install -r requirements.txt`
3. **启动服务器** - 使用`server/`目录下的启动脚本
4. **测试MQTT** - 使用`mqtt/`目录下的测试工具
5. **查看文档** - 参考各目录下的README文件

## MQTT主题格式

所有工具都使用以下MQTT主题格式：

- **数据主题**: `air-quality/hcho/{device_id}/data`
- **状态主题**: `air-quality/hcho/{device_id}/status`
- **响应主题**: `air-quality/hcho/{device_id}/response`
- **配置主题**: `air-quality/hcho/{device_id}/config`
- **命令主题**: `air-quality/hcho/{device_id}/command`

## 环境要求

### Python环境
- Python 3.7+
- 依赖包：paho-mqtt, colorama

### Go环境
- Go 1.19+
- 依赖：github.com/mochi-mqtt/server/v2

### 系统要求
- 端口1883（MQTT）
- 端口8080（HTTP API）
- 数据库连接（MySQL/PostgreSQL）
- Redis连接（可选）

## 故障排除

### 常见问题

1. **MQTT连接失败**
   - 检查MQTT服务器是否启动
   - 确认端口1883未被占用
   - 验证网络连接

2. **Python依赖问题**
   ```bash
   pip install --upgrade pip
   pip install -r requirements.txt
   ```

3. **Go编译错误**
   ```bash
   go mod tidy
   go build -o bin/air-quality-server cmd/air-quality-server/main.go
   ```

4. **端口被占用**
   ```bash
   # Windows
   netstat -an | findstr :8080
   taskkill /PID <进程ID> /F
   
   # Linux/Mac
   lsof -i :8080
   kill -9 <进程ID>
   ```

## 支持

如果遇到问题，请：
1. 查看相关目录的README文件
2. 参考`docs/QUICKSTART.md`中的故障排除部分
3. 检查日志输出
4. 确认环境配置正确

## 许可证

本项目使用MIT许可证。详见LICENSE文件。