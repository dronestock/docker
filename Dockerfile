FROM ccr.ccs.tencentyun.com/storezhang/alpine


LABEL author="storezhang<华寅>"
LABEL email="storezhang@gmail.com"
LABEL qq="160290688"
LABEL wechat="storezhang"
LABEL description="Drone持续集成Docker插件"


RUN apk add docker
# 复制文件
COPY docker /
COPY dockerd /bin


RUN set -ex \
    \
    \
    \
    # && apk update \
    # && apk --no-cache add docker \
    \
    \
    \
    # 增加执行权限
    && chmod +x /bin/dockerd \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /bin/dockerd
