package main

import (
	"sync"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (p *plugin) login() (undo bool, err error) {
	if undo = 0 == len(p.Registries); undo {
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(p.Registries))
	for _, _registry := range p.Registries {
		go p.loginRegistry(_registry, wg, &err)
	}

	// 等待所有任务执行完成
	wg.Wait()

	return
}

func (p *plugin) loginRegistry(registry registry, wg *sync.WaitGroup, err *error) {
	args := []interface{}{
		`login`,
		`--username`, registry.Username,
		`--password`, registry.Password,
		registry.Hostname,
	}

	options := drone.NewExecOptions(
		drone.Args(args...),
		drone.Contains(registry.Mark),
		drone.Async(),
		drone.Dir(p.Context),
	)
	loginErr := p.Exec(exe, options...)

	fields := gox.Fields{
		field.String(`registry`, registry.Hostname),
		field.Strings(`username`, registry.Username),
	}
	if nil != loginErr && registry.Required {
		*err = loginErr
		p.Info(`登录镜像仓库失败`, fields.Connect(field.Error(*err))...)
	} else if 1 < len(p.Registries) {
		p.Info(`登录镜像仓库成功`, fields...)
	}

	// 减少等待个数
	wg.Done()
}
