package main

import (
	`strings`

	`github.com/dronestock/drone`
)

func (p *plugin) login() (undo bool, err error) {
	if undo = `` == strings.TrimSpace(p.Username) || `` == strings.TrimSpace(p.Password); undo {
		return
	}

	// 组装参数
	args := []string{
		`login`,
		`--username`, p.Username,
		`--password`, p.Password,
		p.Registry,
	}

	// 执行代码检查命令
	err = p.Exec(daemonExe, drone.Args(args...), drone.Contains(loginSuccessMark), drone.Async(), drone.Dir(p.Context))

	return
}
