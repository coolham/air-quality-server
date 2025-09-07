# 测试与调试脚本

本目录包含各种测试和调试脚本。

## 脚本说明

- **run-tests.bat** - 运行Go单元测试
- **test-config.go** - 测试配置文件路径解析
- **test-mqtt-storage.bat** - 测试MQTT服务器数据存储功能
- **test-web-data.bat** - 测试Web接口数据功能

## 使用流程

### 运行单元测试
```cmd
run-tests.bat
```

### 测试配置
```cmd
go run test-config.go
```

### 测试MQTT功能
```cmd
test-mqtt-storage.bat
```

### 测试Web功能
```cmd
test-web-data.bat
```
