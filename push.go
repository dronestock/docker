package main

import (
	"fmt"
	"strings"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (p *plugin) push() (undo bool, err error) {
	if undo = `` == strings.TrimSpace(p.Repository); undo {
		return
	}

	for _, tag := range p.tags() {
		for _, _registry := range p.Registries {
			target := fmt.Sprintf(`%s/%s:%s`, _registry.Hostname, p.Repository, tag)
			if err = p.Exec(exe, drone.Args(`tag`, p.tag(), target)); nil != err {
				// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
				target = p.tag()
			}

			fields := gox.Fields{
				field.String(`registry`, _registry.Hostname),
				field.Strings(`tag`, tag),
			}
			pushErr := p.Exec(exe, drone.Args(`push`, target))
			if nil != pushErr && _registry.Required {
				err = pushErr
				p.Info(`推送镜像失败`, fields.Connect(field.Error(err))...)
			} else if 1 < len(p.Registries) {
				p.Info(`推送镜像成功`, fields...)
			}
			if nil != err {
				return
			}
		}
	}

	return
}
