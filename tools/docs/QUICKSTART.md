# MQTT测试工具快速开始指南

## 🚀 快速开始

### 1. 安装依赖

```bash
# 安装Python依赖
pip install -r requirements.txt
```

### 2. 启动MQTT Broker

确保您的MQTT Broker正在运行。如果使用Docker：

```bash
# 启动MQTT Broker
docker-compose -f docker-compose.mqtt.yml up -d
```

### 3. 运行测试

#### 方式1：使用演示脚本（推荐新手）

```bash
python demo.py
```

#### 方式2：直接运行基础测试

```bash
# 发送10条测试消息
python mqtt_test.py --count 10

# 持续模拟设备数据上报
python mqtt_test.py --simulate --interval 30
```

#### 方式3：使用便捷脚本

**Windows:**
```bash
run_test.bat
```

**Linux/Mac:**
```bash
chmod +x run_test.sh
./run_test.sh
```

## 📋 常用命令

### 基础测试

```bash
# 快速测试
python mqtt_test.py --count 5

# 指定设备ID
python mqtt_test.py --device-id hcho_001 --count 10

# 持续模拟
python mqtt_test.py --simulate --interval 60
```

### 高级测试

```bash
# 测试配置下发
python mqtt_advanced_test.py --test-type config

# 测试命令控制
python mqtt_advanced_test.py --test-type command --command calibrate

# 完整测试
python mqtt_advanced_test.py --test-type all
```

### 配置驱动测试

```bash
# 交互模式
python config_driven_test.py --interactive

# 命令行模式
python config_driven_test.py --broker 1 --scenario 1
```

## 🔧 配置说明

### MQTT主题格式

- **数据主题**: `air-quality/hcho/{device_id}/data`
- **状态主题**: `air-quality/hcho/{device_id}/status`
- **配置主题**: `air-quality/hcho/{device_id}/config`
- **命令主题**: `air-quality/hcho/{device_id}/command`
- **响应主题**: `air-quality/hcho/{device_id}/response`

### 默认配置

- **MQTT Broker**: localhost:1883
- **用户名**: admin
- **密码**: password
- **设备ID**: hcho_001

## 🧪 测试场景

### 1. 基础数据上报测试

测试设备数据上报功能，验证MQTT消息格式和传输。

```bash
python mqtt_test.py --count 10 --device-id test_device_001
```

### 2. 持续模拟测试

模拟真实设备持续上报数据。

```bash
python mqtt_test.py --simulate --interval 30 --device-id hcho_001
```

### 3. 配置下发测试

测试服务器向设备下发配置的功能。

```bash
python mqtt_advanced_test.py --test-type config --device-id hcho_001
```

### 4. 命令控制测试

测试服务器向设备发送控制命令的功能。

```bash
python mqtt_advanced_test.py --test-type command --command calibrate --device-id hcho_001
```

### 5. 多设备测试

同时模拟多个设备上报数据。

```bash
# 终端1
python mqtt_test.py --simulate --device-id hcho_001 --interval 30

# 终端2
python mqtt_test.py --simulate --device-id hcho_002 --interval 45

# 终端3
python mqtt_test.py --simulate --device-id hcho_003 --interval 60
```

## 📊 数据格式示例

### 传感器数据消息

```json
{
  "device_id": "hcho_001",
  "device_type": "hcho",
  "timestamp": 1694000000,
  "data": {
    "formaldehyde": 0.08,
    "temperature": 25.5,
    "humidity": 60.2,
    "battery": 85
  },
  "location": {
    "latitude": 39.9042,
    "longitude": 116.4074,
    "address": "北京市朝阳区"
  },
  "quality": {
    "signal_strength": -65,
    "data_quality": "good"
  }
}
```

### 设备状态消息

```json
{
  "device_id": "hcho_001",
  "device_type": "hcho",
  "timestamp": 1694000000,
  "status": {
    "online": true,
    "battery_level": 85,
    "signal_strength": -65,
    "last_data_time": 1693999990,
    "error_code": 0,
    "error_message": ""
  },
  "firmware": {
    "version": "1.2.3",
    "build_date": "2024-01-15"
  }
}
```

## 🛠️ 故障排除

### 连接问题

1. **检查MQTT Broker是否运行**
   ```bash
   # 检查端口是否开放
   netstat -an | grep 1883
   ```

2. **检查网络连接**
   ```bash
   # 测试连接
   telnet localhost 1883
   ```

3. **检查认证信息**
   - 确认用户名和密码正确
   - 检查MQTT Broker的认证配置

### 依赖问题

1. **Python版本**
   ```bash
   python --version  # 需要Python 3.6+
   ```

2. **依赖包**
   ```bash
   pip install -r requirements.txt
   ```

### 权限问题

1. **Linux/Mac执行权限**
   ```bash
   chmod +x run_test.sh
   ```

## 📚 更多信息

- 详细文档: [README.md](README.md)
- 配置文件: [test_config.json](test_config.json)
- 演示脚本: [demo.py](demo.py)

## 🆘 获取帮助

如果遇到问题，可以：

1. 查看详细文档: `README.md`
2. 运行演示脚本: `python demo.py`
3. 检查配置文件: `test_config.json`
4. 查看命令行帮助: `python mqtt_test.py --help`
