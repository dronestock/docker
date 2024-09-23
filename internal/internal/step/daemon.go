package step

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/constant"
	"github.com/goexl/args"
	"github.com/goexl/gox/field"
)

type Daemon struct {
	config  *config.Docker
	command *command.Docker

	defaultMirrors []string
}

func NewDaemon(command *command.Docker, config *config.Docker) *Daemon {
	return &Daemon{
		command: command,
		config:  config,

		defaultMirrors: []string{
			"https://jockerhub.com",
			"https://hub.uuuadc.top",
			"https://docker.anyhub.us.kg",
			"https://dockerhub.jobcher.com",
			"https://dockerhub.icu",
			"https://docker.ckyl.me",
			"https://docker.awsl9527.cn",
			"https://hub.20240220.xyz",
		},
	}
}

func (d *Daemon) Runnable() bool {
	return constant.Unix == d.config.Protocol
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
			d.command.Info(
				"等待Docker启动完成",
				field.New("timeout", d.command.Timeout.Truncate(time.Second).String()),
				field.New("elapsed", d.command.Elapsed()),
				field.New("times", times),
				field.New("address", d.address()),
			)
		}
		time.Sleep(1 * time.Second)
	}

	return
}

func (d *Daemon) startup(ctx *context.Context) (err error) {
	arguments := args.New().Build()
	arguments.Argument("data-root", d.config.Data)
	arguments.Argument("host", d.address())

	if _, se := os.Stat("/etc/config/default.json"); nil == se {
		arguments.Argument("seccomp-profile", "/etc/config/default.json")
	}

	// 驱动
	if "" != d.config.Driver {
		arguments.Argument("storage-driver", d.config.Driver)
	}
	// 镜像加速
	for _, mirror := range d.mirrors() {
		arguments.Argument("registry-mirror", mirror)
	}

	// 启用实验性功能
	if d.config.Experimental {
		arguments.Flag("experimental")
	}

	// 启动后台进程
	err = d.command.Daemon(ctx, arguments.Build(), d.config.Mark)

	return
}

func (d *Daemon) check(ctx *context.Context) error {
	return d.command.Exec(*ctx, args.New().Build().Subcommand("info").Build())
}

func (d *Daemon) mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(d.defaultMirrors))
	if d.command.Default() {
		mirrors = append(mirrors, d.defaultMirrors...)
	}
	mirrors = append(mirrors, d.config.Mirrors...)

	return
}

func (d *Daemon) address() string {
	builder := new(strings.Builder)
	builder.WriteString(d.config.Protocol)
	builder.WriteString(constant.Colon)
	builder.WriteString(constant.Slash)
	builder.WriteString(constant.Slash)
	builder.WriteString(d.config.Host)

	return builder.String()
}
