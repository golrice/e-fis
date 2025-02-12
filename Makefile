# 设置 Go 变量
APP_NAME = distributed-cache
GO = go
BUILD_DIR = bin

# 可执行文件名称
SERVER_BIN = $(BUILD_DIR)/server
CLIENT_BIN = $(BUILD_DIR)/client

# 设置编译参数
LDFLAGS = -s -w  # 去掉调试信息，减小体积
GCFLAGS =       # 编译优化参数
TAGS =          # 需要的构建标签

# 目标编译的平台（可选）
OS ?= $(shell uname -s | tr A-Z a-z)
ARCH ?= $(shell uname -m)
GOOS ?= $(OS)
GOARCH ?= $(ARCH)
BUILD_FLAGS = GOOS=$(GOOS) GOARCH=$(GOARCH)

# 默认目标
.DEFAULT_GOAL := help

## 构建项目
.PHONY: build
build: clean
	@mkdir -p $(BUILD_DIR)
	@echo "🔨 Building Server..."
	$(GO) build -ldflags "$(LDFLAGS)" -o $(SERVER_BIN) ./cmd/server
	@echo "🔨 Building Client..."
	$(GO) build -ldflags "$(LDFLAGS)" -o $(CLIENT_BIN) ./cmd/client
	@echo "✅ Build completed!"

## 运行服务器
.PHONY: run-server
run-server: build
	@echo "🚀 Starting Server..."
	@$(SERVER_BIN) --config=config/config.yaml

## 运行客户端
.PHONY: run-client
run-client: build
	@echo "🚀 Starting Client..."
	@$(CLIENT_BIN) --help

## 运行测试
.PHONY: test
test:
	@echo "🧪 Running tests..."
	@$(GO) test -v ./...

## 代码格式化
.PHONY: fmt
fmt:
	@echo "🖌 Formatting code..."
	@$(GO) fmt ./...

## 代码静态检查
.PHONY: lint
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

## 清理构建的二进制文件
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

## 显示帮助信息
.PHONY: help
help:
	@echo "Golang 分布式缓存系统 - Makefile"
	@echo ""
	@echo "可用命令："
	@echo "  make build         构建服务端和客户端"
	@echo "  make run-server    运行服务器"
	@echo "  make run-client    运行客户端"
	@echo "  make test          运行单元测试"
	@echo "  make fmt           格式化代码"
	@echo "  make lint          运行代码静态检查"
	@echo "  make clean         清理编译产物"

