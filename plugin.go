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
	Dockerfile string `default:"${DOCKERFILE=Dockerfile}" validate:"required"`
	// 上下文
	Context string `default:"${CONTEXT=.}"`
	// 主机
	Host string `default:"${HOST=unix:///var/run/docker.sock}" validate:"required"`
	// 镜像列表
	Mirrors []string `default:"${MIRRORS}"`
	// 标签
	Tag string `default:"${TAG=${DRONE_TAG=0.0.${DRONE_BUILD_NUMBER}}}"`
	// 自动标签
	AutoTag bool `default:"${AUTO_TAG=true}"`
	// 名称
	Name string `default:"${NAME=${DRONE_COMMIT_SHA=latest}}"`

	// 启用实验性功能
	Experimental bool `default:"${EXPERIMENTAL=true}"`
	// 精减镜像层数
	Squash bool `default:"${SQUASH=true}"`
	// 压缩镜像
	Compress bool `default:"${COMPRESS=true}"`
	// 标签列表
	Labels []string `default:"${LABELS}"`

	// 仓库地址
	Remote string `default:"${REMOTE=${DRONE_REMOTE_URL=https://github.com/dronestock/docker}}"`
	// 镜像链接
	// nolint:lll
	Link string `default:"${LINK=${PLUGIN_REPO_LINK=${DRONE_REPO_LINK=https://github.com/dronestock/docker}}}"`

	// 数据目录
	DataRoot string `default:"${DATA_ROOT=/var/lib/docker}"`
	// 驱动
	StorageDriver string `default:"${STORAGE_DRIVER}"`

	// 仓库列表
	Registries []registry `default:"${REGISTRIES}"`
	// 仓库
	Repository string `default:"${REPOSITORY}"`

	// 加速
	Boost boost `default:"${BOOST}"`
}

func newPlugin() drone.Plugin {
	return new(plugin)
}

func (p *plugin) Config() drone.Config {
	return p
}

func (p *plugin) Steps() drone.Steps {
	return drone.Steps{
		drone.NewStep(newBoostStep(p)).Name("加速").Build(),
		drone.NewStep(newDaemonStep(p)).Name("守护").Build(),
		drone.NewStep(newInfoStep(p)).Name("检查").Build(),
		drone.NewStep(newLoginStep(p)).Name("登录").Build(),
		drone.NewStep(newBuildStep(p)).Name("编译").Build(),
		drone.NewStep(newPushStep(p)).Name("推送").Build(),
	}
}

func (p *plugin) Fields() gox.Fields[any] {
	return gox.Fields[any]{
		field.New("dockerfile", p.Dockerfile),
		field.New("context", p.Context),
		field.New("host", p.Host),
		field.New("mirrors", p.Mirrors),
		field.New("tag", p.Tag),
		field.New("tag.auto", p.AutoTag),
		field.New("name", p.Name),

		field.New("experimental", p.Experimental),
		field.New("squash", p.Squash),
		field.New("compress", p.Compress),
		field.New("labels", p.Labels),

		field.New("remote", p.Remote),
		field.New("link", p.Link),

		field.New("repository", p.Repository),
	}
}

func (p *plugin) mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(defaultMirrors))
	if p.Default() {
		mirrors = append(mirrors, defaultMirrors...)
	}
	mirrors = append(mirrors, p.Mirrors...)

	return
}

func (p *plugin) labels() (labels []string) {
	labels = make([]string, 0, 4+len(p.Labels))
	if p.Default() {
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

	autos := strings.Split(p.Tag, common)
	for index := range autos {
		tag := strings.Join(autos[0:index], common)
		tags[tag] = tag
	}
	tags["latest"] = "latest"

	return
}

func (p *plugin) tag() string {
	return fmt.Sprintf("%s:%s", p.Repository, p.Name)
}

func (p *plugin) squash() bool {
	return p.Experimental && p.Squash
}

func (p *plugin) context() (context string) {
	if "" == p.Context {
		context = filepath.Dir(p.Dockerfile)
	} else {
		context = p.Context
	}

	return
}

func (p *plugin) boostEnabled() bool {
	return nil != p.Boost.Enabled && *p.Boost.Enabled
}
