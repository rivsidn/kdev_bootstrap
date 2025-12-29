# kdev_bootstrap

用于构建内核调试环境的一套工具，参见 [基本原理](doc/基本原理.md).


## 工程编译

```bash
# 编译
chmod +x build.sh; bash build.sh
```

## 使用示例

以ubuntu-16.04 为例子.

### 配置文件

```ini
[ubuntu-16.04]

# 发行版信息
distribution = ubuntu
version = 16.04
arch_supported = i386,amd64

# 镜像源(使用阿里云镜像)
mirror = http://mirrors.aliyun.com/ubuntu/

# 系统配置脚本(需要指定路径)
setup_script = ubuntu-16.04-setup.sh

# 安装包
kbuild_packages = make,gcc
```

### 构建系统

```bash
# 构建根文件系统
sudo ./kboot_build_bootfs -a amd64 -f ../configs/ubuntu-16.04.conf
# 创建docker 镜像
sudo ./kboot_build_docker -b ubuntu-16.04-amd64-bootfs/
# 创建qemu 根文件系统
sudo ./kboot_build_qemu -b ubuntu-16.04-amd64-bootfs/

```

## 代码调试

通过docker 镜像编译，通过qemu 调试内核.

### 内核编译

```bash
docker run -it --rm -v $KERNEL_PATH:/workspace -w /workspace --hostname kernel-dev $KDEV_DOCKER_IMAGE bash
```

| 变量              | 说明           |
|-------------------|----------------|
| KERNEL_PATH       | 内核代码路径   |
| KDEV_DOCKER_IMAGE | docker镜像名称 |

### 内核调试

通过`qemu` 调试linux 内核.

## 附录

### 目录结构


```
src/
├── cmd/                    # 命令行工具
│   ├── kboot_build_bootfs/
│   ├── kboot_build_docker/
│   └── kboot_build_qemu/
├── pkg/                    # 核心库
│   ├── config/             # 配置解析
│   ├── builder/            # 构建器实现
│   └── utils/              # 工具函数
├── configs/                # 示例配置文件
├── samples/                # 内核调试脚本示例
├── build.sh                # 构建脚本
└── Makefile                # Make 构建文件
```

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

