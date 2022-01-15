package main

import (
	`path/filepath`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func build(conf *config, logger simaqian.Logger) (err error) {
	args := []string{
		`build`,
		`--rm=true`,
		`--file`, conf.Dockerfile,
		`--tag`, conf.Name,
	}

	// 编译上下文
	args = append(args, conf.context())

	// 精减导数
	if conf.squash() {
		args = append(args, `--squash`)
	}
	// 压缩
	if conf.Compress {
		args = append(args, `--compress`)
	}
	// 添加标签
	for _, label := range conf.labels() {
		args = append(args, `--label`, label)
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`dockerfile`, conf.Dockerfile),
		field.String(`context`, conf.context()),
		field.String(`name`, conf.Name),
		field.Bool(`squash`, conf.squash()),
		field.Bool(`compress`, conf.Compress),
	}
	logger.Info(`开始编译Dockerfile`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.Dir(filepath.Dir(conf.Dockerfile)))
	if _, err = gex.Run(conf.exe, options...); nil != err {
		logger.Error(`编译Dockerfile出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`编译Dockerfile成功`, fields...)

	return
}
