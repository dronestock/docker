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
	Arch string `default:"${ARCH=amd64}" json:"arch,omitempty" validate:"omitempty,oneof=amd64 i386 arm arm64"`
	// 变体
	// nolint:lll
	Variant string `default:"${VARIANT}" json:"variant,omitempty" validate:"omitempty,required_if=Arch arm,oneof=v5 v6 v7"`
}

func (p *Platform) Argument() string {
	builder := gox.StringBuilder()
	if "" != p.Os && "" != p.Arch {
		builder.Append(p.Os).Append(constant.Slash).Append(p.Arch)
	}
	if "" != p.Os && "" != p.Arch && "" != p.Variant {
		builder.Append(constant.Slash).Append(p.Variant)
	}

	return builder.String()
}

func (p *Platform) Qemu() bool {
	return strings.Contains(p.Arch, "arm")
}
