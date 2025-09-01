package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rivsidn/kdev_bootstrap/pkg/config"
	"github.com/rivsidn/kdev_bootstrap/pkg/utils"
)

// QemuBuilder QEMU 镜像构建器
type QemuBuilder struct {
	Config      *config.Config
	BootfsPath  string
	RootfsImage string
	ImageSize   string
}

// NewQemuBuilder 创建新的 QEMU 构建器
func NewQemuBuilder(bootfsPath string, rootfsImage string, imageSize string) (*QemuBuilder, error) {
	// 加载配置文件
	configPath := filepath.Join(bootfsPath, "etc", "bootstrap.conf")
	if !utils.FileExists(configPath) {
		return nil, fmt.Errorf("找不到配置文件: %s", configPath)
	}
	
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	
	// 设置默认大小
	if imageSize == "" {
		imageSize = "2G"
	}
	
	return &QemuBuilder{
		Config:      cfg,
		BootfsPath:  bootfsPath,
		RootfsImage: rootfsImage,
		ImageSize:   imageSize,
	}, nil
}

// Build 构建 QEMU 镜像
func (b *QemuBuilder) Build() error {
	// 1. 检查环境
	if err := b.checkEnvironment(); err != nil {
		return err
	}
	
	// 2. 设置镜像名称
	if b.RootfsImage == "" {
		arch := b.Config.ArchCurrent
		if arch == "" {
			arch = b.inferArch()
		}
		b.RootfsImage = b.Config.GetRootfsName(arch)
	}
	
	// 3. 创建镜像文件
	if err := b.createImage(); err != nil {
		return err
	}
	
	// 4. 格式化镜像
	if err := b.formatImage(); err != nil {
		return err
	}
	
	// 5. 挂载镜像
	mountPoint, err := b.mountImage()
	if err != nil {
		return err
	}
	defer b.unmountImage(mountPoint)
	
	// 6. 复制 rootfs
	if err := b.copyRootfs(mountPoint); err != nil {
		return err
	}
	
	// 7. 安装 bootloader（可选）
	b.installBootloader(mountPoint)
	
	fmt.Printf("\n✅ QEMU 镜像构建成功: %s\n", b.RootfsImage)
	fmt.Printf("   大小: %s\n", b.ImageSize)
	fmt.Printf("   使用方法:\n")
	fmt.Printf("   qemu-system-x86_64 -hda %s -m 1024 -enable-kvm\n", b.RootfsImage)
	
	return nil
}

// checkEnvironment 检查环境
func (b *QemuBuilder) checkEnvironment() error {
	// 检查是否为 root
	if !utils.CheckRoot() {
		return fmt.Errorf("请使用 sudo 或 root 权限运行")
	}
	
	// 检查必要的工具
	deps := []string{"qemu-img", "mkfs.ext3", "mount", "umount"}
	if err := utils.CheckDependencies(deps); err != nil {
		return err
	}
	
	// 检查 bootfs 目录
	if !utils.DirExists(b.BootfsPath) {
		return fmt.Errorf("bootfs 目录不存在: %s", b.BootfsPath)
	}
	
	return nil
}

// inferArch 从路径推断架构
func (b *QemuBuilder) inferArch() string {
	base := filepath.Base(b.BootfsPath)
	if strings.Contains(base, "i386") {
		return "i386"
	}
	if strings.Contains(base, "amd64") {
		return "amd64"
	}
	return "amd64"
}

// createImage 创建镜像文件
func (b *QemuBuilder) createImage() error {
	// 检查镜像是否已存在
	if utils.FileExists(b.RootfsImage) {
		fmt.Printf("镜像文件 %s 已存在\n", b.RootfsImage)
		if !utils.Confirm("是否删除并重新创建？") {
			return fmt.Errorf("用户取消操作")
		}
		if err := os.Remove(b.RootfsImage); err != nil {
			return fmt.Errorf("删除镜像失败: %v", err)
		}
	}
	
	fmt.Printf("\n💾 创建镜像文件: %s (大小: %s)\n", b.RootfsImage, b.ImageSize)
	
	args := []string{
		"create",
		"-f", "raw",
		b.RootfsImage,
		b.ImageSize,
	}
	
	if err := utils.RunCommand("qemu-img", args...); err != nil {
		return fmt.Errorf("创建镜像失败: %v", err)
	}
	
	return nil
}

