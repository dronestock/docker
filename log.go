package main

import (
	`fmt`

	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func log(conf *config, logger simaqian.Logger, err error) {
	var msg string
	if conf.pull() {
		msg = `Git拉取`
	} else {
		msg = `Git推送`
	}

	if nil != err {
		logger.Fatal(fmt.Sprintf(`%s失败`, msg), conf.Fields().Connect(field.Error(err))...)
	} else {
		logger.Info(fmt.Sprintf(`%s成功`, msg), conf.Fields()...)
	}
}
