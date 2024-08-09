package step

import (
	"context"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/key"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type Login struct {
	command *command.Docker
	config  *config.Docker

	registries *config.Registries
	targets    *config.Targets
}

func NewLogin(
	command *command.Docker, config *config.Docker,
	registries *config.Registries, targets *config.Targets,
) *Login {
	return &Login{
		command: command,
		config:  config,

		registries: registries,
		targets:    targets,
	}
}

func (l *Login) Runnable() bool {
	return 0 != len(*l.registries)
}

func (l *Login) Run(ctx *context.Context) (err error) {
	registries := make(config.Registries, 0, len(*l.registries)+len(l.targets.Registries()))
	registries = append(registries, *l.registries...)
	registries = append(registries, l.targets.Registries()...)
	for _, registry := range registries {
		l.login(ctx, registry, &err)
	}

	return
}

func (l *Login) login(ctx *context.Context, registry *config.Registry, err *error) {
	arguments := args.New().Build()
	arguments.Subcommand("login")
	arguments.Argument("username", registry.Username)
	arguments.Argument("password", registry.Password)
	arguments.Subcommand(registry.Hostname)

	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("username", registry.Username),
		field.New("name", registry.Nickname()),
	}
	l.command.Info("准备登录镜像仓库", fields...)
	mark := context.WithValue(*ctx, key.ContextMark, registry.Mark)
	if le := l.command.Exec(mark, arguments.Build()); nil != le {
		*err = le
		l.command.Info("登录镜像仓库失败", fields.Add(field.Error(*err))...)
	} else {
		l.command.Info("登录镜像仓库成功", fields...)
	}
}
