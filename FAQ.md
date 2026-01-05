

## 密钥问题

### 问题现象

根文件系统构建过程中，可能会出现找不到密钥的问题.

```
  Including packages: make, gcc, build-essential, libncurses5-dev, libssl-dev, bc, flex, bison, libelf-dev, systemd, systemd-sysv, dbus, kmod, vim, iproute2, iputils-ping, net-tools,
  wget, curl, openssh-client, isc-dhcp-client, netplan.io, openssh-server, gdb, strace, git, python3, python
  Executing command: debootstrap --arch=amd64 --variant=buildd --components=main,universe --include=make,gcc,build-essential,libncurses5-dev,libssl-dev,bc,flex,bison,libelf-
  dev,systemd,systemd-sysv,dbus,kmod,vim,iproute2,iputils-ping,net-tools,wget,curl,openssh-client,isc-dhcp-client,netplan.io,openssh-server,gdb,strace,git,python3,python bionic ubuntu-
  18.04-amd64-bootfs http://mirrors.aliyun.com/ubuntu/
  I: Retrieving InRelease
  I: Checking Release signature
  E: Release signed by unknown key (key id 3B4FE6ACC0B21F32)
     The specified keyring /usr/share/keyrings/ubuntu-archive-removed-keys.gpg may be incorrect or out of date.
     You can find the latest Debian release key at https://ftp-master.debian.org/keys.html
```

### 解决方案

问题原因是debootstrap 找不到对应的密钥，需要安装密钥.

```bash
# 获取密钥
gpg --keyserver keyserver.ubuntu.com --recv-keys 3B4FE6ACC0B21F32
gpg --export 3B4FE6ACC0B21F32 > ubuntu-bionic-archive-keyring.gpg

# 安装密钥
sudo install -m 644 ubuntu-bionic-archive-keyring.gpg /usr/share/keyrings/ubuntu-archive-removed-keys.gpg

```
