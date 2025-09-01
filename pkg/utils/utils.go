package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

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

func RunCommandOutput(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed %s: %v\noutput: %s", name, err, string(output))
	}
	return string(output), nil
}

func CheckCommand(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func CheckRoot() bool {
	return os.Geteuid() == 0
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

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
	return response == "y" || response == "Y" || response == "yes"
}

