# Scripts 目录说明

本目录包含空气质量监测系统的各种管理脚本，按功能分类组织，便于管理和使用。

## 📁 目录结构

```
scripts/
├── app/                    # 应用程序管理
│   ├── build-app.bat      # 构建应用程序
│   ├── run-app.bat        # 运行已构建的应用程序
│   ├── start-app.bat      # 开发模式启动
│   └── setup.bat          # 一键安装脚本
├── database/              # 数据库管理
│   ├── init-database.bat  # Go迁移工具初始化
│   ├── init-database-sql.bat # SQL脚本初始化
│   ├── check-database.bat # 检查数据库状态
│   └── init.sql           # 数据库初始化SQL脚本
├── docker/                # Docker管理
│   ├── docker-start.bat   # 启动Docker服务
│   ├── docker-start.sh    # Linux版本启动脚本
│   ├── docker-dev-start.bat # 开发环境启动
│   ├── docker-dev-start.sh  # Linux开发环境启动
│   └── docker-stop.sh     # 停止Docker服务
├── firewall/              # 防火墙管理
│   ├── configure_firewall.bat # 配置防火墙规则
│   ├── remove_firewall_rules.bat # 移除防火墙规则
│   ├── check_firewall_rules.ps1 # 检查防火墙规则
│   └── start_with_firewall.bat # 带防火墙启动
├── test/                  # 测试与调试
│   ├── run-tests.bat      # 运行单元测试
│   ├── test-config.go     # 配置路径测试
│   ├── test-mqtt-storage.bat # MQTT存储测试
│   └── test-web-data.bat  # Web数据测试
└── docs/                  # 文档
    └── README.md          # 本说明文档
```

## 🚀 快速开始

### 首次部署

#### 方法1：一键安装（推荐）
```cmd
scripts\app\setup.bat
```

#### 方法2：分步安装
```cmd
# 1. 构建应用程序
scripts\app\build-app.bat

# 2. 初始化数据库
scripts\database\init-database-sql.bat

# 3. 配置防火墙
scripts\firewall\configure_firewall.bat

# 4. 启动应用程序
scripts\app\start-app.bat
```

### 日常开发
```cmd
# 开发模式启动
scripts\app\start-app.bat

# 或带防火墙启动
scripts\firewall\start_with_firewall.bat
```

### 生产部署
```cmd
# 构建应用程序
scripts\app\build-app.bat

# 运行应用程序
scripts\app\run-app.bat
```

## 📋 功能分类详解

### 🖥️ 应用程序管理 (app/)

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `setup.bat` | 一键安装脚本 | 首次部署，自动完成环境检查、依赖下载、数据库初始化、应用构建 |
| `build-app.bat` | 构建应用程序 | 编译Go代码生成可执行文件 |
| `start-app.bat` | 开发模式启动 | 直接运行Go代码，适合开发调试 |
| `run-app.bat` | 运行应用程序 | 运行已构建的可执行文件，适合生产环境 |

### 🗄️ 数据库管理 (database/)

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `init-database-sql.bat` | SQL脚本初始化 | 使用SQL脚本创建数据库和表结构（推荐） |
| `init-database.bat` | Go迁移工具初始化 | 使用Go迁移工具创建表结构 |
| `check-database.bat` | 检查数据库状态 | 验证数据库初始化结果，显示表结构和数据统计 |
| `init.sql` | 数据库初始化SQL脚本 | 包含所有表结构和初始数据的SQL脚本 |

### 🐳 Docker管理 (docker/)

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `docker-start.bat` | 启动Docker服务 | 生产环境Docker部署 |
| `docker-start.sh` | Linux启动脚本 | Linux环境Docker部署 |
| `docker-dev-start.bat` | 开发环境启动 | 开发环境Docker部署 |
| `docker-dev-start.sh` | Linux开发环境启动 | Linux开发环境Docker部署 |
| `docker-stop.sh` | 停止Docker服务 | 停止所有Docker容器 |

