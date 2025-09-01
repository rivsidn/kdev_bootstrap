package main

import (
	"fmt"
	"os"

	"github.com/rivsidn/kdev_bootstrap/pkg/builder"
	"github.com/rivsidn/kdev_bootstrap/pkg/config"
	"github.com/rivsidn/kdev_bootstrap/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	configFile string
	arch       string
	outputDir  string
)

var rootCmd = &cobra.Command{
	Use:   "kboot_build_bootfs",
	Short: "构建内核调试环境的根文件系统",
	Long: `kboot_build_bootfs 根据配置文件构建 bootfs（根文件系统）。

这个工具使用 debootstrap 创建一个最小化的 Ubuntu 根文件系统，
用于后续构建 Docker 镜像和 QEMU 镜像。`,
	RunE: runBuild,
}

func init() {
	// 设置命令行参数
	rootCmd.Flags().StringVarP(&configFile, "file", "f", "", "配置文件路径（必需）")
	rootCmd.Flags().StringVarP(&arch, "arch", "a", "", "目标架构（如：i386, amd64）")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "输出目录（默认为当前目录）")

	rootCmd.MarkFlagRequired("file")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 加载配置文件
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}

	fmt.Printf("配置信息:\n")
	fmt.Printf("   发行版: %s %s\n", cfg.Distribution, cfg.Version)
	fmt.Printf("   支持架构: %v\n", cfg.ArchSupported)
	fmt.Printf("   镜像源: %s\n", cfg.Mirror)

	// 如果没有指定架构，使用配置文件中的或默认架构
	if arch == "" {
		if cfg.ArchCurrent != "" {
			arch = cfg.ArchCurrent
		} else {
			// 获取默认架构
			arch = utils.GetDefaultArch()
			// 检查是否支持
			if !cfg.ValidateArch(arch) {
				if len(cfg.ArchSupported) > 0 {
					arch = cfg.ArchSupported[0]
				} else {
					return fmt.Errorf("无法确定架构")
				}
			}
		}
	}

	fmt.Printf("   目标架构: %s\n", arch)

	// 创建构建器
	builder := builder.NewBootfsBuilder(cfg, arch, outputDir)

	// 执行构建
	if err := builder.Build(); err != nil {
		return fmt.Errorf("构建失败: %v", err)
	}

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}
