kboot_build_docker 命令实现.

通过根文件系统生成docker镜像.

## 支持参数

| 短参数        | 长参数                | 说明             | 是否必须                                   |
|---------------|-----------------------|------------------|--------------------------------------------|
| -b DIR        | --bootfs DIR          | 根文件系统路径   | 是                                         |
| -f DOCKERFILE | --dockfile DOCKERFILE | Dockerfile文件名 | 否                                         |
|               | --image IMAGE:TAG     | 制定镜像名称     | 否，如果不存在根据/etc/bootstrap.conf 生成 |
| -h            | --help                | 显示帮助信息     | 否                                         |

## 示例

```bash
kboot_build_docker -b /tmp/bootfs/
```

## TODO

- 如何将 Dockerfile 包含到程序中

