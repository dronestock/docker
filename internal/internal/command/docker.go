package command

import (
	"context"

	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/key"
	"github.com/dronestock/drone"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

type Docker struct {
	*drone.Base

	binary *config.Binary
	config *config.Docker

	mirrors []string
}

func NewDocker(base *drone.Base, binary *config.Binary, config *config.Docker) *Docker {
	return &Docker{
		Base: base,

		binary: binary,
		config: config,

		mirrors: []string{
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

func (d *Docker) Exec(ctx context.Context, arguments *args.Arguments) (err error) {
	command := d.Command(d.binary.Docker).Context(ctx)
	command.Arguments(arguments)
	// 检查是否要通过输出判断退出
	if mark := ctx.Value(key.ContextMark); nil != mark {
		command.Async().Check().Contains(mark.(string)).Build()
	}
	// 确定运行上下文
	if dir := ctx.Value(key.ContextDir); nil != dir {
		command.Dir(dir.(string))
	}

	fields := gox.Fields[any]{
		field.New("binary", d.binary.Docker),
		field.New("arguments", arguments.Cli()),
	}
	d.Debug("命令执行开始", fields...)
	if _, err = command.Build().Exec(); nil != err {
		d.Warn("命令执行出错", fields.Add(field.Error(err))...)
	} else {
		d.Debug("命令执行成功", fields...)
	}

	return
}

func (d *Docker) Daemon(ctx *context.Context, arguments *args.Arguments, mark string) (err error) {
	if d.Verbose {
		arguments = arguments.Rebuild().Flag("debug").Option("log-level", "debug").Build()
	}

	command := d.Command(d.binary.Daemon).Context(*ctx).Async()
	command.Arguments(arguments)
	// 检查是否完成
	command.Check().Contains(mark).Build()

	fields := gox.Fields[any]{
		field.New("binary", d.binary.Daemon),
		field.New("arguments", arguments.Cli()),
	}
	d.Debug("命令执行开始", fields...)
	if _, err = command.Build().Exec(); nil != err {
		d.Warn("命令执行出错", fields.Add(field.Error(err))...)
	} else {
		d.Debug("命令执行成功", fields...)
	}

	return
}

func (d *Docker) Mirrors() (mirrors []string) {
	mirrors = make([]string, 0, len(d.mirrors))
	if d.Default() {
		mirrors = append(mirrors, d.mirrors...)
	}
	mirrors = append(mirrors, d.config.Mirrors...)

	return
}
