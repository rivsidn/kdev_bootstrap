#!/bin/bash
# Ubuntu 16.04 系统配置脚本
# 在系统启动后运行此脚本完成网络、用户、SSH等配置

echo "Starting Ubuntu 16.04 system configuration..."

# 配置网络
setup_network() {
    echo "Configuring network..."
    
    # 创建网络配置文件
    cat > /etc/network/interfaces << 'EOF'
# interfaces(5) file used by ifup(8) and ifdown(8)
auto lo
iface lo inet loopback

# QEMU 网络配置 - 自动获取 IP
auto eth0
iface eth0 inet dhcp
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
    echo "root:passwd" | chpasswd
    sync
}

# 配置 SSH 服务
setup_ssh() {
    echo "Configuring SSH for passwordless root login..."
    
    # 检查 SSH 配置文件是否存在
    if [ -f /etc/ssh/sshd_config ]; then
        # 修改 SSH 配置允许 root 登录和空密码
        sed -i 's/^#*PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config
        sed -i 's/^#*PermitEmptyPasswords.*/PermitEmptyPasswords yes/' /etc/ssh/sshd_config
        sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication yes/' /etc/ssh/sshd_config
        sed -i 's/^#*UsePAM.*/UsePAM no/' /etc/ssh/sshd_config
        
        # 如果配置项不存在，添加它们
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
        # 如果 SSH 配置文件不存在，创建基本配置
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
echo "Ubuntu 16.04 system configuration completed!"
echo "System is now configured with:"
echo "  - Network: DHCP enabled (will get IP 10.0.2.15 in QEMU)"
echo "  - Root login: no password required (just type 'root')"
echo "  - SSH: configured for passwordless root login"
echo ""
echo "You can delete this script: rm /root/setup.sh"
