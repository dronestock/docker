package main

import (
	"context"
	"strings"
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
	args := []any{
		"build",
		"--rm=true",
		"--file", b.Dockerfile,
		"--tag", b.tag(),
	}

	// 编译上下文
	args = append(args, b.context())

	// 精减导数
	if b.squash() {
		args = append(args, "--squash")
	}
	// 压缩
	if b.Compress {
		args = append(args, "--compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	args = append(args, "--label", strings.Join(b.labels(), " "))

	// 使用本地网络
	args = append(args, "--network", "host")

	// 执行代码检查命令
	err = b.Command(exe).Args(args...).Dir(b.context()).Exec()

	return
}
