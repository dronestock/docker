package step

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/docker/internal/internal/key"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/guc"
)

type Build struct {
	command *command.Docker
	config  *config.Docker

	targets    *config.Targets
	project    *config.Project
	registries *config.Registries
}

func NewBuild(
	command *command.Docker, config *config.Docker,
	targets *config.Targets, project *config.Project, registries *config.Registries,
) *Build {
	return &Build{
		command: command,
		config:  config,

		project:    project,
		targets:    targets,
		registries: registries,
	}
}

func (b *Build) Runnable() bool {
	return 0 != len(*b.targets)
}

func (b *Build) Run(ctx *context.Context) (err error) {
	wg := new(guc.WaitGroup)
	wg.Add(len(*b.targets))
	for _, target := range *b.targets {
		go b.run(ctx, target, wg, &err)
	}
	// 等待所有任务执行完成
	wg.Wait()

	return
}

func (b *Build) run(ctx *context.Context, target *config.Target, wg *guc.WaitGroup, err *error) {
	defer wg.Done()

	directory := target.Dir()
	tags := target.Tags(b.registries, b.config)
	arguments := args.New().Build()

	arguments.Subcommand("buildx", "build")
	arguments.Argument("rm", "true")
	arguments.Argument("file", target.Dockerfile)
	for _, tag := range tags {
		arguments.Argument("tag", tag)
	}

	// 编译上下文
	arguments.Add(directory)
	// 精减层数
	if b.squash() {
		arguments.Flag("squash")
	}
	// 压缩
	if b.config.Compress {
		arguments.Flag("compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	arguments.Argument("label", strings.Join(b.labels(target), constant.Space))

	// 多平台编译
	if "" != target.PlatformArgument() {
		arguments.Argument(constant.Platform, target.PlatformArgument())
	}

	// 直接推送
	pushable := target.Pushable(b.registries, b.config)
	if pushable {
		arguments.Flag("push")
	}

	fields := gox.Fields[any]{
		field.New("tags", tags),
		field.New("target", target),
		field.New("push", pushable),
	}
	dir := context.WithValue(*ctx, key.ContextDir, directory)
	b.command.Info("编译镜像开始", fields...)
	if be := b.command.Exec(dir, arguments.Build()); nil != be {
		*err = be
		b.command.Warn("编译镜像出错", fields.Add(field.Error(be))...)
	} else {
		b.command.Info("编译镜像成功", fields...)
	}
}

func (b *Build) labels(target *config.Target) (labels []string) {
	labels = make([]string, 0, 4+len(b.config.Labels))
	if b.command.Default() {
		labels = append(labels, fmt.Sprintf("created=%s", time.Now().Format(time.RFC3339)))
		labels = append(labels, fmt.Sprintf("revision=%s", target.Name))
		labels = append(labels, fmt.Sprintf("source=%s", b.project.Remote))
		labels = append(labels, fmt.Sprintf("url=%s", b.project.Link))
	}
	labels = append(labels, b.config.Labels...)

	return
}

func (b *Build) squash() bool {
	return b.config.Experimental && b.config.Squash
}
