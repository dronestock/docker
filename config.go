package main

import (
	`fmt`
	`path/filepath`
	`strings`
	`time`

	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/mengpo`
	`github.com/storezhang/validatorx`
)

type config struct {
	// 配置文件
	Dockerfile string `default:"${PLUGIN_DOCKERFILE=${DOCKERFILE=Dockerfile}}" validate:"required"`
	// 上下文
	Context string `default:"${PLUGIN_CONTEXT=${CONTEXT}}"`
	// 主机
	Host string `default:"${PLUGIN_HOST=${HOST=unix:///var/run/docker.sock}}" validate:"required"`
	// 镜像列表
	Mirrors []string `default:"${PLUGIN_MIRRORS=${MIRRORS}}"`
	// 标签
	Tag string `default:"${PLUGIN_TAG=${TAG=latest}}"`
	// 自动标签
	AutoTag bool `default:"${PLUGIN_AUTO_TAG=${AUTO_TAG=true}}"`
	// 名称
	Name string `default:"${PLUGIN_NAME=${NAME=${DRONE_COMMIT_SHA=latest}}}"`

	// 启用实验性功能
	Experimental bool `default:"${PLUGIN_EXPERIMENTAL=${EXPERIMENTAL=true}}"`
	// 精减镜像导数
	Squash bool `default:"${PLUGIN_SQUASH=${SQUASH=true}}"`
	// 压缩镜像
	Compress bool `default:"${PLUGIN_COMPRESS=${COMPRESS=true}}"`
	// 标签列表
	Labels []string `default:"${PLUGIN_LABELS=${LABELS}}"`

	// 仓库地址
	Remote string `default:"${PLUGIN_REMOTE=${REMOTE=${DRONE_REMOTE_URL=https://github.com/dronestock/docker}}}"`
	// 镜像链接
	// nolint:lll
	Link string `default:"${PLUGIN_LINK=${LINK=${PLUGIN_REPO_LINK=${DRONE_REPO_LINK=https://github.com/dronestock/docker}}}}"`

	// 数据目录
	DataRoot string `default:"${PLUGIN_DATA_ROOT=${DATA_ROOT=/var/lib/docker}}"`
	// 驱动
	StorageDriver string `default:"${PLUGIN_STORAGE_DRIVER=${STORAGE_DRIVER}}"`

	// 仓库地址
	Registry string `default:"${PLUGIN_REGISTRY=${REGISTRY=https://index.docker.io/v2}}"`
	// 用户名
	Username string `default:"${PLUGIN_USERNAME=${USERNAME}}"`
	// 密码
	Password string `default:"${PLUGIN_PASSWORD=${PASSWORD}}"`
	// 仓库
	Repository string `default:"${PLUGIN_REPOSITORY=${REPOSITORY}}"`

	// 是否启用默认值
	Defaults bool `default:"${PLUGIN_DEFAULTS=${DEFAULTS=true}}"`
	// 是否显示调试信息
	Verbose bool `default:"${PLUGIN_VERBOSE=${VERBOSE=false}}"`

	defaultMirrors    []string
	exe               string
	daemon            string
	outsideDockerfile string

	daemonSuccessMark string
	loginSuccessMark  string
}

func (c *config) Fields() gox.Fields {
	return []gox.Field{
		field.String(`dockerfile`, c.Dockerfile),
		field.String(`context`, c.Context),
		field.String(`host`, c.Host),
		field.Strings(`mirrors`, c.Mirrors...),
		field.String(`tag`, c.Tag),
		field.Bool(`auto.tag`, c.AutoTag),
		field.String(`name`, c.Name),

		field.Bool(`experimental`, c.Experimental),
		field.Bool(`squash`, c.Squash),
		field.Bool(`compress`, c.Compress),
		field.Strings(`labels`, c.Labels...),

		field.String(`remote`, c.Remote),
		field.String(`link`, c.Link),

		field.String(`registry`, c.Registry),
		field.String(`username`, c.Username),
		field.String(`repository`, c.Repository),

		field.Bool(`defaults`, c.Defaults),
		field.Bool(`verbose`, c.Verbose),
	}
}

func (c *config) load() (err error) {
	if err = mengpo.Set(c); nil != err {
		return
	}
	if err = validatorx.Struct(c); nil != err {
		return
	}
	c.init()

	return
}

func (c *config) init() {
	c.defaultMirrors = []string{
		`https://mirror.baidubce.com`,
		`https://hub.daocloud.io`,
		`https://mirror.ccs.tencentyun.com`,
		`https://docker.mirrors.ustc.edu.cn`,
	}

	c.exe = `/usr/bin/docker`
	c.daemon = `/usr/bin/dockerd`
	c.outsideDockerfile = `/var/run/docker.sock`

	c.daemonSuccessMark = `API listen on /var/run/docker.sock`
	c.loginSuccessMark = `Login Succeeded`
}

func (c *config) mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(c.defaultMirrors))
	if c.Defaults {
		mirrors = append(mirrors, c.defaultMirrors...)
	}
	mirrors = append(mirrors, c.Mirrors...)

	return
}

func (c *config) labels() (labels []string) {
	labels = make([]string, 0, 4)
	if c.Defaults {
		labels = append(labels, fmt.Sprintf("created=%s", time.Now().Format(time.RFC3339)))
		labels = append(labels, fmt.Sprintf("revision=%s", c.Name))
		labels = append(labels, fmt.Sprintf("source=%s", c.Remote))
		labels = append(labels, fmt.Sprintf("url=%s", c.Link))
	}
	labels = append(labels, c.Labels...)

	return
}

func (c *config) tags() (tags map[string]string) {
	tags = make(map[string]string, 3)
	tags[c.Tag] = c.Tag
	if !c.AutoTag {
		return
	}

	autos := strings.Split(c.Tag, `.`)
	_len := len(autos)
	if 1 == _len {
		tags[autos[0]] = autos[0]
	} else if 2 == _len {
		tags[autos[0]] = autos[0]
		second := fmt.Sprintf(`%s.%s`, autos[0], autos[1])
		tags[second] = second
	} else if 3 <= _len {
		tags[autos[0]] = autos[0]
		second := fmt.Sprintf(`%s.%s`, autos[0], autos[1])
		tags[second] = second
		third := fmt.Sprintf(`%s.%s.%s`, autos[0], autos[1], autos[2])
		tags[third] = third
	}
	tags[`latest`] = `latest`

	return
}

func (c *config) squash() bool {
	return c.Experimental && c.Squash
}

func (c *config) context() (context string) {
	if `` == c.Context || `.` == c.Context {
		context = filepath.Dir(c.Dockerfile)
	} else {
		context = c.Context
	}

	return
}
