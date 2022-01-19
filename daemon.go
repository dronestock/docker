package main

import (
	`fmt`
	`os`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func (p *plugin) daemon(logger simaqian.Logger) (undo bool, err error) {
	// 不必要不启动守护进程
	if _, statErr := os.Stat(outsideDockerfile); nil == statErr {
		undo = true
	}
	if undo {
		return
	}

	args := []string{
		"--data-root", p.config.DataRoot,
		fmt.Sprintf(`--host=%s`, p.config.Host),
	}

	if _, statErr := os.Stat("/etc/docker/default.json"); nil == statErr {
		args = append(args, "--seccomp-profile=/etc/docker/default.json")
	}

	// 驱动
	if `` != p.config.StorageDriver {
		args = append(args, "storage-driver", p.config.StorageDriver)
	}
	// 镜像加速
	for _, mirror := range p.config.mirrors() {
		args = append(args, "--registry-mirror", mirror)
	}

	// 启用实验性功能
	if p.config.Experimental {
		args = append(args, "--experimental")
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`host`, p.config.Host),
		field.Strings(`mirrors`, p.config.mirrors()...),
	}
	logger.Info(`开始启动Docker守护进程`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.ContainsChecker(daemonSuccessMark), gex.Async())
	if !p.config.Debug {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(daemonExe, options...); nil != err {
		logger.Error(`启动Docker守护进程出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`启动Docker守护进程成功`, fields...)

	return
}
