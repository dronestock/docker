package main

import (
	"context"
	"strings"

	"github.com/goexl/gox/args"
)

type stepBuild struct {
	*plugin
}

func newBuildStep(plugin *plugin) *stepBuild {
	return &stepBuild{
		plugin: plugin,
	}
}

func (b *stepBuild) Runnable() bool {
	return 0 != len(b.Registries)
}

func (b *stepBuild) Run(_ context.Context) (err error) {
	ba := args.New().Build()
	ba.Subcommand("build")
	ba.Arg("rm", "true")
	ba.Arg("file", b.Dockerfile)
	ba.Arg("tag", b.tag())

	// 编译上下文
	ba.Add(b.context())

	// 精减导数
	if b.squash() {
		ba.Flag("squash")
	}
	// 压缩
	if b.Compress {
		ba.Flag("compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	ba.Arg("label", strings.Join(b.labels(), space))

	// 使用本地网络
	ba.Arg("network", "host")

	// 执行代码检查命令
	_, err = b.Command(exe).Args(ba.Build()).Dir(b.context()).Build().Exec()

	return
}
