package main

import (
	`github.com/storezhang/simaqian`
)

func push(conf *config, logger simaqian.Logger) (err error) {
	// 设置默认分支
	if err = git(conf, logger, `config`, `--global`, `init.defaultBranch`, `master`); nil != err {
		return
	}
	// 设置用户名
	if err = git(conf, logger, `config`, `--global`, `user.name`, conf.Author); nil != err {
		return
	}
	// 设置邮箱
	if err = git(conf, logger, `config`, `--global`, `user.email`, conf.Email); nil != err {
		return
	}
	// 初始化目录
	if err = git(conf, logger, `init`); nil != err {
		return
	}
	// 添加当前目录到Git中
	if err = git(conf, logger, `add`, `.`); nil != err {
		return
	}
	// 提交
	if err = git(conf, logger, `commit`, `.`, `--message`, conf.Message); nil != err {
		return
	}
	// 添加远程仓库地址
	if err = git(conf, logger, `remote`, `add`, `origin`, conf.Remote); nil != err {
		return
	}
	// 如果有标签，推送标签
	if `` != conf.Tag {
		if err = git(conf, logger, `tag`, `--annotate`, conf.Tag, `--message`, conf.Message); nil != err {
			return
		}
		if err = git(conf, logger, `push`, `--set-upstream`, `origin`, conf.Tag, conf.gitForce()); nil != err {
			return
		}
	}
	// 推送
	if err = git(conf, logger, `push`, `--set-upstream`, `origin`, conf.Branch, conf.gitForce()); nil != err {
		return
	}

	return
}
