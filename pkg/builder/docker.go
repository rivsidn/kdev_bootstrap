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
		return nil, fmt.Errorf("找不到配置文件: %s", configPath)
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
			// 尝试从 bootfs 路径推断架构
			arch = b.inferArch()
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
	
	fmt.Printf("\n✅ Docker 镜像构建成功: %s\n", b.ImageName)
	fmt.Printf("   使用方法: docker run -it --rm %s /bin/bash\n", b.ImageName)
	
	return nil
}

// checkEnvironment 检查环境
func (b *DockerBuilder) checkEnvironment() error {
	// 检查是否为 root
	if !utils.CheckRoot() {
		return fmt.Errorf("请使用 sudo 或 root 权限运行")
	}
	
	// 检查 Docker
	if !utils.CheckCommand("docker") {
		return fmt.Errorf("未安装 Docker，请先安装: sudo apt-get install docker.io")
	}
	
	// 检查 bootfs 目录
	if !utils.DirExists(b.BootfsPath) {
		return fmt.Errorf("bootfs 目录不存在: %s", b.BootfsPath)
	}
	
	return nil
}

// inferArch 从路径推断架构
func (b *DockerBuilder) inferArch() string {
	base := filepath.Base(b.BootfsPath)
	if strings.Contains(base, "i386") {
		return "i386"
	}
	if strings.Contains(base, "amd64") {
		return "amd64"
	}
	if strings.Contains(base, "arm64") {
		return "arm64"
	}
	return "amd64" // 默认
}

// createDockerfile 创建 Dockerfile
func (b *DockerBuilder) createDockerfile() error {
	if b.DockerfilePath != "" && utils.FileExists(b.DockerfilePath) {
		fmt.Printf("使用现有 Dockerfile: %s\n", b.DockerfilePath)
		return nil
	}
	
	// 创建临时 Dockerfile
	tmpDir := filepath.Dir(b.BootfsPath)
	b.DockerfilePath = filepath.Join(tmpDir, "Dockerfile.tmp")
	
	arch := b.Config.ArchCurrent
	if arch == "" {
		arch = b.inferArch()
	}
	
	dockerfileContent := `FROM scratch

# 直接添加当前构建上下文（rootfs 目录内容），避免设备文件权限问题
ADD . /

# 设置环境变量
ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
ENV DEBIAN_FRONTEND=noninteractive
ENV ARCH=%s
ENV DISTRIBUTION=%s
ENV VERSION=%s

# 设置工作目录
WORKDIR /root

# 默认命令
CMD ["/bin/bash"]
`
	content := fmt.Sprintf(dockerfileContent, arch, b.Config.Distribution, b.Config.Version)
	
	if err := os.WriteFile(b.DockerfilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("创建 Dockerfile 失败: %v", err)
	}
	
	fmt.Printf("创建临时 Dockerfile: %s\n", b.DockerfilePath)
	return nil
}


// buildImage 构建 Docker 镜像
func (b *DockerBuilder) buildImage() error {
	fmt.Printf("\n🐳 构建 Docker 镜像: %s\n", b.ImageName)
	
	// 构建上下文直接是 bootfs 目录，这样 ADD . / 会添加 bootfs 的内容
	buildContext := b.BootfsPath
	
	args := []string{
		"build",
		"-t", b.ImageName,
		"-f", b.DockerfilePath,
		buildContext,
	}
	
	if err := utils.RunCommand("docker", args...); err != nil {
		return fmt.Errorf("构建 Docker 镜像失败: %v", err)
	}
	
	// 显示镜像信息
	utils.RunCommand("docker", "images", b.ImageName)
	
	return nil
}