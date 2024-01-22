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
	base       *drone.Base
	registries []*config.Registry
	docker     *config.Docker
	logger     log.Logger
}

func NewLogin(base *drone.Base, docker *config.Docker, registries []*config.Registry, logger log.Logger) *Login {
	return &Login{
		base:       base,
		docker:     docker,
		registries: registries,
		logger:     logger,
	}
}

func (l *Login) Runnable() bool {
	return 0 != len(l.registries)
}

func (l *Login) Run(ctx *context.Context) (err error) {
	for _, registry := range l.registries {
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
	}
	dir := (*ctx).Value(constant.KeyDir).(string)
	_, le := l.base.Command(constant.Exe).Args(la.Build()).Checker().Contains(registry.Mark).Dir(dir).Build().Exec()
	if nil != le && registry.Required {
		*err = le
		l.logger.Info("登录镜像仓库失败", fields.Add(field.Error(*err))...)
	} else {
		l.logger.Info("登录镜像仓库成功", fields...)
	}
}
