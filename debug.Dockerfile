FROM ccr.ccs.tencentyun.com/storezhang/alpine:3.20.0


RUN apk update
RUN apk add openssh
RUN apk add docker
RUN mkdir /tmp


# 复制执行文件
COPY docker /
ARG TARGETPLATFORM
COPY dist/${TARGETPLATFORM}/dockerd /usr/local/bin/
RUN chmod +x /usr/local/bin/*


# 执行命令
ENTRYPOINT /usr/local/bin/dockerd
