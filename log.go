package main

import (
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func log(conf *config, logger simaqian.Logger, err error) {
	if nil != err {
		logger.Fatal(`Docker打包失败`, conf.Fields().Connect(field.Error(err))...)
	} else {
		logger.Info(`Docker打包成功`, conf.Fields()...)
	}
}
