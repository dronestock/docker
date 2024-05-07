package step

import (
	"context"
	"strings"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/rs/xid"
)

type Setup struct {
	command *command.Docker
	config  *config.Docker
	targets *config.Targets
}

func NewSetup(command *command.Docker, config *config.Docker, targets *config.Targets) *Setup {
	return &Setup{
		command: command,
		config:  config,
		targets: targets,
	}
}

func (s *Setup) Runnable() bool {
	return true
}

func (s *Setup) Run(ctx *context.Context) (err error) {
	if qe := s.binfmt(ctx); nil != qe {
		err = qe
	} else if de := s.driver(ctx); nil != de {
		err = de
	}

	return
}

func (s *Setup) binfmt(ctx *context.Context) (err error) {
	if !s.targets.Binfmt(s.config) {
		return
	}

	image := "tonistiigi/binfmt"
	arguments := args.New().Build()
	arguments.Subcommand("run")
	arguments.Flag("privileged")
	arguments.Flag("rm")
	arguments.Subcommand(image)
	arguments.Flag("install")
	arguments.Subcommand("all")

	fields := gox.Fields[any]{
		field.New("image", image),
	}
	s.command.Info("准备安装Qemu环境", fields...)
	if err = s.command.Exec(*ctx, arguments.Build()); nil != err {
		s.command.Warn("安装Qemu环境失败", fields.Add(field.Error(err))...)
	} else {
		s.command.Info("安装Qemu环境成功", fields...)
	}

	return
}

func (s *Setup) driver(ctx *context.Context) (err error) {
	platforms := s.targets.Platforms()
	if 1 >= strings.Count(platforms, constant.Comma) { // 只有同时要打包多个平台才需要创建多平台编译驱动
		return
	}

	name := xid.New().String()
	arguments := args.New().Build()
	arguments.Subcommand("buildx")
	arguments.Argument("name", name)
	arguments.Flag("use")
	arguments.Argument("platform", platforms)

	fields := gox.Fields[any]{
		field.New("name", name),
		field.New("platform", platforms),
	}
	s.command.Info("准备创建多平台编译驱动", fields...)
	if err = s.command.Exec(*ctx, arguments.Build()); nil != err {
		s.command.Warn("创建多平台编译驱动失败", fields.Add(field.Error(err))...)
	} else {
		s.command.Info("创建多平台编译驱动成功", fields...)
	}

	return
}
