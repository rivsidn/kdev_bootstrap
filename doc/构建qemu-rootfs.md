kboot_build_qemu 命令实现.

## 支持参数

| 短参数    | 长参数          | 说明                | 是否必须                                        |
|-----------|-----------------|---------------------|-------------------------------------------------|
| -b DIR    | --bootfs DIR    | 指定bootfs 路径     | 是                                              |
| -r ROOTFS | --rootfs ROOTFS | 指定rootfs.img 名称 | 否，如果没指定会根据/etc/bootstrap.conf自动生成 |
| -s SIZE   | --size SIZE     | 指定rootfs 镜像大小 | 否                                              |
| -h        | --help          | 显示帮助信息        | 否                                              |



## 示例

```bash
kboot_build_qemu -b /tmp/bootfs/
```

