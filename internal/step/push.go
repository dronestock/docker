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
	registries config.Registries
	targets    config.Targets
	logger     log.Logger
}

func NewPush(
	base *drone.Base,
	docker *config.Docker, registries config.Registries, targets config.Targets,
	logger log.Logger,
) *Push {
	return &Push{
		base:       base,
		registries: registries,
		docker:     docker,
		targets:    targets,
		logger:     logger,
	}
}

func (p *Push) Runnable() bool {
	return 0 != len(p.registries) || p.targets.Runnable()
}

func (p *Push) Run(ctx *context.Context) (err error) {
	for _, target := range p.targets {
		p.run(ctx, target, &err)
	}

	return
}

func (p *Push) run(ctx *context.Context, target *config.Target, err *error) {
	tags := p.tags(target)
	wg := new(sync.WaitGroup)
	wg.Add((len(p.registries) + len(target.AllRegistries())) * len(tags))
	for _, tag := range tags {
		final := gox.StringBuilder(target.Prefix)
		if "" != tag {
			final.Append(target.Middle)
		}
		final.Append(tag)
		final.Append(target.Suffix)

		localTag := target.LocalTag()
		remoteTag := final.String()
		for _, registry := range p.registries {
			go p.push(ctx, registry, localTag, remoteTag, wg, err)
		}
		for _, registry := range target.AllRegistries() {
			go p.push(ctx, registry, localTag, remoteTag, wg, err)
		}
	}
	// 等待所有任务执行完成
	wg.Wait()
}

func (p *Push) push(
	ctx *context.Context,
	registry *config.Registry,
	localTag string, remoteTag string,
	wg *sync.WaitGroup, err *error,
) {
	// 任何情况下，都必须调用完成方法
	defer wg.Done()

	image := fmt.Sprintf("%s/%s:%s", registry.Hostname, p.docker.Repository, remoteTag)
	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("repository", p.docker.Repository),
		field.New("remoteTag", remoteTag),
		field.New("image", image),
	}

	ta := args.New().Build().Subcommand("tag").Add(localTag, image).Build()
	if _, te := p.base.Command(p.docker.Exe).Context(*ctx).Args(ta).Build().Exec(); nil != te {
		// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
		image = localTag
	} else { // ! 清理打包好的镜像（垃圾文件，不清理会导致磁盘空间占用过大）
		ida := args.New().Build().Subcommand("image", "rm", remoteTag)
		p.base.Cleanup().Command(p.docker.Exe).Args(ida.Build()).Build().Name(fmt.Sprintf("删除镜像：%s", remoteTag)).Build()
	}

	pa := args.New().Build().Subcommand("push").Add(image).Build()
	_, pe := p.base.Command(p.docker.Exe).Context(*ctx).Args(pa).Build().Exec()
	if nil != pe && registry.Required {
		*err = pe
		p.logger.Info("推送镜像失败", fields.Add(field.Error(*err))...)
	} else {
		p.logger.Info("推送镜像成功", fields...)
	}
}

func (p *Push) tags(target *config.Target) (tags map[string]string) {
	tags = make(map[string]string, 3)
	tags[target.Tag] = target.Tag
	if !target.Auto {
		return
	}

	autos := strings.Split(target.Tag, constant.Common)
	for index := range autos {
		tag := strings.Join(autos[0:index+1], constant.Common)
		tags[tag] = tag
	}

	if "" != target.Prefix || "" != target.Suffix {
		tags["latest"] = ""
	} else {
		tags["latest"] = "latest"
	}

	return
}
