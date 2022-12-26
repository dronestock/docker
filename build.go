package main

import (
	"path/filepath"
	"strings"

	"github.com/dronestock/drone"
)

func (p *plugin) build() (undo bool, err error) {
	args := []interface{}{
		"build",
		"--rm=true",
		"--file", p.Dockerfile,
		"--tag", p.tag(),
	}

	// 编译上下文
	args = append(args, p.context())

	// 精减导数
	if p.squash() {
		args = append(args, "--squash")
	}
	// 压缩
	if p.Compress {
		args = append(args, "--compress")
	}

	// 添加标签
	// 通过只添加一个复合标签来减少层
	args = append(args, "--label", strings.Join(p.labels(), " "))

	// 使用本地网络
	args = append(args, "--network", "host")

	// 执行代码检查命令
	err = p.Exec(exe, drone.Args(args...), drone.Dir(filepath.Dir(p.Dockerfile)))

	return
}
