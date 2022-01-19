package main

import (
	`github.com/storezhang/gex`
	`github.com/storezhang/simaqian`
)

func (p *plugin) info(logger simaqian.Logger) (undo bool, err error) {
	options := gex.NewOptions(gex.Args(`info`))
	if !p.config.Debug {
		options = append(options, gex.Quiet())
	}
	if _, err = gex.Run(exe, options...); nil != err {
		logger.Error(`获得Docker信息出错`)
	}

	return
}
