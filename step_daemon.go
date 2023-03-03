package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dronestock/drone"
	"github.com/goexl/gox"
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
	return unix == d.Protocol
}

func (d *stepDaemon) Run(_ context.Context) (err error) {
	args := gox.Args{
		"--data-root", d.DataRoot,
		fmt.Sprintf("--host=%s", d.unix()),
	}

	if _, se := os.Stat("/etc/docker/default.json"); nil == se {
		args.Add("--seccomp-profile=/etc/docker/default.json")
	}

	// 驱动
	if "" != d.StorageDriver {
		args.Add("storage-driver", d.StorageDriver)
	}
	// 镜像加速
	for _, mirror := range d.mirrors() {
		args.Add("--registry-mirror", mirror)
	}

	// 启用实验性功能
	if d.Experimental {
		args.Add("--experimental")
	}

	// 使用阿里DNS
	args.Add("--dns", "223.5.5.5")

	// 执行代码检查命令
	err = d.Command(daemonExe).Args(args...).Checker(drone.Contains(daemonSuccessMark)).Async().Dir(d.context()).Build().Exec()

	return
}
