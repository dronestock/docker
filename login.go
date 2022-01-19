package main

import (
	`strings`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func (p *plugin) login(logger simaqian.Logger) (undo bool, err error) {
	undo = `` == strings.TrimSpace(p.config.Username) || `` == strings.TrimSpace(p.config.Password)
	if undo {
		return
	}

	args := []string{
		`login`,
		`--username`, p.config.Username,
		`--password`, p.config.Password,
		p.config.Registry,
	}

	// 记录启动日志，方便调试
	fields := gox.Fields{
		field.String(`registry`, p.config.Registry),
		field.String(`username`, p.config.Username),
	}
	logger.Info(`开始登录Docker仓库`, fields...)

	// 执行命令
	options := gex.NewOptions(gex.Args(args...), gex.ContainsChecker(loginSuccessMark), gex.Async())
	if !p.config.Debug {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(exe, options...); nil != err {
		logger.Error(`登录Docker仓库出错`, fields.Connect(field.Error(err))...)
	}
	if nil != err {
		return
	}

	// 记录启动功能日志，方便调试
	logger.Info(`登录Docker仓库成功`, fields...)

	return
}
