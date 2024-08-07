package step

import (
	"context"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/guc"
)

type Push struct {
	command *command.Docker
	config  *config.Docker

	registries *config.Registries
	targets    *config.Targets
}

func NewPush(
	command *command.Docker, config *config.Docker,
	targets *config.Targets, registries *config.Registries,
) *Push {
	return &Push{
		command: command,
		config:  config,

		targets:    targets,
		registries: registries,
	}
}

func (p *Push) Runnable() bool {
	return p.targets.Runnable(p.registries, p.config)
}

func (p *Push) Run(ctx *context.Context) (err error) {
	for _, target := range *p.targets {
		p.run(ctx, target, &err)
	}

	return
}

func (p *Push) run(ctx *context.Context, target *config.Target, err *error) {
	tags := target.Tags(p.registries, p.config)
	waiter := guc.New().Wait().Group(len(tags))
	for _, tag := range tags {
		go p.push(ctx, target, tag, waiter, err)
	}
	// 等待所有任务执行完成
	waiter.Wait()
}

func (p *Push) push(ctx *context.Context, target *config.Target, tag string, waiter guc.Waiter, err *error) {
	// 任何情况下，都必须调用完成方法
	defer waiter.Done()

	if 2 <= len(target.AllPlatforms()) {
		return
	}

	fields := gox.Fields[any]{
		field.New("repository", p.config.Repository),
		field.New("tag", tag),
		field.New("target", target),
	}
	if pe := p.command.Exec(*ctx, args.New().Build().Subcommand("push").Add(tag).Build()); nil != pe {
		*err = pe
		p.command.Info("推送镜像失败", fields.Add(field.Error(*err))...)
	} else {
		p.command.Info("推送镜像成功", fields...)
	}
}
