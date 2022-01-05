package main

import (
	`os`
	`os/exec`
	`path/filepath`

	`github.com/storezhang/gox/field`
	`github.com/storezhang/simaqian`
)

func git(conf *config, logger simaqian.Logger, args ...string) (err error) {
	cmd := exec.Command(`git`, args...)
	if cmd.Dir, err = filepath.Abs(conf.Folder); nil != err {
		return
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, conf.envs...)
	if conf.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	fields := conf.Fields().Connect(field.Strings(`args`, args...))
	if err = cmd.Run(); nil != err {
		logger.Error(`执行Git命令出错`, fields.Connect(field.Error(err))...)
	} else {
		logger.Info(`执行Git命令成功`, fields...)
	}

	return
}
