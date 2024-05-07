package config

import (
	"path/filepath"
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/gox"
	"github.com/rs/xid"
)

type Target struct {
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

	// 平台
	Platform Platform `default:"${PLATFORM}" json:"platform,omitempty"`
	// 平台列表
	Platforms Platforms `default:"${PLATFORMS}" json:"platforms,omitempty"`

	// 仓库
	Registry *Registry `json:"registry,omitempty"`
	// 仓库列表
	Registries []*Registry `json:"registries,omitempty"`

	tag string
}

func (t *Target) LocalTag() string {
	if "" == t.tag {
		t.tag = xid.New().String()
	}

	return t.tag
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
