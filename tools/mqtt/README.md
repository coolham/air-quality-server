# MQTT 测试工具

这个目录包含了用于测试MQTT功能的Python工具。

## 文件说明

### 基础测试工具
- **`basic_test.py`** - 基础MQTT测试工具
  - 发送HCHO传感器数据
  - 发送设备状态信息
  - 支持单次测试和持续模拟

### 高级测试工具
- **`advanced_test.py`** - 高级MQTT测试工具
  - 支持配置消息发布
  - 支持命令消息发送
  - 支持设备响应模拟

### 配置驱动测试
- **`config_driven_test.py`** - 基于配置文件的测试工具
  - 使用`test_config.json`配置文件
  - 支持多种测试场景
  - 支持批量测试

### 演示程序
- **`demo.py`** - 交互式演示程序
  - 引导用户使用各种测试工具
  - 支持连续模拟模式
  - 适合学习和演示

### 配置文件
- **`test_config.json`** - 测试配置文件
  - 定义MQTT Broker配置
  - 定义测试场景
  - 定义设备模板

## 使用方法

### 安装依赖
```bash
pip install -r ../requirements.txt
```

### 基础测试
```bash
# 发送3条测试消息
python basic_test.py --count 3

# 持续模拟（每30秒发送一次）
python basic_test.py --simulate --interval 30
```

### 高级测试
```bash
# 发送配置消息
python advanced_test.py --config

# 发送命令消息
python advanced_test.py --command calibrate
```

### 配置驱动测试
```bash
# 使用默认配置运行
python config_driven_test.py

# 指定配置文件
python config_driven_test.py --config custom_config.json
```

### 演示程序
```bash
# 启动交互式演示
python demo.py
```

## 主题格式

所有工具都使用以下MQTT主题格式：

- **数据主题**: `air-quality/hcho/{device_id}/data`
- **状态主题**: `air-quality/hcho/{device_id}/status`
- **响应主题**: `air-quality/hcho/{device_id}/response`
- **配置主题**: `air-quality/hcho/{device_id}/config`
- **命令主题**: `air-quality/hcho/{device_id}/command`

## 注意事项

1. 确保MQTT服务器已启动（端口1883）
2. 默认设备ID为`hcho_001`
3. 所有消息都使用JSON格式
4. 支持QoS 1级别
