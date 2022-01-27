package main

import (
	`github.com/dronestock/drone`
)

func (p *plugin) build() (undo bool, err error) {
	args := []string{
		`build`,
		`--rm=true`,
		`--file`, p.Dockerfile,
		`--tag`, p.Name,
	}

	// 编译上下文
	args = append(args, p.context())

	// 精减导数
	if p.squash() {
		args = append(args, `--squash`)
	}
	// 压缩
	if p.Compress {
		args = append(args, `--compress`)
	}
	// 添加标签
	for _, label := range p.labels() {
		args = append(args, `--label`, label)
	}

	// 执行代码检查命令
	err = p.Exec(exe, drone.Args(args...), drone.Dir(p.Context))

	return
}
