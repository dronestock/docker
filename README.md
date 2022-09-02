# docker

Drone持续集成Docker插件

## 功能

- 自动标签
- 默认镜像
- 自动重试
- 重试背压

## 使用

使用`docker`插件非常简单，只需要基础配置

```yaml
name: 打包Docker到中央仓库
  image: dronestock/docker
  pull: if-not-exists
  settings:
    repository: dronestock/docker
    registries:
      - username: dronestock
        password:
          from_secret: token_docker
      - hostname: ccr.ccs.tencentyun.com
        username: 160290688
        password:
          from_secret: password_tencentyun_ccr
        required: true
```

## 感谢Jetbrains

本项目通过`Jetbrains开源许可IDE`编写源代码，特此感谢
[![Jetbrains图标](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png)](https://www.jetbrains.com/?from=dronestock/docker)
