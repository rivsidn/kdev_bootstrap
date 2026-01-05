# Repository Guidelines

## 项目结构与模块职责
仓库以 Go 模块管理，核心 CLI 位于 `cmd/`，分别包含 `kboot_build_bootfs`、`kboot_build_docker`、`kboot_build_qemu`。公共逻辑集中在 `pkg/`：`config/` 解析 ini 配置，`builder/` 负责根文件系统与镜像创建，`utils/` 提供共享工具。示例配置存放在 `configs/`，对照文档在 `doc/`，辅助脚本位于 `scripts/`，构建产物进入 `bin/`。添加新子命令时请在 `cmd/<name>` 建立入口，并将其依赖封装到相应 `pkg` 子包。

## 构建、测试与开发命令
`make build` 会下载依赖并在 `bin/` 生成全部二进制，提交前须可重复执行。`make clean` 清理产物，避免旧版文件混入。调试阶段可直接 `go run ./cmd/kboot_build_bootfs -a amd64 -f configs/ubuntu-16.04.conf` 以验证参数。构建出的工具常配合 `sudo ./bin/kboot_build_docker -b <bootfs-dir>` 等命令创建镜像。

## 代码风格与命名
使用 Go 1.21，保持 `gofmt` 结果一致，结构体字段与配置键名保持蛇形命名以贴合 INI 文件。包命名使用短小的单词（如 `builder`），新增 Cobra 命令遵循 `New<Feature>Cmd` 工厂模式，并在 `cmd/root.go` 中注册。提交前执行 `go fmt ./... && go vet ./...`。

## 测试准则
尽管当前测试较少，新增功能必须附带 `*_test.go`。优先覆盖配置解析与构建流程，推荐 `go test ./pkg/... -run TestConfig` 逐包验证，并在总线上运行 `go test ./... -cover`，目标是修改文件行覆盖率不低于 80%。集成脚本需提供最小可运行样例或伪造配置以便评审者重现。

## 提交与 PR 规范
历史提交多为中文动宾结构（示例：`同步最新的kdev_env 文件`、`解决ubuntu-22.04 没有/sbin/init 程序问题`），请在 50 字符内概述变更。PR 描述应包括：问题背景、改动要点、相关配置（如 `configs/ubuntu-22.04.conf`）以及验证方式（命令或日志）。涉及脚本或 root 权限操作时附上风险说明与回滚提示。
