# kdev_bootstrap - 内核调试环境构建工具

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

## 附录

### 地址

### Ubuntu版本内核对应

| Ubuntu版本 | 代号              | 内核版本号 |
| ---------- | ---------------   | ---------- |
| 23.10      | Mantic Minotaur   | 6.5        |
| 23.04      | Lunar Lobster     | 6.2        |
| 22.10      | Kinetic Kudu      | 5.19       |
| 22.04      | Jammy Jellyfish   | 5.15       |
| 21.10      | Impish Indri      | 5.13       |
| 21.04      | Hirsute Hippo     | 5.11       |
| 20.10      | Groovy Gorilla    | 5.8        |
| 20.04      | Focal Fossa       | 5.4        |
| 19.10      | Eoan Ermine       | 5.3        |
| 19.04      | Disco Dingo       | 5.0        |
| 18.10      | Cosmic Cuttlefish | 4.18       |
| 18.04      | Bionic Beaver     | 4.15       |
| 17.10      | Artful Aardvark   | 4.13       |
| 17.04      | Zesty Zapus       | 4.10       |
| 16.10      | Yakkety Yak       | 4.8        |
| 16.04      | Xenial Xerus      | 4.4        |
| 15.10      | Wily Werewolf     | 4.2        |
| 15.04      | Vivid Vervet      | 3.19       |
| 14.10      | Utopic Unicorn    | 3.16       |
| 14.04      | Trusty Tahr       | 3.13       |
| 13.10      | Saucy Salamander  | 3.11       |
| 13.04      | Raring Ringtail   | 3.8        |
| 12.10      | Quantal Quetzal   | 3.5        |
| 12.04      | Precise Pangolin  | 3.2+       |
| 11.10      | Oneiric Ocelot    | 3.0        |
| 11.04      | Natty Narwhal     | 2.6.38     |
| 10.10      | Maverick Meerkat  | 2.6.35     |
| 10.04      | Lucid Lynx        | 2.6.32     |
| 09.10      | Karmic Koala      | 2.6.31     |
| 09.04      | Jaunty Jackalope  | 2.6.28     |
| 08.10      | Intrepid Ibex     | 2.6.27     |
| 08.04      | Hardy Heron       | 2.6.24     |
| 07.10      | Gutsy Gibbon      | 2.6.22     |
| 07.04      | Feisty Fawn       | 2.6.20     |
| 06.10      | Edgy Eft          | 2.6.17     |
| 06.06      | Dapper Drake      | 2.6.15     |
| 05.10      | Breezy Badger     | 2.6.12     |
| 05.04      | Hoary Hedgehog    | 2.6.10     |
| 04.10      | Warty Warthog     | 2.6.8      |

## 许可证

MIT License

