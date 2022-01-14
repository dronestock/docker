package main

import (
	`github.com/storezhang/simaqian`
)

func setup(conf *config, logger simaqian.Logger) (err error) {
	// 启动守护进程
	if err = daemon(conf, logger); nil != err {
		return
	}

	// 打印当前Docker信息
	if err = info(conf, logger); nil != err {
		return
	}

	// 登录
	if err = login(conf, logger); nil != err {
		return
	}

	return
}
