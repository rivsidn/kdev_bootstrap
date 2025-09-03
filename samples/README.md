

## QEMU 启动网络模式

### QMEU USER 模式

用于访问互联网.

#### 启动命令

```bash
			qemu-system-x86_64 -kernel $KDEV_KERNEL -hda $KDEV_QEMU_IMAGE \
				-netdev user,id=net0,hostfwd=tcp::2222-:22,hostfwd=tcp::8080-:80 \
				-device e1000,netdev=net0 \
				-append "root=/dev/sda rw init=/sbin/init console=tty0" -m 256M
```

#### 参数解析

| 参数    | 解析                       |
|---------|----------------------------|
| hostfwd | 将宿主机的端口映射到虚拟机 |


#### 虚拟机配置

**手动配置**

```bash
# IP地址设置
ifconfig eth0 10.0.2.15 netmask 255.255.255.0

# 添加默认网关
route add default gw 10.0.2.2

# 配置 DNS
echo "nameserver 10.0.2.3" > /etc/resolv.conf
```

****

```bash
```



## 附录

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

