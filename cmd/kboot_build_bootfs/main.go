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
	Short: "Build root filesystem for kernel debugging environment",
	Long: `kboot_build_bootfs builds bootfs (root filesystem) based on configuration file.

This tool uses debootstrap to create a minimal Ubuntu root filesystem,
which is used for building Docker images and QEMU images.`,
	RunE: runBuild,
}

func init() {
	// 参数解析
	rootCmd.Flags().StringVarP(&configFile, "file", "f", "", "Configuration file path (required)")
	rootCmd.Flags().StringVarP(&arch, "arch", "a", "", "Target architecture (e.g., i386, amd64)")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Output directory (default: current directory)")

	rootCmd.MarkFlagRequired("file")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 加载配置文件
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return err
	}

	fmt.Printf("Configuration:\n")
	fmt.Printf("   Distribution: %s %s\n", cfg.Distribution, cfg.Version)
	fmt.Printf("   Supported architectures: %v\n", cfg.ArchSupported)
	fmt.Printf("   Mirror: %s\n", cfg.Mirror)

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
					return fmt.Errorf("cannot determine architecture")
				}
			}
		}
	}

	fmt.Printf("   Target architecture: %s\n", arch)

	// 创建构建器
	builder := builder.NewBootfsBuilder(cfg, arch, outputDir)

	// 执行构建
	if err := builder.Build(); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}

	return nil
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
