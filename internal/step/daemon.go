package step

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
	"https://dockerproxy.com",
	"https://ustc-edu-cn.mirror.aliyuncs.com",
	"https://mirror.baidubce.com",
	"https://hub.daocloud.io",
	"https://mirror.ccs.tencentyun.com",
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
		time.Sleep(5 * time.Second)
	}
	if nil == err { // 注入上下文，供后续步骤使用
		*ctx = context.WithValue(*ctx, constant.KeyDir, d.dir())
	}

	return
}

func (d *Daemon) startup(_ *context.Context) (err error) {
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
	_, err = d.base.Command(constant.DaemonExe).Args(da.Build()).Checker().Contains(d.docker.Mark).Dir(d.dir()).Build().Exec()

	return
}

func (d *Daemon) check(_ *context.Context) (err error) {
	ia := args.New().Build().Subcommand("info").Build()
	_, err = d.base.Command(constant.Exe).Args(ia).Dir(d.dir()).Build().Exec()

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

func (d *Daemon) dir() (context string) {
	if "" == d.docker.Context {
		context = filepath.Dir(d.docker.Dockerfile)
	} else {
		context = d.docker.Context
	}

	return
}
