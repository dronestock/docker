package step

import (
	"context"
	_ "embed"
	"encoding/json"

	"github.com/dronestock/docker/internal/internal/command"
	"github.com/dronestock/docker/internal/internal/config"
	"github.com/dronestock/docker/internal/internal/step/internal/data"
	"github.com/goexl/args"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/gox/template"
	"github.com/goexl/log"
	"github.com/rs/xid"
)

//go:embed internal/template/buildkitd.toml.gohtml
var kitTemplate string

type Setup struct {
	command *command.Docker
	config  *config.Docker
	targets *config.Targets
}

func NewSetup(command *command.Docker, config *config.Docker, targets *config.Targets) *Setup {
	return &Setup{
		command: command,
		config:  config,
		targets: targets,
	}
}

func (s *Setup) Runnable() bool {
	return s.targets.BinfmtNeedSetup(s.config) || s.targets.DriverNeedSetup()
}

func (s *Setup) Run(ctx *context.Context) (err error) {
	if qe := s.binfmt(ctx); nil != qe {
		err = qe
	} else if de := s.driver(ctx); nil != de {
		err = de
	}

	return
}

func (s *Setup) binfmt(ctx *context.Context) (err error) {
	if !s.targets.BinfmtNeedSetup(s.config) {
		return
	}

	image := "tonistiigi/binfmt"
	arguments := args.New().Build()
	arguments.Subcommand("run")
	arguments.Flag("privileged")
	arguments.Flag("rm")
	arguments.Subcommand(image)
	arguments.Flag("install")
	arguments.Subcommand("all")

	fields := gox.Fields[any]{
		field.New("image", image),
	}
	s.command.Info("准备安装Qemu环境", fields...)
	if err = s.command.Exec(*ctx, arguments.Build()); nil != err {
		s.command.Warn("安装Qemu环境失败", fields.Add(field.Error(err))...)
	} else {
		s.command.Info("安装Qemu环境成功", fields...)
	}

	return
}

func (s *Setup) driver(ctx *context.Context) (err error) {
	if !s.targets.DriverNeedSetup() {
		return
	}

	platforms := s.targets.Platforms()
	name := xid.New().String()
	arguments := args.New().Build()
	arguments.Subcommand("buildx", "create")
	arguments.Argument("name", name)
	arguments.Flag("use")
	arguments.Argument("platform", platforms)
	arguments.Argument("driver", "docker-container")
	arguments.Argument("driver-opt", "network=host")

	fields := gox.Fields[any]{
		field.New("name", name),
		field.New("platform", platforms),
	}
	s.command.Info("准备创建多平台编译驱动", fields...)
	if wce := s.writeConfig(ctx, arguments.Build()); nil != wce {
		err = wce
	} else if ee := s.command.Exec(*ctx, arguments.Build()); nil != ee {
		err = ee
		s.command.Warn("创建多平台编译驱动失败", fields.Add(field.Error(err))...)
	} else {
		s.command.Info("创建多平台编译驱动成功", fields...)
	}

	return
}

func (s *Setup) writeConfig(_ *context.Context, arguments *args.Arguments) (err error) {
	kit := new(data.Buildkitd)
	kit.Debug = s.command.Enabled(log.LevelDebug)

	path := "buildkitd.toml"
	if bytes, me := json.Marshal(s.command.Mirrors()); nil != me {
		err = me
	} else {
		kit.Mirrors = string(bytes)
		err = template.New(kitTemplate).Data(kit).Build().ToFile(path)
	}
	if nil == err {
		arguments.Rebuild().Argument("config", path)
		s.command.Cleanup().File(path).Build()
	}

	return
}
