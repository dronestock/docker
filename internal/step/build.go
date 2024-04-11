package step

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox/args"
	"github.com/goexl/log"
)

type Build struct {
	base    *drone.Base
	targets config.Targets
	docker  *config.Docker
	project *config.Project
	logger  log.Logger
}

func NewBuild(
	base *drone.Base,
	docker *config.Docker, project *config.Project, targets config.Targets,
	logger log.Logger,
) *Build {
	return &Build{
		base:    base,
		docker:  docker,
		project: project,
		targets: targets,
		logger:  logger,
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
	dir := target.Dir()
	tag := target.LocalTag()
	ba := args.New().Build()

	ba.Subcommand("build")
	ba.Arg("rm", "true")
	ba.Arg("file", target.Dockerfile)
	ba.Arg("tag", tag)

	// 编译上下文
	ba.Add(dir)

	// 精减层数
	if b.squash() {
		// ba.Flag("squash")
	}
	// 压缩
	if b.docker.Compress {
		ba.Flag("compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	ba.Arg("label", strings.Join(b.labels(target), constant.Space))

	// 使用本地网络
	ba.Arg("network", "host")

	// 执行代码检查命令
	if _, ee := b.base.Command(b.docker.Exe).Context(*ctx).Args(ba.Build()).Dir(dir).Build().Exec(); nil == ee {
		// ! 清理打包好的镜像（垃圾文件，不清理会导致磁盘空间占用过大）
		ida := args.New().Build().Subcommand("image", "rm", tag)
		b.base.Cleanup().Command(b.docker.Exe).Args(ida.Build()).Build().Name(fmt.Sprintf("删除镜像：%s", tag)).Build()
	} else {
		*err = ee
	}
}

func (b *Build) labels(target *config.Target) (labels []string) {
	labels = make([]string, 0, 4+len(b.docker.Labels))
	if b.base.Default() {
		labels = append(labels, fmt.Sprintf("created=%s", time.Now().Format(time.RFC3339)))
		labels = append(labels, fmt.Sprintf("revision=%s", target.Name))
		labels = append(labels, fmt.Sprintf("source=%s", b.project.Remote))
		labels = append(labels, fmt.Sprintf("url=%s", b.project.Link))
	}
	labels = append(labels, b.docker.Labels...)

	return
}

func (b *Build) squash() bool {
	return b.docker.Experimental && b.docker.Squash
}
