# 应用程序管理脚本

本目录包含应用程序的构建、运行和管理脚本。

## 脚本说明

- **setup.bat** - 一键安装脚本，完成环境检查、依赖下载、数据库初始化、应用构建
- **build-app.bat** - 构建应用程序，生成可执行文件
- **start-app.bat** - 开发模式启动，直接运行Go代码
- **run-app.bat** - 运行已构建的应用程序

## 使用流程

### 首次部署
```cmd
setup.bat
```

### 日常开发
```cmd
start-app.bat
```

### 生产部署
```cmd
build-app.bat
run-app.bat
```
