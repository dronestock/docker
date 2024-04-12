package step

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dronestock/docker/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/dronestock/drone"
	"github.com/goexl/gox/args"
	"github.com/goexl/gox/field"
	"github.com/goexl/log"
)

var defaultMirrors = []string{
	"https://docker.nju.edu.cn",               // 南京大学镜像站
	"https://docker.m.daocloud.io",            // DaoCloud镜像站
	"https://dockerproxy.com",                 // Docker镜像代理
	"https://mirror.iscas.ac.cn",              // 中科院软件所镜像站
	"https://docker.mirrors.sjtug.sjtu.edu.c", // 上海交大镜像站
	"https://mirror.baidubce.com",             // 百度云
}

type Daemon struct {
	base   *drone.Base
	docker *config.Docker
	logger log.Logger
}

func NewDaemon(base *drone.Base, docker *config.Docker, logger log.Logger) *Daemon {
	return &Daemon{
		base:   base,
		docker: docker,
		logger: logger,
	}
}

func (d *Daemon) Runnable() bool {
	return constant.Unix == d.docker.Protocol
}

func (d *Daemon) Run(ctx *context.Context) (err error) {
	times := 0
	for {
		times++
		if err = d.startup(ctx); nil != err { // 启动不成功
			break
		} else if err = d.check(ctx); nil == err || errors.Is(err, context.DeadlineExceeded) { // 超时或者检查不成功
			break
		} else {
			d.logger.Info(
				"等待Docker启动完成",
				field.New("timeout", d.base.Timeout), field.New("elapsed", d.base.Elapsed()),
				field.New("times", times),
				field.New("address", d.address()),
			)
		}
		time.Sleep(1 * time.Second)
	}

	return
}

func (d *Daemon) startup(ctx *context.Context) (err error) {
	da := args.New().Build()
	da.Arg("data-root", d.docker.Data)
	da.Arg("host", d.address())

	if _, se := os.Stat("/etc/docker/default.json"); nil == se {
		da.Arg("seccomp-profile", "/etc/docker/default.json")
	}

	// 驱动
	if "" != d.docker.Driver {
		da.Arg("storage-driver", d.docker.Driver)
	}
	// 镜像加速
	for _, mirror := range d.mirrors() {
		da.Arg("registry-mirror", mirror)
	}

	// 启用实验性功能
	if d.docker.Experimental {
		da.Flag("experimental")
	}

	// 使用阿里DNS
	da.Arg("dns", "223.5.5.5")
	// 执行代码检查命令
	mark := d.docker.Mark
	_, err = d.base.Command(d.docker.Daemon).Context(*ctx).Args(da.Build()).Checker().Contains(mark).Build().Exec()

	return
}

func (d *Daemon) check(ctx *context.Context) (err error) {
	ia := args.New().Build().Subcommand("info").Build()
	_, err = d.base.Command(d.docker.Exe).Context(*ctx).Args(ia).Build().Exec()

	return
}

func (d *Daemon) mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(defaultMirrors))
	if d.base.Default() {
		mirrors = append(mirrors, defaultMirrors...)
	}
	mirrors = append(mirrors, d.docker.Mirrors...)

	return
}

func (d *Daemon) address() string {
	builder := new(strings.Builder)
	builder.WriteString(d.docker.Protocol)
	builder.WriteString(constant.Colon)
	builder.WriteString(constant.Slash)
	builder.WriteString(constant.Slash)
	builder.WriteString(d.docker.Host)

	return builder.String()
}
