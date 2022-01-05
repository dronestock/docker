package main

import (
	`github.com/storezhang/simaqian`
)

func setup(conf *config, logger simaqian.Logger) (err error) {
	// 加速Github
	if err = github(conf, logger); nil != err {
		return
	}
	// 清理目录
	if err = clear(conf, logger); nil != err {
		return
	}
	// 配置SSH
	if err = ssh(conf, logger); nil != err {
		return
	}

	return
}
