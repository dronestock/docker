package main

import (
	`fmt`
	`strings`

	`github.com/dronestock/drone`
)

func (p *plugin) push() (undo bool, err error) {
	if undo = `` == strings.TrimSpace(p.Repository); undo {
		return
	}

	for _, tag := range p.tags() {
		target := fmt.Sprintf(`%s/%s:%s`, p.Registry, p.Repository, tag)
		if err = p.Exec(exe, drone.Args(`tag`, p.tag(), target)); nil != err {
			// 如果命令失败，退化成推送已经打好的镜像，不指定仓库
			target = p.tag()
		}
		if err = p.Exec(exe, drone.Args(`push`, target)); nil != err {
			return
		}
	}

	return
}
