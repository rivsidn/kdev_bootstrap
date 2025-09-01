package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rivsidn/kdev_bootstrap/pkg/config"
	"github.com/rivsidn/kdev_bootstrap/pkg/utils"
)

// DockerBuilder Docker 镜像构建器
type DockerBuilder struct {
	Config         *config.Config
	BootfsPath     string
	DockerfilePath string
	ImageName      string
}

// NewDockerBuilder 创建新的 Docker 构建器
func NewDockerBuilder(bootfsPath string, dockerfilePath string, imageName string) (*DockerBuilder, error) {
	// 加载配置文件
	configPath := filepath.Join(bootfsPath, "etc", "bootstrap.conf")
	if !utils.FileExists(configPath) {
		return nil, fmt.Errorf("configuration file not found: %s", configPath)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return &DockerBuilder{
		Config:         cfg,
		BootfsPath:     bootfsPath,
		DockerfilePath: dockerfilePath,
		ImageName:      imageName,
	}, nil
}

// Build 构建 Docker 镜像
func (b *DockerBuilder) Build() error {
	// 1. 检查环境
	if err := b.checkEnvironment(); err != nil {
		return err
	}

	// 2. 设置镜像名称
	if b.ImageName == "" {
		arch := b.Config.ArchCurrent
		if arch == "" {
			return fmt.Errorf("Can not find the valid arch");
		}
		b.ImageName = b.Config.GetImageName(arch)
	}

	// 3. 创建 Dockerfile
	if err := b.createDockerfile(); err != nil {
		return err
	}

	// 4. 构建 Docker 镜像
	if err := b.buildImage(); err != nil {
		return err
	}

	// 5. 清理临时文件
	if b.DockerfilePath != "" && strings.HasPrefix(filepath.Base(b.DockerfilePath), "Dockerfile.tmp") {
		os.Remove(b.DockerfilePath)
	}

	fmt.Printf("\nDocker image build successful: %s\n", b.ImageName)
	fmt.Printf("   Usage: docker run -it --rm %s /bin/bash\n", b.ImageName)

	return nil
}

// checkEnvironment 检查环境
func (b *DockerBuilder) checkEnvironment() error {
	// 检查 bootfs 目录
	if !utils.DirExists(b.BootfsPath) {
		return fmt.Errorf("bootfs directory does not exist: %s", b.BootfsPath)
	}

	// 检查是否为 root
	if !utils.CheckRoot() {
		return fmt.Errorf("please run with sudo or root privileges")
	}

	return nil
}

// createDockerfile 创建 Dockerfile
func (b *DockerBuilder) createDockerfile() error {
	if b.DockerfilePath != "" && utils.FileExists(b.DockerfilePath) {
		fmt.Printf("Using existing Dockerfile: %s\n", b.DockerfilePath)
		return nil
	}

	// 创建临时 Dockerfile
	tmpDir := filepath.Dir(b.BootfsPath)
	b.DockerfilePath = filepath.Join(tmpDir, "Dockerfile.tmp")

	arch := b.Config.ArchCurrent
	if arch == "" {
		return fmt.Errorf("Can not find the valid arch");
	}

	dockerfileContent := `FROM scratch

# 直接添加当前构建上下文（rootfs 目录内容），避免设备文件权限问题
ADD . /

ARG DEBIAN_FRONTEND=noninteractive

# 设置环境变量
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV ARCH=%s
ENV DISTRIBUTION=%s
ENV VERSION=%s
ENV LANG=C.UTF-8
ENV LC_ALL=C.UTF-8

# 设置工作目录
WORKDIR /root

# 默认命令
CMD ["/bin/bash"]
`
	content := fmt.Sprintf(dockerfileContent, arch, b.Config.Distribution, b.Config.Version)

	if err := os.WriteFile(b.DockerfilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create Dockerfile: %v", err)
	}

	fmt.Printf("Created temporary Dockerfile: %s\n", b.DockerfilePath)
	return nil
}


// buildImage 构建 Docker 镜像
func (b *DockerBuilder) buildImage() error {
	fmt.Printf("\nBuilding Docker image: %s\n", b.ImageName)

	// 构建上下文直接是 bootfs 目录，这样 ADD . / 会添加 bootfs 的内容
	buildContext := b.BootfsPath

	args := []string{
		"build",
		"-t", b.ImageName,
		"-f", b.DockerfilePath,
		buildContext,
	}

	if err := utils.RunCommand("docker", args...); err != nil {
		return fmt.Errorf("failed to build Docker image: %v", err)
	}

	// 显示镜像信息
	utils.RunCommand("docker", "images", b.ImageName)

	return nil
}