// formatImage 格式化镜像
func (b *QemuBuilder) formatImage() error {
	fmt.Println("🔧 格式化镜像为 ext3...")
	
	// 创建 loop 设备
	output, err := utils.RunCommandOutput("losetup", "-f")
	if err != nil {
		return fmt.Errorf("获取空闲 loop 设备失败: %v", err)
	}
	loopDevice := strings.TrimSpace(output)
	
	// 关联镜像到 loop 设备
	if err := utils.RunCommand("losetup", loopDevice, b.RootfsImage); err != nil {
		return fmt.Errorf("关联 loop 设备失败: %v", err)
	}
	defer utils.RunCommand("losetup", "-d", loopDevice)
	
	// 格式化为 ext3
	if err := utils.RunCommand("mkfs.ext3", "-F", loopDevice); err != nil {
		return fmt.Errorf("格式化失败: %v", err)
	}
	
	return nil
}

// mountImage 挂载镜像
func (b *QemuBuilder) mountImage() (string, error) {
	fmt.Println("📁 挂载镜像...")
	
	// 创建临时挂载点
	mountPoint := fmt.Sprintf("/tmp/qemu-mount-%d", os.Getpid())
	if err := utils.CreateDir(mountPoint); err != nil {
		return "", err
	}
	
	// 挂载镜像
	if err := utils.RunCommand("mount", "-o", "loop", b.RootfsImage, mountPoint); err != nil {
		os.RemoveAll(mountPoint)
		return "", fmt.Errorf("挂载镜像失败: %v", err)
	}
	
	return mountPoint, nil
}

// unmountImage 卸载镜像
func (b *QemuBuilder) unmountImage(mountPoint string) {
	fmt.Println("📁 卸载镜像...")
	utils.RunCommand("umount", mountPoint)
	os.RemoveAll(mountPoint)
}

// copyRootfs 复制根文件系统
func (b *QemuBuilder) copyRootfs(mountPoint string) error {
	fmt.Printf("📋 复制根文件系统到镜像...\n")
	
	// 使用 rsync 或 cp 复制文件
	if utils.CheckCommand("rsync") {
		args := []string{
			"-av",
			"--exclude=/proc/*",
			"--exclude=/sys/*",
			"--exclude=/dev/*",
			"--exclude=/tmp/*",
			b.BootfsPath + "/",
			mountPoint + "/",
		}
		if err := utils.RunCommand("rsync", args...); err != nil {
			return fmt.Errorf("复制文件失败: %v", err)
		}
	} else {
		args := []string{
			"-a",
			b.BootfsPath + "/.",
			mountPoint + "/",
		}
		if err := utils.RunCommand("cp", args...); err != nil {
			return fmt.Errorf("复制文件失败: %v", err)
		}
	}
	
	// 创建必要的目录
	dirs := []string{"proc", "sys", "dev", "tmp", "run"}
	for _, dir := range dirs {
		dirPath := filepath.Join(mountPoint, dir)
		if !utils.DirExists(dirPath) {
			utils.CreateDir(dirPath)
		}
	}
	
	// 设置权限
	os.Chmod(filepath.Join(mountPoint, "tmp"), 0777)
	
	return nil
}

// installBootloader 安装 bootloader（可选）
func (b *QemuBuilder) installBootloader(mountPoint string) {
	// 这里可以安装 GRUB 或其他 bootloader
	// 目前跳过，用户可以手动安装或使用 -kernel 参数启动
	fmt.Println("ℹ️  跳过 bootloader 安装，使用 -kernel 参数启动 QEMU")
}

// ParseSize 解析大小字符串（如 "2G", "512M"）为字节数
func ParseSize(size string) (int64, error) {
	size = strings.ToUpper(strings.TrimSpace(size))
	if size == "" {
		return 0, fmt.Errorf("大小不能为空")
	}
	
	// 提取数字和单位
	var numStr string
	var unit string
	
	for i, c := range size {
		if c >= '0' && c <= '9' || c == '.' {
			continue
		}
		numStr = size[:i]
		unit = size[i:]
		break
	}
	
	if numStr == "" {
		numStr = size
	}
	
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("无效的大小: %s", size)
	}
	
	// 转换为字节
	var multiplier float64 = 1
	switch unit {
	case "K", "KB":
		multiplier = 1024
	case "M", "MB":
		multiplier = 1024 * 1024
	case "G", "GB":
		multiplier = 1024 * 1024 * 1024
	case "T", "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "", "B":
		multiplier = 1
	default:
		return 0, fmt.Errorf("未知的单位: %s", unit)
	}
	
	return int64(num * multiplier), nil
}