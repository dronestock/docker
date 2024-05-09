package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/gox"
)

type Target struct {
	Platform `default:"${PLATFORM}"` // 平台

	// 配置文件
	Dockerfile string `default:"${DOCKERFILE=Dockerfile}" json:"dockerfile,omitempty" validate:"required"`
	// 上下文
	Context string `default:"${CONTEXT=.}" json:"context,omitempty"`

	// 标签
	Tag string `default:"${TAG=${DRONE_TAG=0.0.${DRONE_BUILD_NUMBER}}}" json:"tag,omitempty" validate:"required"`
	// 前缀
	Prefix string `default:"${PREFIX}" json:"prefix,omitempty"`
	// 中间
	Middle string `default:"${MIDDLE}" json:"middle,omitempty"`
	// 后缀
	Suffix string `default:"${SUFFIX}" json:"suffix,omitempty"`
	// 自动标签
	Auto bool `default:"${AUTO=true}" json:"auto,omitempty"`
	// 名称
	Name string `default:"${NAME=${DRONE_COMMIT_SHA=latest}}" json:"name,omitempty"`

	// 平台列表
	Platforms Platforms `default:"${PLATFORMS}" json:"platforms,omitempty"`

	// 仓库
	Registry *Registry `json:"registry,omitempty"`
	// 仓库列表
	Registries []*Registry `json:"registries,omitempty"`
}

func (t *Target) AllRegistries() (registries []*Registry) {
	registries = make([]*Registry, 0, len(t.Registries)+1)
	if nil != t.Registry {
		registries = append(registries, t.Registry)
	}
	registries = append(registries, t.Registries...)

	return
}

func (t *Target) AllPlatforms() (platforms Platforms) {
	platforms = make(Platforms, 0, len(t.Platforms)+1)
	if "" != t.Platform.Argument() {
		platforms = append(platforms, &t.Platform)
	}
	platforms = append(platforms, t.Platforms...)

	return
}

func (t *Target) PlatformArgument() (argument string) {
	allPlatforms := t.AllPlatforms()
	platforms := make([]string, 0, len(allPlatforms))
	for _, platform := range allPlatforms {
		platforms = gox.Ift("" != platform.Argument(), append(platforms, platform.Argument()), platforms)
	}

	// 避免出现只有连接符号的字符串
	if 0 != len(platforms) {
		argument = strings.Join(platforms, constant.Comma)
	}

	return
}

func (t *Target) Dir() (context string) {
	if "" == t.Context {
		context = filepath.Dir(t.Dockerfile)
	} else {
		context = t.Context
	}

	return
}

func (t *Target) Qemu() (qemu bool) {
	for _, platform := range t.AllPlatforms() {
		qemu = platform.Qemu()
		if qemu {
			break
		}
	}

	return
}

func (t *Target) BuildWithPush(registries *Registries, docker *Docker) bool {
	return 0 != len(*registries) && "" != docker.Repository && 1 < len(t.AllPlatforms())
}

func (t *Target) Tags(registries *Registries, docker *Docker) (tags []string) {
	autos := t.autos()
	tags = make([]string, 0, (len(*registries)+len(t.AllRegistries()))*len(autos))
	for _, auto := range autos {
		final := gox.StringBuilder(t.Prefix)
		if "" != auto {
			final.Append(t.Middle)
		}
		final.Append(auto)
		final.Append(t.Suffix)

		remote := final.String()
		for _, registry := range *registries {
			tags = append(tags, t.name(registry, docker, remote))
		}
		for _, registry := range t.AllRegistries() {
			tags = append(tags, t.name(registry, docker, remote))
		}
	}

	return
}

func (t *Target) name(registry *Registry, docker *Docker, remote string) (name string) {
	if registry.DockerHub() { // 缩短名字长度
		name = fmt.Sprintf("%s:%s", docker.Repository, remote)
	} else {
		name = fmt.Sprintf("%s/%s:%s", registry.Hostname, docker.Repository, remote)
	}

	return
}

func (t *Target) autos() (tags map[string]string) {
	tags = make(map[string]string, 3)
	tags[t.Tag] = t.Tag
	if !t.Auto {
		return
	}

	autos := strings.Split(t.Tag, constant.Common)
	for index := range autos {
		tag := strings.Join(autos[0:index+1], constant.Common)
		tags[tag] = tag
	}

	if "" != t.Prefix || "" != t.Suffix {
		tags["latest"] = ""
	} else {
		tags["latest"] = "latest"
	}

	return
}
