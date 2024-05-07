package step

import (
	"context"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
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
