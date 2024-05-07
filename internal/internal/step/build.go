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
)

type Build struct {
	targets config.Targets
	config  *config.Docker
	project *config.Project
	command *command.Docker
}

func NewBuild(command *command.Docker, config *config.Docker, project *config.Project, targets config.Targets) *Build {
	return &Build{
		command: command,
		config:  config,
		project: project,
		targets: targets,
	}
}

func (b *Build) Runnable() bool {
	return 0 != len(b.targets)
}

func (b *Build) Run(ctx *context.Context) (err error) {
	for _, target := range b.targets {
		b.run(ctx, target, &err)
	}

	return
}

func (b *Build) run(ctx *context.Context, target *config.Target, err *error) {
	directory := target.Dir()
	tag := target.LocalTag()
	arguments := args.New().Build()

	arguments.Subcommand("build")
	arguments.Argument("rm", "true")
	arguments.Argument("file", target.Dockerfile)
	arguments.Argument("tag", tag)

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

	// 使用本地网络
	arguments.Argument("network", "host")

	// 多平台编译
	if "" != target.PlatformArgument() {
		arguments.Argument(constant.Platform, target.PlatformArgument())
	}

	// 执行代码检查命令
	dir := context.WithValue(*ctx, key.ContextDir, directory)
	*err = b.command.Exec(dir, arguments.Build())
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
