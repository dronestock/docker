package main

import (
	`github.com/storezhang/gex`
	`github.com/storezhang/simaqian`
)

func info(conf *config, logger simaqian.Logger) (err error) {
	// 执行命令
	if _, err = gex.Run(conf.exe, gex.Args(`info`)); nil != err {
		logger.Error(`获得Docker信息出错`)
	}

	return
}
