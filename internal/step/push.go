package step

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

type Push struct {
	base       *drone.Base
	docker     *config.Docker
	registries []*config.Registry
	logger     log.Logger
}

func NewPush(base *drone.Base, docker *config.Docker, registries []*config.Registry, logger log.Logger) *Push {
	return &Push{
		base:       base,
		docker:     docker,
		registries: registries,
		logger:     logger,
	}
}

func (p *Push) Runnable() bool {
	return "" != strings.TrimSpace(p.docker.Repository)
}

func (p *Push) Run(ctx *context.Context) (err error) {
	tags := p.tags()
	wg := new(sync.WaitGroup)
	wg.Add(len(p.registries) * len(tags))
	for _, tag := range tags {
		for _, registry := range p.registries {
			go p.push(ctx, registry, tag, wg, &err)
		}
	}
	// 等待所有任务执行完成
	wg.Wait()

	return
}

func (p *Push) push(ctx *context.Context, registry *config.Registry, tag string, wg *sync.WaitGroup, err *error) {
	// 任何情况下，都必须调用完成方法
	defer wg.Done()

	image := fmt.Sprintf("%s/%s:%s", registry.Hostname, p.docker.Repository, tag)
	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("repository", p.docker.Repository),
		field.New("tag", tag),
		field.New("image", image),
	}

	original := (*ctx).Value(constant.KeyTag).(string)
	ta := args.New().Build().Subcommand("tag").Add(original, image).Build()
	if _, te := p.base.Command(constant.Exe).Args(ta).Build().Exec(); nil != te {
		// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
		image = original
	} else { // ! 清理打包好的镜像（垃圾文件，不清理会导致磁盘空间占用过大）
		ida := args.New().Build().Subcommand("image", "rm", tag)
		p.base.Cleanup().Command(constant.Exe).Args(ida.Build()).Build().Build()
	}

	pa := args.New().Build().Subcommand("push").Add(image).Build()
	_, pe := p.base.Command(constant.Exe).Args(pa).Build().Exec()
	if nil != pe && registry.Required {
		*err = pe
		p.logger.Info("推送镜像失败", fields.Add(field.Error(*err))...)
	} else {
		p.logger.Info("推送镜像成功", fields...)
	}
}

func (p *Push) tags() (tags map[string]string) {
	tags = make(map[string]string, 3)
	tags[p.docker.Tag] = p.docker.Tag
	if !p.docker.Auto {
		return
	}

	autos := strings.Split(p.docker.Tag, constant.Common)
	for index := range autos {
		tag := strings.Join(autos[0:index+1], constant.Common)
		tags[tag] = tag
	}
	tags["latest"] = "latest"

	return
}
