#!/bin/bash

# 安装脚本 - 安装 kdev_bootstrap 工具到系统

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 安装目录
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/usr/local/share/kdev-bootstrap/configs"

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

# 检查权限
check_permission() {
    if [[ $EUID -ne 0 ]]; then
        print_error "此脚本需要 root 权限运行"
        echo "请使用: sudo ./install.sh"
        exit 1
    fi
}

# 检查依赖
check_dependencies() {
    print_info "检查系统依赖..."
    
    local deps=("debootstrap" "docker" "qemu-img")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            missing+=($dep)
            print_warn "  缺少 $dep"
        else
            print_info "  ✓ $dep"
        fi
    done
    
    if [ ${#missing[@]} -gt 0 ]; then
        print_warn "缺少一些依赖，某些功能可能无法使用"
        print_info "安装建议:"
        echo "  sudo apt-get update"
        echo "  sudo apt-get install -y debootstrap docker.io qemu-utils"
        echo ""
        read -p "是否继续安装? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

# 构建工具
build_tools() {
    print_info "构建工具..."
    
    # 检查是否已构建
    if [ ! -d "../bin" ]; then
        ./build.sh
    else
        print_info "检测到已构建的二进制文件"
    fi
}

# 安装二进制文件
install_binaries() {
    print_info "安装二进制文件到 $INSTALL_DIR..."
    
    local binaries=(
        "kboot_build_bootfs"
        "kboot_build_docker"
        "kboot_build_qemu"
    )
    
    for binary in "${binaries[@]}"; do
        if [ -f "../bin/$binary" ]; then
            cp "../bin/$binary" "$INSTALL_DIR/"
            chmod +x "$INSTALL_DIR/$binary"
            print_info "  已安装 $binary"
        else
            print_error "  找不到 $binary"
        fi
    done
}

# 安装配置文件
install_configs() {
    print_info "安装配置文件到 $CONFIG_DIR..."
    
    # 创建配置目录
    mkdir -p "$CONFIG_DIR"
    
    # 复制配置文件
    if [ -d "configs" ]; then
        cp configs/*.conf "$CONFIG_DIR/" 2>/dev/null || true
        print_info "  已安装配置文件"
    else
        print_warn "  未找到配置文件目录"
    fi
}

# 创建命令别名
create_aliases() {
    print_info "创建命令别名..."
    
    # 创建简短别名脚本
    cat > "$INSTALL_DIR/kboot" << 'EOF'
#!/bin/bash

# kboot - kdev_bootstrap 工具的统一入口

case "$1" in
    bootfs)
        shift
        exec kboot_build_bootfs "$@"
        ;;
    docker)
        shift
        exec kboot_build_docker "$@"
        ;;
    qemu)
        shift
        exec kboot_build_qemu "$@"
        ;;
    *)
        echo "kboot - 内核调试环境构建工具"
        echo ""
        echo "使用方法:"
        echo "  kboot <命令> [选项]"
        echo ""
        echo "可用命令:"
        echo "  bootfs  - 构建根文件系统"
        echo "  docker  - 构建 Docker 镜像"
        echo "  qemu    - 构建 QEMU 镜像"
        echo ""
        echo "示例:"
        echo "  sudo kboot bootfs -f ubuntu-22.04.conf -a amd64"
        echo "  sudo kboot docker -b ubuntu-22.04-amd64-bootfs"
        echo "  sudo kboot qemu -b ubuntu-22.04-amd64-bootfs"
        exit 1
        ;;
esac
EOF
    
    chmod +x "$INSTALL_DIR/kboot"
    print_info "  已创建 kboot 命令"
}

# 显示安装信息
show_info() {
    echo ""
    echo -e "${BLUE}=================================${NC}"
    echo -e "${GREEN}✅ kdev_bootstrap 安装成功！${NC}"
    echo -e "${BLUE}=================================${NC}"
    echo ""
    print_info "已安装的工具:"
    echo "  - kboot_build_bootfs : 构建根文件系统"
    echo "  - kboot_build_docker : 构建 Docker 镜像"
    echo "  - kboot_build_qemu   : 构建 QEMU 镜像"
    echo "  - kboot              : 统一命令入口"
    echo ""
    print_info "配置文件位置:"
    echo "  $CONFIG_DIR"
    echo ""
    print_info "快速开始:"
    echo "  1. 查看帮助:"
    echo "     kboot"
    echo ""
    echo "  2. 构建 Ubuntu 22.04 环境:"
    echo "     sudo kboot bootfs -f $CONFIG_DIR/ubuntu-22.04.conf -a amd64"
    echo "     sudo kboot docker -b ubuntu-22.04-amd64-bootfs"
    echo "     sudo kboot qemu -b ubuntu-22.04-amd64-bootfs"
}

# 主函数
main() {
    print_info "开始安装 kdev_bootstrap..."
    
    # 检查权限
    check_permission
    
    # 检查依赖
    check_dependencies
    
    # 构建工具
    build_tools
    
    # 安装二进制文件
    install_binaries
    
    # 安装配置文件
    install_configs
    
    # 创建命令别名
    create_aliases
    
    # 显示安装信息
    show_info
}

# 执行主函数
main "$@"
