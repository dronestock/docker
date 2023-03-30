package main

import (
	"context"
	"os"

	"github.com/goexl/gox/args"
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
	da := args.New().Build()
	da.Arg("data-root", d.DataRoot)
	da.Arg("host", d.unix())

	if _, se := os.Stat("/etc/docker/default.json"); nil == se {
		da.Arg("seccomp-profile", "/etc/docker/default.json")
	}

	// 驱动
	if "" != d.StorageDriver {
		da.Arg("storage-driver", d.StorageDriver)
	}
	// 镜像加速
	for _, mirror := range d.mirrors() {
		da.Arg("registry-mirror", mirror)
	}

	// 启用实验性功能
	if d.Experimental {
		da.Flag("experimental")
	}

	// 使用阿里DNS
	da.Arg("dns", "223.5.5.5")
	// 执行代码检查命令
	_, err = d.Command(daemonExe).Args(da.Build()).Checker().Contains(daemonSuccessMark).Dir(d.context()).Build().Exec()

	return
}
