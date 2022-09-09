package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

var defaultMirrors = []string{
	`https://ustc-edu-cn.mirror.aliyuncs.com`,
	`https://mirror.baidubce.com`,
	`https://hub.daocloud.io`,
	`https://mirror.ccs.tencentyun.com`,
}

type plugin struct {
	drone.Base

	// 配置文件
	Dockerfile string `default:"${PLUGIN_DOCKERFILE=${DOCKERFILE=Dockerfile}}" validate:"required"`
	// 上下文
	Context string `default:"${PLUGIN_CONTEXT=${CONTEXT=.}}"`
	// 主机
	Host string `default:"${PLUGIN_HOST=${HOST=unix:///var/run/docker.sock}}" validate:"required"`
	// 镜像列表
	Mirrors []string `default:"${PLUGIN_MIRRORS=${MIRRORS}}"`
	// 标签
	Tag string `default:"${PLUGIN_TAG=${TAG=${DRONE_TAG=latest}}}"`
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

	// 仓库列表
	Registries []registry `default:"${PLUGIN_REGISTRIES=${REGISTRIES}}"`
	// 仓库
	Repository string `default:"${PLUGIN_REPOSITORY=${REPOSITORY}}"`
}

func newPlugin() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() []*drone.Step {
	return []*drone.Step{
		drone.NewStep(p.daemon, drone.Name(`启动守护进程`)),
		drone.NewStep(p.info, drone.Name(`查看Docker信息`)),
		drone.NewStep(p.login, drone.Name(`登录仓库`)),
		drone.NewStep(p.build, drone.Name(`编译镜像`)),
		drone.NewStep(p.push, drone.Name(`推送镜像`)),
	}
}

func (p *plugin) Fields() gox.Fields {
	return gox.Fields{
		field.String(`dockerfile`, p.Dockerfile),
		field.String(`context`, p.Context),
		field.String(`host`, p.Host),
		field.Strings(`mirrors`, p.Mirrors...),
		field.String(`tag`, p.Tag),
		field.Bool(`auto.tag`, p.AutoTag),
		field.String(`name`, p.Name),

		field.Bool(`experimental`, p.Experimental),
		field.Bool(`squash`, p.Squash),
		field.Bool(`compress`, p.Compress),
		field.Strings(`labels`, p.Labels...),

		field.String(`remote`, p.Remote),
		field.String(`link`, p.Link),

		field.String(`repository`, p.Repository),
	}
}

func (p *plugin) mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(defaultMirrors))
	if p.Defaults {
		mirrors = append(mirrors, defaultMirrors...)
	}
	mirrors = append(mirrors, p.Mirrors...)

	return
}

func (p *plugin) labels() (labels []string) {
	labels = make([]string, 0, 4+len(p.Labels))
	if p.Defaults {
		labels = append(labels, fmt.Sprintf("created=%s", time.Now().Format(time.RFC3339)))
		labels = append(labels, fmt.Sprintf("revision=%s", p.Name))
		labels = append(labels, fmt.Sprintf("source=%s", p.Remote))
		labels = append(labels, fmt.Sprintf("url=%s", p.Link))
	}
	labels = append(labels, p.Labels...)

	return
}

func (p *plugin) tags() (tags map[string]string) {
	tags = make(map[string]string, 3)
	tags[p.Tag] = p.Tag
	if !p.AutoTag {
		return
	}

	autos := strings.Split(p.Tag, `.`)
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

func (p *plugin) tag() string {
	return fmt.Sprintf(`%s:%s`, p.Repository, p.Name)
}

func (p *plugin) squash() bool {
	return p.Experimental && p.Squash
}

func (p *plugin) context() (context string) {
	if `` == p.Context {
		context = filepath.Dir(p.Dockerfile)
	} else {
		context = p.Context
	}

	return
}
