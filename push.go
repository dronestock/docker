package main

import (
	`fmt`
	`strings`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func (p *plugin) push(logger simaqian.Logger) (err error) {
	if `` == strings.TrimSpace(p.config.Repository) {
		return
	}

	for _, _tag := range p.config.tags() {
		target := fmt.Sprintf(`%s/%s:%s`, p.config.Registry, p.config.Repository, _tag)
		tagArgs := []string{
			`tag`,
			p.config.Name,
			target,
		}

		// 记录启动日志，方便调试
		fields := gox.Fields{
			field.String(`name`, p.config.Name),
			field.String(`registry`, p.config.Registry),
			field.String(`repository`, p.config.Repository),
			field.String(`tag`, _tag),
		}

		tagOptions := gex.NewOptions(gex.Args(tagArgs...))
		if !p.config.Verbose {
			tagOptions = append(tagOptions, gex.Quiet())
		}
		if _, err = gex.Run(p.config.exe, tagOptions...); nil != err {
			logger.Error(`镜像打标签出错`, fields.Connect(field.Error(err))...)
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
		if _, err = gex.Run(p.config.exe, gex.Args(pushArgs...)); nil != err {
			logger.Error(`推送Docker镜像到仓库出错`, fields.Connect(field.Error(err))...)
		}
		if nil != err {
			return
		}
		logger.Info(`推送镜像到仓库成功`, fields...)
	}

	return
}
