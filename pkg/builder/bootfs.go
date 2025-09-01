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

	// 5. 执行 debootstrap
	if err := b.runDebootstrap(); err != nil {
		return err
	}

	// 6. 安装额外的包
	if err := b.installPackages(); err != nil {
		return err
	}

	// 7. 配置系统
	if err := b.configureSystem(); err != nil {
		return err
	}

	// 8. 保存配置文件
	b.Config.ArchCurrent = b.Arch
	if err := b.Config.SaveToBootfs(b.BootfsPath); err != nil {
		return err
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

	// 检查依赖
	deps := []string{"debootstrap"}
	if err := utils.CheckDependencies(deps); err != nil {
		return err
	}

	// 验证架构
	if !b.Config.ValidateArch(b.Arch) {
		return fmt.Errorf("unsupported architecture: %s, supported architectures: %s", 
			b.Arch, strings.Join(b.Config.ArchSupported, ", "))
	}

	return nil
}

// setBootfsPath 设置 bootfs 路径
func (b *BootfsBuilder) setBootfsPath() {
	if b.OutputDir == "" {
		b.OutputDir = "."
	}
	
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
	mirror := b.Config.Mirror
	
	args := []string{
		"--arch=" + b.Arch,
		"--variant=buildd",
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

// installPackages 安装额外的包
func (b *BootfsBuilder) installPackages() error {
	packages := b.Config.GetAllPackages()
	if len(packages) == 0 {
		return nil
	}
	
	fmt.Printf("\nInstalling additional packages: %s\n", strings.Join(packages, ", "))
	
	// 更新包列表
	if err := b.chrootRun("apt-get", "update"); err != nil {
		fmt.Printf("Failed to update package list, continuing installation...\n")
	}
	
	// 安装包
	args := []string{"install", "-y", "--no-install-recommends"}
	args = append(args, packages...)
	
	if err := b.chrootRun("apt-get", args...); err != nil {
		fmt.Printf("Some packages failed to install: %v\n", err)
	}
	
	// 清理
	b.chrootRun("apt-get", "clean")
	
	return nil
}

// configureSystem 配置系统
func (b *BootfsBuilder) configureSystem() error {
	fmt.Println("\nConfiguring system...")
	
	// 设置 hostname
	hostnamePath := filepath.Join(b.BootfsPath, "etc", "hostname")
	hostname := fmt.Sprintf("%s-%s", b.Config.Distribution, b.Config.Version)
	if err := os.WriteFile(hostnamePath, []byte(hostname+"\n"), 0644); err != nil {
		fmt.Printf("Failed to set hostname: %v\n", err)
	}
	
	// 设置 hosts
	hostsPath := filepath.Join(b.BootfsPath, "etc", "hosts")
	hostsContent := fmt.Sprintf(`127.0.0.1	localhost
127.0.1.1	%s

# IPv6
::1		localhost ip6-localhost ip6-loopback
ff02::1		ip6-allnodes
ff02::2		ip6-allrouters
`, hostname)
	if err := os.WriteFile(hostsPath, []byte(hostsContent), 0644); err != nil {
		fmt.Printf("Failed to set hosts: %v\n", err)
	}
	
	// 设置 root 密码为空（用于开发环境）
	b.chrootRun("passwd", "-d", "root")
	
	// 设置 DNS
	resolvPath := filepath.Join(b.BootfsPath, "etc", "resolv.conf")
	resolvContent := `nameserver 8.8.8.8
nameserver 8.8.4.4
`
	if err := os.WriteFile(resolvPath, []byte(resolvContent), 0644); err != nil {
		fmt.Printf("Failed to set DNS: %v\n", err)
	}
	
	// 创建必要的目录
	dirs := []string{
		filepath.Join(b.BootfsPath, "root", ".ssh"),
		filepath.Join(b.BootfsPath, "var", "log"),
		filepath.Join(b.BootfsPath, "tmp"),
	}
	for _, dir := range dirs {
		utils.CreateDir(dir)
	}
	
	return nil
}

// chrootRun 在 chroot 环境中运行命令
func (b *BootfsBuilder) chrootRun(name string, args ...string) error {
	chrootArgs := []string{b.BootfsPath, name}
	chrootArgs = append(chrootArgs, args...)
	return utils.RunCommand("chroot", chrootArgs...)
}
