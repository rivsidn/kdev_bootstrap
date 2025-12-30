#!/bin/bash
# Ubuntu 22.04 系统配置脚本
# 在系统启动后运行此脚本完成网络、用户、SSH等配置

echo "Starting Ubuntu 22.04 system configuration..."

# 配置网络（使用 netplan）
setup_network() {
    echo "Configuring network..."

    # 创建 netplan 配置文件
    mkdir -p /etc/netplan
    cat > /etc/netplan/01-netcfg.yaml << 'EOF'
network:
  version: 2
  renderer: networkd
  ethernets:
    eth0:
      dhcp4: true
EOF

    # 配置DNS
    cat > /etc/resolv.conf << 'EOF'
# DNS configuration for QEMU
nameserver 10.0.2.3
nameserver 114.114.114.114
nameserver 8.8.8.8
EOF

    echo "Network configuration completed"
}

# 配置 root 无密码登录
setup_root_password() {
    echo "Configuring root user (no password)..."

    if [ -f /etc/shadow ]; then
        sed -i 's/^root:[^:]*:/root::/' /etc/shadow
        echo "Root password cleared in shadow file"
    else
        sed -i 's/^root:[^:]*:/root::/' /etc/passwd
        echo "Root password cleared in passwd file"
    fi

    echo "Root user configured for passwordless login"
}

# 配置 SSH 服务
setup_ssh() {
    echo "Configuring SSH for passwordless root login..."

    if [ -f /etc/ssh/sshd_config ]; then
        sed -i 's/^#*PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
        sed -i 's/^#*PermitEmptyPasswords.*/PermitEmptyPasswords yes/' /etc/ssh/sshd_config
        sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config
        sed -i 's/^#*UsePAM.*/UsePAM no/' /etc/ssh/sshd_config

        if ! grep -q "^PermitRootLogin" /etc/ssh/sshd_config; then
            echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
        fi
        if ! grep -q "^PermitEmptyPasswords" /etc/ssh/sshd_config; then
            echo "PermitEmptyPasswords yes" >> /etc/ssh/sshd_config
        fi
        if ! grep -q "^PasswordAuthentication" /etc/ssh/sshd_config; then
            echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config
        fi
        if ! grep -q "^UsePAM" /etc/ssh/sshd_config; then
            echo "UsePAM no" >> /etc/ssh/sshd_config
        fi

        echo "SSH configured for passwordless root login"
    else
        mkdir -p /etc/ssh
        cat > /etc/ssh/sshd_config << 'EOF'
# SSH Server Configuration
# Basic configuration with passwordless root login

Port 22
PermitRootLogin yes
PermitEmptyPasswords yes
PasswordAuthentication yes
UsePAM no
EOF
        echo "SSH basic configuration created"
    fi
}

# 执行所有配置
setup_network
setup_root_password
setup_ssh

echo ""
echo "Ubuntu 22.04 system configuration completed!"
echo "System is now configured with:"
echo "  - Network: DHCP enabled (will get IP 10.0.2.15 in QEMU)"
echo "  - Root login: no password required (just type 'root')"
echo "  - SSH: configured for passwordless root login"
echo ""
echo "You can delete this script: rm /root/setup.sh"
