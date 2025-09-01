# Makefile for kdev_bootstrap

# 变量定义
BINARY_DIR := ../bin
MODULE := github.com/kdev/bootstrap
VERSION := 1.0.0
BUILD_TIME := $(shell date +%F_%T)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go 命令
GO := go
GOFMT := gofmt
GOVET := go vet
GOTEST := go test

# 目标二进制文件
BINARIES := kboot_build_bootfs kboot_build_docker kboot_build_qemu

# 默认目标
.PHONY: all
all: clean build

# 构建所有二进制文件
.PHONY: build
build: deps
	@echo "🔨 构建二进制文件..."
	@mkdir -p $(BINARY_DIR)
	@for binary in $(BINARIES); do \
		echo "  构建 $$binary..."; \
		$(GO) build $(LDFLAGS) -o $(BINARY_DIR)/$$binary ./cmd/$$binary; \
	done
	@echo "✅ 构建完成！二进制文件位于 $(BINARY_DIR)/"

# 安装到系统
.PHONY: install
install: build
	@echo "📦 安装到 /usr/local/bin..."
	@for binary in $(BINARIES); do \
		sudo cp $(BINARY_DIR)/$$binary /usr/local/bin/; \
		sudo chmod +x /usr/local/bin/$$binary; \
		echo "  已安装 $$binary"; \
	done
	@echo "✅ 安装完成！"

# 卸载
.PHONY: uninstall
uninstall:
	@echo "🗑️  卸载..."
	@for binary in $(BINARIES); do \
		sudo rm -f /usr/local/bin/$$binary; \
		echo "  已删除 /usr/local/bin/$$binary"; \
	done
	@echo "✅ 卸载完成！"

# 安装依赖
.PHONY: deps
deps:
	@echo "📦 下载依赖..."
	@$(GO) mod download
	@$(GO) mod tidy

# 格式化代码
.PHONY: fmt
fmt:
	@echo "🎨 格式化代码..."
	@$(GOFMT) -w .

# 代码检查
.PHONY: vet
vet:
	@echo "🔍 检查代码..."
	@$(GOVET) ./...

# 运行测试
.PHONY: test
test:
	@echo "🧪 运行测试..."
	@$(GOTEST) -v ./...

# 清理
.PHONY: clean
clean:
	@echo "🧹 清理..."
	@rm -rf $(BINARY_DIR)
	@$(GO) clean
	@echo "✅ 清理完成！"

# 构建 Docker 镜像（用于开发）
.PHONY: docker-dev
docker-dev:
	@echo "🐳 构建开发 Docker 镜像..."
	@docker build -t kdev-bootstrap-dev:latest -f Dockerfile.dev .

# 运行示例
.PHONY: example
example: build
	@echo "📝 运行示例..."
	@echo "1. 构建 Ubuntu 22.04 bootfs:"
	@echo "   sudo $(BINARY_DIR)/kboot_build_bootfs -f configs/ubuntu-22.04.conf -a amd64"
	@echo ""
	@echo "2. 构建 Docker 镜像:"
	@echo "   sudo $(BINARY_DIR)/kboot_build_docker -b ubuntu-22.04-amd64-bootfs"
	@echo ""
	@echo "3. 构建 QEMU 镜像:"
	@echo "   sudo $(BINARY_DIR)/kboot_build_qemu -b ubuntu-22.04-amd64-bootfs"

# 显示帮助
.PHONY: help
help:
	@echo "kdev_bootstrap - 内核调试环境构建工具"
	@echo ""
	@echo "使用方法:"
	@echo "  make [目标]"
	@echo ""
	@echo "可用目标:"
	@echo "  all       - 清理并构建所有二进制文件（默认）"
	@echo "  build     - 构建所有二进制文件"
	@echo "  install   - 安装到 /usr/local/bin"
	@echo "  uninstall - 从系统卸载"
	@echo "  deps      - 下载并整理依赖"
	@echo "  fmt       - 格式化代码"
	@echo "  vet       - 静态代码检查"
	@echo "  test      - 运行测试"
	@echo "  clean     - 清理构建产物"
	@echo "  example   - 显示使用示例"
	@echo "  help      - 显示此帮助信息"

# 设置默认目标
.DEFAULT_GOAL := help