package main

import (
	`github.com/storezhang/gex`
	`github.com/storezhang/simaqian`
)

func (p *plugin) info(logger simaqian.Logger) (err error) {
	if _, err = gex.Run(p.config.exe, gex.Args(`info`)); nil != err {
		logger.Error(`获得Docker信息出错`)
	}

	return
}
