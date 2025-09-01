package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// Config 配置文件结构
type Config struct {
	Distribution   string            `ini:"distribution"`
	Version        string            `ini:"version"`
	ArchSupported  []string          `ini:"-"`
	ArchCurrent    string            `ini:"arch_current"`
	Mirror         string            `ini:"mirror"`
	Packages       map[string]string `ini:"-"`
	
	// 内部字段
	sectionName      string
	ArchSupportedRaw string `ini:"arch_supported"`
}

// UbuntuSuiteMap Ubuntu版本与suite的映射
var UbuntuSuiteMap = map[string]string{
	"5.10":  "breezy",
	"16.04": "xenial",
	"18.04": "bionic",
	"20.04": "focal",
	"22.04": "jammy",
	"24.04": "noble",
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法加载配置文件 %s: %v", configPath, err)
	}

	// 获取第一个非默认section
	var section *ini.Section
	var sectionName string
	for _, s := range cfg.Sections() {
		if s.Name() != ini.DefaultSection {
			section = s
			sectionName = s.Name()
			break
		}
	}

	if section == nil {
		return nil, fmt.Errorf("配置文件中没有找到有效的section")
	}

	config := &Config{
		sectionName: sectionName,
		Packages:    make(map[string]string),
	}

	// 解析基本字段
	if err := section.MapTo(config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}

	// 解析 arch_supported
	if config.ArchSupportedRaw != "" {
		config.ArchSupported = strings.Split(config.ArchSupportedRaw, ",")
		for i := range config.ArchSupported {
			config.ArchSupported[i] = strings.TrimSpace(config.ArchSupported[i])
		}
	}

	// 解析所有 _packages 结尾的配置
	for _, key := range section.Keys() {
		if strings.HasSuffix(key.Name(), "_packages") {
			config.Packages[key.Name()] = key.Value()
		}
	}

	// 设置默认镜像
	if config.Mirror == "" {
		if strings.HasPrefix(config.Version, "5.") {
			config.Mirror = "http://old-releases.ubuntu.com/ubuntu/"
		} else {
			config.Mirror = "http://mirrors.aliyun.com/ubuntu/"
		}
	}

	return config, nil
}

// SaveToBootfs 将配置保存到 bootfs 的 /etc/bootstrap.conf
func (c *Config) SaveToBootfs(bootfsPath string) error {
	configPath := filepath.Join(bootfsPath, "etc", "bootstrap.conf")
	
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建新的 INI 文件
	cfg := ini.Empty()
	section, err := cfg.NewSection(c.sectionName)
	if err != nil {
		return fmt.Errorf("创建section失败: %v", err)
	}

	// 写入基本字段
	section.NewKey("distribution", c.Distribution)
	section.NewKey("version", c.Version)
	section.NewKey("arch_supported", strings.Join(c.ArchSupported, ","))
	if c.ArchCurrent != "" {
		section.NewKey("arch_current", c.ArchCurrent)
	}
	if c.Mirror != "" {
		section.NewKey("mirror", c.Mirror)
	}

	// 写入 packages
	for key, value := range c.Packages {
		section.NewKey(key, value)
	}

	// 保存文件
	if err := cfg.SaveTo(configPath); err != nil {
		return fmt.Errorf("保存配置文件失败: %v", err)
	}

	return nil
}

// GetSuite 获取Ubuntu的suite名称
func (c *Config) GetSuite() string {
	if suite, ok := UbuntuSuiteMap[c.Version]; ok {
		return suite
	}
	// 如果找不到映射，使用版本号作为suite
	return c.Version
}

// GetAllPackages 获取所有要安装的包
func (c *Config) GetAllPackages() []string {
	var packages []string
	for _, pkgList := range c.Packages {
		for _, pkg := range strings.Split(pkgList, ",") {
			pkg = strings.TrimSpace(pkg)
			if pkg != "" {
				packages = append(packages, pkg)
			}
		}
	}
	return packages
}

// GetImageName 生成镜像名称
func (c *Config) GetImageName(arch string) string {
	return fmt.Sprintf("%s-%s-%s", c.Distribution, c.Version, arch)
}

// GetRootfsName 生成 rootfs 名称
func (c *Config) GetRootfsName(arch string) string {
	return fmt.Sprintf("%s-%s-%s-rootfs.img", c.Distribution, c.Version, arch)
}

// ValidateArch 验证架构是否支持
func (c *Config) ValidateArch(arch string) bool {
	for _, a := range c.ArchSupported {
		if a == arch {
			return true
		}
	}
	return false
}