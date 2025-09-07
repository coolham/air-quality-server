# Windows系统数据库初始化指南

## 概述

本文档专门针对Windows系统用户，说明如何在不使用Makefile的情况下进行数据库初始化。

## 前置条件

1. **MySQL 8.0+** 已安装并运行
2. **Go 1.21+** 开发环境
3. **PowerShell** 或 **命令提示符**
4. **配置文件** 已正确配置数据库连接信息

## 快速开始

### 1. 配置数据库连接

编辑 `config/config.yaml` 文件，确保数据库配置正确：

```yaml
database:
  host: localhost
  port: 3306
  username: air_quality
  password: air_quality123
  database: air_quality
  charset: utf8mb4
  max_idle: 10
  max_open: 100
  max_life: 3600
```

### 2. 创建数据库

#### 方法1：使用MySQL命令行

```cmd
# 打开命令提示符或PowerShell
# 连接到MySQL
mysql -u root -p

# 在MySQL中执行以下命令
CREATE DATABASE air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（可选）
CREATE USER 'air_quality'@'localhost' IDENTIFIED BY 'air_quality123';
GRANT ALL PRIVILEGES ON air_quality.* TO 'air_quality'@'localhost';
FLUSH PRIVILEGES;

# 退出MySQL
EXIT;
```

#### 方法1.1：使用SQL脚本（推荐）

```cmd
# 使用项目提供的SQL脚本创建数据库
mysql -u root -p < scripts/init.sql
```

#### 方法2：使用MySQL Workbench

1. 打开MySQL Workbench
2. 连接到MySQL服务器
3. 执行以下SQL语句：

```sql
CREATE DATABASE air_quality CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'air_quality'@'localhost' IDENTIFIED BY 'air_quality123';
GRANT ALL PRIVILEGES ON air_quality.* TO 'air_quality'@'localhost';
FLUSH PRIVILEGES;
```

### 3. 初始化数据库

#### 方法1：使用PowerShell（推荐）

```powershell
# 打开PowerShell，切换到项目目录
cd D:\Work\esp32\projects\air-quality-server

# 初始化数据库
go run cmd/migrate/main.go -action init

# 查看数据库状态
go run cmd/migrate/main.go -action status

# 查看帮助信息
go run cmd/migrate/main.go -help
```

#### 方法2：使用命令提示符

```cmd
# 打开命令提示符，切换到项目目录
cd D:\Work\esp32\projects\air-quality-server

# 初始化数据库
go run cmd\migrate\main.go -action init

# 查看数据库状态
go run cmd\migrate\main.go -action status

# 查看帮助信息
go run cmd\migrate\main.go -help
```

#### 方法3：使用指定配置文件

```powershell
# 使用指定配置文件初始化
go run cmd/migrate/main.go -config config/config.yaml -action init
```

### 4. 构建和运行应用程序

#### 构建应用程序

```powershell
# 构建应用程序
go build -o bin/air-quality-server.exe ./cmd/air-quality-server

# 运行应用程序
.\bin\air-quality-server.exe
```

#### 直接运行（开发模式）

```powershell
# 直接运行，不构建
go run cmd/air-quality-server/main.go
```

## Windows批处理脚本

为了简化操作，您可以创建批处理脚本：

### 创建 `init-database.bat`

```batch
@echo off
echo 正在初始化数据库...
go run cmd\migrate\main.go -action init
if %errorlevel% equ 0 (
    echo 数据库初始化成功！
) else (
    echo 数据库初始化失败！
    pause
)
```

### 创建 `check-database.bat`

```batch
@echo off
echo 正在检查数据库状态...
go run cmd\migrate\main.go -action status
pause
```

### 创建 `start-app.bat`

```batch
@echo off
echo 正在启动应用程序...
go run cmd\air-quality-server\main.go
pause
```

### 创建 `build-app.bat`

```batch
@echo off
echo 正在构建应用程序...
go build -o bin\air-quality-server.exe .\cmd\air-quality-server
if %errorlevel% equ 0 (
    echo 构建成功！可执行文件位于 bin\air-quality-server.exe
) else (
    echo 构建失败！
)
pause
```

## 使用批处理脚本

