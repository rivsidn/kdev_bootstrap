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
	Short: "Generate QEMU image from root filesystem",
	Long: `kboot_build_qemu creates QEMU disk image using built bootfs (root filesystem).
	
This image is mainly used for kernel debugging and can be started and tested in QEMU virtual machine.`,
	RunE: runBuild,
}

func init() {
	rootCmd.Flags().StringVarP(&bootfsPath, "bootfs", "b", "", "Root filesystem path (required)")
	rootCmd.Flags().StringVarP(&rootfsImage, "rootfs", "r", "", "Output rootfs.img name (optional)")
	rootCmd.Flags().StringVarP(&imageSize, "size", "s", "2G", "Image size (default: 2G)")
	
	rootCmd.MarkFlagRequired("bootfs")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 创建构建器
	builder, err := builder.NewQemuBuilder(bootfsPath, rootfsImage, imageSize)
	if err != nil {
		return err
	}
	
	fmt.Printf("Configuration:\n")
	fmt.Printf("   Distribution: %s %s\n", builder.Config.Distribution, builder.Config.Version)
	fmt.Printf("   Architecture: %s\n", builder.Config.ArchCurrent)
	fmt.Printf("   Bootfs: %s\n", bootfsPath)
	fmt.Printf("   Image size: %s\n", imageSize)
	
	// 执行构建
	if err := builder.Build(); err != nil {
		return fmt.Errorf("build failed: %v", err)
	}
	
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
