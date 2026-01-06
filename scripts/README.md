
## 调试脚本

集合了docker、qemu 启动命令，编译、调试Linux 内核时使用.


### 编译

```bash
source kdev_env

kdev_build
```

### 运行


```bash
source kdev_env

kdev_run terminal
```

| 模式     | 说明                                       |
|----------|--------------------------------------------|
| terminal | 终端启动                                   |
| debug    | 调试模式，需要链接GDB 设备才能正常启动     |
| internet | 网络模式，可以访问互联网                   |
| bridge   | 桥模式，不仅可以访问互联网，还可以与PC互联 |


### 文件传输

```bash

# 默认将文件传输到/root/

kdev_push modules.ko 

```

## 附录

### 常见问题

#### terminal模式模块卸载问题

现象.

```
root@(none):/root# rmmod waitqueue_sample
rmmod: ERROR: ../libkmod/libkmod.c:514 lookup_builtin_file() could not open builtin file '/lib/modules/4.4.115/modules.builtin.bin'
rmmod: ERROR: Module waitqueue_sample is not currently loaded
```

**terminal模式启动调试内核的时候，需要手动挂载proc、sys 文件系统.**

```bash
mount -t proc  none /proc
mount -t sysfs none /sys
```

#### bridge模式网络访问

如何设置访问互联网.

```bash
#  虚拟机设置
## 网口link
ip link set eth0 up
## 设置IP地址
ip addr add 172.20.0.2/24 dev eth0
## 添加默认路由
ip route add default via 172.20.0.1

# 宿主机设置
sudo iptables -t nat -A POSTROUTING -s 172.20.0.0/24 ! -o tap0 -j MASQUERADE

```

### 地址说明

  QEMU User 模式的固定网络拓扑.

  | IP 地址     | 角色          | 说明                      |
  |-------------|---------------|---------------------------|
  | 10.0.2.0/24 | 网络段        | User 模式的虚拟网络       |
  | 10.0.2.2    | 网关/宿主机   | 虚拟机访问宿主机用这个 IP |
  | 10.0.2.3    | DNS 服务器    | QEMU 内置的 DNS 转发服务  |
  | 10.0.2.4    | SMB 服务器    | 文件共享（如果启用）      |
  | 10.0.2.15   | 虚拟机默认 IP | DHCP 分配的第一个 IP      |

这些 IP 地址是 QEMU User 模式网络的固定内置地址，不是随意设置的，而是 QEMU 硬编码的默认值。

