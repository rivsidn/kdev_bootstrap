# Makefile for kdev_bootstrap

# 变量定义
BINARY_DIR := bin

# Go 命令
GO := go

# 目标二进制文件
BINARIES := kboot_build_bootfs kboot_build_docker kboot_build_qemu

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help:
	@echo "kdev_bootstrap - Linux kernel development environment builder"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  build - Build all binaries"
	@echo "  clean - Remove build artifacts"
	@echo "  help  - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make clean"
	@echo "  make build"

# 构建所有二进制文件
.PHONY: build
build:
	@echo "Checking Go installation..."
	@which $(GO) > /dev/null || (echo "error: Go not installed" >&2 && exit 1)
	@echo "Creating binary directory..."
	@mkdir -p $(BINARY_DIR)
	@echo "Preventing go mod tidy from scanning bootfs directory..."
	@echo "module ignore" > $(BINARY_DIR)/go.mod
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "Building binaries..."
	@for binary in $(BINARIES); do \
		echo "  Building $$binary..."; \
		$(GO) build -o $(BINARY_DIR)/$$binary ./cmd/$$binary || exit 1; \
	done
	@echo "Build completed. Binaries are in $(BINARY_DIR)/"

# 清理
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@for binary in $(BINARIES); do \
		rm -f $(BINARY_DIR)/$$binary; \
	done
	@$(GO) clean
	@echo "Clean completed"
