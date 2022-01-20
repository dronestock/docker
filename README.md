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
    # registry: ccr.ccs.tencentyun.com
    repository: dronestock/docker
    username: dronestock
    password:
      from_secret: token_docker
```
