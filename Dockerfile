# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai

# 设置Go代理（解决中国大陆网络访问问题）
# 支持通过构建参数覆盖
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn
ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o air-quality-server \
    ./cmd/air-quality-server

# 运行阶段
FROM alpine:latest

# 安装必要的工具
RUN apk add --no-cache wget ca-certificates tzdata

# 从构建阶段复制必要的文件
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# 设置时区
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/air-quality-server .

# 复制配置文件和Web模板
COPY --from=builder /app/config ./config
COPY --from=builder /app/web ./web

# 创建日志目录并设置权限
RUN mkdir -p /app/logs && \
    chmod 755 /app/logs && \
    chown -R root:root /app/config && \
    chmod -R 755 /app/config && \
    chown -R root:root /app/web && \
    chmod -R 755 /app/web && \
    chown root:root /app/air-quality-server && \
    chmod 755 /app/air-quality-server

# 暴露端口
EXPOSE 8080 1883

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD ["wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"] || exit 1

# 启动应用
CMD ["./air-quality-server"]
