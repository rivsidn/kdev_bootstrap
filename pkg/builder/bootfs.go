package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivsidn/kdev_bootstrap/pkg/config"
	"github.com/rivsidn/kdev_bootstrap/pkg/utils"
)

// BootfsBuilder bootfs 构建器
type BootfsBuilder struct {
	Config     *config.Config
	Arch       string
	OutputDir  string
	BootfsPath string
}

// NewBootfsBuilder 创建新的 bootfs 构建器
func NewBootfsBuilder(cfg *config.Config, arch, outputDir string) *BootfsBuilder {
	return &BootfsBuilder{
		Config:    cfg,
		Arch:      arch,
		OutputDir: outputDir,
	}
}

// Build 构建 bootfs
func (b *BootfsBuilder) Build() error {
	// 1. 检查环境
	if err := b.checkEnvironment(); err != nil {
		return err
	}

	// 2. 设置 bootfs 路径
	b.setBootfsPath()

	// 3. 检查是否已存在
	if utils.DirExists(b.BootfsPath) {
		fmt.Printf("Directory %s already exists\n", b.BootfsPath)
		if !utils.Confirm("Delete and recreate?") {
			return fmt.Errorf("operation cancelled by user")
		}
		if err := os.RemoveAll(b.BootfsPath); err != nil {
			return fmt.Errorf("failed to remove directory: %v", err)
		}
	}

	// 4. 创建目录
	if err := utils.CreateDir(b.BootfsPath); err != nil {
		return err
	}

	// 5. 执行 debootstrap（包含额外的包）
	if err := b.runDebootstrap(); err != nil {
		return err
	}

	// 6. 保存配置文件
	b.Config.ArchCurrent = b.Arch
	if err := b.Config.SaveToBootfs(b.BootfsPath); err != nil {
		return err
	}

	// 7. 配置网络
	if err := b.configureNetwork(); err != nil {
		return fmt.Errorf("failed to configure network: %v", err)
	}

	fmt.Printf("\nBootfs build successful: %s\n", b.BootfsPath)
	return nil
}

// checkEnvironment 检查环境
func (b *BootfsBuilder) checkEnvironment() error {
	// 检查是否为 root
	if !utils.CheckRoot() {
		return fmt.Errorf("please run with sudo or root privileges")
	}

	return nil
}

// setBootfsPath 设置 bootfs 路径
func (b *BootfsBuilder) setBootfsPath() {
	if b.OutputDir != "" {
		b.BootfsPath = b.OutputDir
		return
	}

	b.OutputDir = "."

	dirName := fmt.Sprintf("%s-%s-%s-bootfs",
		strings.ToLower(b.Config.Distribution),
		b.Config.Version,
		b.Arch)

	b.BootfsPath = filepath.Join(b.OutputDir, dirName)
}

// runDebootstrap 执行 debootstrap
func (b *BootfsBuilder) runDebootstrap() error {
	fmt.Println("\nRunning debootstrap...")

	suite := b.Config.GetSuite()
	if suite == "" {
		return fmt.Errorf("Not find the valid suite, add first")
	}

	mirror := b.Config.Mirror

	args := []string{
		"--arch=" + b.Arch,
		"--variant=buildd",
	}

	// 获取所有要安装的包
	packages := b.Config.GetAllPackages()
	if len(packages) > 0 {
		fmt.Printf("Including packages: %s\n", strings.Join(packages, ", "))
		args = append(args, "--include="+strings.Join(packages, ","))
	}

	// 对于旧版本 Ubuntu，添加特殊参数
	if strings.HasPrefix(b.Config.Version, "5.") {
		args = append(args, "--no-check-gpg")
	}

	args = append(args, suite, b.BootfsPath, mirror)

	if err := utils.RunCommand("debootstrap", args...); err != nil {
		return fmt.Errorf("debootstrap failed: %v", err)
	}

	return nil
}

// configureNetwork 配置网络
func (b *BootfsBuilder) configureNetwork() error {
	fmt.Println("Configuring network...")

	// 创建 /etc/network 目录
	networkDir := filepath.Join(b.BootfsPath, "etc", "network")
	if err := os.MkdirAll(networkDir, 0755); err != nil {
		return fmt.Errorf("failed to create network directory: %v", err)
	}

	// 创建 /etc/network/interfaces 文件
	interfacesPath := filepath.Join(networkDir, "interfaces")
	interfacesContent := `# interfaces(5) file used by ifup(8) and ifdown(8)
auto lo
iface lo inet loopback

# QEMU 网络配置 - 自动获取 IP
auto eth0
iface eth0 inet dhcp
`

	if err := os.WriteFile(interfacesPath, []byte(interfacesContent), 0644); err != nil {
		return fmt.Errorf("failed to create interfaces file: %v", err)
	}

	// 创建 /etc/resolv.conf 作为备份
	resolvPath := filepath.Join(b.BootfsPath, "etc", "resolv.conf")
	resolvContent := `# DNS configuration for QEMU
nameserver 10.0.2.3
nameserver 114.114.114.114
nameserver 8.8.8.8
`
	if err := os.WriteFile(resolvPath, []byte(resolvContent), 0644); err != nil {
		return fmt.Errorf("failed to create resolv.conf: %v", err)
	}

	fmt.Println("Network configuration completed")
	return nil
}

