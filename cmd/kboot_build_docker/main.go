package main

import (
	"fmt"
	"os"

	"github.com/rivsidn/kdev_bootstrap/pkg/builder"
	"github.com/spf13/cobra"
)

var (
	bootfsPath     string
	dockerfilePath string
	imageName      string
)

var rootCmd = &cobra.Command{
	Use:   "kboot_build_docker",
	Short: "Generate Docker image from root filesystem",
	Long: `kboot_build_docker creates Docker image using built bootfs (root filesystem).
	
This image is mainly used for kernel compilation environment, including basic tools and libraries needed for kernel building.`,
	RunE: runBuild,
}

func init() {
	rootCmd.Flags().StringVarP(&bootfsPath, "bootfs", "b", "", "Root filesystem path (required)")
	rootCmd.Flags().StringVarP(&dockerfilePath, "dockerfile", "f", "", "Dockerfile file path (optional)")
	rootCmd.Flags().StringVar(&imageName, "image", "", "Image name (format: name:tag, optional)")
	
	rootCmd.MarkFlagRequired("bootfs")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 创建构建器
	builder, err := builder.NewDockerBuilder(bootfsPath, dockerfilePath, imageName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Configuration:\n")
	fmt.Printf("   Distribution: %s %s\n", builder.Config.Distribution, builder.Config.Version)
	fmt.Printf("   Architecture: %s\n", builder.Config.ArchCurrent)
	fmt.Printf("   Bootfs: %s\n", bootfsPath)
	
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