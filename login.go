package main

import (
	`strings`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func login(conf *config, logger simaqian.Logger) (err error) {
	if `` == strings.TrimSpace(conf.Username) || `` == strings.TrimSpace(conf.Password) {
		return
	}

	args := []string{
		`login`,
		`--username`, conf.Username,
		`--password`, conf.Password,
		conf.Registry,
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`registry`, conf.Registry),
		field.String(`username`, conf.Username),
	}
	logger.Info(`开始登录Docker仓库`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.ContainsChecker(conf.loginSuccessMark), gex.Async())
	if !conf.Verbose {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(conf.exe, options...); nil != err {
		logger.Error(`登录Docker仓库出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`登录Docker仓库成功`, fields...)

	return
}
