package step

import (
	"context"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

type Login struct {
	base    *drone.Base
	targets []*config.Target
	docker  *config.Docker
	logger  log.Logger
}

func NewLogin(base *drone.Base, docker *config.Docker, targets []*config.Target, logger log.Logger) *Login {
	return &Login{
		base:    base,
		docker:  docker,
		targets: targets,
		logger:  logger,
	}
}

func (l *Login) Runnable() (runnable bool) {
	for _, target := range l.targets {
		if nil != target.Registry || 0 != len(target.Registries) {
			runnable = true
		}
		if runnable {
			break
		}
	}

	return
}

func (l *Login) Run(ctx *context.Context) (err error) {
	for _, target := range l.targets {
		l.run(ctx, target, &err)
	}

	return
}

func (l *Login) run(ctx *context.Context, target *config.Target, err *error) {
	registries := make([]*config.Registry, 0, len(target.Registries)+1)
	if nil != target.Registry {
		registries = append(registries, target.Registry)
	}
	registries = append(registries, target.Registries...)

	for _, registry := range registries {
		l.login(ctx, registry, err)
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
	}
	l.logger.Info("准备登录镜像仓库", fields...)

	dir := (*ctx).Value(constant.KeyDir).(string)
	_, le := l.base.Command(constant.Exe).Args(la.Build()).Checker().Contains(registry.Mark).Dir(dir).Build().Exec()
	if nil != le && registry.Required {
		*err = le
		l.logger.Info("登录镜像仓库失败", fields.Add(field.Error(*err))...)
	} else {
		l.logger.Info("登录镜像仓库成功", fields...)
	}
}
