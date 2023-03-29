FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2 AS builder

# 复制执行文件
COPY dockerd /docker/usr/local/bin/
# 复制脚本
COPY docker /docker



FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.17.2


LABEL author="storezhang<华寅>" \
    email="storezhang@gmail.com" \
    qq="160290688" \
    wechat="storezhang" \
    description="Drone持续集成Docker插件，增加以下功能：1、多镜像仓库支持；2、镜像推送；3、镜像编译；4、多镜像仓库登录"


# 复制文件，多个复制命令都合并成一个
COPY --from=builder /docker /


RUN set -ex \
    \
    \
    \
    && apk update \
    # 安装SSH客户端
    && apk --no-cache add openssh \
    \
    # 安装Docker客户端
    && apk --no-cache add docker \
    \
    \
    \
    # 增加执行权限
    && chmod +x /usr/local/bin/* \
    \
    \
    \
    && rm -rf /var/cache/apk/*


# 执行命令
ENTRYPOINT /usr/local/bin/dockerd
