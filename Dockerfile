FROM storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
email="storezhang@gmail.com" \
qq="160290688" \
wechat="storezhang" \
description="Drone持续集成Docker插件，增加以下功能：1、多镜像仓库支持；2、镜像推送；3、镜像编译；4、多镜像仓库登录"


# 复制文件
COPY docker /bin


RUN set -ex \
    \
    \
    \
    && apk update \
    && apk --no-cache add docker \
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
