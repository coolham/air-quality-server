# 快速解决Docker构建超时问题

## 🚨 问题症状
```
ERROR [builder 5/7] RUN go mod download
go: github.com/bytedance/sonic@v1.9.1: Get "https://proxy.golang.org/github.com/bytedance/sonic/@v/v1.9.1.mod": dial tcp 142.250.69.177:443: i/o timeout
```

## ⚡ 快速解决方案

### 方案1：使用专用构建脚本（推荐）

**Windows用户：**
```cmd
# 构建生产环境
scripts\docker\docker-build-china.bat

# 构建开发环境  
scripts\docker\docker-build-china.bat dev
```

**Linux/macOS用户：**
```bash
# 构建生产环境
./scripts/docker/docker-build-china.sh

# 构建开发环境
./scripts/docker/docker-build-china.sh dev
```

### 方案2：手动设置环境变量

```cmd
# Windows PowerShell
$env:GOPROXY="https://goproxy.cn,direct"
$env:GOSUMDB="sum.golang.google.cn"
docker-compose up --build -d
```

```bash
# Linux/macOS
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
docker-compose up --build -d
```

### 方案3：直接使用docker build

```cmd
docker build --build-arg GOPROXY=https://goproxy.cn,direct --build-arg GOSUMDB=sum.golang.google.cn -f Dockerfile -t air-quality-server:latest .
```

## 🔍 测试代理是否有效

运行测试脚本：
```cmd
scripts\docker\test-go-proxy.bat
```

## 📋 已修改的文件

- ✅ `Dockerfile` - 添加了Go代理配置
- ✅ `Dockerfile.dev` - 添加了Go代理配置  
- ✅ `scripts/docker/docker-build-china.bat` - Windows构建脚本
- ✅ `scripts/docker/docker-build-china.sh` - Linux构建脚本
- ✅ `scripts/docker/test-go-proxy.bat` - 代理测试脚本
- ✅ `docs/docker_china_guide.md` - 详细解决方案文档

## 🎯 推荐使用流程

1. **首次使用**：运行测试脚本确认代理有效
   ```cmd
   scripts\docker\test-go-proxy.bat
   ```

2. **构建镜像**：使用专用构建脚本
   ```cmd
   scripts\docker\docker-build-china.bat
   ```

3. **启动服务**：使用docker-compose
   ```cmd
   docker-compose up -d
   ```

## 🆘 如果仍然失败

1. **尝试备用代理**：
   ```cmd
   set GOPROXY=https://goproxy.io,direct
   ```

2. **检查网络连接**：
   ```cmd
   ping goproxy.cn
   ```

3. **使用VPN**：如果网络环境限制严重，建议使用VPN

## 📞 获取帮助

- 查看详细文档：`docs/docker_china_guide.md`
- 检查Docker状态：`docker info`
- 查看构建日志：`docker-compose logs`
