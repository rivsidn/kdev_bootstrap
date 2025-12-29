#!/bin/bash

# 构建脚本 - 快速构建所有工具

set -e

cmd_build() {
	if ! command -v go &> /dev/null; then
		echo "error: Go not installed" >&2
		exit 1
	fi

	mkdir -p bin

	# 阻止 go mod tidy 扫描 bin 目录下的 bootfs 文件系统
	echo "module ignore" > bin/go.mod

	go mod download
	go mod tidy

	for cmd in kboot_build_bootfs kboot_build_docker kboot_build_qemu; do
		go build -o bin/$cmd ./cmd/$cmd
	done
}

cmd_clean() {
	rm -rf bin
}

# 默认为build命令
case "${1:-build}" in
	build)
		cmd_build
		;;
	clean)
		cmd_clean
		;;
	*)
		echo "usage: $0 {build|clean}" >&2
		exit 1
		;;
esac

