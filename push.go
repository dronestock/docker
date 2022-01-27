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

	for _, _tag := range p.tags() {
		target := fmt.Sprintf(`%s/%s:%s`, p.Registry, p.Repository, _tag)
		if err = p.Exec(exe, drone.Args(`tag`, p.Name, target), drone.Dir(p.Context)); nil != err {
			return
		}
		if err = p.Exec(exe, drone.Args(`push`, target)); nil != err {
			return
		}
	}

	return
}
