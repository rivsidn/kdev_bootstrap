package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCommand 执行系统命令
func RunCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	fmt.Printf("Executing command: %s %s\n", name, strings.Join(args, " "))
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed %s: %v", name, err)
	}
	
	return nil
}

// RunCommandOutput 执行命令并返回输出
func RunCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed %s: %v\noutput: %s", name, err, string(output))
	}
	return string(output), nil
}

// CheckCommand 检查命令是否存在
func CheckCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

// CheckRoot 检查是否以 root 权限运行
func CheckRoot() bool {
	return os.Geteuid() == 0
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists 检查目录是否存在
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CreateDir 创建目录
func CreateDir(path string) error {
	if !DirExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}
	return nil
}

// Confirm 用户确认
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

// CheckDependencies 检查系统依赖
func CheckDependencies(deps []string) error {
	var missing []string
	for _, dep := range deps {
		if !CheckCommand(dep) {
			missing = append(missing, dep)
		}
	}
	
	if len(missing) > 0 {
		return fmt.Errorf("missing dependencies: %s\nplease run: sudo apt-get install %s", 
			strings.Join(missing, ", "), strings.Join(missing, " "))
	}
	
	return nil
}

// GetDefaultArch 获取默认架构
func GetDefaultArch() string {
	output, err := RunCommandOutput("dpkg", "--print-architecture")
	if err != nil {
		return "amd64"
	}
	return strings.TrimSpace(output)
}