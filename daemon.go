package main

import (
	`fmt`
	`os`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func daemon(conf *config, logger simaqian.Logger) (err error) {
	args := []string{
		fmt.Sprintf(`--host=%s`, conf.Host),
	}

	if _, statErr := os.Stat("/etc/docker/default.json"); nil == statErr {
		args = append(args, "--seccomp-profile=/etc/docker/default.json")
	}

	// 镜像加速
	for _, mirror := range conf.mirrors() {
		args = append(args, "--registry-mirror", mirror)
	}

	// 启用实验性功能
	if conf.Experimental {
		args = append(args, "--experimental")
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`host`, conf.Host),
		field.Strings(`mirrors`, conf.mirrors()...),
	}
	logger.Info(`开始启动Docker守护进程`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.ContainsChecker(conf.daemonSuccessMark), gex.Async())
	if !conf.Verbose {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(conf.daemon, options...); nil != err {
		logger.Error(`启动Docker守护进程出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`启动Docker守护进程成功`, fields...)

	return
}
