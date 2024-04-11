package step

import (
	"context"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

type Login struct {
	base       *drone.Base
	docker     *config.Docker
	registries config.Registries
	targets    config.Targets
	logger     log.Logger
}

func NewLogin(
	base *drone.Base,
	docker *config.Docker, registries config.Registries, targets config.Targets,
	logger log.Logger,
) *Login {
	return &Login{
		base:       base,
		docker:     docker,
		registries: registries,
		targets:    targets,
		logger:     logger,
	}
}

func (l *Login) Runnable() bool {
	return 0 != len(l.registries) || l.targets.Runnable()
}

func (l *Login) Run(ctx *context.Context) (err error) {
	registries := make(config.Registries, 0, len(l.registries)+len(l.targets.Registries()))
	registries = append(l.registries, l.targets.Registries()...)
	for _, registry := range registries {
		l.login(ctx, registry, &err)
	}

	return
}

func (l *Login) login(ctx *context.Context, registry *config.Registry, err *error) {
	la := args.New().Build()
	la.Subcommand("login")
	la.Arg("username", registry.Username)
	la.Arg("password", registry.Password)
	la.Add(registry.Hostname)

	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("username", registry.Username),
		field.New("name", registry.Nickname()),
	}
	l.logger.Info("准备登录镜像仓库", fields...)

	_, le := l.base.Command(l.docker.Exe).Context(*ctx).Args(la.Build()).Checker().Contains(registry.Mark).Build().Exec()
	if nil != le && registry.Required {
		*err = le
		l.logger.Info("登录镜像仓库失败", fields.Add(field.Error(*err))...)
	} else {
		l.logger.Info("登录镜像仓库成功", fields...)
	}
}
