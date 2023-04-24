package main

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
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

func (d *stepDaemon) Run(ctx context.Context) (err error) {
	times := 0
	for {
		times++
		if err = d.startup(ctx); nil != err { // 启动不成功
			break
		} else if err = d.check(ctx); nil == err || errors.Is(err, context.DeadlineExceeded) { // 超时或者检查不成功
			break
		} else {
			d.Info(
				"等待Docker启动完成",
				field.New("timeout", d.Timeout), field.New("elapsed", d.Elapsed()),
				field.New("times", times),
				field.New("address", d.address()),
			)
		}
		time.Sleep(5 * time.Second)
	}

	return
}

func (d *stepDaemon) startup(_ context.Context) (err error) {
	da := args.New().Build()
	da.Arg("data-root", d.Data)
	da.Arg("host", d.address())

	if _, se := os.Stat("/etc/docker/default.json"); nil == se {
		da.Arg("seccomp-profile", "/etc/docker/default.json")
	}

	// 驱动
	if "" != d.Driver {
		da.Arg("storage-driver", d.Driver)
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
	_, err = d.Command(daemonExe).Args(da.Build()).Checker().Contains(d.Mark).Dir(d.context()).Build().Exec()

	return
}

func (d *stepDaemon) check(_ context.Context) (err error) {
	ia := args.New().Build().Subcommand("info").Build()
	_, err = d.Command(exe).Args(ia).Dir(d.context()).Build().Exec()

	return
}
