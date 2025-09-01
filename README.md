# kdev_bootstrap - 内核调试环境构建工具

基于《构建调试环境.md》设计文档的 Go 语言实现，用于快速构建内核编译和调试环境。

## 功能特性

- 🚀 **一键构建**：通过配置文件自动构建完整的调试环境
- 🐳 **Docker 支持**：生成用于内核编译的 Docker 镜像
- 💿 **QEMU 支持**：生成用于内核调试的 QEMU 磁盘镜像
- 🔧 **多架构支持**：支持 i386、amd64、arm64 等架构
- 📦 **版本管理**：支持多个 Ubuntu 版本（5.10、16.04、20.04、22.04）

## 系统架构

```
配置文件 (.conf)
    ↓
kboot_build_bootfs (构建根文件系统)
    ↓
    ├── kboot_build_docker (生成 Docker 镜像)
    └── kboot_build_qemu (生成 QEMU 镜像)
```

## 快速开始

### 安装

```bash
cd src
chmod +x build.sh install.sh
./build.sh          # 构建工具
sudo ./install.sh   # 安装到系统
```

### 使用示例

1. **构建 Ubuntu 22.04 调试环境**

```bash
# 构建根文件系统
sudo kboot bootfs -f /usr/local/share/kdev-bootstrap/configs/ubuntu-22.04.conf -a amd64

# 构建 Docker 镜像（用于编译）
sudo kboot docker -b ubuntu-22.04-amd64-bootfs

# 构建 QEMU 镜像（用于调试）
sudo kboot qemu -b ubuntu-22.04-amd64-bootfs
```

2. **使用 Docker 镜像编译内核**

```bash
docker run -it --rm -v /path/to/kernel:/kernel ubuntu-22.04-amd64 /bin/bash
# 在容器内编译内核
cd /kernel
make menuconfig
make -j$(nproc)
```

3. **使用 QEMU 调试内核**

```bash
qemu-system-x86_64 \
    -kernel /path/to/bzImage \
    -hda ubuntu-22.04-amd64-rootfs.img \
    -append "root=/dev/sda rw console=ttyS0" \
    -m 2048 \
    -enable-kvm \
    -nographic
```

## 命令详解

### kboot_build_bootfs

构建根文件系统，使用 debootstrap 创建最小化 Ubuntu 环境。

```bash
kboot_build_bootfs [选项]
  -f, --file FILE    配置文件路径（必需）
  -a, --arch ARCH    目标架构（i386/amd64/arm64）
  -o, --output DIR   输出目录
  -h, --help         显示帮助
```

### kboot_build_docker

从根文件系统构建 Docker 镜像。

```bash
kboot_build_docker [选项]
  -b, --bootfs DIR        根文件系统路径（必需）
  -f, --dockerfile FILE   Dockerfile 路径（可选）
  --image NAME:TAG        镜像名称（可选）
  -h, --help             显示帮助
```

### kboot_build_qemu

从根文件系统构建 QEMU 磁盘镜像。

```bash
kboot_build_qemu [选项]
  -b, --bootfs DIR     根文件系统路径（必需）
  -r, --rootfs FILE    输出镜像名称（可选）
  -s, --size SIZE      镜像大小（默认 2G）
  -h, --help          显示帮助
```

## 配置文件格式

```ini
[ubuntu-22.04]
distribution = ubuntu
version = 22.04
arch_supported = amd64,arm64
mirror = http://mirrors.aliyun.com/ubuntu/

# 软件包组
kbuild_packages = make,gcc,build-essential,libncurses-dev
module_packages = kmod
debug_packages = gdb,strace
network_packages = wget,curl,openssh-client
```

## 目录结构

```
src/
├── cmd/                    # 命令行工具
│   ├── kboot_build_bootfs/
│   ├── kboot_build_docker/
│   └── kboot_build_qemu/
├── pkg/                    # 核心库
│   ├── config/            # 配置解析
│   ├── builder/           # 构建器实现
│   └── utils/             # 工具函数
├── configs/               # 示例配置文件
├── build.sh              # 构建脚本
├── install.sh            # 安装脚本
└── Makefile              # Make 构建文件
```

## 系统要求

- Go 1.21 或更高版本
- Ubuntu/Debian 系统
- root 权限（用于 debootstrap 和挂载操作）
- 依赖工具：debootstrap、docker、qemu-utils

## 开发

```bash
# 获取代码
git clone <repository>
cd kdev_bootstrap/src

# 构建
make build

# 运行测试
make test

# 安装到系统
sudo make install

# 清理
make clean
```

## 许可证

MIT License
