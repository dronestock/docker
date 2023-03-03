package main

import (
	"context"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type stepLogin struct {
	*plugin
}

func newLoginStep(plugin *plugin) *stepLogin {
	return &stepLogin{
		plugin: plugin,
	}
}

func (l *stepLogin) Runnable() bool {
	return 0 != len(l.Registries)
}

func (l *stepLogin) Run(_ context.Context) (err error) {
	for _, _registry := range l.Registries {
		l.login(_registry, &err)
	}

	return
}

func (l *stepLogin) login(registry registry, err *error) {
	args := gox.Args{
		"login",
		"--username", registry.Username,
		"--password", registry.Password,
		registry.Hostname,
	}

	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("username", registry.Username),
	}
	le := l.Command(exe).Args(args...).Checker(drone.Contains(registry.Mark)).Async().Dir(l.context()).Build().Exec()
	if nil != le && registry.Required {
		*err = le
		l.Info("登录镜像仓库失败", fields.Add(field.Error(*err))...)
	} else {
		l.Info("登录镜像仓库成功", fields...)
	}
}
