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
steps:
  - name: 打包Docker到中央仓库
    image: dronestock/docker
    settings:
      repository: dronestock/docker
      registries:
        - username: dronestock
          password: password_docker
        - hostname: ccr.ccs.tencentyun.com
          username: "160290688"
          password: password_ccr
          required: true
```

更多使用教程，请参考[使用文档](https://www.dronestock.tech/plugin/stock/docker)

## 交流

![微信群](https://www.dronestock.tech/communication/wxwork.jpg)

## 捐助

![支持宝](https://github.com/storezhang/donate/raw/master/alipay-small.jpg)
![微信](https://github.com/storezhang/donate/raw/master/weipay-small.jpg)

## 感谢Jetbrains

本项目通过`Jetbrains开源许可IDE`编写源代码，特此感谢

[![Jetbrains图标](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)](https://www.jetbrains.com/?from=dronestock/docker)
