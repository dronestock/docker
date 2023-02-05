package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dronestock/drone"
)

type stepDaemon struct {
	*plugin
}

func newDaemonStep(plugin *plugin) *stepDaemon {
	return &stepDaemon{
		plugin: plugin,
	}
}

func (d *stepDaemon) Runnable() (runnable bool) {
	if _, se := os.Stat(outsideDockerfile); nil != se && os.IsExist(se) {
		runnable = true
	}

	return
}

func (d *stepDaemon) Run(_ context.Context) (err error) {
	args := []any{
		"--data-root", d.DataRoot,
		fmt.Sprintf("--host=%s", d.Host),
	}

	if _, se := os.Stat("/etc/docker/default.json"); nil == se {
		args = append(args, "--seccomp-profile=/etc/docker/default.json")
	}

	// 驱动
	if "" != d.StorageDriver {
		args = append(args, "storage-driver", d.StorageDriver)
	}
	// 镜像加速
	for _, mirror := range d.mirrors() {
		args = append(args, "--registry-mirror", mirror)
	}

	// 启用实验性功能
	if d.Experimental {
		args = append(args, "--experimental")
	}

	// 使用阿里DNS
	args = append(args, "--dns", "223.5.5.5")

	// 执行代码检查命令
	err = d.Command(daemonExe).Args(args...).Checker(drone.Contains(daemonSuccessMark)).Async().Dir(d.context()).Exec()

	return
}
