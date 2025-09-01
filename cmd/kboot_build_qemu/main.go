package main

import (
	"fmt"
	"os"

	"github.com/rivsidn/kdev_bootstrap/pkg/builder"
	"github.com/spf13/cobra"
)

var (
	bootfsPath  string
	rootfsImage string
	imageSize   string
)

var rootCmd = &cobra.Command{
	Use:   "kboot_build_qemu",
	Short: "通过根文件系统生成 QEMU 镜像",
	Long: `kboot_build_qemu 使用已构建的 bootfs（根文件系统）创建 QEMU 磁盘镜像。
	
这个镜像主要用于内核调试，可以在 QEMU 虚拟机中启动并测试内核。`,
	RunE: runBuild,
}

func init() {
	rootCmd.Flags().StringVarP(&bootfsPath, "bootfs", "b", "", "根文件系统路径（必需）")
	rootCmd.Flags().StringVarP(&rootfsImage, "rootfs", "r", "", "输出的 rootfs.img 名称（可选）")
	rootCmd.Flags().StringVarP(&imageSize, "size", "s", "2G", "镜像大小（默认 2G）")
	
	rootCmd.MarkFlagRequired("bootfs")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 创建构建器
	builder, err := builder.NewQemuBuilder(bootfsPath, rootfsImage, imageSize)
	if err != nil {
		return err
	}
	
	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("   发行版: %s %s\n", builder.Config.Distribution, builder.Config.Version)
	fmt.Printf("   架构: %s\n", builder.Config.ArchCurrent)
	fmt.Printf("   Bootfs: %s\n", bootfsPath)
	fmt.Printf("   镜像大小: %s\n", imageSize)
	
	// 执行构建
	if err := builder.Build(); err != nil {
		return fmt.Errorf("构建失败: %v", err)
	}
	
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 错误: %v\n", err)
		os.Exit(1)
	}
}