1. 将上述脚本保存到项目根目录
2. 双击运行相应的 `.bat` 文件
3. 或者在命令提示符中运行：

```cmd
# 初始化数据库
init-database.bat

# 检查数据库状态
check-database.bat

# 构建应用程序
build-app.bat

# 启动应用程序
start-app.bat
```

## PowerShell脚本（可选）

如果您更喜欢PowerShell，可以创建 `.ps1` 脚本：

### 创建 `init-database.ps1`

```powershell
Write-Host "正在初始化数据库..." -ForegroundColor Green
go run cmd/migrate/main.go -action init
if ($LASTEXITCODE -eq 0) {
    Write-Host "数据库初始化成功！" -ForegroundColor Green
} else {
    Write-Host "数据库初始化失败！" -ForegroundColor Red
    Read-Host "按任意键继续"
}
```

### 创建 `start-app.ps1`

```powershell
Write-Host "正在启动应用程序..." -ForegroundColor Green
go run cmd/air-quality-server/main.go
```

## 验证安装

### 1. 检查数据库状态

```powershell
go run cmd/migrate/main.go -action status
```

### 2. 连接数据库验证

```cmd
# 连接到数据库
mysql -u air_quality -p air_quality

# 查看所有表
SHOW TABLES;

# 检查用户数据
SELECT * FROM users;

# 检查角色数据
SELECT * FROM roles;
```

### 3. 启动应用程序

```powershell
# 方法1：直接运行
go run cmd/air-quality-server/main.go

# 方法2：先构建再运行
go build -o bin/air-quality-server.exe ./cmd/air-quality-server
.\bin\air-quality-server.exe
```

访问 `http://localhost:8080/health` 检查应用是否正常启动。

## 常见问题解决

### 1. Go命令找不到

**错误信息**：`'go' 不是内部或外部命令`

**解决方案**：
- 确保Go已正确安装
- 检查环境变量PATH是否包含Go的bin目录
- 重新打开命令提示符或PowerShell

### 2. MySQL连接失败

**错误信息**：`连接数据库失败`

**解决方案**：
- 检查MySQL服务是否运行
- 验证数据库连接配置
- 确认用户权限

### 3. 权限问题

**错误信息**：`Access denied`

**解决方案**：
- 以管理员身份运行命令提示符
- 检查数据库用户权限
- 确保用户有CREATE、INSERT、SELECT权限

### 4. 路径问题

**错误信息**：`找不到文件`

**解决方案**：
- 确保在正确的项目目录中
- 使用绝对路径
- 检查文件是否存在

## 环境变量配置

如果需要使用环境变量配置，可以设置：

```cmd
# 设置配置文件路径
set AIR_QUALITY_CONFIG=config\config.yaml

# 设置数据库配置
set DB_HOST=localhost
set DB_PORT=3306
set DB_USERNAME=air_quality
set DB_PASSWORD=air_quality123
set DB_DATABASE=air_quality
```

## 开发环境设置

### 使用VS Code

1. 安装Go扩展
2. 打开项目文件夹
3. 使用集成终端运行命令

### 使用GoLand

1. 打开项目
2. 配置运行配置
3. 使用内置终端

## 生产环境部署

### 构建生产版本

```powershell
# 构建生产版本
go build -ldflags "-s -w" -o bin/air-quality-server.exe ./cmd/air-quality-server

# 创建服务（需要管理员权限）
sc create "AirQualityServer" binpath="D:\path\to\air-quality-server.exe" start=auto
```

### 使用NSSM（推荐）

1. 下载NSSM：https://nssm.cc/download
2. 安装服务：

```cmd
nssm install AirQualityServer "D:\path\to\air-quality-server.exe"
nssm start AirQualityServer
```

## 总结

Windows系统下的数据库初始化步骤：

1. **配置数据库**：编辑 `config/config.yaml`
2. **创建数据库**：使用MySQL命令行或Workbench
3. **初始化数据库**：`go run cmd/migrate/main.go -action init`
4. **验证安装**：`go run cmd/migrate/main.go -action status`
5. **启动应用**：`go run cmd/air-quality-server/main.go`

使用批处理脚本可以简化操作流程，提高开发效率。

---

**文档版本**: v1.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