### 🔥 防火墙管理 (firewall/)

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `configure_firewall.bat` | 配置防火墙规则 | 为应用程序和MQTT端口添加防火墙规则 |
| `remove_firewall_rules.bat` | 移除防火墙规则 | 清理之前添加的防火墙规则 |
| `check_firewall_rules.ps1` | 检查防火墙规则 | 查看现有的防火墙规则状态 |
| `start_with_firewall.bat` | 带防火墙启动 | 自动配置防火墙并启动应用程序 |

### 🧪 测试与调试 (test/)

| 脚本 | 功能 | 使用场景 |
|------|------|----------|
| `run-tests.bat` | 运行单元测试 | 执行Go单元测试，验证代码质量 |
| `test-config.go` | 配置路径测试 | 测试配置文件路径解析，调试配置问题 |
| `test-mqtt-storage.bat` | MQTT存储测试 | 测试MQTT服务器数据存储功能 |
| `test-web-data.bat` | Web数据测试 | 测试Web接口数据功能 |

## 🔧 使用流程

### 开发环境设置
```cmd
# 1. 环境检查与安装
scripts\app\setup.bat

# 2. 开发模式启动
scripts\app\start-app.bat
```

### 生产环境部署
```cmd
# 1. 构建应用程序
scripts\app\build-app.bat

# 2. 初始化数据库
scripts\database\init-database-sql.bat

# 3. 配置防火墙
scripts\firewall\configure_firewall.bat

# 4. 运行应用程序
scripts\app\run-app.bat
```

### Docker环境部署
```cmd
# 1. 启动Docker服务
scripts\docker\docker-start.bat

# 2. 停止Docker服务
scripts\docker\docker-stop.sh
```

### 测试与调试
```cmd
# 1. 运行单元测试
scripts\test\run-tests.bat

# 2. 测试配置
go run scripts\test\test-config.go

# 3. 测试MQTT功能
scripts\test\test-mqtt-storage.bat
```

## ⚠️ 注意事项

### 权限要求
- **防火墙脚本**：需要管理员权限
- **Docker脚本**：需要Docker Desktop运行
- **数据库脚本**：需要MySQL服务运行

### 环境依赖
- **Go环境**：Go 1.21或更高版本
- **MySQL**：MySQL 8.0或更高版本
- **Docker**：Docker Desktop（可选）

### 配置文件
- 确保 `config/config.yaml` 中的数据库配置正确
- Docker部署需要 `config/config.docker.yaml`

## 🐛 故障排除

### 常见问题

#### 1. 数据库连接失败
```cmd
# 检查数据库状态
scripts\database\check-database.bat

# 重新初始化数据库
scripts\database\init-database-sql.bat
```

#### 2. 防火墙问题
```cmd
# 检查防火墙规则
powershell -ExecutionPolicy Bypass -File scripts\firewall\check_firewall_rules.ps1

# 重新配置防火墙
scripts\firewall\configure_firewall.bat
```

#### 3. 应用程序启动失败
```cmd
# 检查配置
go run scripts\test\test-config.go

# 重新构建
scripts\app\build-app.bat
```

#### 4. Docker问题
```cmd
# 检查Docker状态
docker info

# 重启Docker服务
scripts\docker\docker-stop.sh
scripts\docker\docker-start.bat
```

## 📚 技术说明

### 脚本特性
- **跨平台支持**：提供Windows (.bat) 和Linux (.sh) 版本
- **中文支持**：使用UTF-8编码，支持中文显示
- **错误处理**：提供详细的错误信息和解决建议
- **用户友好**：交互式界面，清晰的执行反馈

### 路径处理
- **自动检测**：自动检测项目根目录
- **相对路径**：支持相对路径和绝对路径
- **目录创建**：自动创建必要的目录结构

### 安全考虑
- **权限检查**：检查必要的系统权限
- **输入验证**：验证配置文件和依赖
- **错误恢复**：提供错误恢复机制

---

**文档版本**: v2.0  
**创建日期**: 2024-09-07  
**最后更新**: 2024-09-07  
**作者**: 空气质量监测系统开发团队
