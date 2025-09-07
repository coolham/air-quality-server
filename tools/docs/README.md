# 工具文档

这个目录包含了tools目录下所有工具的详细文档。

## 文档列表

### 快速开始指南
- **`QUICKSTART.md`** - 快速开始指南
  - 环境准备
  - 依赖安装
  - 基本使用
  - 常见问题

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
└── README.md               # 根目录说明
```

## 工具分类

### MQTT测试工具
用于测试MQTT消息收发功能，包括：
- 基础消息发送测试
- 高级功能测试（配置、命令）
- 配置驱动的批量测试
- 交互式演示程序

### 服务器启动脚本
用于启动空气质量监测服务器：
- Windows批处理脚本
- Linux/Mac Shell脚本
- 环境变量配置
- 错误处理

### 文档
提供详细的使用说明：
- 快速开始指南
- 工具使用说明
- 故障排除指南

## 使用流程

1. **环境准备** - 参考`QUICKSTART.md`
2. **安装依赖** - 运行`pip install -r requirements.txt`
3. **启动服务器** - 使用`server/`目录下的启动脚本
4. **测试MQTT** - 使用`mqtt/`目录下的测试工具
5. **查看文档** - 参考各目录下的README文件

## 支持

如果遇到问题，请：
1. 查看相关目录的README文件
2. 参考`QUICKSTART.md`中的故障排除部分
3. 检查日志输出
4. 确认环境配置正确
