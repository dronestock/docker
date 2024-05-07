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
}

func NewDocker(base *drone.Base, binary *config.Binary) *Docker {
	return &Docker{
		Base: base,

		binary: binary,
	}
}

func (d *Docker) Exec(ctx context.Context, arguments *args.Arguments) (err error) {
	command := d.Command(d.binary.Docker).Context(ctx)
	command.Args(arguments)
	// 检查是否要通过输出判断退出
	if mark := ctx.Value(key.ContextMark); nil != mark {
		command.Async().Checker().Contains(mark.(string))
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
	command.Args(arguments)
	// 检查是否完成
	command.Checker().Contains(mark)

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
