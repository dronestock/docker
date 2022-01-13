FROM storezhang/alpine


LABEL author="storezhang<华寅>"
LABEL email="storezhang@gmail.com"
LABEL qq="160290688"
LABEL wechat="storezhang"
LABEL architecture="AMD64/x86_64" version="latest" build="2021-12-31"
LABEL Description="Drone持续集成Git插件，增加标签功能以及Github加速功能。同时支持推拉模式"


# 复制文件
COPY docker /bin


RUN set -ex \
    \
    \
    \
    && apk update \
    && apk add docker \
    \
    \
    \
    # 增加执行权限
    && chmod +x /bin/docker \
    \
    \
    \
    && rm -rf /var/cache/apk/*



# 执行命令
ENTRYPOINT /bin/docker
