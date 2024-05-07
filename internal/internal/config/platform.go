package config

import (
	"strings"

	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/gox"
)

type Platform struct {
	// 操作系统
	Os string `default:"${OS=linux}" json:"os,omitempty" validate:"omitempty,oneof=linux windows"`
	// 架构
	Arch string `default:"${ARCH}" json:"arch,omitempty" validate:"omitempty,oneof=amd64 i386 arm/v7 arm64"`
}

func (p *Platform) Argument() (platform string) {
	if "" != p.Os && "" != p.Arch {
		platform = gox.StringBuilder(p.Os, constant.Slash, p.Arch).String()
	}

	return
}

func (p *Platform) Qemu() bool {
	return strings.Contains(p.Arch, "arm")
}
