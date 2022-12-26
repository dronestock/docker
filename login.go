package main

import (
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (p *plugin) login() (undo bool, err error) {
	if undo = 0 == len(p.Registries); undo {
		return
	}

	for _, _registry := range p.Registries {
		p.loginRegistry(_registry, &err)
	}

	return
}

func (p *plugin) loginRegistry(registry registry, err *error) {
	args := []interface{}{
		"login",
		"--username", registry.Username,
		"--password", registry.Password,
		registry.Hostname,
	}

	options := drone.NewExecOptions(
		drone.Args(args...),
		drone.Contains(registry.Mark),
		drone.Async(),
		drone.Dir(p.Context),
	)
	loginErr := p.Exec(exe, options...)

	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("username", registry.Username),
	}
	if nil != loginErr {
		if registry.Required {
			*err = loginErr
		}
		p.Info("登录镜像仓库失败", fields.Connect(field.Error(*err))...)
	} else {
		p.Info("登录镜像仓库成功", fields...)
	}
}
