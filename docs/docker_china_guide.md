# Docker构建中国大陆网络问题解决方案

## 问题描述

在中国大陆构建Docker镜像时，经常会遇到Go模块下载超时的问题：

```
ERROR [builder 5/7] RUN go mod download
go: github.com/bytedance/sonic@v1.9.1: Get "https://proxy.golang.org/github.com/bytedance/sonic/@v/v1.9.1.mod": dial tcp 142.250.69.177:443: i/o timeout
```

这是因为Go默认的代理服务器`proxy.golang.org`在中国大陆访问受限导致的。

## 解决方案

### 方案1：使用专门的构建脚本（推荐）

我们提供了专门为中国大陆网络环境优化的构建脚本：

#### Linux/macOS
```bash
# 构建生产环境镜像
./scripts/docker/docker-build-china.sh

# 构建开发环境镜像
./scripts/docker/docker-build-china.sh dev
```

#### Windows
```cmd
REM 构建生产环境镜像
scripts\docker\docker-build-china.bat

REM 构建开发环境镜像
scripts\docker\docker-build-china.bat dev
```

### 方案2：手动设置环境变量

#### 使用docker-compose构建
```bash
# 设置环境变量
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 构建并启动
docker-compose up --build -d
```

#### 使用docker build命令
```bash
docker build \
  --build-arg GOPROXY=https://goproxy.cn,direct \
  --build-arg GOSUMDB=sum.golang.google.cn \
  -f Dockerfile \
  -t air-quality-server:latest .
```

### 方案3：修改本地Go配置

如果您在本地开发环境中也遇到类似问题，可以设置全局Go代理：

```bash
# 设置Go代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 验证设置
go env GOPROXY
go env GOSUMDB
```

## 代理服务器说明

### 主要代理服务器

1. **goproxy.cn** - 七牛云提供的Go模块代理
   - 地址：https://goproxy.cn
   - 特点：稳定、快速、免费

2. **goproxy.io** - 另一个Go模块代理
   - 地址：https://goproxy.io
   - 特点：备用选择

3. **athens** - 自建代理服务器
   - 适用于企业内网环境

### 配置说明

```bash
# 推荐的代理配置
GOPROXY=https://goproxy.cn,direct
GOSUMDB=sum.golang.google.cn
```

- `GOPROXY`：指定模块代理服务器，`direct`表示如果代理失败则直接从源下载
- `GOSUMDB`：指定校验和数据库，用于验证模块完整性

## 故障排除

### 1. 仍然超时
如果使用goproxy.cn仍然超时，可以尝试：

```bash
# 使用备用代理
export GOPROXY=https://goproxy.io,direct

# 或者使用多个代理
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct
```

### 2. 私有模块问题
如果项目包含私有模块，需要配置：

```bash
# 设置私有模块不走代理
export GOPRIVATE=github.com/your-org/*,gitlab.com/your-org/*
```

### 3. 网络环境检测
可以使用以下命令测试网络连接：

```bash
# 测试代理连接
curl -I https://goproxy.cn

# 测试Go模块下载
go list -m -versions github.com/gin-gonic/gin
```

## 最佳实践

1. **使用构建脚本**：推荐使用提供的构建脚本，它们已经预配置了最佳的网络设置

2. **缓存优化**：在Dockerfile中先复制`go.mod`和`go.sum`，再运行`go mod download`，这样可以利用Docker的层缓存

3. **多阶段构建**：使用多阶段构建减少最终镜像大小

4. **网络重试**：在CI/CD环境中，可以设置重试机制

## 相关文件

- `Dockerfile` - 生产环境Dockerfile
- `Dockerfile.dev` - 开发环境Dockerfile
- `scripts/docker/docker-build-china.sh` - Linux/macOS构建脚本
- `scripts/docker/docker-build-china.bat` - Windows构建脚本
- `docker-compose.yml` - 生产环境配置
- `docker-compose.dev.yml` - 开发环境配置

## 更新日志

- 2024-01-XX：添加Go代理配置到Dockerfile
- 2024-01-XX：创建专门的构建脚本
- 2024-01-XX：添加详细的故障排除指南
