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
	"github.com/rs/xid"
)

type Build struct {
	base       *drone.Base
	registries []*config.Registry
	docker     *config.Docker
	project    *config.Project
	logger     log.Logger
}

func NewBuild(
	base *drone.Base,
	docker *config.Docker, project *config.Project, registries []*config.Registry,
	logger log.Logger,
) *Build {
	return &Build{
		base:       base,
		docker:     docker,
		project:    project,
		registries: registries,
		logger:     logger,
	}
}

func (b *Build) Runnable() bool {
	return 0 != len(b.registries)
}

func (b *Build) Run(ctx *context.Context) (err error) {
	dir := (*ctx).Value(constant.KeyDir).(string)
	tag := b.tag()
	ba := args.New().Build()

	ba.Subcommand("build")
	ba.Arg("rm", "true")
	ba.Arg("file", b.docker.Dockerfile)
	ba.Arg("tag", tag)

	// 编译上下文
	ba.Add(dir)

	// 精减导数
	if b.squash() {
		ba.Flag("squash")
	}
	// 压缩
	if b.docker.Compress {
		ba.Flag("compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	ba.Arg("label", strings.Join(b.labels(), constant.Space))

	// 使用本地网络
	ba.Arg("network", "host")

	// 执行代码检查命令
	if _, err = b.base.Command(constant.Exe).Args(ba.Build()).Dir(dir).Build().Exec(); nil == err {
		*ctx = context.WithValue(*ctx, constant.KeyTag, tag)
		// ! 清理打包好的镜像（垃圾文件，不清理会导致磁盘空间占用过大）
		ida := args.New().Build().Subcommand("image", "rm", tag)
		b.base.Cleanup().Command(constant.Exe).Args(ida.Build()).Build().Name(fmt.Sprintf("删除镜像：%s", tag)).Build()
	}

	return
}

func (b *Build) labels() (labels []string) {
	labels = make([]string, 0, 4+len(b.docker.Labels))
	if b.base.Default() {
		labels = append(labels, fmt.Sprintf("created=%s", time.Now().Format(time.RFC3339)))
		labels = append(labels, fmt.Sprintf("revision=%s", b.docker.Name))
		labels = append(labels, fmt.Sprintf("source=%s", b.project.Remote))
		labels = append(labels, fmt.Sprintf("url=%s", b.project.Link))
	}
	labels = append(labels, b.docker.Labels...)

	return
}

func (b *Build) tag() string {
	return xid.New().String()
}

func (b *Build) squash() bool {
	return b.docker.Experimental && b.docker.Squash
}
