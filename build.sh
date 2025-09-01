#!/bin/bash

# 构建脚本 - 快速构建所有工具

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Go 环境
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go 未安装，请先安装 Go 语言环境"
        echo "访问 https://golang.org/dl/ 下载安装"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    print_info "检测到 Go 版本: $GO_VERSION"
}

# 主函数
main() {
    print_info "开始构建 kdev_bootstrap 工具..."
    
    # 检查环境
    check_go
    
    # 创建 bin 目录
    mkdir -p bin
    
    # 下载依赖
    print_info "下载依赖包..."
    go mod download
    go mod tidy
    
    # 构建二进制文件
    print_info "构建二进制文件..."
    
    BINARIES=(
        "kboot_build_bootfs"
        "kboot_build_docker"
        "kboot_build_qemu"
    )
    
    for binary in "${BINARIES[@]}"; do
        print_info "  构建 $binary..."
        go build -o bin/$binary ./cmd/$binary
        chmod +x bin/$binary
    done
    
    print_info "✅ 构建完成！"
    echo ""
    print_info "二进制文件位于: bin/"
    echo ""
    print_info "使用示例:"
    echo "  1. 构建 bootfs:"
    echo "     sudo bin/kboot_build_bootfs -f configs/ubuntu-22.04.conf -a amd64"
    echo ""
    echo "  2. 构建 Docker 镜像:"
    echo "     sudo bin/kboot_build_docker -b ubuntu-22.04-amd64-bootfs"
    echo ""
    echo "  3. 构建 QEMU 镜像:"
    echo "     sudo bin/kboot_build_qemu -b ubuntu-22.04-amd64-bootfs"
}

# 执行主函数
main "$@"
