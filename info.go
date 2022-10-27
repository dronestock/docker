package main

import (
	"github.com/dronestock/drone"
)

func (p *plugin) info() (undo bool, err error) {
	// 执行代码检查命令
	err = p.Exec(exe, drone.Args(`info`), drone.Dir(p.Context))

	return
}
