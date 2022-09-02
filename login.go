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

	// 登录仓库
	for _, _registry := range p.Registries {
		args := []interface{}{
			`login`,
			`--username`, _registry.Username,
			`--password`, _registry.Password,
			_registry.Hostname,
		}

		options := drone.NewExecOptions(
			drone.Args(args...),
			drone.Contains(loginSuccessMark),
			drone.Async(),
			drone.Dir(p.Context),
		)
		loginErr := p.Exec(exe, options...)

		fields := gox.Fields{
			field.String(`registry`, _registry.Hostname),
			field.Strings(`username`, _registry.Username),
		}
		if nil != loginErr && _registry.Required {
			err = loginErr
			p.Info(`登录镜像仓库失败`, fields.Connect(field.Error(err))...)
		} else if 1 < len(p.Registries) {
			p.Info(`登录镜像仓库成功`, fields...)
		}
		if nil != err {
			return
		}
	}

	return
}
