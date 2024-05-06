package config

import (
	"path/filepath"

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

	// 操作系统
	Os string `default:"${OS=linux}" json:"os,omitempty" validate:"omitempty,oneof=linux windows"`
	// 架构
	Arch string `default:"${ARCH}" json:"arch,omitempty" validate:"omitempty,oneof=amd64 i386 arm32v7 arm64v8"`

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

func (t *Target) Dir() (context string) {
	if "" == t.Context {
		context = filepath.Dir(t.Dockerfile)
	} else {
		context = t.Context
	}

	return
}

func (t *Target) Platform() (platform string) {
	if "" != t.Os && "" != t.Arch {
		platform = gox.StringBuilder(t.Os, constant.Slash, t.Arch).String()
	}

	return
}
