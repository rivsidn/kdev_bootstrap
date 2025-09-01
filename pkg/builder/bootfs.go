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
	if b.outputDir != "" {
		b.BootfsPath = b.outputDir
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

// chrootRun 在 chroot 环境中运行命令
func (b *BootfsBuilder) chrootRun(name string, args ...string) error {
	chrootArgs := []string{b.BootfsPath, name}
	chrootArgs = append(chrootArgs, args...)
	return utils.RunCommand("chroot", chrootArgs...)
}
