package main

import (
	`fmt`

	`github.com/storezhang/simaqian`
)

func pull(conf *config, logger simaqian.Logger) (err error) {
	// 克隆项目
	cloneArgs := []string{`clone`, conf.Remote}
	if conf.Submodules {
		cloneArgs = append(cloneArgs, `--remote-submodules`, `--recurse-submodules`)
	}
	if 0 != conf.Depth {
		cloneArgs = append(cloneArgs, `--depth`, fmt.Sprintf(`%d`, conf.Depth))
	}
	cloneArgs = append(cloneArgs, conf.Folder)
	if err = git(conf, logger, cloneArgs...); nil != err {
		return
	}
	// 检出提交的代码
	err = git(conf, logger, `checkout`, conf.checkout())

	return
}
