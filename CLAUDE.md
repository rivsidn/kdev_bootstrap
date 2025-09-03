# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

kdev_bootstrap 是一个Linux内核开发环境构建工具，用于自动化创建内核编译的Docker容器和内核调试的QEMU虚拟机环境。

## 常用命令

### 构建项目
```bash
# 构建所有二进制文件
make build

# 快速构建（使用build.sh）
./build.sh

# 安装到系统
make install

# 运行测试（注意：目前无测试文件）
make test
```

### 使用示例
```bash
# 构建根文件系统
./bin/kboot_build_bootfs -c configs/ubuntu-22.04.conf -o /tmp/rootfs

# 从根文件系统构建Docker镜像
./bin/kboot_build_docker -r /tmp/rootfs -t myimage:latest

# 构建QEMU磁盘镜像
./bin/kboot_build_qemu -r /tmp/rootfs -o /tmp/disk.img
```

## 代码架构

### 核心模块结构
```
pkg/
├── config/       # 配置管理
│   ├── config.go    # 配置文件解析（INI格式）
│   └── suites.go    # 软件包套件管理
└── builder/      # 构建器实现
    ├── bootfs.go    # 根文件系统构建（debootstrap）
    ├── docker.go    # Docker镜像构建
    └── qemu.go      # QEMU磁盘镜像构建

cmd/              # 命令行工具
├── kboot_build_bootfs/   # 构建根文件系统
├── kboot_build_docker/   # 构建Docker镜像
└── kboot_build_qemu/     # 构建QEMU镜像
```

### 构建流程
1. 读取配置文件（configs/*.conf）定义发行版、架构、软件包
2. 使用debootstrap创建最小化根文件系统
3. 基于根文件系统构建Docker镜像或QEMU磁盘镜像

### 关键设计模式
- 配置驱动：通过INI配置文件控制构建参数
- 分层构建：根文件系统作为Docker和QEMU镜像的共同基础
- 模块化软件包管理：内核构建包、调试工具、网络工具等分组

## 开发注意事项

### 依赖要求
- Go 1.21+
- 系统工具：debootstrap, docker, qemu-utils
- 需要sudo权限（debootstrap和chroot操作）

### 配置文件格式
配置文件使用INI格式，包含以下关键节：
- `[distro]`: 发行版信息（名称、版本、架构）
- `[mirror]`: APT镜像源配置
- `[packages.*]`: 软件包组定义（kernel、debug、network等）

### 已知问题（来自TODO.md）
- Docker镜像构建后标签问题
- 需要sudo权限运行的安全性考虑
- rsync vs cp的使用差异
- 架构相关的编译问题

### 代码规范遵循
- Go代码：标准Go格式，注释用中文
- Bash脚本：尽早返回原则，避免嵌套，仅在失败时输出错误
