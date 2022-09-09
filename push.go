package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (p *plugin) push() (undo bool, err error) {
	if undo = `` == strings.TrimSpace(p.Repository); undo {
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(len(p.Registries))
	for _, tag := range p.tags() {
		for _, _registry := range p.Registries {
			go p.pushToRegistry(_registry, tag, wg, &err)
		}
	}

	// 等待所有任务执行完成
	wg.Wait()

	return
}

func (p *plugin) pushToRegistry(registry registry, tag string, wg *sync.WaitGroup, err *error) {
	target := fmt.Sprintf(`%s/%s:%s`, registry.Hostname, p.Repository, tag)
	fields := gox.Fields{
		field.String(`registry`, registry.Hostname),
		field.Strings(`tag`, tag),
	}

	if tagErr := p.Exec(exe, drone.Args(`tag`, p.tag(), target)); nil != tagErr {
		// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
		target = p.tag()
	}

	pushErr := p.Exec(exe, drone.Args(`push`, target))
	if nil != pushErr {
		if registry.Required {
			*err = pushErr
		}
		p.Info(`推送镜像失败`, fields.Connect(field.Error(*err))...)
	} else {
		p.Info(`推送镜像成功`, fields...)
	}

	// 减少等待个数
	wg.Done()
}
