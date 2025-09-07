# 空气质量监测服务端系统 Makefile

.PHONY: help build test clean docker-build docker-run dev

# 默认目标
help:
	@echo "空气质量监测服务端系统构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  build          构建所有服务"
	@echo "  test           运行测试"
	@echo "  clean          清理构建文件"
	@echo "  docker-build   构建Docker镜像"
	@echo "  docker-run     运行Docker容器"
	@echo "  dev            启动开发环境"
	@echo "  lint           代码检查"
	@echo "  fmt            代码格式化"
	@echo "  proto          生成protobuf代码"
	@echo "  migrate        初始化数据库"
	@echo "  migrate-status 查看数据库状态"

# 构建应用
build:
	@echo "构建应用..."
	@go build -o bin/air-quality-server ./cmd/air-quality-server
	@echo "构建完成!"

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/
	@go clean

# 代码检查
lint:
	@echo "运行代码检查..."
	@golangci-lint run

# 代码格式化
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 生成protobuf代码
proto:
	@echo "生成protobuf代码..."
	@protoc --go_out=. --go-grpc_out=. ./api/proto/*.proto

# 构建Docker镜像
docker-build:
	@echo "构建Docker镜像..."
	@docker-compose build

# 运行Docker容器
docker-run:
	@echo "启动Docker容器..."
	@docker-compose up -d

# 启动开发环境
dev:
	@echo "启动开发环境..."
	@docker-compose -f docker-compose.dev.yml up -d

# 停止开发环境
dev-stop:
	@echo "停止开发环境..."
	@docker-compose -f docker-compose.dev.yml down

# 生成API文档
docs:
	@echo "生成API文档..."
	@swag init -g ./cmd/api-gateway/main.go -o ./docs/swagger

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy

# 初始化数据库
migrate:
	@echo "初始化数据库..."
	@go run ./cmd/migrate/main.go -action init

# 查看数据库状态
migrate-status:
	@echo "查看数据库状态..."
	@go run ./cmd/migrate/main.go -action status

# 启动应用
start:
	@echo "启动应用..."
	@make build
	@./bin/air-quality-server

# 停止应用
stop:
	@echo "停止应用..."
	@pkill -f air-quality-server
	@echo "应用已停止"
