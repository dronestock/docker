package main

import (
	`path/filepath`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func (p *plugin) build(logger simaqian.Logger) (undo bool, err error) {
	args := []string{
		`build`,
		`--rm=true`,
		`--file`, p.config.Dockerfile,
		`--tag`, p.config.Name,
	}

	// 编译上下文
	args = append(args, p.config.context())

	// 精减导数
	if p.config.squash() {
		args = append(args, `--squash`)
	}
	// 压缩
	if p.config.Compress {
		args = append(args, `--compress`)
	}
	// 添加标签
	for _, label := range p.config.labels() {
		args = append(args, `--label`, label)
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`dockerfile`, p.config.Dockerfile),
		field.String(`context`, p.config.context()),
		field.String(`name`, p.config.Name),
		field.Bool(`squash`, p.config.squash()),
		field.Bool(`compress`, p.config.Compress),
	}
	logger.Info(`开始编译Dockerfile`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.Dir(filepath.Dir(p.config.Dockerfile)))
	if !p.config.Debug {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(p.exe, options...); nil != err {
		logger.Error(`编译Dockerfile出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`编译Dockerfile成功`, fields...)

	return
}
