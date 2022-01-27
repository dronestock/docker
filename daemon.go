package main

import (
	`fmt`
	`os`

	`github.com/dronestock/drone`
)

func (p *plugin) daemon() (undo bool, err error) {
	// 不必要不启动守护进程
	if _, statErr := os.Stat(outsideDockerfile); nil == statErr {
		undo = true
	}
	if undo {
		return
	}

	args := []string{
		"--data-root", p.DataRoot,
		fmt.Sprintf(`--host=%s`, p.Host),
	}

	if _, statErr := os.Stat("/etc/docker/default.json"); nil == statErr {
		args = append(args, "--seccomp-profile=/etc/docker/default.json")
	}

	// 驱动
	if `` != p.StorageDriver {
		args = append(args, "storage-driver", p.StorageDriver)
	}
	// 镜像加速
	for _, mirror := range p.mirrors() {
		args = append(args, "--registry-mirror", mirror)
	}

	// 启用实验性功能
	if p.Experimental {
		args = append(args, "--experimental")
	}

	// 执行代码检查命令
	err = p.Exec(daemonExe, drone.Args(args...), drone.Contains(daemonSuccessMark), drone.Async(), drone.Dir(p.Context))

	return
}
