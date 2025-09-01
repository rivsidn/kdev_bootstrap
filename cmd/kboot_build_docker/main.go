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
	Short: "通过根文件系统生成 Docker 镜像",
	Long: `kboot_build_docker 使用已构建的 bootfs（根文件系统）创建 Docker 镜像。
	
这个镜像主要用于内核编译环境，包含了构建内核所需的基本工具和库。`,
	RunE: runBuild,
}

func init() {
	rootCmd.Flags().StringVarP(&bootfsPath, "bootfs", "b", "", "根文件系统路径（必需）")
	rootCmd.Flags().StringVarP(&dockerfilePath, "dockerfile", "f", "", "Dockerfile 文件路径（可选）")
	rootCmd.Flags().StringVar(&imageName, "image", "", "镜像名称（格式：name:tag，可选）")
	
	rootCmd.MarkFlagRequired("bootfs")
}

func runBuild(cmd *cobra.Command, args []string) error {
	// 创建构建器
	builder, err := builder.NewDockerBuilder(bootfsPath, dockerfilePath, imageName)
	if err != nil {
		return err
	}
	
	fmt.Printf("📋 配置信息:\n")
	fmt.Printf("   发行版: %s %s\n", builder.Config.Distribution, builder.Config.Version)
	fmt.Printf("   架构: %s\n", builder.Config.ArchCurrent)
	fmt.Printf("   Bootfs: %s\n", bootfsPath)
	
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