package step

import (
	"context"
	"fmt"
	"strings"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/guc"
)

type Push struct {
	command    *command.Docker
	config     *config.Docker
	registries config.Registries
	targets    config.Targets
}

func NewPush(
	command *command.Docker, config *config.Docker,
	registries config.Registries, targets config.Targets,
) *Push {
	return &Push{
		command:    command,
		registries: registries,
		config:     config,
		targets:    targets,
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
	wg := new(guc.WaitGroup)
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
	local string, remote string,
	wg *guc.WaitGroup, err *error,
) {
	// 任何情况下，都必须调用完成方法
	defer wg.Done()

	image := fmt.Sprintf("%s/%s:%s", registry.Hostname, p.config.Repository, remote)
	fields := gox.Fields[any]{
		field.New("registry", registry.Hostname),
		field.New("repository", p.config.Repository),
		field.New("tag", remote),
		field.New("image", image),
	}

	if te := p.command.Exec(*ctx, args.New().Build().Subcommand("tag").Add(local, image).Build()); nil != te {
		// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
		image = local
	}

	pe := p.command.Exec(*ctx, args.New().Build().Subcommand("push").Add(image).Build())
	if nil != pe && registry.Required {
		*err = pe
		p.command.Info("推送镜像失败", fields.Add(field.Error(*err))...)
	} else {
		p.command.Info("推送镜像成功", fields...)
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
