package main

import (
	`fmt`
	`strings`

	`github.com/storezhang/gex`
	`github.com/storezhang/gox`
	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func (p *plugin) push(logger simaqian.Logger) (undo bool, err error) {
	undo = `` == strings.TrimSpace(p.config.Repository)
	if undo {
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
		if !p.config.Debug {
			tagOptions = append(tagOptions, gex.Quiet())
		}
		if _, err = gex.Run(exe, tagOptions...); nil != err {
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

		// 记录日志
		logger.Info(`开始推送镜像到仓库`, fields...)

		pushOptions := gex.NewOptions(gex.Args(pushArgs...))
		if !p.config.Debug {
			pushOptions = append(pushOptions, gex.Quiet())
		}
		if _, err = gex.Run(exe, pushOptions...); nil != err {
			logger.Error(`推送Docker镜像到仓库出错`, fields.Connect(field.Error(err))...)
		}
		if nil != err {
			return
		}
		logger.Info(`推送镜像到仓库成功`, fields...)
	}

	return
}
