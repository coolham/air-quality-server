# 空气质量监测服务端系统

## 项目概述

这是一个企业级的空气质量监测服务端系统，用于接收、存储、处理和分析来自ESP32空气质量监测设备的数据。系统采用简化的微服务架构，使用Go语言开发，具备高可用、高性能、易部署的特点。

## 系统特性

### 核心功能
- **数据接收**: 支持HTTP/WebSocket协议接收设备数据
- **数据存储**: 使用MySQL存储所有数据，Redis作为缓存和消息队列
- **实时处理**: 基于Redis Pub/Sub的实时数据处理和告警
- **可视化**: 提供RESTful API接口进行数据查询和可视化

### 技术特性
- **简化架构**: 去除了复杂的时序数据库和消息队列，使用MySQL+Redis的轻量级方案
- **单体应用**: 模块化设计，单一部署，易于维护
- **高性能**: 支持高并发数据接收和处理
- **易部署**: 使用Docker Compose一键部署

## 系统架构

### 整体架构
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ESP32设备     │    │   移动端APP     │    │   Web管理端     │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP/WebSocket       │ HTTP/WebSocket       │ HTTP/WebSocket
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │    空气质量监测服务端      │
                    │   (单体应用架构)          │
                    └─────────────┬─────────────┘
                                 │
        ┌────────────────────────┼────────────────────────┐
        │                        │                        │
┌───────▼────────┐    ┌─────────▼─────────┐    ┌─────────▼─────────┐
│   消息队列      │    │   关系数据库      │    │   Redis缓存       │
│  (Redis Pub/Sub)│    │   (MySQL)         │    │   (Redis)         │
└────────────────┘    └───────────────────┘    └───────────────────┘
```

### 核心功能模块
1. **设备管理模块** - 设备注册、状态管理、配置下发
2. **数据接收模块** - 接收设备数据、数据验证、格式转换
3. **数据处理模块** - 实时数据处理、统计分析、异常检测
4. **数据查询模块** - 历史数据查询、报表生成、数据导出
5. **告警模块** - 告警规则管理、告警触发、通知发送
6. **用户管理模块** - 用户认证、权限管理、组织架构
7. **配置管理模块** - 系统配置、设备配置、业务配置

## 技术栈

### 后端技术
- **语言**: Go 1.21+
- **框架**: Gin (HTTP) + gRPC
- **数据库**: MySQL 8.0 + Redis 7.0
- **消息队列**: Redis Pub/Sub
- **监控**: Prometheus + Grafana
- **日志**: 结构化日志 (JSON格式)

### 部署和运维
- **容器化**: Docker + Docker Compose
- **反向代理**: Nginx
- **CI/CD**: GitLab CI/CD
- **监控**: Prometheus + Grafana

## 快速开始

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+
- Go 1.21+ (开发环境)

### 部署步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd air-quality-server
```

2. **启动服务**
```bash
# 启动所有服务
docker-compose up -d
```

3. **验证部署**
```bash
# 检查服务状态
docker-compose ps

# 测试API
curl http://localhost:8080/health
```

### 开发环境

1. **安装依赖**
```bash
go mod download
```

2. **配置环境**
```bash
# 使用默认配置文件
export CONFIG_FILE=config/config.yaml

# 或者使用环境变量配置
export DB_HOST=localhost
export DB_PORT=3306
export REDIS_HOST=localhost
export REDIS_PORT=6379
```

3. **启动开发环境**
```bash
make dev
```

4. **运行测试**
```bash
make test
```

### 配置文件说明

系统支持多种配置加载方式：

1. **配置文件**：优先使用 `config/config.yaml`
2. **环境变量**：通过 `CONFIG_FILE` 指定配置文件路径
3. **自动查找**：系统会自动查找以下路径的配置文件：
   - `config/config.yaml`
   - `config.yaml`
   - 项目根目录下的配置文件

配置文件路径支持相对路径和绝对路径，相对路径会自动相对于项目根目录解析。

## 项目结构

```
air-quality-server/
├── cmd/                # 应用入口
│   └── server/
│       └── main.go     # 主程序入口
├── internal/           # 内部包
│   ├── models/        # 数据模型
│   ├── services/      # 业务逻辑层
│   ├── repositories/  # 数据访问层
│   ├── handlers/      # HTTP处理器
│   ├── middleware/    # 中间件
│   ├── config/        # 配置管理
│   └── utils/         # 工具包
├── config/            # 配置文件
├── scripts/           # 脚本文件
├── docs/              # 文档
├── docker-compose.yml # Docker编排文件
├── Dockerfile         # Docker构建文件
├── Makefile           # 构建脚本
└── go.mod             # Go模块文件
```

## 文档

- [系统设计文档](docs/system_design.md) - 详细的系统架构设计
- [模块接口文档](docs/module_interfaces.md) - 各模块接口定义
- [数据库设计文档](docs/database_design.md) - 数据库表结构设计
- [部署指南](docs/deployment_guide.md) - 详细的部署和运维指南

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 项目Issues: [GitHub Issues](https://github.com/your-repo/air-quality-server/issues)
- 邮箱: xingshizhai@gmail.com



