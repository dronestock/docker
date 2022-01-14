package main

import (
	`fmt`
	`strings`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func push(conf *config, logger simaqian.Logger) (err error) {
	if `` == strings.TrimSpace(conf.Repository) {
		return
	}

	for _, _tag := range conf.tags() {
		target := fmt.Sprintf(`%s:%s`, conf.Repository, _tag)
		tagArgs := []string{
			`tag`,
			conf.Name,
			target,
		}

		// 记录启动日志，方便调试
		fields := gox.Fields{
			field.String(`name`, conf.Name),
			field.String(`registry`, conf.Registry),
			field.String(`repository`, conf.Repository),
			field.String(`tag`, _tag),
		}

		tagOptions := gex.NewOptions(gex.Args(tagArgs...))
		if !conf.Verbose {
			tagOptions = append(tagOptions, gex.Quiet())
		}
		if _, err = gex.Run(conf.exe, tagOptions...); nil != err {
			logger.Error(`给Docker打标签出错`, fields.Connect(field.Error(err))...)
		}
		if nil != err {
			return
		}

		// 推关到仓库
		pushArgs := []string{
			`push`,
			target,
		}

		logger.Info(`开始推送镜像到仓库`, fields...)
		if _, err = gex.Run(conf.exe, gex.Args(pushArgs...)); nil != err {
			logger.Error(`推送Docker镜像到仓库出错`, fields.Connect(field.Error(err))...)
		}
		if nil != err {
			return
		}
		logger.Info(`推送镜像到仓库成功`, fields...)
	}

	return
}
