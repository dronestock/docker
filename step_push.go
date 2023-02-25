package main

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type stepPush struct {
	*plugin
}

func newPushStep(plugin *plugin) *stepPush {
	return &stepPush{
		plugin: plugin,
	}
}

func (p *stepPush) Runnable() bool {
	return "" != strings.TrimSpace(p.Repository)
}

func (p *stepPush) Run(_ context.Context) (err error) {
	tags := p.tags()
	wg := new(sync.WaitGroup)
	wg.Add(len(p.Registries) * len(tags))
	for _, tag := range tags {
		for _, _registry := range p.Registries {
			go p.push(_registry, tag, wg, &err)
		}
	}

	// 等待所有任务执行完成
	wg.Wait()

	return
}

func (p *stepPush) push(registry registry, tag string, wg *sync.WaitGroup, err *error) {
	// 任何情况下，都必须调用完成方法
	defer wg.Done()

	target := fmt.Sprintf("%s/%s:%s", registry.Hostname, p.Repository, tag)
	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("tag", tag),
	}

	if te := p.Command(exe).Args("tag", p.tag(), target).Exec(); nil != te {
		// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
		target = p.tag()
	}

	pe := p.Command(exe).Args("push", target).Exec()
	if nil != pe {
		if registry.Required {
			*err = pe
		}
		p.Info("推送镜像失败", fields.Add(field.Error(*err))...)
	} else {
		p.Info("推送镜像成功", fields...)
	}
}
